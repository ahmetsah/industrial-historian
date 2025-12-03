use crate::config::CalculatedTagConfig;
use evalexpr::{eval_with_context, ContextWithMutableVariables, HashMapContext};
use historian_core::SensorData;
use std::time::{SystemTime, UNIX_EPOCH};
use tracing::{error, info};
//use tracing::error;

pub struct Engine {
    context: HashMapContext,
    calculations: Vec<CalculatedTagConfig>,
}

impl Engine {
    pub fn new(calculations: Vec<CalculatedTagConfig>) -> Self {
        let context = HashMapContext::new();

        Self {
            context,
            calculations,
        }
    }

    pub fn process(&mut self, data: SensorData) -> Vec<SensorData> {
        // 1. Update Context with new data
        // We store values as float in the context
        if let Err(e) = self
            .context
            .set_value(data.sensor_id.clone(), data.value.into())
        {
            error!("Failed to set context value for {}: {}", data.sensor_id, e);
            return vec![data]; // Return original data even if context update fails
        }

        let mut results = vec![data]; // Start with the raw data

        // 2. Evaluate all calculated tags
        // Note: This does not handle chained calculations (C = A + B, D = C + 1) in a single pass if order is wrong.
        // For MVP, we assume topological order or single-level depth.
        for calc in &self.calculations {
            match eval_with_context(&calc.expression, &self.context) {
                Ok(val) => {
                    if let Ok(float_val) = val.as_float() {
                        // Create new SensorData
                        let calc_data = SensorData {
                            sensor_id: calc.name.clone(),
                            value: float_val,
                            timestamp_ms: SystemTime::now()
                                .duration_since(UNIX_EPOCH)
                                .unwrap_or_default()
                                .as_millis() as i64,
                            quality: 1,
                        };

                        // Update context with this new calculated value so subsequent calculations can use it
                        let _ = self.context.set_value(calc.name.clone(), float_val.into());

                        results.push(calc_data);
                    }
                }
                Err(_e) => {
                    // It's normal to fail if dependencies are missing (e.g. startup)
                    // info!("Skipping calculation for {}: {}", calc.name, _e);
                }
            }
        }

        results
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_simple_calculation() {
        let config = vec![CalculatedTagConfig {
            name: "Efficiency".to_string(),
            expression: "Output / Input".to_string(),
        }];
        let mut engine = Engine::new(config);

        // 1. Send Input
        let input = SensorData {
            sensor_id: "Input".to_string(),
            value: 10.0,
            timestamp_ms: 100,
            quality: 1,
        };
        let res1 = engine.process(input);
        assert_eq!(res1.len(), 1); // Only Input, Efficiency fails (Output missing)

        // 2. Send Output
        let output = SensorData {
            sensor_id: "Output".to_string(),
            value: 80.0,
            timestamp_ms: 100,
            quality: 1,
        };
        let res2 = engine.process(output);
        assert_eq!(res2.len(), 2); // Output + Efficiency

        let eff = &res2[1];
        assert_eq!(eff.sensor_id, "Efficiency");
        assert_eq!(eff.value, 8.0);
    }
}
