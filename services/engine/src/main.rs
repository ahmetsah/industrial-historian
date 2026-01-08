mod downsample;
mod export;
mod ingest;
mod metadata;
mod query;
mod query_rest;
mod storage;

use anyhow::Result;
use axum::Router;
use std::sync::Arc;
use storage::rocksdb::RocksDBStorage;
use crate::storage::tiered::{s3::S3Client, policy::TieringPolicy, job::TieringJob};

#[tokio::main]
async fn main() -> Result<()> {
    tracing_subscriber::fmt::init();
    tracing::info!("Starting Engine Service");

    let nats_url =
        std::env::var("NATS_URL").unwrap_or_else(|_| "nats://localhost:4222".to_string());
    let db_path = std::env::var("DB_PATH").unwrap_or_else(|_| "/tmp/historian-db".to_string());
    let subject = std::env::var("NATS_SUBJECT").unwrap_or_else(|_| "data.>".to_string());
    let grpc_addr = std::env::var("GRPC_ADDR")
        .unwrap_or_else(|_| "0.0.0.0:50051".to_string())
        .parse()?;
    let http_port = std::env::var("HTTP_PORT")
        .unwrap_or_else(|_| "8080".to_string())
        .parse::<u16>()?;

    let client = ingest::connect_nats(&nats_url).await?;
    tracing::info!("Connected to NATS at {}", nats_url);

    // S3 & Tiering Configuration
    let s3_endpoint = std::env::var("S3_ENDPOINT").ok();
    let s3_bucket = std::env::var("S3_BUCKET").ok();
    let s3_access_key = std::env::var("S3_ACCESS_KEY").ok();
    let s3_secret_key = std::env::var("S3_SECRET_KEY").ok();
    let tiering_enabled = std::env::var("TIERING_ENABLED").unwrap_or("false".to_string()) == "true";
    let tiering_min_age_ms = std::env::var("TIERING_MIN_AGE_MS")
        .unwrap_or("3600000".to_string()) // 1 hour default
        .parse::<u64>().unwrap_or(3600000);

    let mut s3_client = None;
    if let (Some(endpoint), Some(bucket), Some(access), Some(secret)) = (s3_endpoint, s3_bucket, s3_access_key, s3_secret_key) {
        if tiering_enabled {
            tracing::info!("Initializing S3 Client for bucket: {}", bucket);
            match S3Client::new(&endpoint, &bucket, &access, &secret) {
                Ok(client) => s3_client = Some(client),
                Err(e) => tracing::error!("Failed to create S3 client: {}", e),
            }
        }
    }

    let storage = Arc::new(RocksDBStorage::new(&db_path, s3_client.clone())?);
    
    if let Some(client) = s3_client {
        let policy = TieringPolicy {
            max_age_ms: tiering_min_age_ms as i64,

        };
        let job = TieringJob::new(storage.clone(), client, policy, 60000); // Run every 60s
        tokio::spawn(async move {
            job.run().await;
        });
        tracing::info!("Tiering Job started. Old data (>{}ms) will be moved to S3.", tiering_min_age_ms);
    }
    tracing::info!("Initialized RocksDB at {}", db_path);

    let query_service = query::QueryService::new(storage.clone());
    let metadata_index = Arc::new(metadata::MetadataIndex::new());
    metadata_index.load_from_storage(storage.as_ref()).await;
    tracing::info!("Initialized Metadata Index");

    // Combine routers
    let app = Router::new()
        .merge(export::router(storage.clone()))
        .merge(metadata::router(metadata_index.clone()))
        .merge(query_rest::router(storage.clone()))
        .layer(tower_http::cors::CorsLayer::permissive());

    let http_addr = std::net::SocketAddr::from(([0, 0, 0, 0], http_port));
    tracing::info!("Starting ingestion on subject {}", subject);
    tracing::info!("Starting gRPC server on {}", grpc_addr);
    tracing::info!("Starting HTTP server on {}", http_addr);

    let listener = tokio::net::TcpListener::bind(http_addr).await?;

    tokio::select! {
        res = ingest::start_ingestion(client, storage.clone(), &subject, metadata_index.clone()) => {
            if let Err(e) = res {
                tracing::error!("Ingestion failed: {}", e);
            }
        }
        res = tonic::transport::Server::builder()
            .add_service(historian_core::historian_query_server::HistorianQueryServer::new(query_service))
            .serve(grpc_addr) => {
                if let Err(e) = res {
                    tracing::error!("gRPC server failed: {}", e);
                }
            }
        res = axum::serve(listener, app) => {
            if let Err(e) = res {
                tracing::error!("HTTP server failed: {}", e);
            }
        }
        _ = tokio::signal::ctrl_c() => {
            tracing::info!("Shutting down Engine Service");
        }
    }
    Ok(())
}
