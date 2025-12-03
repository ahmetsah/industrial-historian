# historian - Epic Breakdown

**Author:** Ahmet
**Date:** 2025-12-02
**Project Level:** Industrial IoT Platform / SaaS / Embedded
**Target Scale:** >50k events/sec, Edge Deployment

---

## Overview

This document provides the complete epic and story breakdown for historian, decomposing the requirements from the [PRD](./PRD.md) into implementable stories.

**Living Document Notice:** This is the initial version. It will be updated after UX Design and Architecture workflows add interaction and technical details to stories.

### Epics Summary

*   **Epic 0: Foundation & Infrastructure (Scaffolding):** Setting up the Polyglot Monorepo, Docker Compose, NATS, and CI/CD pipelines.
*   **Epic 1: Secure Access & Identity (Auth):** Implementing the Go-based Auth service with JWT, RBAC, and FDA-compliant Re-authentication.
*   **Epic 2: High-Performance Ingestion (Ingestor):** Rust Ingestor implementation with Modbus/OPC UA adapters and Store & Forward logic.
*   **Epic 3: Efficient Storage Engine (Engine):** Rust Engine implementation with LSM-Tree, Gorilla Compression, and gRPC Query interface.
*   **Epic 4: Real-Time Visualization (Viz):** React/Vite frontend with uPlot, Zustand state management, and NATS WebSocket integration.
*   **Epic 5: Compliance & Safety (Audit & Alarm):** Implementing the Go-based Audit service (Immutable Logs) and Alarm service (ISA 18.2 State Machine).
*   **Epic 6: Predictive Intelligence (Sim):** Python-based Digital Twin integration for real-time inference and anomaly detection.

---

## Functional Requirements Inventory

### Data Ingestion & Management
*   **FR-ING-01 (Protocol Support):** Support Modbus TCP, OPC UA, and Siemens S7 protocols.
*   **FR-ING-02 (Store & Forward):** RAM buffering and Spill-to-Disk for zero data loss during network outages.
*   **FR-ING-03 (Calculated Tags):** Real-time mathematical and logical operations for virtual tags.

### Visualization & Analysis
*   **FR-VIS-01 (Trend Analysis):** Interactive charts (Zoom/Pan) for 1 hour to 10 years data range.
*   **FR-VIS-02 (High-Perf Export):** Download <1M rows in CSV/XLSX in <5 seconds.
*   **FR-VIS-03 (Confidence Interval):** Visual distinction for predicted data (dashed lines, shaded areas).

### Alarm Management (ISA 18.2)
*   **FR-ALM-01 (Alarm Lifecycle):** Full ISA 18.2 state machine (Unack/Active, Ack/Active, etc.).
*   **FR-ALM-02 (Alarm Shelving):** Temporary suppression of alarms by operators.

### Digital Twin & Prediction
*   **FR-DT-01 (Inference Engine):** Real-time anomaly detection and value prediction using pre-trained models.
*   **FR-DT-02 (Model Accuracy):** Display confidence scores based on historical performance.

### Audit & Compliance (FDA Part 11)
*   **FR-AUD-01 (Immutable Logging):** Chained hash logging for all critical user actions.
*   **FR-AUD-02 (Re-Authentication):** Electronic signature (password re-entry) for critical changes.

### System Administration & Security
*   **FR-SYS-01 (Tiered Storage):** Automated movement of old data to S3/MinIO.
*   **FR-SYS-02 (RBAC):** Role-based access control (Operator, Engineer, Admin, Auditor, Service Account).

---

## Epic Structure & Technical Context

### Epic 0: Foundation & Infrastructure (Scaffolding)
*   **Goal:** Establish the "Walking Skeleton" of the system, enabling all other services to be built, deployed, and communicate.
*   **User Value:** Developers can build features; Ops can deploy the system.
*   **Technical Context:**
    *   **Monorepo:** Initialize `historian` root with `crates/`, `services/`, `viz/`, `ops/`.
    *   **Messaging:** NATS JetStream setup in Docker Compose with correct subject hierarchy configuration (`enterprise.>`).
    *   **CI/CD:** GitHub Actions for multi-arch builds (ARM64/AMD64).
    *   **Shared Libs:** `historian-core` crate with Protobuf definitions.

### Epic 1: Secure Access & Identity (Auth)
*   **Goal:** Implement a secure, FDA-compliant authentication system from Day 1.
*   **User Value:** Users can securely log in; System meets basic access control regulations.
*   **PRD Coverage:** FR-SYS-02 (RBAC), FR-AUD-02 (Re-Authentication).
*   **Technical Context:**
    *   **Service:** `go-services/auth`.
    *   **Tech:** Go, JWT (RS256), bcrypt.
    *   **Data:** Users/Roles in relational DB (Postgres/SQLite).
    *   **Integration:** Publishes `user.login` events to NATS.

### Epic 2: High-Performance Ingestion (Ingestor)
*   **Goal:** Enable high-speed data collection from industrial devices with zero data loss guarantees.
*   **User Value:** System captures real-world data reliably, even during network glitches.
*   **PRD Coverage:** FR-ING-01 (Protocols), FR-ING-02 (Store & Forward), FR-ING-03 (Calculated Tags).
*   **Technical Context:**
    *   **Service:** `services/ingestor` (Rust).
    *   **Tech:** `tokio` for async I/O, `modbus` crate.
    *   **Pattern:** "Hybrid Spill-to-Disk" Ring Buffer.
    *   **Integration:** Publishes to `enterprise.site.area.line.device.sensor` on NATS.

### Epic 3: Efficient Storage Engine (Engine)
*   **Goal:** Store massive amounts of time-series data efficiently and retrieve it instantly.
*   **User Value:** Users can access years of history without waiting.
*   **PRD Coverage:** FR-SYS-01 (Tiered Storage), NFR-PERF-02 (Query Latency).
*   **Technical Context:**
    *   **Service:** `services/engine` (Rust).
    *   **Tech:** LSM-Tree (custom or RocksDB wrapper), Gorilla Compression.
    *   **API:** gRPC for internal query, GraphQL for frontend.
    *   **Data:** `data/tsdb/` directory.

### Epic 4: Real-Time Visualization (Viz)
*   **Goal:** Provide a responsive, modern interface for monitoring and analysis.
*   **User Value:** Operators and Engineers can visualize trends and make decisions.
*   **PRD Coverage:** FR-VIS-01 (Trend Analysis), FR-VIS-02 (Export), FR-VIS-03 (Confidence Interval).
*   **Technical Context:**
    *   **App:** `viz` (React + Vite).
    *   **Tech:** `uPlot` for charts, `Zustand` for state.
    *   **Pattern:** WebSocket connection to NATS (via Gateway) for live data.

### Epic 5: Compliance & Safety (Audit & Alarm)
*   **Goal:** Satisfy critical industrial regulations for safety and data integrity.
*   **User Value:** Quality team passes audits; Operators manage alarms effectively.
*   **PRD Coverage:** FR-ALM-01 (Lifecycle), FR-ALM-02 (Shelving), FR-AUD-01 (Immutable Logging).
*   **Technical Context:**
    *   **Services:** `go-services/audit`, `go-services/alarm`.
    *   **Tech:** Go.
    *   **Audit:** Chained Hash implementation for log integrity.
    *   **Alarm:** State Machine implementation (ISA 18.2).

### Epic 6: Predictive Intelligence (Sim)
*   **Goal:** Transform the system from reactive to proactive with Digital Twin capabilities.
*   **User Value:** Operators receive early warnings before failures occur.
*   **PRD Coverage:** FR-DT-01 (Inference), FR-DT-02 (Accuracy).
*   **Technical Context:**
    *   **Service:** `services/sim` (Python).
    *   **Tech:** GEKKO, NumPy.
    *   **Pattern:** Subscribes to live data, publishes `...predicted` topics.

---

## FR Coverage Map

*   **FR-ING-01:** Epic 2 (Story 2.1, 2.2)
*   **FR-ING-02:** Epic 2 (Story 2.3)
*   **FR-ING-03:** Epic 2 (Story 2.4)
*   **FR-VIS-01:** Epic 4 (Story 4.2)
*   **FR-VIS-02:** Epic 4 (Story 4.3)
*   **FR-VIS-03:** Epic 4 (Story 4.2)
*   **FR-ALM-01:** Epic 5 (Story 5.2)
*   **FR-ALM-02:** Epic 5 (Story 5.2)
*   **FR-DT-01:** Epic 6 (Story 6.1)
*   **FR-DT-02:** Epic 6 (Story 6.2)
*   **FR-AUD-01:** Epic 5 (Story 5.1)
*   **FR-AUD-02:** Epic 1 (Story 1.2)
*   **FR-SYS-01:** Epic 3 (Story 3.2)
*   **FR-SYS-02:** Epic 1 (Story 1.1)

---

## Epic 0: Foundation & Infrastructure (Scaffolding)

**Goal:** Establish the "Walking Skeleton" of the system, enabling all other services to be built, deployed, and communicate. This includes the Polyglot Monorepo structure, NATS JetStream infrastructure, and shared Protocol Buffer definitions.

### Story 0.1: Monorepo & Polyglot Workspace Initialization

**As a** Developer,
**I want** a structured Polyglot Monorepo (Rust, Go, TS),
**So that** I can develop multiple services in a unified environment with shared tooling.

**Acceptance Criteria:**
*   **Given** a fresh git repository
*   **When** I initialize the project structure
*   **Then** the following directories exist: `crates/`, `services/`, `go-services/`, `viz/`, `ops/`
*   **And** `Cargo.toml` is configured as a Rust Workspace including `crates/*` and `services/*`
*   **And** `go-services/auth` and `go-services/audit` are initialized as Go modules
*   **And** `viz/` is initialized as a Vite + React + TS project
*   **And** a root `Makefile` exists with commands to build all components (`make build-rust`, `make build-go`, `make build-viz`)

**Technical Notes:**
*   Follow Architecture Directory Structure exactly.
*   Rust Edition: 2024 (or latest stable).
*   Go Version: 1.21+.
*   Node Version: 20+ (LTS).

### Story 0.2: Infrastructure as Code (Docker Compose & NATS)

**As a** DevOps Engineer,
**I want** a Docker Compose environment with NATS JetStream,
**So that** I can run the entire distributed system locally for development.

**Acceptance Criteria:**
*   **Given** the `ops/` directory
*   **When** I run `docker-compose up -d`
*   **Then** NATS JetStream is running on port 4222
*   **And** NATS Management UI (if available) or CLI tool can connect
*   **And** A Stream named `EVENTS` is created with subject `enterprise.>`
*   **And** MinIO (S3 compatible) is running for future object storage needs
*   **And** PostgreSQL is running for the Auth service

**Technical Notes:**
*   Use `nats:latest` image with JetStream enabled (`-js`).
*   Configure NATS to use file-based storage for persistence.
*   Define NATS configuration in `ops/nats.conf`.

### Story 0.3: Core Library & Protobuf Schema Setup

**As a** Backend Developer,
**I want** a shared library with Protobuf definitions,
**So that** Rust and Go services can communicate using a strictly typed schema.

**Acceptance Criteria:**
*   **Given** the `crates/historian-core` library
*   **When** I define `.proto` files in `crates/historian-core/src/proto/`
*   **Then** `build.rs` automatically generates Rust structs
*   **And** I can import these structs in `services/ingestor`
*   **And** A script `scripts/gen-go-proto.sh` generates Go structs for `go-services/`
*   **And** The schema includes basic messages: `SensorData`, `LogEntry`, `UserAction`

**Technical Notes:**
*   Use `prost` and `tonic-build` for Rust generation.
*   Use `protoc-gen-go` for Go generation.
*   Define package name as `historian.v1`.

### Story 0.4: CI/CD Pipeline Foundation

**As a** Team Lead,
**I want** a GitHub Actions pipeline,
**So that** every commit is automatically built and linted to prevent regressions.

**Acceptance Criteria:**
*   **Given** a push to `main` or a PR
*   **When** the CI pipeline runs
*   **Then** it builds all Rust binaries (`cargo build`)
*   **And** it runs Rust tests (`cargo test`) and linter (`cargo clippy`)
*   **And** it builds Go binaries (`go build`) and runs `go vet`
*   **And** it builds the Frontend (`npm run build`)
*   **And** the pipeline fails if any step fails

**Technical Notes:**
*   Use caching for Cargo registry and `node_modules` to speed up builds.
*   Separate jobs for Rust, Go, and Viz to run in parallel.

---

## Epic 1: Secure Access & Identity (Auth)

**Goal:** Implement a secure, FDA-compliant authentication system from Day 1. This ensures that all subsequent development happens within a secure context.

### Story 1.1: Auth Service & JWT Implementation

**As a** System Administrator,
**I want** a centralized authentication service,
**So that** I can manage users and secure API access.

**Acceptance Criteria:**
*   **Given** the `go-services/auth` service
*   **When** I POST to `/api/v1/login` with valid credentials
*   **Then** I receive a JWT (RS256 signed) with `sub`, `role`, and `exp` claims
*   **And** the login event is published to NATS subject `sys.auth.login`
*   **And** invalid credentials return 401 Unauthorized

**Technical Notes:**
*   Use `golang-jwt/jwt` for token generation.
*   Store users in PostgreSQL with `bcrypt` hashed passwords.
*   Generate RSA keys on startup if missing (for dev) or load from secrets (prod).

### Story 1.2: FDA Re-Authentication (Electronic Signature)

**As a** Quality Manager,
**I want** a re-authentication mechanism for critical actions,
**So that** the system complies with FDA 21 CFR Part 11 "Electronic Signature" requirements.

**Acceptance Criteria:**
*   **Given** an authenticated user performing a critical action (e.g., changing alarm limits)
*   **When** the UI prompts for re-authentication
*   **And** the user enters their password again
*   **Then** the Auth service verifies the password
*   **And** issues a short-lived "Signing Token" (valid for 1 minute)
*   **And** this token is logged in the Audit Trail as a signature

**Technical Notes:**
*   Endpoint: `POST /api/v1/re-auth`.
*   This is distinct from the initial login session.

### Story 1.3: RBAC & Service Accounts

**As a** Security Officer,
**I want** Role-Based Access Control,
**So that** operators cannot change engineering settings.

**Acceptance Criteria:**
*   **Given** a user with role `OPERATOR`
*   **When** they try to access an `ADMIN` only endpoint
*   **Then** the request is denied (403 Forbidden)
*   **And** I can create "Service Accounts" (API Keys) for Ingestors to authenticate against NATS

**Technical Notes:**
*   Roles: `ADMIN`, `ENGINEER`, `OPERATOR`, `AUDITOR`, `SERVICE`.
*   Implement middleware in Go to check JWT `role` claim.

---

## Epic 2: High-Performance Ingestion (Ingestor)

**Goal:** Enable high-speed data collection from industrial devices with zero data loss guarantees.

### Story 2.1: Modbus TCP Adapter (Rust)

**As a** Control Engineer,
**I want** to collect data from Modbus TCP devices,
**So that** I can monitor legacy PLCs.

**Acceptance Criteria:**
*   **Given** a configuration file listing Modbus registers
*   **When** the Ingestor service starts
*   **Then** it connects to the Modbus TCP server
*   **And** polls registers at the defined interval (e.g., 100ms)
*   **And** converts raw bytes to `SensorData` Protobuf messages
*   **And** pushes them to the internal Ring Buffer

**Technical Notes:**
*   Use `tokio-modbus` crate.
*   Ensure connection retry logic is robust (exponential backoff).

### Story 2.2: Hybrid Store & Forward (Ring Buffer)

**As a** Plant Manager,
**I want** zero data loss during network outages,
**So that** my historical records are complete.

**Acceptance Criteria:**
*   **Given** the Ingestor is disconnected from NATS
*   **When** new data arrives from sensors
*   **Then** it is stored in an in-memory Ring Buffer
*   **And** if RAM fills up, it spills to a local disk file (WAL)
*   **When** connection is restored
*   **Then** buffered data is flushed to NATS in chronological order

**Technical Notes:**
*   Implement a custom `Buffer` struct in Rust.
*   Use a separate async task for flushing to NATS.

### Story 2.3: Calculated Tags Engine

**As a** Process Engineer,
**I want** to define virtual tags based on math formulas,
**So that** I can monitor derived values (e.g., Efficiency = Output / Input).

**Acceptance Criteria:**
*   **Given** a config defining `Tag C = Tag A + Tag B`
*   **When** `Tag A` or `Tag B` changes
*   **Then** `Tag C` is automatically recalculated
*   **And** published to NATS as a new data point

**Technical Notes:**
*   Use a lightweight expression parser (e.g., `evalexpr` crate).
*   Processing must happen in microsecond range.

---

## Epic 3: Efficient Storage Engine (Engine)

**Goal:** Store massive amounts of time-series data efficiently and retrieve it instantly.

### Story 3.1: LSM-Tree Storage Engine (Rust)

**As a** Database Administrator,
**I want** a write-optimized storage engine,
**So that** I can handle >50k events/second without disk I/O bottlenecks.

**Acceptance Criteria:**
*   **Given** a stream of `SensorData` from NATS
*   **When** the Engine receives data
*   **Then** it writes to an in-memory MemTable
*   **And** flushes to SSTables on disk when full
*   **And** applies Gorilla (XOR) compression to floating point values

**Technical Notes:**
*   Use `rocksdb` binding or a pure Rust LSM implementation (e.g., `agatedb` or custom).
*   Key format: `[SensorID][Timestamp]`.

### Story 3.2: Tiered Storage Manager

**As a** CFO,
**I want** old data moved to cheaper storage,
**So that** we don't waste expensive SSD space on 5-year-old logs.

**Acceptance Criteria:**
*   **Given** data older than 1 month
*   **When** the Tiering Job runs
*   **Then** it moves SSTables from local SSD to MinIO (S3) bucket
*   **And** updates the index to point to the remote location
*   **And** queries seamlessly fetch data from either source

**Technical Notes:**
*   Use `rust-s3` crate.
*   Implement transparent query federation.

### Story 3.3: gRPC Query API

**As a** Frontend Developer,
**I want** a fast API to query historical data,
**So that** I can populate charts.

**Acceptance Criteria:**
*   **Given** a gRPC request for `Sensor X` from `T1` to `T2`
*   **When** the Engine processes the request
*   **Then** it performs automatic downsampling (LTTB algorithm) if the point count > 1000
*   **And** returns the data stream in < 100ms

**Technical Notes:**
*   Define `service HistorianQuery` in Protobuf.
*   Implement LTTB (Largest-Triangle-Three-Buckets) for visual downsampling.

---

## Epic 4: Real-Time Visualization (Viz)

**Goal:** Provide a responsive, modern interface for monitoring and analysis.

### Story 4.1: Real-time Dashboard Framework

**As a** Operator,
**I want** a customizable dashboard,
**So that** I can arrange charts relevant to my machine.

**Acceptance Criteria:**
*   **Given** the React application
*   **When** I drag and drop a "Chart Widget"
*   **Then** I can configure it to listen to a specific NATS subject
*   **And** the layout is saved to local storage

**Technical Notes:**
*   Use `react-grid-layout`.
*   Use `Zustand` for managing the layout state.

### Story 4.2: Trend Chart Component (uPlot)

**As a** Engineer,
**I want** high-performance charts,
**So that** I can zoom into millisecond-level data without the browser freezing.

**Acceptance Criteria:**
*   **Given** a chart displaying 100,000 points
*   **When** I zoom or pan
*   **Then** the interaction is smooth (60 FPS)
*   **And** new real-time points append to the right side instantly

**Technical Notes:**
*   Wrap `uPlot` in a React component.
*   Ensure memory management (destroy chart instances on unmount).

### Story 4.3: Data Export Service

**As a** Analyst,
**I want** to export data to Excel,
**So that** I can perform offline analysis.

**Acceptance Criteria:**
*   **Given** a selected time range on a chart
*   **When** I click "Export CSV"
*   **Then** the browser downloads a file containing the raw data
*   **And** the generation takes < 5 seconds for 1M rows

**Technical Notes:**
*   Stream data from gRPC directly to a browser stream if possible, or generate on backend.

---

## Epic 5: Compliance & Safety (Audit & Alarm)

**Goal:** Satisfy critical industrial regulations for safety and data integrity.

### Story 5.1: Immutable Audit Service (Go)

**As a** Auditor,
**I want** proof that logs haven't been tampered with,
**So that** we pass the FDA inspection.

**Acceptance Criteria:**
*   **Given** a user action (e.g., "Changed Setpoint")
*   **When** the Audit service receives the log
*   **Then** it calculates a hash including the *previous* log's hash (Chained Hash)
*   **And** stores it in the database
*   **And** provides an API to verify the chain integrity

**Technical Notes:**
*   Use SHA-256.
*   Store in `audit_logs` table.

### Story 5.2: Alarm State Machine (ISA 18.2)

**As a** Operator,
**I want** a standard alarm system,
**So that** I don't miss critical failures.

**Acceptance Criteria:**
*   **Given** an alarm definition (High Temp > 100)
*   **When** the value exceeds 100
*   **Then** the alarm state transitions to `Unack/Active`
*   **And** when I click "Acknowledge", it moves to `Ack/Active`
*   **And** when value drops below 100, it moves to `Ack/RTN` (Return to Normal)

**Technical Notes:**
*   Implement full ISA 18.2 state diagram.
*   Persist state in Postgres/SQLite to survive restarts.

---

## Epic 6: Predictive Intelligence (Sim)

**Goal:** Transform the system from reactive to proactive with Digital Twin capabilities.

### Story 6.1: Digital Twin Service (Python)

**As a** Engineer,
**I want** to run a physics-based simulation,
**So that** I can compare actual vs. ideal performance.

**Acceptance Criteria:**
*   **Given** a GEKKO model of a reactor
*   **When** the Sim service starts
*   **Then** it subscribes to real-time inputs from NATS
*   **And** runs the model solver
*   **And** publishes the calculated "Ideal Output" back to NATS

**Technical Notes:**
*   Use `gekko` Python library.
*   Run in a separate Docker container (`services/sim`).

### Story 6.2: Inference & Anomaly Detection

**As a** Maintenance Tech,
**I want** to know when a machine is behaving abnormally,
**So that** I can fix it before it breaks.

**Acceptance Criteria:**
*   **Given** the real value and the simulated "Ideal" value
*   **When** the difference (residual) exceeds a threshold
*   **Then** an "Anomaly" event is published
*   **And** a confidence score is calculated

**Technical Notes:**
*   Simple residual analysis initially.
*   Publish to `sys.analytics.anomaly`.

---

## FR Coverage Matrix

| FR ID | Description | Covered By |
| :--- | :--- | :--- |
| **FR-ING-01** | Protocol Support | Story 2.1, 2.2 |
| **FR-ING-02** | Store & Forward | Story 2.2 |
| **FR-ING-03** | Calculated Tags | Story 2.3 |
| **FR-VIS-01** | Trend Analysis | Story 4.1, 4.2 |
| **FR-VIS-02** | High-Perf Export | Story 4.3 |
| **FR-VIS-03** | Confidence Interval | Story 4.2 |
| **FR-ALM-01** | Alarm Lifecycle | Story 5.2 |
| **FR-ALM-02** | Alarm Shelving | Story 5.2 |
| **FR-DT-01** | Inference Engine | Story 6.1, 6.2 |
| **FR-DT-02** | Model Accuracy | Story 6.2 |
| **FR-AUD-01** | Immutable Logging | Story 5.1 |
| **FR-AUD-02** | Re-Authentication | Story 1.2 |
| **FR-SYS-01** | Tiered Storage | Story 3.2 |
| **FR-SYS-02** | RBAC | Story 1.3 |

---

## Summary

**Total Epics:** 7 (0-6)
**Total Stories:** 19
**Readiness:** Ready for Sprint Planning.

_This document is now complete and ready for the Development Team._
