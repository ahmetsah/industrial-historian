#![allow(dead_code)]
use crate::storage::rocksdb::RocksDBStorage;
use crate::storage::tiered::policy::TieringPolicy;
use crate::storage::tiered::s3::S3Client;
use anyhow::Result;
use std::sync::Arc;
use std::time::Duration;
use tokio::time;

pub struct TieringJob {
    storage: Arc<RocksDBStorage>,
    s3_client: S3Client,
    policy: TieringPolicy,
    interval_ms: u64,
}

impl TieringJob {
    pub fn new(
        storage: Arc<RocksDBStorage>,
        s3_client: S3Client,
        policy: TieringPolicy,
        interval_ms: u64,
    ) -> Self {
        Self {
            storage,
            s3_client,
            policy,
            interval_ms,
        }
    }

    pub async fn run(&self) {
        let mut interval = time::interval(Duration::from_millis(self.interval_ms));

        loop {
            interval.tick().await;
            if let Err(e) = self.execute_tiering().await {
                tracing::error!("Tiering job failed: {}", e);
            }
        }
    }

    async fn execute_tiering(&self) -> Result<()> {
        let now = std::time::SystemTime::now()
            .duration_since(std::time::UNIX_EPOCH)?
            .as_millis() as i64;
        let threshold = now - self.policy.max_age_ms;

        // Process in batches
        let batch_size = 100;
        let old_data = self.storage.get_old_data(threshold, batch_size)?;

        if old_data.is_empty() {
            return Ok(());
        }

        tracing::info!("Found {} chunks to tier", old_data.len());

        for (sensor_id, timestamp, data) in old_data {
            let key = format!("{}/{}.bin", sensor_id, timestamp);

            match self.s3_client.put_object(&key, &data).await {
                Ok(_) => {
                    // CRITICAL: Record metadata BEFORE deleting local data to ensure reachability
                    if let Err(e) = self
                        .storage
                        .record_tiered_metadata(&sensor_id, timestamp, &key)
                    {
                        tracing::error!("Failed to record metadata for {}: {}", key, e);
                        // Do not delete if metadata record failed
                        continue;
                    }

                    if let Err(e) = self.storage.delete_data(&sensor_id, timestamp) {
                        tracing::error!("Failed to delete tiered data for {}: {}", key, e);
                    } else {
                        tracing::debug!("Tiered {}", key);
                    }
                }
                Err(e) => {
                    tracing::error!("Failed to upload {}: {}", key, e);
                }
            }
        }

        Ok(())
    }
}
