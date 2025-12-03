use anyhow::Result;
use historian_core::SensorData;
use prost::Message;
use std::collections::VecDeque;
use std::path::PathBuf;
use tokio::fs::{self, OpenOptions};
use tokio::io::{AsyncReadExt, AsyncWriteExt};

pub struct RingBuffer<T> {
    buffer: VecDeque<T>,
    capacity: usize,
}

impl<T> RingBuffer<T> {

    #[allow(dead_code)]
    pub fn new(capacity: usize) -> Self {
        Self {
            buffer: VecDeque::with_capacity(capacity),
            capacity,
        }
    }

    pub fn push(&mut self, item: T) -> Option<T> {
        let mut popped = None;
        if self.buffer.len() >= self.capacity {
            popped = self.buffer.pop_front();
        }
        self.buffer.push_back(item);
        popped
    }

    pub fn pop(&mut self) -> Option<T> {
        self.buffer.pop_front()
    }

    #[allow(dead_code)]
    pub fn is_empty(&self) -> bool {
        self.buffer.is_empty()
    }

    #[allow(dead_code)]
    pub fn len(&self) -> usize {
        self.buffer.len()
    }
}

pub struct DiskBuffer {
    file_path: PathBuf,
}

impl DiskBuffer {
    pub fn new(path: impl Into<PathBuf>) -> Self {
        Self {
            file_path: path.into(),
        }
    }

    pub async fn push(&self, item: &SensorData) -> Result<()> {
        let mut file = OpenOptions::new()
            .create(true)
            .append(true)
            .open(&self.file_path)
            .await?;

        let mut buf = Vec::new();
        item.encode(&mut buf)?;

        let len = buf.len() as u32;
        file.write_u32(len).await?;
        file.write_all(&buf).await?;
        Ok(())
    }

    pub async fn read_all_and_clear(&self) -> Result<Vec<SensorData>> {
        // TODO: Implement chunked reading or streaming to avoid loading entire file into memory
        // if the WAL file becomes very large (e.g., > 100MB).
        if !self.file_path.exists() {
            return Ok(Vec::new());
        }

        let mut file = fs::File::open(&self.file_path).await?;
        let mut content = Vec::new();
        file.read_to_end(&mut content).await?;

        let mut items = Vec::new();
        let mut cursor = std::io::Cursor::new(content);
        //use std::io::Read;

        while cursor.position() < cursor.get_ref().len() as u64 {
            // Read length (u32)
            let mut len_buf = [0u8; 4];
            if std::io::Read::read_exact(&mut cursor, &mut len_buf).is_err() {
                break;
            }
            let len = u32::from_be_bytes(len_buf) as usize;

            // Read protobuf message
            let mut msg_buf = vec![0u8; len];
            if std::io::Read::read_exact(&mut cursor, &mut msg_buf).is_err() {
                break;
            }

            let item = SensorData::decode(&msg_buf[..])?;
            items.push(item);
        }

        fs::remove_file(&self.file_path).await?;

        Ok(items)
    }
}

pub struct HybridBuffer {
    mem_buffer: RingBuffer<SensorData>,
    disk_buffer: DiskBuffer,
}

impl HybridBuffer {
    pub fn new(mem_capacity: usize, disk_path: impl Into<PathBuf>) -> Self {
        Self {
            mem_buffer: RingBuffer::new(mem_capacity),
            disk_buffer: DiskBuffer::new(disk_path),
        }
    }

    pub async fn push(&mut self, item: SensorData) -> Result<()> {
        // Push to memory first
        if let Some(overflow_item) = self.mem_buffer.push(item) {
            // If memory is full, the oldest item pops out.
            // We write this oldest item to disk.
            self.disk_buffer.push(&overflow_item).await?;
        }
        Ok(())
    }

    pub fn pop_mem(&mut self) -> Option<SensorData> {
        self.mem_buffer.pop()
    }

    pub async fn flush_disk(&self) -> Result<Vec<SensorData>> {
        self.disk_buffer.read_all_and_clear().await
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_ring_buffer_overflow() {
        let mut rb = RingBuffer::new(2);
        assert_eq!(rb.push(1), None);
        assert_eq!(rb.push(2), None);
        assert_eq!(rb.push(3), Some(1)); // Should pop 1
        assert_eq!(rb.len(), 2);
        assert_eq!(rb.pop(), Some(2));
        assert_eq!(rb.pop(), Some(3));
    }

    #[tokio::test]
    async fn test_disk_buffer() {
        let path = "test_buffer.wal";
        let db = DiskBuffer::new(path);

        // Clean up before test
        let _ = fs::remove_file(path).await;

        let item1 = SensorData {
            sensor_id: "s1".to_string(),
            value: 1.0,
            timestamp_ms: 100,
            quality: 1,
        };
        let item2 = SensorData {
            sensor_id: "s2".to_string(),
            value: 2.0,
            timestamp_ms: 200,
            quality: 1,
        };

        db.push(&item1).await.unwrap();
        db.push(&item2).await.unwrap();

        let items = db.read_all_and_clear().await.unwrap();
        assert_eq!(items.len(), 2);
        assert_eq!(items[0].sensor_id, "s1");
        assert_eq!(items[1].sensor_id, "s2");

        assert!(!std::path::Path::new(path).exists());
    }

    #[tokio::test]
    async fn test_hybrid_buffer_spill() {
        let path = "test_hybrid.wal";
        let _ = fs::remove_file(path).await;

        // Capacity 1. Push 2 items. 1st should spill to disk.
        let mut hb = HybridBuffer::new(1, path);

        let item1 = SensorData {
            sensor_id: "s1".to_string(),
            value: 1.0,
            timestamp_ms: 100,
            quality: 1,
        };
        let item2 = SensorData {
            sensor_id: "s2".to_string(),
            value: 2.0,
            timestamp_ms: 200,
            quality: 1,
        };

        hb.push(item1.clone()).await.unwrap(); // Buffer: [s1]
        hb.push(item2.clone()).await.unwrap(); // Buffer: [s2], Disk: [s1]

        // Check memory
        let popped = hb.pop_mem().unwrap();
        assert_eq!(popped.sensor_id, "s2");
        assert!(hb.pop_mem().is_none());

        // Check disk
        let disk_items = hb.flush_disk().await.unwrap();
        assert_eq!(disk_items.len(), 1);
        assert_eq!(disk_items[0].sensor_id, "s1");

        let _ = fs::remove_file(path).await;
    }
}
