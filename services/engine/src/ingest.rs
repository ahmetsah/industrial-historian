use crate::storage::StorageEngine;
use anyhow::{Context, Result};
use async_nats::Client;
use futures::StreamExt;
use historian_core::SensorData;
use prost::Message;
use std::sync::Arc;

pub async fn connect_nats(url: &str) -> Result<Client> {
    let client = async_nats::connect(url)
        .await
        .context("Failed to connect to NATS")?;
    Ok(client)
}

pub async fn start_ingestion(
    client: Client,
    storage: Arc<dyn StorageEngine>,
    subject: &str,
    metadata_index: Arc<crate::metadata::MetadataIndex>,
) -> Result<()> {
    tracing::info!("Subscribing to {}", subject);
    let mut subscriber = client
        .subscribe(subject.to_string())
        .await
        .context("Failed to subscribe")?;

    while let Some(msg) = subscriber.next().await {
        let data = match SensorData::decode(msg.payload) {
            Ok(d) => d,
            Err(e) => {
                tracing::warn!("Failed to deserialize message: {}", e);
                continue;
            }
        };

        metadata_index.ensure_sensor_exists(&data.sensor_id).await;

        if let Err(e) = storage.write(&data) {
            tracing::error!("Failed to write to storage: {}", e);
        }
    }
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_nats_connection_failure() {
        // This test expects to fail connecting to a non-existent NATS server
        let result = connect_nats("nats://non-existent-host:4222").await;
        assert!(result.is_err());
    }
}
