# Story 0.2: Infrastructure as Code (Docker Compose & NATS)

Status: done

## Story

As a DevOps Engineer,
I want a Docker Compose environment with NATS JetStream,
so that I can run the entire distributed system locally for development.

## Acceptance Criteria

1. **Given** the `ops/` directory
2. **When** I run `docker-compose up -d`
3. **Then** NATS JetStream is running on port 4222
4. **And** NATS Management UI (if available) or CLI tool can connect
5. **And** A Stream named `EVENTS` is created with subject `enterprise.>`
6. **And** MinIO (S3 compatible) is running for future object storage needs
7. **And** PostgreSQL is running for the Auth service

## Tasks / Subtasks

- [x] Create NATS Configuration
  - [x] Create `ops/nats.conf`
  - [x] Enable JetStream (`jetstream { store_dir: "/data/jetstream" }`)
  - [x] Configure file-based storage
- [x] Create Docker Compose File
  - [x] Define `nats` service (image: `nats:latest`, ports: 4222, 8222, command: `-c /etc/nats/nats.conf`)
  - [x] Define `minio` service (image: `minio/minio`, ports: 9000, 9001)
  - [x] Define `postgres` service (image: `postgres:16-alpine`, ports: 5432)
  - [x] Define volumes for persistence (`nats_data`, `minio_data`, `pg_data`)
- [x] Automate Stream Creation
  - [x] Create a setup script or use `nats-box` sidecar in compose to initialize the `EVENTS` stream
  - [x] Command: `nats stream add EVENTS --subjects "enterprise.>" --storage file --retention limits`
- [x] Verify Environment
  - [x] Test `docker-compose up`
  - [x] Verify NATS connection
  - [x] Verify MinIO UI access
  - [x] Verify Postgres connection

## Dev Notes

### Technical Stack Versions
- **NATS:** Latest (JetStream enabled)
- **MinIO:** Latest
- **PostgreSQL:** 16-alpine
- **Docker Compose:** V2

### Configuration Notes
- **NATS:**
  - Enable JetStream.
  - Store directory: `/data/jetstream` (mapped to volume).
  - Subject hierarchy: `enterprise.>` (as defined in Architecture).
- **MinIO:**
  - Default credentials (dev): `minioadmin` / `minioadmin`.
  - Bucket creation: Optional for now, but good to have a setup script.
- **Postgres:**
  - Default credentials (dev): `postgres` / `postgres`.
  - DB Name: `historian`.

### Architecture Compliance
- **Infrastructure:** All defined in `ops/` directory.
- **Persistence:** Use Docker volumes for data persistence across restarts.
- **Network:** All services on default bridge network, accessible via localhost ports.

### References
- [Architecture Infrastructure](docs/architecture.md#Infrastructure--Deployment)
- [Epic 0 Details](docs/epics.md#Epic-0-Foundation--Infrastructure-Scaffolding)

## Dev Agent Record

### Context Reference
- **Story ID:** 0.2
- **Story Key:** 0-2-infrastructure-as-code-docker-compose-nats

### Agent Model Used
- Gemini 2.0 Flash

### Completion Notes List
- Ultimate context engine analysis completed - comprehensive developer guide created
- Created NATS configuration file `ops/nats.conf` with JetStream enabled.
- Created Docker Compose file `ops/docker-compose.yml` with NATS, MinIO, and Postgres.
- Created Stream Setup script `ops/setup_streams.sh` and integrated it into Docker Compose.
- Verified environment: NATS, MinIO, and Postgres are running and accessible. Stream EVENTS created successfully.

## File List
- ops/nats.conf
- ops/docker-compose.yml
- ops/setup_streams.sh
- tests/infra/verify_nats_conf.sh
- tests/infra/verify_docker_compose.sh
- tests/infra/verify_stream_setup.sh

## Change Log
- 2025-12-03: Verified environment and stream creation.
- 2025-12-03: Created Stream Setup script and integrated with Docker Compose.
- 2025-12-03: Created Docker Compose configuration and verification script.
- 2025-12-03: Created Docker Compose configuration and verification script.
- 2025-12-03: Created NATS configuration and verification script.
