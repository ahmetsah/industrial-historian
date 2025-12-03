# Story 1.2: FDA Re-Authentication (Electronic Signature)

Status: done

## Story

As a Quality Manager,
I want a re-authentication mechanism for critical actions,
so that the system complies with FDA 21 CFR Part 11 "Electronic Signature" requirements.

## Acceptance Criteria

1. **Given** an authenticated user performing a critical action (e.g., changing alarm limits)
2. **When** the UI prompts for re-authentication
3. **And** the user enters their password again
4. **Then** the Auth service verifies the password
5. **And** issues a short-lived "Signing Token" (valid for 1 minute)
6. **And** this token is logged in the Audit Trail as a signature (simulated for now via NATS event)

## Tasks / Subtasks

- [x] Define Signing Token Structure
  - [x] Update `token_service.go` to support "signing" token type
  - [x] Set short expiration (1 minute)
  - [x] Include specific claim `scope: signature`
- [x] Implement Re-Auth Handler
  - [x] Create `POST /api/v1/re-auth` endpoint
  - [x] Verify user's current session (Access Token)
  - [x] Verify provided password against DB hash
  - [x] Generate and return Signing Token
- [x] Publish Signature Event
  - [x] Publish `sys.auth.signature_issued` event to NATS
  - [x] Payload: UserID, Timestamp, Reason (optional)
- [x] Testing
  - [x] Unit test: Generate signing token with correct claims
  - [x] Integration test: Full re-auth flow with valid/invalid password
  - [x] Verify token expiration (1 minute)

## Dev Notes

### FDA 21 CFR Part 11 Compliance
- **Electronic Signatures:** Must be unique to the individual and executed for each signing act.
- **Re-auth:** This is NOT a session refresh. It is a deliberate action to "sign" a record.
- **Audit:** The issuance of this token is a critical audit event.

### Technical Implementation
- **Endpoint:** `POST /api/v1/re-auth`
- **Input:** `{"password": "..."}` (User is identified by their Bearer token)
- **Output:** `{"signing_token": "..."}`

### Architecture Compliance
- **Service:** `go-services/auth`
- **NATS Subject:** `sys.auth.signature_issued`

### References
- [Epic 1 Details](docs/epics.md#Epic-1-Secure-Access--Identity-Auth)
- [FDA 21 CFR Part 11](https://www.accessdata.fda.gov/scripts/cdrh/cfdocs/cfcfr/CFRSearch.cfm?CFRPart=11)

## Dev Agent Record

### Context Reference
- **Story ID:** 1.2
- **Story Key:** 1-2-fda-re-authentication-electronic-signature

### Agent Model Used
- Gemini 2.0 Flash

### Completion Notes List
- Ultimate context engine analysis completed - comprehensive developer guide created
