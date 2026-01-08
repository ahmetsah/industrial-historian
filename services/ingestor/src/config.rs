use config::{Config, ConfigError, File};
use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct ModbusConfig {
    pub ip: String,
    pub port: u16,
    pub unit_id: u8,
    pub poll_interval_ms: u64,
    pub registers: Vec<RegisterConfig>,
}

#[derive(Debug, Deserialize)]
pub struct RegisterConfig {
    pub address: u16,
    pub name: String,
    pub data_type: String, // e.g., "Float32", "Int16"
}

#[derive(Debug, Deserialize)]
pub struct NatsConfig {
    pub url: String,
    #[serde(default = "default_nats_subject")]
    pub subject: String,
}

fn default_nats_subject() -> String {
    "data.raw".to_string()
}

#[derive(Debug, Deserialize)]
pub struct BufferConfig {
    pub memory_capacity: usize,
    pub disk_path: String,
}

#[derive(Debug, Deserialize)]
pub struct CalculatedTagConfig {
    pub name: String,
    pub expression: String,
}

#[derive(Debug, Deserialize)]
pub struct Settings {
    pub modbus_devices: Vec<ModbusConfig>, // Changed from single to array
    pub nats: NatsConfig,
    pub buffer: BufferConfig,
    #[serde(default)]
    pub calculated_tags: Vec<CalculatedTagConfig>,
}

impl Settings {
    pub async fn load() -> Result<Self, anyhow::Error> {
        // 1. Check CONFIG_URL (API-based config)
        if let Ok(url) = std::env::var("CONFIG_URL") {
            println!("🌐 Loading config from URL: {}", url);
            let client = reqwest::Client::builder()
                .danger_accept_invalid_certs(true) // For internal nats/https if needed
                .build()?;

            let settings = client.get(&url).send().await?.json::<Settings>().await?;
            return Ok(settings);
        }

        // 2. Fallback to CONFIG_FILE (File-based config)
        let config_file =
            std::env::var("CONFIG_FILE").unwrap_or_else(|_| "config/default".to_string());

        println!("📂 Loading config from file: {}", config_file);

        let s = Config::builder()
            .add_source(File::with_name(&config_file).required(true))
            .build()?;

        Ok(s.try_deserialize()?)
    }

    /// Load from specific file path (for testing or manual override)
    pub fn from_file(path: &str) -> Result<Self, ConfigError> {
        let s = Config::builder()
            .add_source(File::with_name(path).required(true))
            .build()?;
        s.try_deserialize()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_load_config_from_string() {
        let config_str = r#"
            [[modbus_devices]]
            ip = "127.0.0.1"
            port = 502
            unit_id = 1
            poll_interval_ms = 100
            
            [[modbus_devices.registers]]
            address = 0
            name = "Temperature"
            data_type = "Float32"

            [nats]
            url = "nats://localhost:4222"

            [buffer]
            memory_capacity = 10000
            disk_path = "buffer.wal"
        "#;

        let s = Config::builder()
            .add_source(File::from_str(config_str, config::FileFormat::Toml))
            .build()
            .unwrap();

        let settings: Settings = s.try_deserialize().unwrap();
        assert_eq!(settings.modbus_devices.len(), 1);
        assert_eq!(settings.modbus_devices[0].ip, "127.0.0.1");
        assert_eq!(settings.nats.url, "nats://localhost:4222");
        assert_eq!(settings.buffer.memory_capacity, 10000);
        assert_eq!(settings.modbus_devices[0].port, 502);
        assert_eq!(settings.modbus_devices[0].registers.len(), 1);
        assert_eq!(settings.modbus_devices[0].registers[0].name, "Temperature");
    }
}
