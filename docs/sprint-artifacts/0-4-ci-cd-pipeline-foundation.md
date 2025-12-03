# Story 0.4: CI/CD Pipeline Foundation

Status: done

## Story

As a Team Lead,
I want a GitHub Actions pipeline,
so that every commit is automatically built and linted to prevent regressions.

## Acceptance Criteria

1. **Given** a push to `main` or a PR
2. **When** the CI pipeline runs
3. **Then** it builds all Rust binaries (`cargo build`)
4. **And** it runs Rust tests (`cargo test`) and linter (`cargo clippy`)
5. **And** it builds Go binaries (`go build`) and runs `go vet`
6. **And** it builds the Frontend (`npm run build`)
7. **And** the pipeline fails if any step fails

## Tasks / Subtasks

- [x] Create GitHub Actions Workflow
  - [x] Create `.github/workflows/ci.yml`
  - [x] Define triggers: `push` to `main`, `pull_request` to `main`
- [x] Define Rust Job
  - [x] Use `actions-rs/toolchain` or `rust-toolchain` action
  - [x] Cache `~/.cargo/registry` and `target/`
  - [x] Run `cargo check`, `cargo clippy`, `cargo test`
  - [x] Run `cargo build --release`
- [x] Define Go Job
  - [x] Use `actions/setup-go`
  - [x] Cache `~/go/pkg/mod`
  - [x] Run `go vet ./...`
  - [x] Run `go test ./...`
  - [x] Run `go build ./...`
- [x] Define Frontend Job
  - [x] Use `actions/setup-node`
  - [x] Cache `node_modules`
  - [x] Run `npm ci`
  - [x] Run `npm run lint` (if configured)
  - [x] Run `npm run build`
- [x] Verify Pipeline
  - [x] (Optional) Use `act` to test locally if available, or rely on pushing to repo.

## Dev Notes

### Technical Stack Versions
- **GitHub Actions:** Latest stable actions
- **Rust:** Stable (1.91.1+)
- **Go:** 1.22+ (Match local version)
- **Node:** 20+ (LTS)

### Pipeline Optimization
- **Parallelism:** Run Rust, Go, and Frontend jobs in parallel.
- **Caching:** Critical for Rust builds. Use `Swatinem/rust-cache` for easy caching.
- **Fail Fast:** Ensure jobs fail immediately on error.

### Architecture Compliance
- **CI/CD:** Defined in `.github/workflows/`.
- **Monorepo:** Pipeline must handle the root directory context correctly.

### References
- [Epic 0 Details](docs/epics.md#Epic-0-Foundation--Infrastructure-Scaffolding)

## Dev Agent Record

### Context Reference
- **Story ID:** 0.4
- **Story Key:** 0-4-ci-cd-pipeline-foundation

### Agent Model Used
- Gemini 2.0 Flash

### Completion Notes List
- Ultimate context engine analysis completed - comprehensive developer guide created
- Created GitHub Actions workflow `.github/workflows/ci.yml` with parallel jobs for Rust, Go, and Frontend.
- Skipped local verification (`act` not installed). Pipeline will be verified upon push to `main`.

## File List
- .github/workflows/ci.yml

## Change Log
- 2025-12-03: Created CI/CD pipeline foundation.
