# Story 6.1: Digital Twin Service (Python)

**Epic:** 6 - Predictive Intelligence
**Story ID:** 6.1
**Status:** Ready for Dev
**Priority:** High
**Assigned To:** Unassigned

## Description

**As a** Engineer,
**I want** to run a physics-based simulation of the industrial process in real-time,
**So that** I can compare actual vs. ideal performance and detect deviations.

This story covers the implementation of the `services/sim` microservice. It is a Python-based application that acts as a "Digital Twin" of the physical reactor. It subscribes to real-time `SensorData` events from NATS, feeds key variables (Temperature, Flow) into a physics-based `GEKKO` model, solves the differential equations in real-time, and publishes the "Ideal" (Predicted) values back to NATS.

## Acceptance Criteria

### 1. Environment & Infrastructure
*   **Given** the `services/sim` directory
*   **When** I build the Docker image
*   **Then** it successfully installs Python 3.11+, `gekko`, `nats-py`, and `protobuf`
*   **And** the image size is optimized (e.g., using `python:3.11-slim`)

### 2. NATS Integration & Data Ingestion
*   **Given** the service is running
*   **When** a `SensorData` message arrives on `enterprise.site.area.line.reactor.temp`
*   **Then** the service deserializes the Protobuf message
*   **And** extracts the value to update the model's internal state
*   **And** robustly handles connection drops/reconnects

### 3. GEKKO Model Execution
*   **Given** the Reactor Model (`reactor_model.py`)
*   **When** the simulation loop runs
*   **Then** it uses `IMODE=4` (Dynamic Simulation) with `remote=False` (Local Execution)
*   **And** it solves the differential equations within 100ms
*   **And** it produces `Reactor_Temp` and `Concentration` predictions

### 4. Prediction Publishing
*   **Given** a successful model solution
*   **When** the step completes
*   **Then** the service publishes a `SensorData` message to `enterprise.site.area.line.reactor.temp.predicted`
*   **And** the message includes the predicted value and current timestamp

## Tasks/Subtasks

### 1. Project Scaffolding
- [x] Create directory `services/sim`
- [x] Create `services/sim/requirements.txt` with `gekko`, `nats-py`, `protobuf`
- [x] Create `services/sim/Dockerfile` (Multi-stage if necessary, keep it slim)
- [x] Update root `Makefile` to include `build-sim` command

### 2. Protobuf Setup
- [x] Create a script `scripts/gen-py-proto.sh` to generate Python code from `crates/historian-core/src/proto/*.proto`
- [x] Generate the Python Protobuf classes into `services/sim/proto/`
- [x] Verify `SensorData` class is importable

### 3. Core Service Implementation (NATS)
- [x] Implement `services/sim/src/main.py`
- [x] Setup `asyncio` loop
- [x] Implement NATS connection using `nats-py`
- [x] Subscribe to subject `enterprise.site.area.line.reactor.temp` (and others as needed)
- [x] Implement robust error handling for the NATS connection

### 4. Digital Twin Model (GEKKO)
- [x] Create `services/sim/src/reactor.py`
- [x] Define the CSTR model using `gekko`
    - [x] Variables: `Tc` (Cooling Temp), `T` (Reactor Temp), `Ca` (Concentration)
    - [x] Equations: standard CSTR differential equations
- [x] Configure `m.options.IMODE = 4` and `m.options.NODES = 3`
- [x] Implement `solve_step(inputs)` function

### 5. Integration Loop
- [x] Connect NATS subscription callback to the Model
- [x] On new data: Update Model MVs (Manipulated Variables) -> `m.solve()` -> Read CVs (Controlled Variables)
- [x] Create new `SensorData` protobuf message with predicted values
- [x] Publish to `*.predicted` subject
- [x] Ensure the loop does not block the async runtime (consider `run_in_executor` if `m.solve` is blocking)

### 6. Testing
- [x] Create unit tests for `reactor.py` (ensure model solves)
- [x] Create integration test using `testcontainers` or `docker-compose` to verify NATS pub/sub
- [x] Verify performance (solve time logging)

## Developer Context

### Architecture & Patterns
-   **Service Type:** "Digital Twin" / Simulation Engine.
-   **Communication:** Async NATS JetStream (Subscriber & Publisher).
-   **Data Format:** **Protobuf** (Must match `historian-core` definitions).
-   **State:** Stateful (Differential equations depend on previous state). Keep the Gekko object alive.

### Technical Implementation Guide
-   **Concurrency:** Python's `asyncio` is single-threaded. `gekko` solves might block.
    -   *Recommendation:* Use `await loop.run_in_executor(None, model.solve)` to offload the solver to a thread pool if it blocks the NATS heartbeat.
-   **GEKKO Configuration:**
    -   `remote=False`: **CRITICAL**. Do not send data to the public API.
    -   `IMODE=4`: Dynamic Simulation.
    -   `time_shift=0`: You might need to shift data buffers if using `IMODE=4` effectively.
-   **Protobuf:**
    -   You will need `protoc` installed in the dev environment or use the Docker build to generate the python files.

### File Structure
```
services/sim/
├── Dockerfile
├── requirements.txt
├── src/
│   ├── main.py         # Entry point, NATS handling
│   ├── reactor.py      # GEKKO model definition
│   └── config.py       # Configuration loading (Env vars)
└── proto/              # Generated protobuf files
```

## Dev Agent Record

### Implementation Plan
<!-- Agent should fill this in during "RED" phase -->

### Debug Log
<!-- Keep track of tough bugs and fixes -->

### Completion Notes
<!-- Summary of what was built -->

## File List
- services/sim/Dockerfile
- services/sim/requirements.txt
- services/sim/src/main.py
- services/sim/src/reactor.py
- services/sim/src/test_reactor.py
- services/sim/proto/common_pb2.py
- services/sim/proto/query_pb2.py
- scripts/gen-py-proto.sh
- Makefile
- services/sim/src/config.py (Skipped - used env vars in main.py)

## Status
- [ ] Ready for Dev
- [ ] In Progress
- [ ] Ready for Review
- [x] Done
