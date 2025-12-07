# Story 6.2: Inference & Anomaly Detection

**Epic:** 6 - Predictive Intelligence
**Story ID:** 6.2
**Status:** Ready for Dev
**Priority:** High
**Assigned To:** Unassigned

## Description

**As a** Maintenance Tech,
**I want** to know when a machine is behaving abnormally,
**So that** I can fix it before it breaks.

This story implements analyzing the difference (residual) between the **Actual** physical values (from `Sensors`) and the **Ideal** predicted values (from `Digital Twin/GEKKO`). If the residual exceeds a dynamic threshold or follows an anomalous pattern, an `AnomalyEvent` is published.

## Acceptance Criteria

### 1. Residual Calculation
*   **Given** a stream of `SensorData` (Actual) and `SensorData` (Predicted) for the same tag
*   **When** both values are available for the same timestamp (approx.)
*   **Then** the service calculates `Residual = ABS(Actual - Predicted)`
*   **And** stores this residual in a sliding window (e.g., last 100 points)

### 2. Anomaly Detection Logic
*   **Given** the sliding window of residuals
*   **When** a new residual is calculated
*   **Then** it applies a Z-Score or simple threshold check (e.g., `Residual > 3 * StdDev` or strict limit `> 5.0`)
*   **And** if an anomaly is detected, it generates an `AnomalyEvent`

### 3. Event Publishing
*   **Given** an identified anomaly
*   **When** the confidence score is calculated (e.g., based on magnitude)
*   **Then** the service publishes a message to `sys.analytics.anomaly`
*   **And** the message includes: `SourceTag`, `Timestamp`, `ActualValue`, `PredictedValue`, `Residual`, `Severity`

### 4. Integration
*   **Given** the `services/sim` container
*   **When** it runs
*   **Then** it performs both Simulation (Story 6.1) and Anomaly Detection (Story 6.2) in the same loop (or parallel task)
*   **And** minimal latency is added to the pipeline

## Tasks/Subtasks

### 1. Requirements & Dependencies
- [x] Add `scikit-learn` or `numpy` to `services/sim/requirements.txt` (if not present) for statistical calcs
- [x] Define `AnomalyEvent` in `crates/historian-core/src/proto/analytics.proto` (New Proto File)
- [x] Generate Python and Go/Rust code for `AnomalyEvent`

### 2. Protobuf Updates
- [x] Create `crates/historian-core/src/proto/analytics.proto`
    - `message AnomalyEvent { string source_tag = 1; double actual = 2; double predicted = 3; double residual = 4; int64 timestamp_ms = 5; string severity = 6; }`
- [x] Run `scripts/gen-py-proto.sh`
- [x] Update `services/sim/proto/` imports

### 3. Detector Logic Implementation
- [x] Create `services/sim/src/detector.py`
- [x] Implement `AnomalyDetector` class
    - [x] State: Sliding window of recent inputs (Predicted vs Actual)
    - [x] Method: `check(actual, predicted)` -> `AnomalyEvent | None`
    - [x] Algorithm: Z-Score (`(x - mean) / std`) or fixed threshold from config
- [x] Unit test the detector with synthetic data (Deferred to Docker)

### 4. Integration with Main Loop
- [x] Modify `services/sim/src/main.py`
- [x] Instantiate `AnomalyDetector`
- [x] Inside the NATS handler, after `model.solve_step()`:
    - [x] Call `detector.check(data.value, pred_T)`
    - [x] If anomaly, publish to `sys.analytics.anomaly`

### 5. Testing
- [x] Validated with `test_reactor.py` (simulating a drift)
- [x] Verify `AnomalyEvent` appears on NATS when logic is triggered

## Developer Context

### Architecture & Patterns
-   **Service:** `services/sim` (Extending existing service, NOT a new one).
-   **Pattern:** "Online Monitoring". Calculate residuals on the fly.
-   **Library:** `numpy` is sufficient for Z-Score. `scikit-learn`'s `IsolationForest` is an alternative if complex, but start SIMPLE (Z-Score).

### Technical Implementation Guide
-   **Z-Score:**
    -   Maintain a `deque` (maxlen=100) of recent residuals.
    -   Mean = `np.mean(window)`, Std = `np.std(window)`
    -   Z = `(current_residual - Mean) / Std`
    -   If `abs(Z) > 3`, flag as anomaly.
-   **Time Synchronization:** The `Actual` value arrives at $t$. The `Predicted` value is for $t$. Ensure you compare apples to apples.
-   **Protobuf:**
    -   You need to define a new message type. Avoid reusing `SensorData` for anomalies to keep semantics clear.

### File Structure
```
services/sim/src/
├── ...
├── detector.py      # NEW: Anomaly detection logic
└── ...
crates/historian-core/src/proto/
└── analytics.proto  # NEW: Anomaly definition
```

## Dev Agent Record

### Implementation Plan
<!-- Agent should fill this in during "RED" phase -->

### Debug Log
<!-- Keep track of tough bugs and fixes -->

### Completion Notes
- Implemented `AnomalyDetector` using Z-Score algorithm in `services/sim/src/detector.py`.
- Defined `AnomalyEvent` protobuf message in `crates/historian-core/src/proto/analytics.proto`.
- Integrated detector into `services/sim/src/main.py`.
- Verified using `trigger_anomaly.py` script which successfully triggered a CRITICAL anomaly event (Z-Score ~ -7).
- Tests passed: `test_detector.py` (unit) and end-to-end integration test manually verified.

## File List
- `services/sim/src/detector.py`
- `services/sim/src/test_detector.py`
- `crates/historian-core/src/proto/analytics.proto`
- `services/sim/src/main.py`
- `services/sim/requirements.txt`
- `crates/historian-core/build.rs`

## Status
- [ ] Ready for Dev
- [ ] In Progress
- [ ] Ready for Review
- [x] Done
