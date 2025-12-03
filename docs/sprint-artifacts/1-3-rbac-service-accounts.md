# Story 1.3: RBAC & Service Accounts

Status: done

## Story

As a Security Officer,
I want Role-Based Access Control,
so that operators cannot change engineering settings and services can authenticate securely.

## Acceptance Criteria

1. **Given** a user with role `OPERATOR`
2. **When** they try to access an `ADMIN` only endpoint
3. **Then** the request is denied (403 Forbidden)
4. **And** I can create "Service Accounts" (API Keys or long-lived JWTs) for Ingestors
5. **And** Service Accounts have the `SERVICE` role
6. **And** Middleware enforces role checks on protected routes

## Tasks / Subtasks

- [x] Define Roles
  - [x] Update `user_repo.go` to enforce valid roles: `ADMIN`, `ENGINEER`, `OPERATOR`, `AUDITOR`, `SERVICE`
- [x] Implement RBAC Middleware
  - [x] Create `internal/middleware/rbac.go`
  - [x] Implement `RequireRole(role string)` middleware function
  - [x] Check `role` claim from JWT context
- [x] Implement Service Accounts
  - [x] Create `POST /api/v1/service-accounts` endpoint (Admin only)
  - [x] Generate long-lived JWTs (e.g., 1 year) for services
  - [x] Store service account metadata in DB (name, created_by)
- [x] Protect Endpoints
  - [x] Apply `RequireRole("ADMIN")` to user management endpoints
  - [x] Apply `RequireRole("SERVICE")` (or similar) to ingestion endpoints (if any in Auth service, otherwise just verify token generation)
- [x] Testing
  - [x] Unit test: Middleware allows/denies based on role
  - [x] Integration test: Operator cannot create users
  - [x] Integration test: Admin can create service account

## Dev Notes

### Roles Definition
- **ADMIN:** Full access (User management, System config).
- **ENGINEER:** Configuration access (Tags, Alarms), but not User management.
- **OPERATOR:** Read-only operational data, Acknowledge alarms.
- **AUDITOR:** Read-only Audit logs.
- **SERVICE:** Machine-to-machine access (Ingestors).

### Service Accounts
- For MVP, Service Accounts can be just Users with `role=SERVICE` and a long-lived JWT.
- Future: specific API Key table if needed, but JWT is fine for now.

### Architecture Compliance
- **Middleware:** Go middleware pattern.
- **Security:** Fail closed (default deny).

### References
- [Epic 1 Details](docs/epics.md#Epic-1-Secure-Access--Identity-Auth)

## Dev Agent Record

### Context Reference
- **Story ID:** 1.3
- **Story Key:** 1-3-rbac-service-accounts

### Agent Model Used
- Gemini 2.0 Flash

### Completion Notes List
- Ultimate context engine analysis completed - comprehensive developer guide created
