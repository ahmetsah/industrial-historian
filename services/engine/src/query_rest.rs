use crate::storage::StorageEngine;
use axum::{
    extract::{Query, State},
    response::IntoResponse,
    routing::get,
    Json, Router,
};
use serde::{Deserialize, Serialize};
use std::sync::Arc;

#[derive(Deserialize)]
pub struct QueryParams {
    sensor_id: String,
    start_ts: i64,
    end_ts: i64,
    #[serde(default = "default_max_points")]
    max_points: usize,
}

fn default_max_points() -> usize {
    1000
}

#[derive(Serialize)]
pub struct QueryResponse {
    pub points: Vec<DataPoint>,
}

#[derive(Serialize)]
pub struct DataPoint {
    pub timestamp: i64,
    pub value: f64,
}

pub fn router(storage: Arc<dyn StorageEngine>) -> Router {
    Router::new()
        .route("/api/v1/query", get(query_handler))
        .with_state(storage)
}

async fn query_handler(
    State(storage): State<Arc<dyn StorageEngine>>,
    Query(params): Query<QueryParams>,
) -> impl IntoResponse {
    match storage
        .scan(&params.sensor_id, params.start_ts, params.end_ts)
        .await
    {
        Ok(points) => {
            // Downsample if needed
            let final_points = if points.len() > params.max_points {
                downsample(points, params.max_points)
            } else {
                points
            };

            let response = QueryResponse {
                points: final_points
                    .into_iter()
                    .map(|p| DataPoint {
                        timestamp: p.timestamp_ms,
                        value: p.value,
                    })
                    .collect(),
            };

            Json(response).into_response()
        }
        Err(e) => {
            tracing::error!("Query error for sensor {}: {}", params.sensor_id, e);
            (
                axum::http::StatusCode::INTERNAL_SERVER_ERROR,
                format!("Query error: {}", e),
            )
                .into_response()
        }
    }
}

// Simple LTTB-like downsampling
fn downsample(
    points: Vec<historian_core::SensorData>,
    max_points: usize,
) -> Vec<historian_core::SensorData> {
    if points.len() <= max_points || max_points < 3 {
        return points;
    }

    let bucket_size = (points.len() as f64) / (max_points as f64);
    let mut result = Vec::with_capacity(max_points);

    // Always include first point
    result.push(points[0].clone());

    for i in 1..(max_points - 1) {
        let bucket_start = ((i as f64) * bucket_size) as usize;
        let bucket_end = (((i + 1) as f64) * bucket_size) as usize;
        let bucket_end = bucket_end.min(points.len());

        if bucket_start < bucket_end {
            // Find point in bucket with max value (simple approach)
            let mut max_idx = bucket_start;
            let mut max_val = points[bucket_start].value;
            for j in bucket_start..bucket_end {
                if points[j].value > max_val {
                    max_val = points[j].value;
                    max_idx = j;
                }
            }
            result.push(points[max_idx].clone());
        }
    }

    // Always include last point
    result.push(points[points.len() - 1].clone());

    result
}
