# Story 3.3: gRPC Query API

Status: done

## Story

As a Frontend Developer,
I want a fast API to query historical data,
so that I can populate charts without freezing the browser.

## Acceptance Criteria

1.  **Given** a gRPC request for `Sensor X` from `T1` to `T2` with `max_points=1000`
    *   **When** the Engine processes the request
    *   **Then** it scans the storage (RocksDB) for the time range
    *   **And** if the point count > 1000, it applies LTTB downsampling
    *   **And** returns the data stream in < 100ms (for reasonable ranges)
2.  **Given** a request for a non-existent sensor
    *   **When** processed
    *   **Then** it returns a `NOT_FOUND` gRPC status
3.  **Given** a large time range (e.g., 1 year)
    *   **When** requested
    *   **Then** the response size is bounded by `max_points` (e.g., ~1000 points), ensuring low network overhead.

## Tasks / Subtasks

- [x] Define Protobuf Service (`crates/historian-core/src/proto/query.proto`)
    - [x] Define `HistorianQuery` service
    - [x] Define `QueryRequest` (sensor_id, start_ts, end_ts, max_points)
    - [x] Define `QueryResponse` (stream of `SensorData` or `DataChunk`)
    - [x] Re-generate code (`cargo build`)
- [x] Implement Query Service (`services/engine/src/query.rs`)
    - [x] Implement `HistorianQuery` trait using `tonic`
    - [x] Inject `StorageEngine` reference
- [x] Implement Downsampling Logic (`services/engine/src/downsample.rs`)
    - [x] Add `lttb` crate dependency
    - [x] Implement function `downsample(points: Vec<SensorData>, threshold: usize) -> Vec<SensorData>`
    - [x] *Optimization:* If possible, apply downsampling while streaming from RocksDB to avoid loading all raw points into RAM first (though LTTB usually requires all points in a bucket). For MVP, loading into RAM for the requested range is acceptable if range is not massive. If massive, we rely on Tiering/Pre-aggregation (future).
- [x] Wire up Main (`services/engine/src/main.rs`)
    - [x] Start gRPC server on port 50051 (or configured)
    - [x] Run alongside the Ingestion loop (use `tokio::spawn`)

## Dev Notes

- **Architecture Patterns:**
    - **Protocol:** gRPC (HTTP/2).
    - **Library:** `tonic` (Rust gRPC).
    - **Algorithm:** LTTB (Largest-Triangle-Three-Buckets) for visual downsampling.
    - **Data Flow:** `gRPC Request` -> `Engine` -> `RocksDB Scan` -> `LTTB` -> `gRPC Response`.

- **Dependencies:**
    - `tonic`
    - `prost`
    - `lttb`

- **Performance Consideration:**
    - Scanning 1 year of raw data (e.g., 31M points at 1Hz) to downsample to 1000 points *on the fly* is too slow for <100ms latency.
    - *Constraint:* For this story, assume queries are for reasonable ranges (e.g., last 24h) or that the data density allows it.
    - *Future:* Pre-computed aggregates (rollups) will be needed for fast multi-year queries. This story focuses on the *API* and *Raw/Downsampled* path.

### References

- [Architecture Decision: API & Communication](docs/architecture.md#api--communication-patterns)
- [Epic 3: Efficient Storage Engine](docs/epics.md#epic-3-efficient-storage-engine-engine)

## Dev Agent Record

### Context Reference

- **Epic:** 3
- **Story:** 3
- **Key:** 3-3-grpc-query-api

### Agent Model Used

- Model: Gemini 2.0 Flash Experimental
- Date: 2025-12-03

### Completion Notes List

- Defined gRPC service in `crates/historian-core/src/proto/query.proto`.
- Implemented `scan` method in `StorageEngine` and `RocksDBStorage`.
- Implemented `QueryService` using `tonic`.
- Implemented LTTB downsampling using `lttb` crate.
- Wired up gRPC server in `main.rs`.
- Verified compilation and tests.

### File List

- `crates/historian-core/src/proto/query.proto`
- `crates/historian-core/build.rs`
- `services/engine/Cargo.toml`
- `services/engine/src/storage/mod.rs`
- `services/engine/src/storage/rocksdb.rs`
- `services/engine/src/query.rs`
- `services/engine/src/downsample.rs`
- `services/engine/src/main.rs`
