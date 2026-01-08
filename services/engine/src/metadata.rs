use axum::{extract::State, routing::get, Json, Router};
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use tokio::sync::RwLock;

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub struct SensorMetadata {
    pub id: String,
    pub desc: String,
    pub factory: String,
    pub line: String,
    pub machine: String,
    #[serde(rename = "type")]
    pub sensor_type: String,
    pub unit: String,
}

#[derive(Debug, Default)]
pub struct MetadataIndex {
    sensors: Arc<RwLock<Vec<SensorMetadata>>>,
}

impl MetadataIndex {
    pub fn new() -> Self {
        Self {
            sensors: Arc::new(RwLock::new(Vec::new())),
        }
    }

    pub async fn add_sensor(&self, sensor: SensorMetadata) {
        let mut sensors = self.sensors.write().await;
        // Basic check to avoid duplicates by ID? For now simpler is better.
        sensors.push(sensor);
    }

    pub async fn get_all(&self) -> Vec<SensorMetadata> {
        self.sensors.read().await.clone()
    }

    pub async fn ensure_sensor_exists(&self, sensor_id: &str) {
        {
            let sensors = self.sensors.read().await;
            if sensors.iter().any(|s| s.id == sensor_id) {
                return;
            }
        }

        let mut sensors = self.sensors.write().await;
        if sensors.iter().any(|s| s.id == sensor_id) {
            return;
        }

        let parts: Vec<&str> = sensor_id.split('.').collect();
        let (factory, line, machine, sensor_type, id, desc, unit) = if parts.len() >= 5 {
             (
                parts[0].to_string(),
                parts[1].to_string(),
                parts[2].to_string(),
                parts[3].to_string(),
                sensor_id.to_string(),
                format!("{} - {}", parts[3], parts[4]),
                Self::infer_unit(parts[3]),
            )
        } else {
            (
                "Unknown".to_string(),
                "Unknown".to_string(),
                "Unknown".to_string(),
                "Generic".to_string(),
                sensor_id.to_string(),
                sensor_id.to_string(),
                "".to_string(),
            )
        };

        sensors.push(SensorMetadata {
            id,
            desc,
            factory,
            line,
            machine,
            sensor_type,
            unit,
        });
        tracing::info!("Auto-registered new sensor: {}", sensor_id);
    }

    /// Load sensor metadata from storage
    /// Expects sensor IDs in format: factory.line.machine.type.id
    /// Example: F1.L2.M3.Temp.S001
    pub async fn load_from_storage(&self, storage: &dyn crate::storage::StorageEngine) {
        let mut sensors = self.sensors.write().await;

        // Get all sensor IDs from storage
        match storage.get_all_sensor_ids() {
            Ok(sensor_ids) => {
                for sensor_id in sensor_ids {
                    // Parse sensor ID to extract hierarchy
                    // Expected format: factory.line.machine.type.id
                    let parts: Vec<&str> = sensor_id.split('.').collect();

                    let (factory, line, machine, sensor_type, id, desc, unit) = if parts.len() >= 5
                    {
                        // Full hierarchical format
                        (
                            parts[0].to_string(),
                            parts[1].to_string(),
                            parts[2].to_string(),
                            parts[3].to_string(),
                            sensor_id.clone(),
                            format!("{} - {}", parts[3], parts[4]),
                            Self::infer_unit(parts[3]),
                        )
                    } else {
                        // Fallback: use sensor_id as-is with defaults
                        (
                            "Unknown".to_string(),
                            "Unknown".to_string(),
                            "Unknown".to_string(),
                            "Generic".to_string(),
                            sensor_id.clone(),
                            sensor_id.clone(),
                            "".to_string(),
                        )
                    };

                    sensors.push(SensorMetadata {
                        id,
                        desc,
                        factory,
                        line,
                        machine,
                        sensor_type,
                        unit,
                    });
                }

                tracing::info!("Loaded {} sensors from storage", sensors.len());
            }
            Err(e) => {
                tracing::error!("Failed to load sensor IDs from storage: {}", e);
            }
        }
    }

    fn infer_unit(sensor_type: &str) -> String {
        match sensor_type.to_lowercase().as_str() {
            "temp" | "temperature" => "°C".to_string(),
            "pressure" => "bar".to_string(),
            "flow" => "m³/h".to_string(),
            "level" => "m".to_string(),
            "speed" => "rpm".to_string(),
            "power" => "kW".to_string(),
            "voltage" => "V".to_string(),
            "current" => "A".to_string(),
            _ => "".to_string(),
        }
    }
}

pub fn router(index: Arc<MetadataIndex>) -> Router {
    Router::new()
        .route("/api/v1/metadata", get(get_metadata))
        .with_state(index)
}

async fn get_metadata(State(index): State<Arc<MetadataIndex>>) -> Json<serde_json::Value> {
    let sensors = index.get_all().await;
    Json(serde_json::json!({ "sensors": sensors }))
}


