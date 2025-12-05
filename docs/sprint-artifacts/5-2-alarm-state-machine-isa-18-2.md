# Story 5.2: Alarm State Machine (ISA 18.2)

Status: Ready for Review

## Story

**As a** Operator,
**I want** a standard alarm system,
**So that** I don't miss critical failures.

## Acceptance Criteria

1.  **Given** an alarm definition (e.g., "High Temp > 100")
2.  **When** the monitored value exceeds the threshold
3.  **Then** the alarm state transitions to `Unack/Active`
4.  **And** an event is published to NATS `sys.alarm.events`
5.  **When** the operator acknowledges the alarm
6.  **Then** the state moves to `Ack/Active`
7.  **When** the value drops below the threshold
8.  **Then** the state moves to `Ack/RTN` (Return to Normal) or `Normal` depending on the sequence
9.  **And** the state is persisted to Postgres to survive service restarts

## Tasks / Subtasks

- [x] **Service Initialization**
  - [x] Initialize `go-services/alarm` module
  - [x] Configure `config` package (Env vars: DB_URL, NATS_URL)
  - [x] Setup `main.go` with graceful shutdown

- [x] **Database Implementation**
  - [x] Design `alarm_definitions` table (id, tag, threshold, type, priority)
  - [x] Design `active_alarms` table (id, definition_id, state, activation_time, ack_time, value)
  - [x] Implement migration (using `golang-migrate` or raw SQL)
  - [x] Implement `Repository` for CRUD operations

- [x] **State Machine Logic (ISA 18.2)**
  - [x] Define States: `Normal`, `UnackActive`, `AckActive`, `UnackRTN`, `Shelved`, `Suppressed`
  - [x] Implement `AlarmFSM` struct with strict transition logic
  - [x] Implement "Shelving" logic (temporary suppression with timeout)

- [x] **NATS Integration**
  - [x] **Consumer:** Subscribe to `enterprise.>` (wildcard) to monitor all sensor values
  - [x] **Evaluator:** Check incoming values against loaded `alarm_definitions`
  - [x] **Producer:** Publish transitions to `sys.alarm.events`

- [x] **API Implementation**
  - [x] `POST /api/v1/alarms/:id/ack` (Acknowledge)
  - [x] `POST /api/v1/alarms/:id/shelve` (Shelve)
  - [x] `GET /api/v1/alarms/active` (List active alarms)
  - [x] `POST /api/v1/alarms/definitions` (CRUD for definitions)

## Dev Notes

### Technical Requirements
*   **Language:** Go 1.21+
*   **Database:** PostgreSQL (use `pgx` driver)
*   **State Machine:** Implement a custom FSM or use a lightweight library. Ensure it supports the full ISA 18.2 cycle including `Unack/RTN`.
*   **Performance:** The Evaluator must be fast. Cache `alarm_definitions` in memory and refresh on change.

### Architecture Compliance
*   **Location:** `go-services/alarm/`
*   **Communication:** NATS for input (sensor data) and output (events). HTTP for Operator actions.
*   **Naming:** Use `sys.alarm.events` for published events.

### Project Structure
```
go-services/alarm/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── core/
│   │   ├── domain.go      # Alarm, AlarmDefinition structs
│   │   └── fsm.go         # ISA 18.2 Logic
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
*   [Architecture: Alarm Service](./../architecture.md#epic-fda-compliance-auth--audit)

## Dev Agent Record

### Context Reference
*   **Architecture:** `docs/architecture.md`
*   **Epics:** `docs/epics.md`
*   **Previous Story:** `docs/sprint-artifacts/5-1-immutable-audit-service-go.md`

### Agent Model Used
Antigravity (Google Deepmind)

### Debug Log References

### Completion Notes List
- Implemented Alarm Service in Go
- Designed Postgres schema with `alarm_definitions` and `active_alarms` tables
- Implemented ISA 18.2 State Machine (Normal, UnackActive, AckActive, UnackRTN, Shelved)
- Implemented Shelving logic with `shelved_until`
- Implemented NATS Consumer for sensor data and Publisher for alarm events
- Implemented HTTP API for Acknowledge, Shelve, List Active, and Create Definition
- Added unit tests for FSM, Evaluator, and Service logic
- Verified all acceptance criteria

### File List
- go-services/alarm/go.mod
- go-services/alarm/go.sum
- go-services/alarm/internal/config/config.go
- go-services/alarm/internal/config/config_test.go
- go-services/alarm/cmd/server/main.go
- go-services/alarm/cmd/server/main_test.go
- go-services/alarm/Dockerfile
- go-services/alarm/migrations/000001_init_schema.up.sql
- go-services/alarm/migrations/000001_init_schema.down.sql
- go-services/alarm/internal/core/domain.go
- go-services/alarm/internal/core/fsm.go
- go-services/alarm/internal/core/fsm_test.go
- go-services/alarm/internal/core/evaluator.go
- go-services/alarm/internal/core/evaluator_test.go
- go-services/alarm/internal/core/service.go
- go-services/alarm/internal/core/service_test.go
- go-services/alarm/internal/transport/nats_consumer.go
- go-services/alarm/internal/transport/http_handler.go
