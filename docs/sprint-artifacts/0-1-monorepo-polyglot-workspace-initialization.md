# Story 0.1: Monorepo & Polyglot Workspace Initialization

Status: done

## Story

As a Developer,
I want a structured Polyglot Monorepo (Rust, Go, TS),
so that I can develop multiple services in a unified environment with shared tooling.

## Acceptance Criteria

1. **Given** a fresh git repository
2. **When** I initialize the project structure
3. **Then** the following directories exist: `crates/`, `services/`, `go-services/`, `viz/`, `ops/`
4. **And** `Cargo.toml` is configured as a Rust Workspace including `crates/*` and `services/*`
5. **And** `go-services/auth` and `go-services/audit` are initialized as Go modules
6. **And** `viz/` is initialized as a Vite + React + TS project
7. **And** a root `Makefile` exists with commands to build all components (`make build-rust`, `make build-go`, `make build-viz`)

## Tasks / Subtasks

- [x] Initialize Project Root & Git
  - [x] Create project directory `historian` (if not exists) and `git init`
  - [x] Create `.gitignore` (Rust, Go, Node, IDEs)
- [x] Initialize Rust Workspace
  - [x] Create root `Cargo.toml` with `[workspace]` members: `crates/*`, `services/*`
  - [x] Create `crates/historian-core` (lib)
  - [x] Create `services/ingestor` (bin)
  - [x] Create `services/engine` (bin)
- [x] Initialize Go Services
  - [x] Create `go-services/auth` and `go mod init`
  - [x] Create `go-services/audit` and `go mod init`
- [x] Initialize Frontend
  - [x] Run `npm create vite@latest viz -- --template react-ts`
  - [x] Install dependencies in `viz/`
- [x] Infrastructure Setup
  - [x] Create `ops/` directory
  - [x] Create `ops/docker-compose.yml` (empty or basic structure)
- [x] Build Automation
  - [x] Create `Makefile` with targets: `build-rust`, `build-go`, `build-viz`, `dev-up`, `dev-down`

## Dev Notes

### Technical Stack Versions
- **Rust:** 1.91.1+ (Edition 2024)
- **Go:** 1.22+
- **Node.js:** 24.11.1+ (LTS)
- **Vite:** 7.2.6+

### Project Structure Notes
- Follow the structure defined in `docs/architecture.md` exactly.
- **Rust Workspace:** Ensure `crates/historian-core` is a library and `services/*` are binaries.
- **Go Modules:** Each service in `go-services/` is an independent module.
- **Frontend:** Use `viz/` as the root for the React application.

### Architecture Compliance
- **Monorepo:** Single repo for all services.
- **Workspace:** Rust workspace to deduplicate dependencies.
- **Polyglot:** Explicit separation of Rust (High Perf), Go (Business Logic), and TS (UI).

### References
- [Architecture Structure](docs/architecture.md#Complete-Project-Directory-Structure)
- [Epic 0 Details](docs/epics.md#Epic-0-Foundation--Infrastructure-Scaffolding)

## Dev Agent Record

### Context Reference
- **Story ID:** 0.1
- **Story Key:** 0-1-monorepo-polyglot-workspace-initialization

### Agent Model Used
- Gemini 2.0 Flash

### Completion Notes List
- Ultimate context engine analysis completed - comprehensive developer guide created
