# Story 4.1: Real-time Dashboard Framework

Status: ready-for-dev

## Story

**As a** Operator,
**I want** a customizable dashboard,
**So that** I can arrange charts relevant to my machine.

## Acceptance Criteria

1.  **Given** the React application
2.  **When** I drag and drop a "Chart Widget"
3.  **Then** I can configure it to listen to a specific NATS subject
4.  **And** the layout is saved to local storage
5.  **And** the layout persists across page reloads

## Tasks / Subtasks

- [ ] **Initialize Dashboard Feature**
  - [ ] Create directory structure: `viz/src/features/dashboard/{components,stores,types}`
  - [ ] Define `DashboardWidget` type (id, type, position, config)

- [ ] **State Management (Zustand)**
  - [ ] Install `zustand` (v5+)
  - [ ] Create `useDashboardStore` in `viz/src/features/dashboard/stores/useDashboardStore.ts`
  - [ ] Implement `addWidget`, `removeWidget`, `updateLayout` actions
  - [ ] Use `persist` middleware to save state to `localStorage`
  - [ ] **Constraint:** Separate actions from state object (best practice)

- [ ] **Grid Layout Implementation**
  - [ ] Install `react-grid-layout` and `@types/react-grid-layout`
  - [ ] Create `DashboardLayout` component in `viz/src/features/dashboard/components/DashboardLayout.tsx`
  - [ ] Implement `ResponsiveReactGridLayout` with `WidthProvider`
  - [ ] Configure breakpoints (`lg`, `md`, `sm`, `xs`, `xxs`) and columns
  - [ ] Handle `onLayoutChange` to update Zustand store

- [ ] **Widget Component**
  - [ ] Create `DashboardWidget` wrapper component
  - [ ] Implement "Edit Mode" toggle (only allow dragging/resizing in edit mode)
  - [ ] Add "Remove" button to widget header
  - [ ] Create a placeholder content for widgets (Real charts come in Story 4.2)

- [ ] **Integration**
  - [ ] Add `DashboardPage` to main routing
  - [ ] Verify persistence by reloading page

## Dev Notes

### Technical Requirements
*   **Library:** `react-grid-layout` v1.5.2+
*   **State:** `zustand` v5.0.9+
*   **Icons:** Use `lucide-react` (standard for this project)
*   **Styling:** TailwindCSS for widget containers

### Architecture Compliance
*   **Store Pattern:** Use "Slice" pattern if store grows, but for now a single `useDashboardStore` is fine.
*   **Persistence:** Use `create(persist(...))` from `zustand/middleware`.
*   **Performance:** Use `useShallow` from `zustand/react/shallow` when selecting multiple state values to prevent unnecessary re-renders.

### React 18 & Grid Layout
*   `react-grid-layout` has known quirks with React 18 Strict Mode.
*   **Workaround:** If you see "cursor desync" or layout jumping, ensure you are using `WidthProvider` and potentially memoize children.
*   **CSS:** You must import the grid layout CSS:
    ```typescript
    import 'react-grid-layout/css/styles.css';
    import 'react-resizable/css/styles.css';
    ```

### Project Structure
```
viz/src/features/dashboard/
├── components/
│   ├── DashboardLayout.tsx
│   ├── DashboardWidget.tsx
│   └── WidgetPlaceholder.tsx
├── stores/
│   └── useDashboardStore.ts
└── types/
    └── index.ts
```

### References
*   [Epics: Epic 4](./../epics.md#epic-4-real-time-visualization-viz)
*   [Architecture: Frontend Architecture](./../architecture.md#frontend-architecture)

## Dev Agent Record

### Context Reference
*   **Architecture:** `docs/architecture.md`
*   **Epics:** `docs/epics.md`

### Agent Model Used
Antigravity (Google Deepmind)

### Completion Notes List
*   [ ] Confirmed `react-grid-layout` installation
*   [ ] Verified persistence works
*   [ ] Checked responsiveness on mobile breakpoint

### File List
*   `viz/package.json`
*   `viz/src/features/dashboard/stores/useDashboardStore.ts`
*   `viz/src/features/dashboard/components/DashboardLayout.tsx`
*   `viz/src/features/dashboard/components/DashboardWidget.tsx`
