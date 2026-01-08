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
use tracing::{error, info};
//use tracing_subscriber;

mod buffer;
mod config;
mod engine;
mod modbus;
mod publisher;

use crate::buffer::HybridBuffer;
use crate::config::Settings;
use crate::modbus::ModbusAdapter;
use crate::publisher::Publisher;
use historian_core::SensorData;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // 1. Ortam değişkenlerini ve loglamayı başlat
    dotenv().ok();
    tracing_subscriber::fmt::init();

    info!("Starting Ingestor Service");

    // 2. Konfigürasyonu yükle
    let settings = Settings::load().await.map_err(|e| {
        error!("Failed to load configuration: {}", e);
        e
    })?;

    info!("Configuration loaded");

    // 3. Hybrid Buffer ve Publisher oluştur
    let buffer = HybridBuffer::new(settings.buffer.memory_capacity, settings.buffer.disk_path);

    let mut publisher = Publisher::new(
        buffer,
        settings.nats.url.clone(),
        settings.nats.subject.clone(),
    );

    // 4. Veri kanalları oluştur (Dinamik buffer boyutu)
    // Modbus -> Engine (tx_raw -> rx_raw)
    // Engine -> Publisher (tx_pub -> rx_pub)
    let device_count = settings.modbus_devices.len();
    let buffer_size = (100 * device_count).max(1000); // Min 1000, cihaz başına 100

    info!(
        "Channel buffer size: {} (for {} devices)",
        buffer_size, device_count
    );

    let (tx_raw, mut rx_raw) = mpsc::channel::<SensorData>(buffer_size);
    let (tx_pub, rx_pub) = mpsc::channel::<SensorData>(buffer_size);

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

    // 7. Modbus Adaptörlerini başlat (Her cihaz için ayrı task)
    let device_count = settings.modbus_devices.len();
    for (idx, modbus_config) in settings.modbus_devices.into_iter().enumerate() {
        let tx_clone = tx_raw.clone();
        let device_name = format!("{}:{}", modbus_config.ip, modbus_config.port);

        tokio::spawn(async move {
            info!(
                "[Device {}/{}] Starting Modbus adapter for {}",
                idx + 1,
                device_count,
                device_name
            );
            let mut adapter = ModbusAdapter::new(modbus_config, tx_clone);

            if let Err(e) = adapter.connect().await {
                error!("[{}] Initial connection failed: {}", device_name, e);
            }
            adapter.poll_loop().await;
        });
    }

    // 8. Ana prosesi ayakta tut
    info!(
        "Ingestor running with {} Modbus device(s). Press Ctrl+C to stop.",
        device_count
    );
    tokio::signal::ctrl_c().await?;
    info!("Shutting down...");

    Ok(())
}
