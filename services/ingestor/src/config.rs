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
    pub modbus: ModbusConfig,
    pub nats: NatsConfig,
    pub buffer: BufferConfig,
    #[serde(default)]
    pub calculated_tags: Vec<CalculatedTagConfig>,
}

impl Settings {
    pub fn new() -> Result<Self, ConfigError> {
        let s = Config::builder()
            .add_source(File::with_name("config/default").required(false))
            .add_source(File::with_name("config/local").required(false))
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
            [modbus]
            ip = "127.0.0.1"
            port = 502
            unit_id = 1
            poll_interval_ms = 100
            
            [[modbus.registers]]
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
        assert_eq!(settings.modbus.ip, "127.0.0.1");
        assert_eq!(settings.nats.url, "nats://localhost:4222");
        assert_eq!(settings.buffer.memory_capacity, 10000);
        assert_eq!(settings.modbus.port, 502);
        assert_eq!(settings.modbus.registers.len(), 1);
        assert_eq!(settings.modbus.registers[0].name, "Temperature");
    }
}
