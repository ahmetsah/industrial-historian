# Retrospective - Epic 2: High-Performance Ingestion (Ingestor)

**Date:** 2025-12-03
**Participants:** Alice (PO), Bob (SM), Charlie (Senior Dev), Dana (QA), Elena (Junior Dev), Ahmet (Project Lead)

## 1. Epic Summary

**Status:** Completed (3/3 Stories)
**Goal:** Enable high-speed data collection from industrial devices with zero data loss guarantees.

**Delivery Metrics:**
- **Completed:** 3/3 stories (100%)
- **Quality:** High. Core requirements (Modbus, Buffering, Calculation) met.
- **Velocity:** Consistent. The team adapted well to the Rust transition.

## 2. What Went Well (Successes)

- **Rust Transition:** The team successfully implemented the Ingestor service in Rust, leveraging `tokio` for high concurrency.
- **Pipeline Architecture:** The modular design (`ModbusAdapter` -> `CalculationEngine` -> `HybridBuffer` -> `Publisher`) proved to be clean and extensible.
- **Data Safety:** The "Hybrid Store & Forward" mechanism (Story 2.2) works as intended, ensuring data is spilled to disk when the memory buffer fills up, guaranteeing zero data loss.
- **Calculation Power:** The `evalexpr` integration (Story 2.3) provides a flexible way to define derived tags without complex hardcoding.

## 3. Challenges & Lessons Learned

- **Memory Management:** We identified a potential risk with `DiskBuffer::read_all_and_clear` loading the entire WAL file into memory.
    - **Lesson:** For "unbounded" resources like disk logs, always prefer streaming or chunked processing from day one.
- **Performance Scalability:** The Calculation Engine currently iterates over *all* formulas for *every* incoming tag update (O(N*M)).
    - **Lesson:** While fine for MVP, this will be a bottleneck at scale. We need to plan for "premature optimization" when the architectural pattern (O(N^2)) is obviously unscalable.
- **Modbus Nuances:** We missed setting the `Unit ID` (Slave ID) in the initial implementation, which is critical for some gateways.
    - **Lesson:** Protocol details matter. Always verify against real hardware specs or comprehensive mocks.

## 4. Action Items

| Action Item | Owner | Priority | Status |
| :--- | :--- | :--- | :--- |
| **Disk Buffer Optimization:** Refactor `DiskBuffer` to use chunked reading or streaming to prevent OOM during recovery from long outages. | Charlie | High | Todo |
| **Calculation Optimization:** Implement a Dependency Graph (DAG) or Reverse Index for the Calculation Engine to only evaluate affected formulas. | Elena | Medium | Todo |
| **Load Testing:** Conduct a load test (e.g., 10k events/sec) to validate the system's stability and identify actual bottlenecks. | Dana | High | Todo |

## 5. Next Epic Readiness (Epic 3: Efficient Storage Engine)

**Status:** Ready to Start

**Dependencies Check:**
- [x] **Ingestion Pipeline:** Ready (Epic 2). We have a stream of data to store.
- [x] **Protobuf:** `SensorData` defined.

**Risks:**
- **Complexity:** Implementing an LSM-Tree (Log-Structured Merge-Tree) from scratch (or integrating a complex crate) is non-trivial.
- **Data Volume:** We need to handle the high volume of data coming from the Ingestor.

**Preparation Plan:**
- Research Rust LSM-tree crates (e.g., `rocksdb`, `sled`, or custom implementation).
- Review storage requirements (retention, query patterns).

---
**Facilitator Notes:**
The shift to Rust was a success. The team is building confidence in low-level systems programming. The "Hybrid Buffer" is a critical asset for system reliability.
