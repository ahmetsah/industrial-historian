mod downsample;
mod ingest;
mod query;
mod storage;

use anyhow::Result;
use std::sync::Arc;
use storage::rocksdb::RocksDBStorage;

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

    let client = ingest::connect_nats(&nats_url).await?;
    tracing::info!("Connected to NATS at {}", nats_url);

    let storage = Arc::new(RocksDBStorage::new(&db_path, None)?);
    tracing::info!("Initialized RocksDB at {}", db_path);

    let query_service = query::QueryService::new(storage.clone());
    tracing::info!("Starting ingestion on subject {}", subject);
    tracing::info!("Starting gRPC server on {}", grpc_addr);

    tokio::select! {
        res = ingest::start_ingestion(client, storage, &subject) => {
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
        _ = tokio::signal::ctrl_c() => {
            tracing::info!("Shutting down Engine Service");
        }
    }
    Ok(())
}
