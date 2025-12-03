use historian_core::SensorData;
use lttb::{lttb, DataPoint};

pub fn downsample(points: Vec<SensorData>, threshold: usize) -> Vec<SensorData> {
    if points.len() <= threshold || threshold < 3 {
        return points;
    }

    let mut data_points = Vec::with_capacity(points.len());
    for p in &points {
        data_points.push(DataPoint {
            x: p.timestamp_ms as f64,
            y: p.value,
        });
    }

    let downsampled_data = lttb(data_points, threshold);

    // Optimization: Since LTTB preserves original points, we can just map back.
    // LTTB returns DataPoint { x, y }. x is timestamp as f64.
    // We can iterate both sorted lists to find matches in O(N).

    let mut result = Vec::with_capacity(downsampled_data.len());
    let mut point_iter = points.into_iter();
    let mut current_point = point_iter.next();

    for dp in downsampled_data {
        let target_ts = dp.x as i64;

        while let Some(p) = current_point.as_ref() {
            if p.timestamp_ms == target_ts {
                result.push(p.clone());
                current_point = point_iter.next(); // Advance for next iteration
                break;
            } else if p.timestamp_ms < target_ts {
                // Skip points that were dropped by LTTB
                current_point = point_iter.next();
            } else {
                // p.timestamp_ms > target_ts: Should not happen if LTTB preserves order and points exist
                // But if it does, we might have missed it or float precision issue?
                // LTTB uses the exact points from input, so exact match on timestamp (converted to f64) should work
                // provided no precision loss for i64 -> f64 (safe for timestamps < 2^53, which is until year 285,428)
                break;
            }
        }
    }

    result
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_downsample() {
        let mut points = Vec::new();
        for i in 0..100 {
            points.push(SensorData {
                sensor_id: "test".to_string(),
                timestamp_ms: i * 1000,
                value: i as f64,
                quality: 1,
            });
        }

        let downsampled = downsample(points.clone(), 10);
        assert!(downsampled.len() <= 10);
        assert_eq!(downsampled.first().unwrap().timestamp_ms, 0);
        assert_eq!(downsampled.last().unwrap().timestamp_ms, 99000);
    }
}
