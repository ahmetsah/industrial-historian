pub mod historian {
    pub mod v1 {
        tonic::include_proto!("historian.v1");
    }
}

pub use historian::v1::*;

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_sensor_data_instantiation() {
        let sensor_data = SensorData {
            sensor_id: "sensor-1".to_string(),
            value: 123.45,
            timestamp_ms: 1678886400000,
            quality: 1,
        };

        assert_eq!(sensor_data.sensor_id, "sensor-1");
        assert_eq!(sensor_data.value, 123.45);
        assert_eq!(sensor_data.quality, 1);
    }
}
