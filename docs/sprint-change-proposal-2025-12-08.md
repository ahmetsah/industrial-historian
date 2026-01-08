# Sprint Change Proposal: Resetting Visualization Module

**System Date:** 2025-12-08
**Author:** Product Manager (John)

## 1. Issue Summary
**Trigger:** Strategic decision to rebuild the visualization module (`viz/`) from scratch.
**Problem Statement:** The current implementation of the Visualization module (Epic 4) requires a complete reset. The user has explicitly requested to delete all contents of the `viz/` directory and restart implementation from Story 4.1.
**Context:** This change intends to clear technical debt or misalignment in the current frontend implementation before proceeding further.

## 2. Impact Analysis
### Epic Impact
- **Epic 4 (Visualization & Dashboards):** Status must be reverted from `done` to `in-progress` (or `contexted/drafted`).
- **Story 4.1, 4.2, 4.3:** Status must be reverted from `done` to `drafted`. Implementation needs to be re-executed.

### Artifact Conflicts
- **Codebase:** `viz/` directory will be purged.
- **Sprint Status:** `docs/sprint-artifacts/sprint-status.yaml` reflects Epic 4 as `done`. This creates a false state vs. reality.

### Technical Impact
- **Dependencies:** No direct backend dependencies identified that would break core ingestion/storage (Epics 0-3, 5, 6). The dashboard is a consumer.
- **Risk:** Low risk to backend integrity. Main risk is timeline delay for frontend availability.

## 3. Recommended Approach
**Selected Path:** **Option 2: Explicit Rollback & Restart**
We will clean the slate for the frontend to allow a fresh, high-quality implementation without legacy constraints.

**Rationale:**
- User explicitly requested a "fresh start".
- Attempting to patch a fundamentally rejected implementation is costlier than rewriting.
- Requirements (Story 4.1, 4.2) appear to remain valid (unless user specifies otherwise during re-implementation).

## 4. Detailed Change Proposals

### A. Artifact Updates
**Target:** `docs/sprint-artifacts/sprint-status.yaml`

**Changes:**
```yaml
development_status:
  epic-4: contexted  # Was: done
  4-1-real-time-dashboard-framework: drafted  # Was: done
  4-2-trend-chart-component-uplot: drafted  # Was: done
  4-3-data-export-service: drafted  # Was: done
  epic-4-retrospective: optional  # Was: completed
```

## 5. Implementation Handoff
**Scope:** **Major** (Resetting an entire Epic).
**Handoff Plan:**
1. **PM (Me):** Update `sprint-status.yaml` to reflect the reset.
2. **System/Dev:** Execute command to delete `viz/`.
3. **Dev Agent:** Pick up Story 4.1 as the next task.

## 6. Approval & Execution
**Status:** **Awaiting Final Confirmation**
The next step is to execute the file updates and deletion command upon user approval.
