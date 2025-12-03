#![allow(dead_code)]
use anyhow::Result;
use historian_core::SensorData;

#[async_trait::async_trait]
pub trait StorageEngine: Send + Sync {
    fn write(&self, data: &SensorData) -> Result<()>;
    async fn read(&self, sensor_id: &str, timestamp: i64) -> Result<Option<f64>>;
    async fn scan(&self, sensor_id: &str, start_ts: i64, end_ts: i64) -> Result<Vec<SensorData>>;
}

pub mod compression;
pub mod rocksdb;
pub mod tiered;
