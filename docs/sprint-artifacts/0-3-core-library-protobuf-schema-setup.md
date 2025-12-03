# Story 0.3: Core Library & Protobuf Schema Setup

Status: done

## Story

As a Backend Developer,
I want a shared library with Protobuf definitions,
so that Rust and Go services can communicate using a strictly typed schema.

## Acceptance Criteria

1. **Given** the `crates/historian-core` library
2. **When** I define `.proto` files in `crates/historian-core/src/proto/`
3. **Then** `build.rs` automatically generates Rust structs
4. **And** I can import these structs in `services/ingestor`
5. **And** A script `scripts/gen-go-proto.sh` generates Go structs for `go-services/`
6. **And** The schema includes basic messages: `SensorData`, `LogEntry`, `UserAction`

## Tasks / Subtasks

- [x] Setup `historian-core` Crate
  - [x] Add dependencies to `crates/historian-core/Cargo.toml`: `prost`, `tonic`, `serde`
  - [x] Add build dependencies: `tonic-build`
  - [x] Create `crates/historian-core/build.rs` to compile protos
- [x] Define Protobuf Schemas
  - [x] Create `crates/historian-core/src/proto/common.proto` (SensorData, LogEntry, UserAction)
  - [x] Define package `historian.v1`
  - [x] Define messages with correct types and field numbers
- [x] Implement Rust Generation
  - [x] Configure `build.rs` to output to `OUT_DIR`
  - [x] Re-export generated modules in `crates/historian-core/src/lib.rs`
- [x] Implement Go Generation
  - [x] Create `scripts/gen-go-proto.sh`
  - [x] Install `protoc-gen-go` and `protoc-gen-go-grpc` tools
  - [x] Script should generate Go code into `go-services/internal/proto` (or similar shared location)
- [x] Verify Integration
  - [x] Create a test in `historian-core` that instantiates a `SensorData` struct
  - [x] Verify Go structs are generated and compile

## Dev Notes

### Technical Stack Versions
- **Rust:** `prost` (latest), `tonic` (latest)
- **Go:** `google.golang.org/protobuf` (latest), `google.golang.org/grpc` (latest)
- **Protoc:** Ensure `protoc` compiler is installed or handled by `tonic-build`

### Schema Definitions
- **SensorData:**
  - `string sensor_id`
  - `double value`
  - `int64 timestamp_ms`
  - `int32 quality` (0=Bad, 1=Good)
- **LogEntry:**
  - `string level`
  - `string message`
  - `int64 timestamp_ms`
  - `string service`
- **UserAction:**
  - `string user_id`
  - `string action`
  - `string resource`
  - `int64 timestamp_ms`

### Architecture Compliance
- **Naming:** Follow `PascalCase` for messages, `snake_case` for fields.
- **Package:** `historian.v1`
- **Location:** All `.proto` files MUST live in `crates/historian-core/src/proto/`.

### References
- [Architecture Protobuf Patterns](docs/architecture.md#Protobuf-Naming)
- [Epic 0 Details](docs/epics.md#Epic-0-Foundation--Infrastructure-Scaffolding)

## Dev Agent Record

### Context Reference
- **Story ID:** 0.3
- **Story Key:** 0-3-core-library-protobuf-schema-setup

### Agent Model Used
- Gemini 2.0 Flash

### Completion Notes List
- Ultimate context engine analysis completed - comprehensive developer guide created
- Configured `historian-core` crate with `prost` and `tonic`.
- Created `crates/historian-core/build.rs` for Protobuf compilation.
- Defined `SensorData`, `LogEntry`, and `UserAction` in `crates/historian-core/src/proto/common.proto`.
- Implemented Rust code generation and re-exported modules in `crates/historian-core/src/lib.rs`.
- Verified `historian-core` builds successfully.
- Created `scripts/gen-go-proto.sh` and generated Go Protobuf code in `go-services/pkg/proto`.
- Verified Rust integration with unit tests in `historian-core`.
- Verified Go integration with `tests/proto_verification/main.go`.

## File List
- crates/historian-core/Cargo.toml
- crates/historian-core/build.rs
- crates/historian-core/src/proto/common.proto
- crates/historian-core/src/lib.rs
- scripts/gen-go-proto.sh
- go-services/pkg/proto/common.pb.go
- tests/proto_verification/main.go
- tests/proto_verification/go.mod

## Change Log
- 2025-12-03: Verified Rust and Go Protobuf integration.
- 2025-12-03: Implemented Go Protobuf generation.
- 2025-12-03: Implemented Rust Protobuf generation and verified build.
- 2025-12-03: Initialized Protobuf setup in `historian-core`.
