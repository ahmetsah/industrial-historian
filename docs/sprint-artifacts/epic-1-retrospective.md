# Retrospective - Epic 1: Secure Access & Identity (Auth)

**Date:** 2025-12-03
**Participants:** Alice (PO), Bob (SM), Charlie (Senior Dev), Dana (QA), Elena (Junior Dev), Ahmet (Project Lead)

## 1. Epic Summary

**Status:** Completed (3/3 Stories)
**Goal:** Implement a secure, FDA-compliant authentication system (JWT, Re-auth, RBAC).

**Delivery Metrics:**
- **Completed:** 3/3 stories (100%)
- **Quality:** High. Core FDA requirements met.
- **Velocity:** Consistent with Epic 0.

## 2. What Went Well (Successes)

- **FDA Compliance:** The "Re-authentication" flow (Story 1.2) was implemented precisely as required, providing a solid foundation for Part 11 compliance.
- **Event-Driven Audit:** Publishing `sys.auth.login` and `sys.auth.signature_issued` events to NATS proved to be a powerful pattern for decoupling the Audit Trail from the Auth Service.
- **Version Discipline:** We stuck to Go 1.22 (as decided in Epic 0 Retro), avoiding the version mismatch headaches we faced previously.
- **RBAC Flexibility:** The middleware design evolved during implementation to correctly allow `ADMIN` override, showing good adaptability.

## 3. Challenges & Lessons Learned

- **Production Readiness:** We initially missed "graceful shutdown" and used hardcoded secrets/strings.
    - **Lesson:** Even for MVP, basic operational hygiene (shutdown, config) should be part of the "Definition of Done" or standard boilerplate.
- **Role Case Sensitivity:** A minor but annoying issue with `admin` vs `ADMIN` caused friction.
    - **Lesson:** Use constants for Enums (Roles, Statuses) everywhere, from the DB layer to the API handler, to prevent string typos.
- **Testing Scope:** We manually tested with `curl`.
    - **Lesson:** We should start adding automated integration tests that run these `curl` scenarios (or use a tool like Postman/Bruno) to prevent regression.

## 4. Action Items

| Action Item | Owner | Priority | Status |
| :--- | :--- | :--- | :--- |
| **Shared Constants:** Create a shared Go package (or use `historian-core` protos) for Role definitions to ensure consistency across services. | Charlie | Medium | Todo |
| **Integration Test Suite:** Create a basic script (e.g., `tests/auth_integration.sh`) that runs the `curl` commands we used manually. | Dana | High | Todo |
| **Secrets Management:** Refactor `main.go` to load secrets (keys, passwords) strictly from environment variables or a secrets manager, removing defaults. | Elena | Low | Todo |

## 5. Next Epic Readiness (Epic 2: Ingestor)

**Status:** Ready to Start

**Dependencies Check:**
- [x] **Service Accounts:** Ready (Story 1.3). Ingestors can now authenticate.
- [x] **Protobuf:** `SensorData` defined (Epic 0).
- [x] **Rust Environment:** Toolchain ready.

**Risks:**
- **Context Switch:** Switching from Go (Auth) to Rust (Ingestor) requires a mental shift.
- **Performance:** Epic 2 is about "High Performance". We need to be careful with memory management and async runtime (Tokio) usage.

**Preparation Plan:**
- Review `tokio-modbus` documentation.
- Ensure local Modbus simulator is available (or use a mock).

---
**Facilitator Notes:**
The team is finding its rhythm. The "Event-Driven" architecture is proving its worth. The shift to Rust for Epic 2 will be the next big test.
