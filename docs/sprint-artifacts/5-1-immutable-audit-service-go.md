# Story 5.1: Immutable Audit Service (Go)

Status: ready-for-dev

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

- [ ] **Service Initialization**
  - [ ] Initialize `go-services/audit` module
  - [ ] Configure `config` package (Env vars: DB_URL, NATS_URL)
  - [ ] Setup `main.go` with graceful shutdown

- [ ] **Database Implementation**
  - [ ] Design `audit_logs` schema: `id` (ULID/UUID), `timestamp`, `actor`, `action`, `details` (JSONB), `prev_hash`, `curr_hash`
  - [ ] Implement migration (using `golang-migrate` or raw SQL on startup)
  - [ ] Implement `Repository` pattern for inserting logs
  - [ ] **Critical:** The insert transaction must lock the last row or use serializable isolation to ensure `prev_hash` is accurate under load.

- [ ] **Chained Hash Logic**
  - [ ] Implement `Hasher` service
  - [ ] Logic: `curr_hash = SHA256(prev_hash + timestamp + actor + action + details)`
  - [ ] Handle "Genesis Block" case (first log entry has zero `prev_hash`)

- [ ] **NATS Consumer**
  - [ ] Connect to NATS JetStream
  - [ ] Subscribe to `sys.auth.login` (from Auth Service)
  - [ ] Subscribe to `sys.audit.>` (Generic audit events)
  - [ ] Map incoming events to `LogEntry` model

- [ ] **Verification API**
  - [ ] Implement `GET /api/v1/audit/verify`
  - [ ] Logic: Re-calculate hashes from start (or a checkpoint) and compare with stored `curr_hash`
  - [ ] Return `valid: true` or `valid: false` with the ID of the first broken link

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
*   [ ] Verified chain integrity with test data
*   [ ] Tested race conditions (concurrent writes)
*   [ ] Verified NATS subscription works

### File List
*   `go-services/audit/go.mod`
*   `go-services/audit/cmd/server/main.go`
*   `go-services/audit/internal/core/hasher.go`
