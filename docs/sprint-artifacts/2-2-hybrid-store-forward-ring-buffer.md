# Story 2.2: Hybrid Store & Forward (Ring Buffer)

Status: done

## Story

As a Plant Manager,
I want zero data loss during network outages,
so that my historical records are complete.

## Acceptance Criteria

1. **Given** the Ingestor is disconnected from NATS
2. **When** new data arrives from sensors
3. **Then** it is stored in an in-memory Ring Buffer (e.g., circular queue)
4. **And** if RAM fills up (or buffer full), it spills to a local disk file (WAL - Write Ahead Log) or simply overwrites oldest (depending on strictness, but "Hybrid" implies disk spill) -> *Clarification: Requirement says "spills to local disk file".*
5. **When** connection is restored
6. **Then** buffered data is flushed to NATS in chronological order

## Tasks / Subtasks

- [x] Implement In-Memory Ring Buffer
  - [x] Create `struct RingBuffer<T>` with fixed capacity
  - [x] Implement `push(item: T)` and `pop() -> Option<T>`
- [x] Implement Disk Spilling (WAL)
  - [x] If buffer full, write to `buffer.wal` file
  - [x] Implement `DiskBuffer` struct
- [x] Implement Hybrid Buffer Manager
  - [x] Combine Memory and Disk buffers
  - [x] Logic: Push to Mem -> if Full -> Push to Disk
- [x] Implement NATS Publisher Task
  - [x] Create async task that reads from Buffer and publishes to NATS
  - [x] Handle NATS connection status (connected/disconnected)
  - [x] Retry logic for publishing
- [x] Integration
  - [x] Connect `ModbusAdapter` output to `Buffer` input
  - [x] Connect `Buffer` output to `NATS Publisher`
- [x] Testing
  - [x] Unit test: Ring Buffer wrap-around logic
  - [x] Unit test: Disk spill and recovery
  - [x] Integration test: Simulate NATS outage and verify data flush

## Dev Notes

### Technical Stack
- **Language:** Rust
- **Crate:** `tokio` (async), `serde` (serialization)
- **Disk I/O:** `tokio::fs` or `std::fs` (careful with blocking)

### Architecture
- **Ingestor Pipeline:** Modbus -> Channel -> Buffer -> NATS
- **Buffer Strategy:**
  - **Memory:** Fast, limited size (e.g., 10,000 points).
  - **Disk:** Slower, larger size. Used only when NATS is down or Mem is full.

### References
- [Epic 2 Details](docs/epics.md#Epic-2-High-Performance-Ingestion-Ingestor)

## File List

- services/ingestor/src/buffer.rs
- services/ingestor/src/publisher.rs
- services/ingestor/src/main.rs
- services/ingestor/src/config.rs
- services/ingestor/Cargo.toml
- config/default.toml

## Change Log

- 2025-12-03: Implemented Hybrid Store & Forward Ring Buffer with Disk Spilling (WAL) and NATS Publisher. (Ahmet)

## Dev Agent Record

### Implementation Plan - Hybrid Buffer

- **Approach:**
  - Implemented `RingBuffer` (in-memory) using `VecDeque`.
  - Implemented `DiskBuffer` (WAL) using `tokio::fs` and `prost` for Protobuf serialization.
  - Implemented `HybridBuffer` to manage spilling: when memory is full, the oldest item is popped and written to disk.
  - Implemented `Publisher` task that consumes from `HybridBuffer` and publishes to NATS. It prioritizes flushing disk buffer when connected.
- **Testing:**
  - Added unit tests for `RingBuffer` overflow logic.
  - Added async tests for `DiskBuffer` read/write.
  - Added async tests for `HybridBuffer` spill logic.
- **Decisions:**
  - Used `prost` for efficient binary serialization to disk.
  - Used `std::io::Read` for synchronous parsing of loaded file content to avoid async complexity with `Cursor`.
  - Implemented "flush all" for disk buffer for MVP simplicity, with a TODO for chunked reading in future.

### Context Reference

- **Story ID:** 2.2
- **Story Key:** 2-2-hybrid-store-forward-ring-buffer

### Agent Model Used
- Gemini 2.0 Flash

### Completion Notes List
- Ultimate context engine analysis completed - comprehensive developer guide created
