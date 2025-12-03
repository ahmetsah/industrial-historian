use axum::{
    body::Body,
    extract::{Query, State},
    response::IntoResponse,
    routing::get,
    Router,
};
use futures::StreamExt;
use serde::Deserialize;
use std::sync::Arc;
use crate::storage::StorageEngine;

#[derive(Deserialize)]
pub struct ExportQuery {
    sensor_id: String,
    start_ts: i64,
    end_ts: i64,
}

pub async fn start_server(storage: Arc<dyn StorageEngine>, port: u16) -> anyhow::Result<()> {
    let app = Router::new()
        .route("/api/v1/export", get(export_handler))
        .with_state(storage)
        .layer(tower_http::cors::CorsLayer::permissive());

    let addr = std::net::SocketAddr::from(([0, 0, 0, 0], port));
    tracing::info!("Starting Export HTTP server on {}", addr);
    let listener = tokio::net::TcpListener::bind(addr).await?;
    axum::serve(listener, app).await?;
    Ok(())
}

async fn export_handler(
    State(storage): State<Arc<dyn StorageEngine>>,
    Query(params): Query<ExportQuery>,
) -> impl IntoResponse {
    let stream_result = storage.scan_stream(&params.sensor_id, params.start_ts, params.end_ts);

    match stream_result {
        Ok(stream) => {
            let csv_stream = stream.map(|res| {
                match res {
                    Ok(point) => {
                        // Format: timestamp_ms,value,quality
                        Ok::<_, std::io::Error>(format!("{},{},{}\n", point.timestamp_ms, point.value, point.quality))
                    }
                    Err(e) => {
                        tracing::error!("Error in stream: {}", e);
                        Err(std::io::Error::new(std::io::ErrorKind::Other, e.to_string()))
                    }
                }
            });
            
            // Add header
            let header = futures::stream::once(async { Ok("timestamp_ms,value,quality\n".to_string()) });
            let body_stream = header.chain(csv_stream);

            // Convert to bytes stream
            let byte_stream = body_stream.map(|res| {
                res.map(|s| axum::body::Bytes::from(s))
            });

            // Set headers for file download
            let headers = [
                (axum::http::header::CONTENT_TYPE, "text/csv"),
                (axum::http::header::CONTENT_DISPOSITION, "attachment; filename=\"export.csv\""),
            ];

            (headers, Body::from_stream(byte_stream)).into_response()
        }
        Err(e) => {
            (axum::http::StatusCode::INTERNAL_SERVER_ERROR, e.to_string()).into_response()
        }
    }
}
