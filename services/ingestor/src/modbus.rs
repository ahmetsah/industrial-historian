use crate::config::ModbusConfig;
use anyhow::Result;
use historian_core::SensorData;
use std::net::SocketAddr;
use std::time::Duration;
use std::time::{SystemTime, UNIX_EPOCH};
use tokio_modbus::client::Context;
use tokio_modbus::prelude::*;

use tokio::sync::mpsc;

pub struct ModbusAdapter {
    config: ModbusConfig,
    ctx: Option<Context>,
    sender: mpsc::Sender<SensorData>,
}

impl ModbusAdapter {
    pub fn new(config: ModbusConfig, sender: mpsc::Sender<SensorData>) -> Self {
        Self {
            config,
            ctx: None,
            sender,
        }
    }

    pub async fn connect(&mut self) -> Result<()> {
        let addr_str = format!("{}:{}", self.config.ip, self.config.port);
        let addr: SocketAddr = addr_str.parse()?;

        // Note: In real implementation we might need to handle Unit ID if using RTU over TCP,
        // but for pure TCP usually Unit ID is handled in the request or connection.
        // tokio-modbus tcp::connect returns a Context.

        let mut ctx = tcp::connect(addr).await?;

        // Set the Unit ID (Slave ID)
        ctx.set_slave(Slave(self.config.unit_id));

        self.ctx = Some(ctx);
        Ok(())
    }

    pub async fn poll_loop(&mut self) {
        let interval = Duration::from_millis(self.config.poll_interval_ms);
        let mut interval_timer = tokio::time::interval(interval);
        let mut backoff = Duration::from_secs(1);

        loop {
            if self.ctx.is_none() {
                tracing::info!("Attempting to connect...");
                match self.connect().await {
                    Ok(_) => {
                        tracing::info!("Connected to Modbus device");
                        backoff = Duration::from_secs(1); // Reset backoff
                    }
                    Err(e) => {
                        tracing::error!("Connection failed: {}. Retrying in {:?}", e, backoff);
                        tokio::time::sleep(backoff).await;
                        backoff = std::cmp::min(backoff * 2, Duration::from_secs(60));
                        continue;
                    }
                }
            }

            interval_timer.tick().await;
            if let Err(e) = self.poll().await {
                tracing::error!("Modbus poll error: {}", e);
                // Force reconnection on error
                self.ctx = None;
            }
        }
    }

    async fn poll(&mut self) -> Result<()> {
        if let Some(ctx) = &mut self.ctx {
            for reg_config in &self.config.registers {
                let count = get_register_count(&reg_config.data_type);

                let data = ctx
                    .read_holding_registers(reg_config.address, count)
                    .await?;
                // In some versions/configurations, read_holding_registers might return a Result<Vec<u16>, ExceptionCode> inside the outer Result?
                // Or maybe the outer Result is the IO error, and the inner is the Modbus exception?
                // Based on the error message, 'data' is Result<Vec<u16>, ExceptionCode>.
                // So we need to handle the inner result.
                let data = data.map_err(|e| anyhow::anyhow!("Modbus Exception: {:?}", e))?;

                let value = convert_registers(&data, &reg_config.data_type)?;

                let sensor_data = SensorData {
                    sensor_id: reg_config.name.clone(),
                    value,
                    timestamp_ms: SystemTime::now().duration_since(UNIX_EPOCH)?.as_millis() as i64,
                    quality: 1,
                };

                tracing::info!("Read sensor: {:?}", sensor_data);
                if let Err(e) = self.sender.send(sensor_data).await {
                    tracing::error!("Failed to send sensor data: {}", e);
                }
            }
        }
        Ok(())
    }
}

/// Get the number of 16-bit registers required for a data type
fn get_register_count(data_type: &str) -> u16 {
    match data_type {
        "Bool" | "Int16" | "UInt16" => 1,
        "Int32" | "UInt32" | "Float32" => 2,
        "Int64" | "UInt64" | "Float64" => 4,
        _ => 1, // Default to 1 register
    }
}

fn convert_registers(regs: &[u16], data_type: &str) -> Result<f64> {
    match data_type {
        "Bool" => {
            if regs.is_empty() {
                return Err(anyhow::anyhow!("Not enough registers for Bool"));
            }
            Ok(if regs[0] != 0 { 1.0 } else { 0.0 })
        }
        "Int16" => {
            if regs.is_empty() {
                return Err(anyhow::anyhow!("Not enough registers for Int16"));
            }
            let val = regs[0] as i16;
            Ok(val as f64)
        }
        "UInt16" => {
            if regs.is_empty() {
                return Err(anyhow::anyhow!("Not enough registers for UInt16"));
            }
            Ok(regs[0] as f64)
        }
        "Int32" => {
            if regs.len() < 2 {
                return Err(anyhow::anyhow!("Not enough registers for Int32"));
            }
            let high = regs[0] as u32;
            let low = regs[1] as u32;
            let combined = (high << 16) | low;
            let val = combined as i32;
            Ok(val as f64)
        }
        "UInt32" => {
            if regs.len() < 2 {
                return Err(anyhow::anyhow!("Not enough registers for UInt32"));
            }
            let high = regs[0] as u32;
            let low = regs[1] as u32;
            let combined = (high << 16) | low;
            Ok(combined as f64)
        }
        "Float32" => {
            if regs.len() < 2 {
                return Err(anyhow::anyhow!("Not enough registers for Float32"));
            }
            let high = regs[0];
            let low = regs[1];
            let combined = ((high as u32) << 16) | (low as u32);
            let val = f32::from_bits(combined);
            Ok(val as f64)
        }
        "Int64" => {
            if regs.len() < 4 {
                return Err(anyhow::anyhow!("Not enough registers for Int64"));
            }
            let combined = ((regs[0] as u64) << 48)
                | ((regs[1] as u64) << 32)
                | ((regs[2] as u64) << 16)
                | (regs[3] as u64);
            let val = combined as i64;
            Ok(val as f64)
        }
        "UInt64" => {
            if regs.len() < 4 {
                return Err(anyhow::anyhow!("Not enough registers for UInt64"));
            }
            let combined = ((regs[0] as u64) << 48)
                | ((regs[1] as u64) << 32)
                | ((regs[2] as u64) << 16)
                | (regs[3] as u64);
            Ok(combined as f64)
        }
        "Float64" => {
            if regs.len() < 4 {
                return Err(anyhow::anyhow!("Not enough registers for Float64"));
            }
            let combined = ((regs[0] as u64) << 48)
                | ((regs[1] as u64) << 32)
                | ((regs[2] as u64) << 16)
                | (regs[3] as u64);
            let val = f64::from_bits(combined);
            Ok(val)
        }
        _ => Err(anyhow::anyhow!("Unsupported data type: {}", data_type)),
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_connect_fail() {
        let (tx, _rx) = tokio::sync::mpsc::channel(10);
        let config = ModbusConfig {
            ip: "127.0.0.1".to_string(),
            port: 50200, // Unlikely to be open
            unit_id: 1,
            poll_interval_ms: 100,
            registers: vec![],
        };
        let mut adapter = ModbusAdapter::new(config, tx);
        let res = adapter.connect().await;
        assert!(res.is_err());
    }

    #[test]
    fn test_convert_registers() {
        let regs = vec![0x42F6, 0xE979];
        let val = convert_registers(&regs, "Float32").unwrap();
        assert!((val - 123.456).abs() < 0.001);
    }
}
