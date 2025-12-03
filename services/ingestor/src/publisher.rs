use crate::buffer::HybridBuffer;
use async_nats::Client;
use historian_core::SensorData;
use prost::Message;
use tokio::sync::mpsc;
use tracing::{error, info, warn};

pub struct Publisher {
    buffer: HybridBuffer,
    nats_client: Option<Client>,
    nats_url: String,
}

impl Publisher {
    pub fn new(buffer: HybridBuffer, nats_url: String) -> Self {
        Self {
            buffer,
            nats_client: None,
            nats_url,
        }
    }

    pub async fn run(&mut self, mut rx: mpsc::Receiver<SensorData>) {
        // Connect to NATS initially
        self.connect_nats().await;

        loop {
            tokio::select! {
                Some(data) = rx.recv() => {
                    if let Err(e) = self.buffer.push(data).await {
                        error!("Failed to buffer data: {}", e);
                    }
                    // Try to flush immediately if connected
                    self.flush().await;
                }
                _ = tokio::time::sleep(tokio::time::Duration::from_millis(100)) => {
                    // Periodic flush retry
                    self.flush().await;
                }
            }
        }
    }

    async fn connect_nats(&mut self) {
        match async_nats::connect(&self.nats_url).await {
            Ok(client) => {
                info!("Connected to NATS at {}", self.nats_url);
                self.nats_client = Some(client);
            }
            Err(e) => {
                // Only log error if we were previously connected or it's a new attempt
                if self.nats_client.is_none() {
                    warn!("Failed to connect to NATS: {}", e);
                }
                self.nats_client = None;
            }
        }
    }

    async fn flush(&mut self) {
        if self.nats_client.is_none() {
            self.connect_nats().await;
            if self.nats_client.is_none() {
                return;
            }
        }

        let client = self.nats_client.as_ref().unwrap().clone();

        // 1. Flush Disk Buffer
        match self.buffer.flush_disk().await {
            Ok(items) => {
                if !items.is_empty() {
                    info!("Flushing {} items from disk buffer", items.len());
                    for (i, item) in items.iter().enumerate() {
                        if let Err(e) = self.publish(&client, item).await {
                            error!("Failed to publish disk item: {}", e);
                            // Re-save remaining items to disk to preserve them
                            // We push them back to the buffer.
                            // Note: pushing them back to HybridBuffer might put them in Mem or Disk depending on state.
                            // Ideally we should prepend them, but RingBuffer doesn't support push_front easily without rotation.
                            // For now, we push them back (they go to end of Mem or Disk). Order is slightly compromised but data is saved.
                            for remaining in items.iter().skip(i) {
                                let _ = self.buffer.push(remaining.clone()).await;
                            }
                            self.nats_client = None;
                            return;
                        }
                    }
                }
            }
            Err(e) => error!("Failed to read disk buffer: {}", e),
        }

        // 2. Flush Memory Buffer
        while let Some(item) = self.buffer.pop_mem() {
            if let Err(e) = self.publish(&client, &item).await {
                error!("Failed to publish mem item: {}", e);
                // Push back
                let _ = self.buffer.push(item).await;
                self.nats_client = None;
                return;
            }
        }
    }

    async fn publish(&self, client: &Client, item: &SensorData) -> anyhow::Result<()> {
        let mut buf = Vec::new();
        item.encode(&mut buf)?;

        // --- DUZELTME ---
        // Eski: let subject = format!("historian.data.{}", item.sensor_id);
        // Yeni: Sabit "data.raw" konusuna basiyoruz (Ingestion standardimiz)
        // Ileride bu dinamik olabilir ama simdilik dinleyicimiz "data.>" bekliyor.
        let subject = "data.raw";

        client.publish(subject.to_string(), buf.into()).await?;
        Ok(())
    }
}
