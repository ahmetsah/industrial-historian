# Story 1.1: Auth Service & JWT Implementation

Status: done

## Story

As a System Administrator,
I want a centralized authentication service,
so that I can manage users and secure API access.

## Acceptance Criteria

1. **Given** the `go-services/auth` service
2. **When** I POST to `/api/v1/login` with valid credentials
3. **Then** I receive a JWT (RS256 signed) with `sub`, `role`, and `exp` claims
4. **And** the login event is published to NATS subject `sys.auth.login`
5. **And** invalid credentials return 401 Unauthorized
6. **And** I can create an initial admin user via CLI or seed script

## Tasks / Subtasks

- [x] Setup Go Module & Dependencies
  - [x] Initialize `go-services/auth` (already done in 0.1)
  - [x] Add dependencies: `github.com/golang-jwt/jwt/v5`, `golang.org/x/crypto/bcrypt`, `github.com/lib/pq` (or `pgx`), `github.com/nats-io/nats.go`
- [x] Implement Database Layer
  - [x] Create `internal/repository/user_repo.go`
  - [x] Define `User` struct (ID, Username, PasswordHash, Role)
  - [x] Implement `CreateUser` and `GetUserByUsername`
  - [x] Create migration script `migrations/001_create_users_table.sql`
- [x] Implement JWT Logic
  - [x] Create `internal/service/token_service.go`
  - [x] Implement RSA key generation/loading (generate `private.pem` if missing)
  - [x] Implement `GenerateToken(user User) string`
- [x] Implement Auth Handler
  - [x] Create `internal/handler/auth_handler.go`
  - [x] Implement `Login(w http.ResponseWriter, r *http.Request)`
  - [x] Validate credentials using bcrypt
  - [x] Publish `sys.auth.login` event to NATS on success
- [x] Wire up Main
  - [x] Configure DB connection (Postgres)
  - [x] Configure NATS connection
  - [x] Start HTTP server on port 8080 (or configured port)
- [x] Testing
  - [x] Write unit tests for `token_service`
  - [x] Write integration test for Login flow

## Dev Notes

### Technical Stack
- **Go:** 1.22
- **Database:** PostgreSQL 16 (via `ops/docker-compose.yml`)
- **Messaging:** NATS JetStream
- **JWT Library:** `github.com/golang-jwt/jwt/v5` (Use v5, not legacy)

### Security Requirements
- **Algorithm:** RS256 (Asymmetric)
- **Password Hashing:** Bcrypt (Cost 12+)
- **Token Expiry:** 1 hour (Access Token), 24 hours (Refresh Token - optional for MVP, stick to Access for now)

### Architecture Compliance
- **Service Location:** `go-services/auth`
- **NATS Subject:** `sys.auth.login` (Strict naming)
- **Database:** Owns `users` table in `historian` DB.

### References
- [Architecture Auth Decision](docs/architecture.md#Authentication--Security)
- [Epic 1 Details](docs/epics.md#Epic-1-Secure-Access--Identity-Auth)

## Dev Agent Record

### Context Reference
- **Story ID:** 1.1
- **Story Key:** 1-1-auth-service-jwt-implementation

### Agent Model Used
- Gemini 2.0 Flash

### Completion Notes List
- Ultimate context engine analysis completed - comprehensive developer guide created
