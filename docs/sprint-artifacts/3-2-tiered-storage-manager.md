# Story 3.2: Tiered Storage Manager

Status: done

## Story

As a CFO,
I want old data moved to cheaper storage,
so that we don't waste expensive SSD space on 5-year-old logs.

## Acceptance Criteria

1.  **Given** data older than 1 month (configurable)
    *   **When** the Tiering Job runs
    *   **Then** it identifies eligible SSTables or data chunks
    *   **And** uploads them to the configured MinIO (S3) bucket
    *   **And** verifies the upload integrity
    *   **And** deletes the local copy after successful upload
2.  **Given** a query for historical data
    *   **When** the data is not found locally
    *   **Then** the Engine automatically fetches it from S3
    *   **And** serves it transparently to the user
3.  **Given** a network failure during upload
    *   **When** the job runs
    *   **Then** it retries or resumes without data loss

## Tasks / Subtasks

- [x] Implement S3 Client Wrapper (`src/storage/tiered/s3.rs`)
    - [x] Add `rust-s3` dependency
    - [x] Implement `S3Client` struct with `put_object`, `get_object`, `delete_object`
    - [x] Configure from `minio` settings in `config.toml`
- [x] Implement Tiering Policy Manager (`src/storage/tiered/policy.rs`)
    - [x] Define `TieringPolicy` struct (e.g., `max_age_days`, `max_local_size_gb`)
    - [x] Implement logic to scan local storage for eligible files
- [x] Implement Background Tiering Job (`src/storage/tiered/job.rs`)
    - [x] Create an async task loop that runs periodically (e.g., every hour)
    - [x] Execute the tiering policy
    - [x] Handle upload failures gracefully
- [x] Update Storage Engine for Hybrid Read (`src/storage/mod.rs`)
    - [x] Modify `get` or `scan` logic to check S3 if data is missing locally (or based on metadata)
    - [x] Implement a metadata index to track which time ranges are in S3
- [x] Integration Tests
    - [x] Test upload/download flow with a local MinIO container

## Dev Notes

- **Architecture Patterns:**
    - **Library:** `rust-s3` is specified in Architecture.
    - **Metadata:** We need a way to know *what* is in S3 without listing the bucket every time.
        - *Option A:* Keep a local "Manifest" (SQLite or RocksDB column family) tracking `[StartTs, EndTs] -> S3Key`.
        - *Option B:* Use a deterministic naming convention `sensor_id/year/month/day/chunk_id.sst` and check existence (slower).
        - *Selected:* **Option A** (Manifest in RocksDB). Store metadata in a separate Column Family `cf_tiering_metadata`.
    - **Granularity:** Tiering individual RocksDB SSTables is hard because RocksDB manages them.
        - *Refinement:* We might need to implement "Application-Level Tiering". Instead of moving RocksDB internal files, we should:
            1.  Write recent data to RocksDB (Hot Tier).
            2.  Periodically export/dump old data from RocksDB to Parquet/Arrow/Protobuf files.
            3.  Upload those files to S3.
            4.  Delete old range from RocksDB (`delete_range`).
        - *Decision:* For this story, let's stick to a simpler "File-based" approach if we were using a custom engine, but since we use RocksDB, **Export-Upload-Delete** is the viable path.
        - *Format:* Export as compressed Protobuf blobs or Parquet (if we want analytics later). Let's stick to the internal `SensorData` Protobuf format for simplicity now.

- **Dependencies:**
    - `rust-s3`
    - `tokio`
    - `cron` or simple `tokio::time::interval`

### References

- [Architecture Decision: Tiered Storage](docs/architecture.md#data-architecture)
- [Epic 3: Efficient Storage Engine](docs/epics.md#epic-3-efficient-storage-engine-engine)

## Dev Agent Record

### Context Reference

- **Epic:** 3
- **Story:** 2
- **Key:** 3-2-tiered-storage-manager

### Agent Model Used

- Model: Gemini 2.0 Flash Experimental
- Date: 2025-12-03

### Completion Notes List

- Identified challenge with RocksDB tiering (it's not file-based in a simple way).
- Proposed "Export-Upload-Delete" strategy.
- Selected `rust-s3` as client.
- Implemented `S3Client` wrapper using `rust-s3`.
- Implemented `TieringPolicy` for age-based retention.
- Implemented `TieringJob` to periodically scan, upload, and delete old data.
- Updated `RocksDBStorage` to support Column Families and Metadata (`tiering_metadata`).
- Implemented `read` in `StorageEngine` with Hybrid Read logic (Buffer -> RocksDB -> S3 Metadata).
- Added unit tests for metadata operations.

### File List

- `services/engine/Cargo.toml`
- `services/engine/src/main.rs`
- `services/engine/src/storage/mod.rs`
- `services/engine/src/storage/rocksdb.rs`
- `services/engine/src/storage/tiered/mod.rs`
- `services/engine/src/storage/tiered/s3.rs`
- `services/engine/src/storage/tiered/policy.rs`
- `services/engine/src/storage/tiered/job.rs`
