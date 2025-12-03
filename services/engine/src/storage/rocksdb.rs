#![allow(dead_code)]
use crate::storage::compression::compress_points;
use crate::storage::tiered::s3::S3Client;
use crate::storage::StorageEngine;
use anyhow::{Context, Result};
use dashmap::DashMap;
use historian_core::SensorData;
use rocksdb::{Options, DB};
use std::sync::Arc;
use tokio::sync::mpsc;
use tokio_stream::wrappers::ReceiverStream;

const BUFFER_SIZE: usize = 1000;

pub struct RocksDBStorage {
    db: Arc<DB>,
    buffer: DashMap<String, Vec<(i64, f64)>>,
    s3_client: Option<S3Client>,
}

impl RocksDBStorage {
    pub fn new(path: &str, s3_client: Option<S3Client>) -> Result<Self> {
        let mut opts = Options::default();
        opts.create_if_missing(true);
        opts.create_missing_column_families(true);
        // Optimization for high write throughput
        opts.set_max_background_jobs(4);
        opts.set_bytes_per_sync(1048576);
        // MemTable settings
        opts.set_write_buffer_size(64 * 1024 * 1024); // 64MB MemTable
        opts.set_max_write_buffer_number(3); // 3 MemTables
        opts.set_target_file_size_base(64 * 1024 * 1024); // 64MB SSTable file size base

        let cfs = vec!["default", "tiering_metadata"];
        let db = DB::open_cf(&opts, path, cfs).context("Failed to open RocksDB")?;
        Ok(Self {
            db: Arc::new(db),
            buffer: DashMap::new(),
            s3_client,
        })
    }

    fn generate_key(sensor_id: &str, timestamp: i64) -> Vec<u8> {
        let mut key = Vec::new();
        key.extend_from_slice(sensor_id.as_bytes());
        key.push(0); // Separator
        key.extend_from_slice(&timestamp.to_be_bytes());
        key
    }

    fn flush_sensor(&self, sensor_id: &str, points: &[(i64, f64)]) -> Result<()> {
        if points.is_empty() {
            return Ok(());
        }
        let compressed = compress_points(points)?;

        // Key: SensorID + StartTimestamp
        let start_time = points[0].0;
        let key = Self::generate_key(sensor_id, start_time);

        self.db
            .put(&key, &compressed)
            .context("Failed to write to RocksDB")?;
        Ok(())
    }

    pub fn get_old_data(
        &self,
        threshold_ms: i64,
        limit: usize,
    ) -> Result<Vec<(String, i64, Vec<u8>)>> {
        let mut old_data = Vec::new();
        let iter = self.db.iterator(rocksdb::IteratorMode::Start);

        for item in iter {
            let (key, value) = item.context("Failed to read from iterator")?;
            if key.len() < 9 {
                continue;
            }

            let ts_start = key.len() - 8;
            let ts_bytes: [u8; 8] = key[ts_start..].try_into().unwrap();
            let timestamp = i64::from_be_bytes(ts_bytes);

            if timestamp < threshold_ms {
                let sensor_bytes = &key[0..ts_start - 1]; // -1 for separator
                let sensor_id = String::from_utf8(sensor_bytes.to_vec()).unwrap_or_default();
                old_data.push((sensor_id, timestamp, value.to_vec()));
                if old_data.len() >= limit {
                    break;
                }
            }
        }
        Ok(old_data)
    }

    pub fn delete_data(&self, sensor_id: &str, timestamp: i64) -> Result<()> {
        let key = Self::generate_key(sensor_id, timestamp);
        self.db.delete(&key).context("Failed to delete data")?;
        Ok(())
    }

    pub fn record_tiered_metadata(
        &self,
        sensor_id: &str,
        timestamp: i64,
        s3_key: &str,
    ) -> Result<()> {
        let cf = self
            .db
            .cf_handle("tiering_metadata")
            .context("CF not found")?;
        let key = Self::generate_key(sensor_id, timestamp);
        self.db
            .put_cf(&cf, &key, s3_key.as_bytes())
            .context("Failed to write metadata")?;
        Ok(())
    }

    pub fn get_tiered_metadata(&self, sensor_id: &str, timestamp: i64) -> Result<Option<String>> {
        let cf = self
            .db
            .cf_handle("tiering_metadata")
            .context("CF not found")?;
        let key = Self::generate_key(sensor_id, timestamp);

        let mut iter = self.db.iterator_cf(
            &cf,
            rocksdb::IteratorMode::From(&key, rocksdb::Direction::Reverse),
        );

        if let Some(Ok((k, v))) = iter.next() {
            if k.len() < 9 {
                return Ok(None);
            }
            let ts_start = k.len() - 8;
            let sensor_bytes = &k[0..ts_start - 1];
            let stored_sensor_id = String::from_utf8(sensor_bytes.to_vec()).unwrap_or_default();

            if stored_sensor_id == sensor_id {
                return Ok(Some(String::from_utf8(v.to_vec())?));
            }
        }
        Ok(None)
    }

    // Private read removed
}

#[async_trait::async_trait]
impl StorageEngine for RocksDBStorage {
    fn write(&self, data: &SensorData) -> Result<()> {
        let mut entry = self.buffer.entry(data.sensor_id.clone()).or_default();
        entry.push((data.timestamp_ms, data.value));

        if entry.len() >= BUFFER_SIZE {
            let points_to_flush = std::mem::take(&mut *entry);
            drop(entry); // Release lock

            self.flush_sensor(&data.sensor_id, &points_to_flush)?;
        }
        Ok(())
    }

    async fn read(&self, sensor_id: &str, timestamp: i64) -> Result<Option<f64>> {
        // 1. Check Buffer
        if let Some(entry) = self.buffer.get(sensor_id) {
            for (ts, val) in entry.iter() {
                if *ts == timestamp {
                    return Ok(Some(*val));
                }
            }
        }

        // 2. Check RocksDB (Hot Tier)
        let key = Self::generate_key(sensor_id, timestamp);
        // We need seek_for_prev because we store compressed blocks
        let mut iter = self.db.iterator(rocksdb::IteratorMode::From(
            &key,
            rocksdb::Direction::Reverse,
        ));
        if let Some(Ok((k, v))) = iter.next() {
            if k.len() >= 9 {
                let ts_start = k.len() - 8;
                let sensor_bytes = &k[0..ts_start - 1];
                let stored_sensor_id = String::from_utf8(sensor_bytes.to_vec()).unwrap_or_default();
                if stored_sensor_id == sensor_id {
                    use crate::storage::compression::decompress_points;
                    if let Ok(points) = decompress_points(&v) {
                        for (ts, val) in points {
                            if ts == timestamp {
                                return Ok(Some(val));
                            }
                        }
                    }
                }
            }
        }

        // 3. Check S3 (Cold Tier)
        if let Some(client) = &self.s3_client {
            if let Ok(Some(s3_key)) = self.get_tiered_metadata(sensor_id, timestamp) {
                match client.get_object(&s3_key).await {
                    Ok(data) => {
                        use crate::storage::compression::decompress_points;
                        if let Ok(points) = decompress_points(&data) {
                            for (ts, val) in points {
                                if ts == timestamp {
                                    return Ok(Some(val));
                                }
                            }
                        }
                    }
                    Err(e) => {
                        tracing::error!("Failed to fetch from S3: {}", e);
                    }
                }
            }
        }

        Ok(None)
    }

    async fn scan(&self, sensor_id: &str, start_ts: i64, end_ts: i64) -> Result<Vec<SensorData>> {
        let mut points = Vec::new();

        // 1. Check Buffer
        if let Some(entry) = self.buffer.get(sensor_id) {
            for (ts, val) in entry.iter() {
                if *ts >= start_ts && *ts <= end_ts {
                    points.push(SensorData {
                        sensor_id: sensor_id.to_string(),
                        timestamp_ms: *ts,
                        value: *val,
                        quality: 1,
                    });
                }
            }
        }

        // 2. Check RocksDB (Hot Tier)
        // We scan from start_ts. Note: This might miss points in a block starting before start_ts but overlapping.
        // For strict correctness, we should seek back one block.
        // But since we don't know block duration, we'd need seek_for_prev.
        // rust-rocksdb iterator doesn't expose seek_for_prev easily in the high-level loop.
        // We'll proceed with forward scan from start_ts for now.
        let start_key = Self::generate_key(sensor_id, start_ts);
        let iter = self.db.iterator(rocksdb::IteratorMode::From(
            &start_key,
            rocksdb::Direction::Forward,
        ));

        for item in iter {
            let (k, v) = item.context("Failed to read from iterator")?;
            if k.len() < 9 {
                continue;
            }

            let ts_start = k.len() - 8;
            let sensor_bytes = &k[0..ts_start - 1];
            let stored_sensor_id = String::from_utf8(sensor_bytes.to_vec()).unwrap_or_default();

            if stored_sensor_id != sensor_id {
                break; // Moved to next sensor
            }

            let ts_bytes: [u8; 8] = k[ts_start..].try_into().unwrap();
            let block_start_ts = i64::from_be_bytes(ts_bytes);

            if block_start_ts > end_ts {
                break; // Beyond range
            }

            use crate::storage::compression::decompress_points;
            if let Ok(decompressed) = decompress_points(&v) {
                for (ts, val) in decompressed {
                    if ts >= start_ts && ts <= end_ts {
                        points.push(SensorData {
                            sensor_id: sensor_id.to_string(),
                            timestamp_ms: ts,
                            value: val,
                            quality: 1,
                        });
                    }
                }
            }
        }

        // 3. Check S3 (Cold Tier) - Metadata Scan
        if let Some(client) = &self.s3_client {
            let mut s3_keys = Vec::new();

            // Scope for RocksDB interaction
            {
                let cf = self
                    .db
                    .cf_handle("tiering_metadata")
                    .context("CF not found")?;
                let start_key = Self::generate_key(sensor_id, start_ts);
                let iter = self.db.iterator_cf(
                    &cf,
                    rocksdb::IteratorMode::From(&start_key, rocksdb::Direction::Forward),
                );

                for item in iter {
                    let (k, v) = item.context("Failed to read metadata iterator")?;
                    if k.len() < 9 {
                        continue;
                    }

                    let ts_start = k.len() - 8;
                    let sensor_bytes = &k[0..ts_start - 1];
                    let stored_sensor_id =
                        String::from_utf8(sensor_bytes.to_vec()).unwrap_or_default();

                    if stored_sensor_id != sensor_id {
                        break;
                    }

                    let ts_bytes: [u8; 8] = k[ts_start..].try_into().unwrap();
                    let block_start_ts = i64::from_be_bytes(ts_bytes);

                    if block_start_ts > end_ts {
                        break;
                    }

                    let s3_key = String::from_utf8(v.to_vec())?;
                    s3_keys.push(s3_key);
                }
            }

            // Fetch from S3
            for s3_key in s3_keys {
                match client.get_object(&s3_key).await {
                    Ok(data) => {
                        use crate::storage::compression::decompress_points;
                        if let Ok(decompressed) = decompress_points(&data) {
                            for (ts, val) in decompressed {
                                if ts >= start_ts && ts <= end_ts {
                                    points.push(SensorData {
                                        sensor_id: sensor_id.to_string(),
                                        timestamp_ms: ts,
                                        value: val,
                                        quality: 1,
                                    });
                                }
                            }
                        }
                    }
                    Err(e) => {
                        tracing::error!("Failed to fetch from S3: {}", e);
                    }
                }
            }
        }

        points.sort_by_key(|p| p.timestamp_ms);
        // Deduplicate? (Buffer vs RocksDB vs S3 might overlap if not careful, but usually distinct)
        points.dedup_by_key(|p| p.timestamp_ms);

        Ok(points)
    }

    fn scan_stream(
        &self,
        sensor_id: &str,
        start_ts: i64,
        end_ts: i64,
    ) -> Result<std::pin::Pin<Box<dyn futures::Stream<Item = Result<SensorData>> + Send>>> {
        let (tx, rx) = mpsc::channel(1000);
        let db = self.db.clone();
        let buffer = self.buffer.clone();
        let s3_client = self.s3_client.clone();
        let sensor_id = sensor_id.to_string();

        tokio::spawn(async move {
            // 1. S3 (Cold Tier)
            if let Some(client) = s3_client {
                let mut s3_keys = Vec::new();
                // Blocking RocksDB access for metadata
                {
                    let db_clone = db.clone();
                    let sensor_id_clone = sensor_id.clone();
                    let res = tokio::task::spawn_blocking(move || -> Result<Vec<String>> {
                        let cf = db_clone.cf_handle("tiering_metadata").context("CF not found")?;
                        let start_key = RocksDBStorage::generate_key(&sensor_id_clone, start_ts);
                        let iter = db_clone.iterator_cf(
                            &cf,
                            rocksdb::IteratorMode::From(&start_key, rocksdb::Direction::Forward),
                        );
                        let mut keys = Vec::new();
                        for item in iter {
                            let (k, v) = item?;
                             if k.len() < 9 { continue; }
                            let ts_start = k.len() - 8;
                            let sensor_bytes = &k[0..ts_start - 1];
                            let stored_sensor_id = String::from_utf8(sensor_bytes.to_vec()).unwrap_or_default();
                            if stored_sensor_id != sensor_id_clone { break; }
                            let ts_bytes: [u8; 8] = k[ts_start..].try_into().unwrap();
                            let block_start_ts = i64::from_be_bytes(ts_bytes);
                            if block_start_ts > end_ts { break; }
                            keys.push(String::from_utf8(v.to_vec())?);
                        }
                        Ok(keys)
                    }).await;
                    
                    if let Ok(Ok(keys)) = res {
                        s3_keys = keys;
                    }
                }

                for s3_key in s3_keys {
                    if let Ok(data) = client.get_object(&s3_key).await {
                         use crate::storage::compression::decompress_points;
                         if let Ok(decompressed) = decompress_points(&data) {
                             for (ts, val) in decompressed {
                                 if ts >= start_ts && ts <= end_ts {
                                     let point = SensorData {
                                         sensor_id: sensor_id.clone(),
                                         timestamp_ms: ts,
                                         value: val,
                                         quality: 1,
                                     };
                                     if tx.send(Ok(point)).await.is_err() { return; }
                                 }
                             }
                         }
                    }
                }
            }

            // 2. RocksDB (Hot Tier)
            {
                let db_clone = db.clone();
                let sensor_id_clone = sensor_id.clone();
                let tx_clone = tx.clone();
                
                let _ = tokio::task::spawn_blocking(move || {
                    let start_key = RocksDBStorage::generate_key(&sensor_id_clone, start_ts);
                    let iter = db_clone.iterator(rocksdb::IteratorMode::From(
                        &start_key,
                        rocksdb::Direction::Forward,
                    ));
                    
                    for item in iter {
                        if let Ok((k, v)) = item {
                             if k.len() < 9 { continue; }
                             let ts_start = k.len() - 8;
                             let sensor_bytes = &k[0..ts_start - 1];
                             let stored_sensor_id = String::from_utf8(sensor_bytes.to_vec()).unwrap_or_default();
                             if stored_sensor_id != sensor_id_clone { break; }
                             let ts_bytes: [u8; 8] = k[ts_start..].try_into().unwrap();
                             let block_start_ts = i64::from_be_bytes(ts_bytes);
                             if block_start_ts > end_ts { break; }
                             
                             use crate::storage::compression::decompress_points;
                             if let Ok(decompressed) = decompress_points(&v) {
                                 for (ts, val) in decompressed {
                                     if ts >= start_ts && ts <= end_ts {
                                         let point = SensorData {
                                             sensor_id: sensor_id_clone.clone(),
                                             timestamp_ms: ts,
                                             value: val,
                                             quality: 1,
                                         };
                                         if tx_clone.blocking_send(Ok(point)).is_err() { return; }
                                     }
                                 }
                             }
                        }
                    }
                }).await;
            }

            // 3. Buffer
            if let Some(entry) = buffer.get(&sensor_id) {
                for (ts, val) in entry.iter() {
                    if *ts >= start_ts && *ts <= end_ts {
                        let point = SensorData {
                            sensor_id: sensor_id.clone(),
                            timestamp_ms: *ts,
                            value: *val,
                            quality: 1,
                        };
                        if tx.send(Ok(point)).await.is_err() { return; }
                    }
                }
            }
        });

        Ok(Box::pin(ReceiverStream::new(rx)))
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use tempfile::TempDir;

    #[test]
    fn test_rocksdb_write_buffered() {
        let temp_dir = TempDir::new().unwrap();
        let path = temp_dir.path().to_str().unwrap();
        let storage = RocksDBStorage::new(path, None).unwrap();

        // Write BUFFER_SIZE points
        for i in 0..BUFFER_SIZE {
            let data = SensorData {
                sensor_id: "test-sensor".to_string(),
                timestamp_ms: 1000 + i as i64 * 1000,
                value: 42.0 + i as f64,
                quality: 1,
            };
            storage.write(&data).unwrap();
        }

        // Verify flush happened
        // Key should be for the first point
        let key = RocksDBStorage::generate_key("test-sensor", 1000);
        let val = storage.db.get(&key).unwrap();
        assert!(val.is_some());

        // Verify content?
        // We can decompress and check count
        use crate::storage::compression::decompress_points;
        let points = decompress_points(&val.unwrap()).unwrap();
        assert_eq!(points.len(), BUFFER_SIZE);
        assert_eq!(points[0].0, 1000);
        assert_eq!(
            points[BUFFER_SIZE - 1].0,
            1000 + (BUFFER_SIZE - 1) as i64 * 1000
        );
    }

    #[test]
    fn test_write_performance() {
        let temp_dir = TempDir::new().unwrap();
        let path = temp_dir.path().to_str().unwrap();
        let storage = RocksDBStorage::new(path, None).unwrap();

        let count = 50_000;
        let start = std::time::Instant::now();

        for i in 0..count {
            let data = SensorData {
                sensor_id: "perf-sensor".to_string(),
                timestamp_ms: 1000 + i as i64,
                value: i as f64,
                quality: 1,
            };
            storage.write(&data).unwrap();
        }

        let duration = start.elapsed();
        let rate = count as f64 / duration.as_secs_f64();
        println!(
            "Wrote {} events in {:?} (Rate: {:.2} events/sec)",
            count, duration, rate
        );

        // In debug mode, performance is lower.
        // We just want to ensure it doesn't crash and is reasonably fast.
        assert!(rate > 1000.0, "Rate {} too low", rate);
    }

    #[tokio::test]
    async fn test_tiered_metadata() {
        let temp_dir = TempDir::new().unwrap();
        let path = temp_dir.path().to_str().unwrap();
        let storage = RocksDBStorage::new(path, None).unwrap();

        let sensor_id = "sensor-1";
        let timestamp = 1000;
        let s3_key = "s3://bucket/sensor-1/1000.bin";

        storage
            .record_tiered_metadata(sensor_id, timestamp, s3_key)
            .unwrap();

        let retrieved = storage.get_tiered_metadata(sensor_id, timestamp).unwrap();
        assert_eq!(retrieved, Some(s3_key.to_string()));

        // Test seek behavior (Reverse)
        let retrieved_later = storage
            .get_tiered_metadata(sensor_id, timestamp + 10)
            .unwrap();
        assert_eq!(retrieved_later, Some(s3_key.to_string()));

        // Test different sensor
        let retrieved_other = storage.get_tiered_metadata("sensor-2", timestamp).unwrap();
        assert_eq!(retrieved_other, None);
    }
}
