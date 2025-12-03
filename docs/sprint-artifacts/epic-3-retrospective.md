# Retrospective - Epic 3: Efficient Storage Engine (Engine)

**Date:** 2025-12-04
**Participants:** Alice (PO), Bob (SM), Charlie (Senior Dev), Dana (QA), Elena (Junior Dev), Ahmet (Project Lead)

## 1. Epic Summary

**Status:** Completed (3/3 Stories)
**Goal:** Store massive amounts of time-series data efficiently and retrieve it instantly.

**Delivery Metrics:**
- **Completed:** 3/3 stories (100%)
- **Quality:** High. Core requirements (LSM-Tree, Tiering, gRPC Query) met.
- **Velocity:** High. Complex Rust features (RocksDB integration, S3 tiering, gRPC) delivered rapidly.

## 2. What Went Well (Successes)

- **LSM-Tree Integration:** Successfully integrated `rust-rocksdb` for high-throughput writes (>50k events/sec validated in unit tests).
- **Tiered Storage:** Implemented a sophisticated "Export-Upload-Delete" strategy for tiering data to S3/MinIO, including a metadata index for transparent hybrid reads.
- **gRPC Performance:** The gRPC API with LTTB downsampling provides a fast query interface. Optimization to O(N) downsampling significantly improved efficiency.
- **Data Safety:** Identified and fixed a critical data loss bug in the tiering job before release.

## 3. Challenges & Lessons Learned

- **RocksDB Complexity:** Tiering individual SSTables proved difficult with RocksDB's internal management.
    - **Lesson:** For managed engines like RocksDB, application-level tiering (export/delete) is often more practical than trying to manipulate internal files.
- **Metadata Management:** We initially missed recording metadata after S3 upload, which would have led to data loss.
    - **Lesson:** In distributed data movement, always ensure the "pointer" to the new location is persisted *before* deleting the source.
- **Query Semantics:** Handling "Not Found" vs "Empty Range" for non-existent sensors is tricky without a central registry.
    - **Lesson:** A Sensor Registry service or table is needed for strict validation, but for high-performance ingestion, "implicit creation" is acceptable if documented.

## 4. Action Items

| Action Item | Owner | Priority | Status |
| :--- | :--- | :--- | :--- |
| **Benchmark Suite:** Create a proper `criterion` benchmark suite to validate >50k events/sec under sustained load (not just unit tests). | Charlie | High | Todo |
| **WAL Recovery Test:** Add an explicit integration test to verify WAL recovery after a hard crash. | Dana | High | Todo |
| **Allocation Optimization:** Optimize `generate_key` in `rocksdb.rs` to avoid memory allocation in the hot write path. | Elena | Medium | Todo |
| **Sensor Registry:** Consider adding a lightweight Sensor Registry to support strict "Not Found" checks. | Ahmet | Low | Backlog |

## 5. Next Epic Readiness (Epic 4: Real-Time Visualization)

**Status:** Ready to Start

**Dependencies Check:**
- [x] **Query API:** Ready (Epic 3). gRPC API is available for the backend-for-frontend (BFF) or direct consumption.
- [x] **Data Source:** Ingestion (Epic 2) and Storage (Epic 3) are flowing.

**Risks:**
- **Frontend Performance:** Rendering 100k+ points in the browser is challenging. `uPlot` is chosen but needs careful integration.
- **WebSocket Load:** Streaming real-time data to many clients might require a dedicated Gateway service.

**Preparation Plan:**
- Initialize React/Vite project (already done in Epic 0).
- Prototype `uPlot` integration.
- Design the WebSocket topic structure for the frontend.

---
**Facilitator Notes:**
This was a technically dense epic. The team handled the complexity of LSM-trees and S3 integration well. The "Adversarial Code Review" caught critical bugs, proving its value.
