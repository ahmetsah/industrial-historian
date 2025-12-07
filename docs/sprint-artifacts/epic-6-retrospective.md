# Retrospective - Epic 6: Predictive Intelligence (Sim)

**Date:** 2025-12-07
**Participants:** Alice (PO), Bob (SM), Charlie (Senior Dev), Dana (QA), Elena (Junior Dev), Ahmet (Project Lead)

## 1. Epic Summary

**Status:** Completed (2/2 Stories)
**Goal:** Transform the system from reactive to proactive with Digital Twin capabilities (GEKKO simulation & Anomaly Detection).

**Delivery Metrics:**
- **Completed:** 2/2 stories (100%)
- **Quality:** Meets MVP requirements. Functional end-to-end pipeline established.
- **Velocity:** Fast. Python ecosystem allowed for rapid prototyping of complex math models.

## 2. What Went Well (Successes)

- **Rapid Prototyping with Python:** Using Python for the simulation service (`services/sim`) was the right choice. Libraries like `GEKKO` and `numpy` made implementing differential equations and statistical checks trivial compared to Rust or Go.
- **NATS Integration:** The `nats-py` client integrated seamlessly with our existing NATS JetStream infrastructure. The pattern of "Subscribe to Real -> Publish Predicted" proved very effective.
- **Protobuf Workflow:** Generating Python code from the shared `historian-core` protobuf definitions worked smoothly, maintaining type safety across the polyglot stack (Rust/Go/Python).
- **Concurrency Management:** We proactively handled the blocking nature of the GEKKO solver by using `loop.run_in_executor`, preventing the NATS heartbeat from timing out during heavy computations.

## 3. Challenges & Lessons Learned

- **"Apples to Oranges" Comparison:** In Story 6.2, we simplified the anomaly detection to compare `Input Tc` vs `Predicted T`. While this verified the software pipeline, it is chemically invalid.
    - **Lesson:** For MVP demos, functional verification (pipes working) often takes precedence over domain correctness, but technical debt must be explicitly documented to prevent confusion later.
- **Docker Path Issues:** We encountered minor issues with script paths (`trigger_anomaly.py`) when running inside Docker.
    - **Lesson:** Always verify file placements in `Dockerfile` COPY instructions, especially for ad-hoc scripts used for testing.
- **Environment Differences:** Initial local testing failed because `numpy` wasn't installed on the host, only in Docker.
    - **Lesson:** "It works on my machine" is solved by "It works in the container". Rely more on `docker-compose` for local dev/test cycles.

## 4. Code Review Findings

**Story 6.2 (Anomaly Detection):**
- **Accepted Debt:** The logic comparing `data.value` (Input) with `pred_T` (Output) was flagged but accepted for the sake of the demo.
- **Simplification:** Using a simple Z-Score algorithm was deemed sufficient for MVP, avoiding the complexity of `IsolationForest` or other ML models for now.

## 5. Action Items

| Action Item | Owner | Priority | Status |
| :--- | :--- | :--- | :--- |
| **Model Refinement:** Update the Simulation to subscribe to *both* `Tc` and `T` streams for valid residual calculation. | Charlie | High | Backlog |
| **Visualization:** Add `AnomalyEvent` visualization to the Frontend (Epic 4 extension). | Elena | Medium | Backlog |
| **Unit Testing:** Improve unit test coverage for `detector.py` to handle edge cases (empty windows, zero std dev). | Dana | Low | Todo |

## 6. Project Completion & Next Steps

**Epic 6 concludes the core MVP functionality defined in the initial Roadmap.**

**System Capabilities Achieved:**
1.  **Foundation:** Polyglot Monorepo, NATS, CI/CD.
2.  **Auth:** Secure FDA-compliant access.
3.  **Ingestion:** High-speed buffering and protocol support.
4.  **Storage:** Efficient LSM-tree storage.
5.  **Viz:** Real-time dashboards.
6.  **Compliance:** Audit logs and Alarm management.
7.  **Intelligence:** Digital Twin and Anomaly Detection.

**Recommendation:**
- Proceed to **Project Retrospective** or **Release Candidate** preparation.
- Address high-priority Technical Debt accumulated across Epics 1-6.
- Plan "Phase 2" features based on stakeholder feedback from the MVP demo.

## 7. Key Insights

1.  **Polyglot Architecture Successful:** We successfully integrated Rust (Performance), Go (Business Logic), and Python (Data Science) using NATS and Protobuf. This proves the architecture is flexible and scalable.
2.  **Protocol Buffers are the Glue:** Shared schema definitions were critical. Without them, keeping data structures aligned across 3 languages would have been a nightmare.
3.  **Simulation adds Value:** Even a simple Digital Twin provides immediate visualization of "what should be happening," which is a powerful tool for operators.

## 8. Team Feedback

**Alice (PO):** "The platform looks complete. We have data flowing from end-to-end, simulated prediction, and compliance. Ready for demo."
**Charlie (Senior Dev):** "The Python service was fun to build. Good break from the strictness of Rust."
**Dana (QA):** "Testing the anomaly detector was interesting. I'd like to see more rigorous data scenarios in the future."
**Ahmet (Project Lead):** "Excellent work. We went from zero to a full industrial IoT platform MVP."
