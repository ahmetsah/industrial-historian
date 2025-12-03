//use historian_core::hello;
/*
mod config;
mod modbus;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();
    tracing::info!("Starting Ingestor Service");
    //hello();
}

*/


use dotenvy::dotenv;
use tokio::sync::mpsc;
use tracing::{info, error};
use tracing_subscriber;

mod config;
mod modbus;
mod buffer;
mod publisher;
mod engine;

use crate::config::Settings;
use crate::modbus::ModbusAdapter;
use crate::buffer::HybridBuffer;
use crate::publisher::Publisher;
use historian_core::SensorData; 

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // 1. Ortam değişkenlerini ve loglamayı başlat
    dotenv().ok();
    tracing_subscriber::fmt::init();

    info!("Starting Ingestor Service");

    // 2. Konfigürasyonu yükle
    let settings = Settings::new().map_err(|e| {
        error!("Failed to load configuration: {}", e);
        e
    })?;
    
    info!("Configuration loaded");

    // 3. Hybrid Buffer ve Publisher oluştur
    let buffer = HybridBuffer::new(
        settings.buffer.memory_capacity,
        settings.buffer.disk_path
    );

    let mut publisher = Publisher::new(buffer, settings.nats.url.clone());

    // 4. Veri kanalları oluştur
    // Modbus -> Engine (tx_raw -> rx_raw)
    // Engine -> Publisher (tx_pub -> rx_pub)
    let (tx_raw, mut rx_raw) = mpsc::channel::<SensorData>(100);
    let (tx_pub, rx_pub) = mpsc::channel::<SensorData>(100);

    // 5. Publisher'ı başlat (Ayrı task)
    tokio::spawn(async move {
        publisher.run(rx_pub).await;
    });

    // 6. Engine'i başlat (Ayrı task)
    use crate::engine::Engine;
    let mut engine = Engine::new(settings.calculated_tags);
    
    tokio::spawn(async move {
        info!("Calculation Engine started");
        while let Some(data) = rx_raw.recv().await {
            let results = engine.process(data);
            for res in results {
                if let Err(e) = tx_pub.send(res).await {
                    error!("Failed to send to publisher: {}", e);
                    break;
                }
            }
        }
    });

    // 7. Modbus Adaptörünü başlat
    let modbus_config = settings.modbus;
    let mut adapter = ModbusAdapter::new(modbus_config, tx_raw);

    tokio::spawn(async move {
        if let Err(e) = adapter.connect().await {
            error!("Initial connection failed: {}", e);
        }
        adapter.poll_loop().await;
    });

    // 7. Ana prosesi ayakta tut
    info!("Ingestor running. Press Ctrl+C to stop.");
    tokio::signal::ctrl_c().await?;
    info!("Shutting down...");

    Ok(())
}
