# Story 5.1: Immutable Audit Service (Go)

Status: Done

## Story

**As a** Auditor,
**I want** proof that logs haven't been tampered with,
**So that** we pass the FDA inspection.

## Acceptance Criteria

1.  **Given** a user action (e.g., "Login", "Changed Setpoint")
2.  **When** the Audit service receives the log event via NATS
3.  **Then** it calculates a SHA-256 hash including the *previous* log's hash (Chained Hash)
4.  **And** stores the entry in the `audit_logs` PostgreSQL table
5.  **And** the chain integrity can be verified via an API endpoint
6.  **And** any tampering with a past row breaks the verification of all subsequent rows

## Tasks / Subtasks

- [x] **Service Initialization**
  - [x] Initialize `go-services/audit` module
  - [x] Configure `config` package (Env vars: DB_URL, NATS_URL)
  - [x] Setup `main.go` with graceful shutdown

- [x] **Database Implementation**
  - [x] Design `audit_logs` schema: `id` (ULID/UUID), `timestamp`, `actor`, `action`, `details` (JSONB), `prev_hash`, `curr_hash`
  - [x] Implement migration (using `golang-migrate` or raw SQL on startup)
  - [x] Implement `Repository` pattern for inserting logs
  - [x] **Critical:** The insert transaction must lock the last row or use serializable isolation to ensure `prev_hash` is accurate under load.

- [x] **Chained Hash Logic**
  - [x] Implement `Hasher` service
  - [x] Logic: `curr_hash = SHA256(prev_hash + timestamp + actor + action + details)`
  - [x] Handle "Genesis Block" case (first log entry has zero `prev_hash`)

- [x] **NATS Consumer**
  - [x] Connect to NATS JetStream
  - [x] Subscribe to `sys.auth.login` (from Auth Service)
  - [x] Subscribe to `sys.audit.>` (Generic audit events)
  - [x] Map incoming events to `LogEntry` model

- [x] **Verification API**
  - [x] Implement `GET /api/v1/audit/verify`
  - [x] Logic: Re-calculate hashes from start (or a checkpoint) and compare with stored `curr_hash`
  - [x] Return `valid: true` or `valid: false` with the ID of the first broken link

## Dev Notes

### Technical Requirements
*   **Language:** Go 1.21+
*   **Database:** PostgreSQL (use `pgx` driver)
*   **Hashing:** `crypto/sha256`
*   **Concurrency:** Be careful with the "Get Last Hash -> Calculate New -> Insert" race condition. Use a database lock or single-threaded writer channel if necessary.

### Architecture Compliance
*   **Location:** `go-services/audit/`
*   **Communication:** NATS for writing, HTTP for verification.
*   **Security:** This service is critical. Ensure SQL injection protection (use parameterized queries).

### Project Structure
```
go-services/audit/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── core/
│   │   ├── domain.go      # LogEntry struct
│   │   └── hasher.go      # SHA256 logic
│   ├── repository/
│   │   └── postgres.go
│   └── transport/
│       ├── nats_consumer.go
│       └── http_handler.go
├── go.mod
└── Dockerfile
```

### References
*   [Epics: Epic 5](./../epics.md#epic-5-compliance--safety-audit--alarm)
*   [Architecture: Audit Service](./../architecture.md#epic-fda-compliance-auth--audit)

## Dev Agent Record

### Context Reference
*   **Architecture:** `docs/architecture.md`
*   **Epics:** `docs/epics.md`

### Agent Model Used
Antigravity (Google Deepmind)

### Completion Notes List
*   [x] Verified chain integrity with test data
*   [x] Tested race conditions (concurrent writes)
*   [x] Verified NATS subscription works

### Implementation Plan
*   **Service Initialization:** Created basic structure, config loading with env vars, and graceful shutdown in main.go.
*   **Database Implementation:** Implemented Postgres repository with Serializable isolation for chained hash integrity. Created migration for audit_logs table.
*   **Chained Hash Logic:** Implemented SHA256Hasher with deterministic hashing of log fields.
*   **NATS Consumer:** Implemented NATS JetStream consumer to ingest audit events and persist them using the repository.
*   **Verification API:** Implemented HTTP endpoint to verify the integrity of the audit log chain.

### File List
*   `go-services/audit/go.mod`
*   `go-services/audit/cmd/server/main.go`
*   `go-services/audit/internal/config/config.go`
*   `go-services/audit/internal/config/config_test.go`
*   `go-services/audit/internal/core/domain.go`
*   `go-services/audit/internal/core/hasher.go`
*   `go-services/audit/internal/core/hasher_test.go`
*   `go-services/audit/internal/core/sha256_hasher.go`
*   `go-services/audit/internal/repository/postgres.go`
*   `go-services/audit/internal/repository/postgres_test.go`
*   `go-services/audit/internal/transport/http_handler.go`
*   `go-services/audit/internal/transport/nats_consumer.go`
*   `go-services/audit/migrations/001_create_audit_logs_table.sql`

### Change Log
*   Implemented Audit Service with Chained Hashing.
*   Added NATS consumer for `sys.audit.>` and `sys.auth.login`.
*   Added HTTP verification endpoint.
*   Added Postgres repository with Serializable isolation.
*   **Fix (AI Review):** Fixed timestamp precision mismatch (Microsecond truncation) to ensure hash verification passes.
*   **Fix (AI Review):** Added retry logic for `40001` serialization failures to handle concurrent writes.
*   **Fix (AI Review):** Added `LIMIT 1000` to `IterateLogs` to prevent potential DoS.
