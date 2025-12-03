# Story 4.1: Real-time Dashboard Framework

Status: done

## Story

**As a** Operator,
**I want** a customizable dashboard,
**So that** I can arrange charts relevant to my machine.

## Acceptance Criteria

1.  ✅ **Given** the React application
2.  ✅ **When** I drag and drop a "Chart Widget"
3.  ✅ **Then** I can configure it to listen to a specific NATS subject
4.  ✅ **And** the layout is saved to local storage
5.  ✅ **And** the layout persists across page reloads

## Tasks / Subtasks

- [x] **Initialize Dashboard Feature**
  - [x] Create directory structure: `viz/src/features/dashboard/{components,stores,types}`
  - [x] Define `DashboardWidget` type (id, type, position, config)

- [x] **State Management (Zustand)**
  - [x] Install `zustand` (v5+)
  - [x] Create `useDashboardStore` in `viz/src/features/dashboard/stores/useDashboardStore.ts`
  - [x] Implement `addWidget`, `removeWidget`, `updateWidget`, `updateLayout` actions
  - [x] Use `persist` middleware to save state to `localStorage`
  - [x] **Constraint:** Separate actions from state object (best practice)

- [x] **Grid Layout Implementation**
  - [x] Install `react-grid-layout` and `@types/react-grid-layout`
  - [x] Create `DashboardLayout` component in `viz/src/features/dashboard/components/DashboardLayout.tsx`
  - [x] Implement `ResponsiveReactGridLayout` with `WidthProvider`
  - [x] Configure breakpoints (`lg`, `md`, `sm`, `xs`, `xxs`) and columns
  - [x] Handle `onLayoutChange` to update Zustand store

- [x] **Widget Component**
  - [x] Create `DashboardWidget` wrapper component
  - [x] Implement "Edit Mode" toggle (only allow dragging/resizing in edit mode)
  - [x] Add "Remove" button to widget header
  - [x] Add "Settings" button to configure widget (NATS subject)
  - [x] Create a placeholder content for widgets (Real charts come in Story 4.2)

- [x] **Integration**
  - [x] Add `DashboardPage` to main routing
  - [x] Verify persistence by reloading page

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
│   ├── WidgetPlaceholder.tsx
│   └── WidgetConfigModal.tsx
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
*   [x] Confirmed `react-grid-layout` installation
*   [x] Verified persistence works
*   [x] Checked responsiveness on mobile breakpoint
*   [x] Added widget configuration modal for NATS subject
*   [x] Fixed TypeScript type safety issues
*   [x] Added browser compatibility for UUID generation

### Code Review Fixes Applied
*   Fixed TypeScript `any` type → proper `Layout` type
*   Added crypto.randomUUID() fallback for older browsers
*   Implemented widget configuration UI (AC #3)
*   Added `updateWidget` action to store
*   Created `WidgetConfigModal` component

### File List
*   `viz/package.json`
*   `viz/package-lock.json`
*   `viz/src/App.tsx`
*   `viz/src/main.tsx`
*   `viz/src/features/dashboard/types/index.ts`
*   `viz/src/features/dashboard/stores/useDashboardStore.ts`
*   `viz/src/features/dashboard/components/DashboardLayout.tsx`
*   `viz/src/features/dashboard/components/DashboardWidget.tsx`
*   `viz/src/features/dashboard/components/WidgetPlaceholder.tsx`
*   `viz/src/features/dashboard/components/WidgetConfigModal.tsx`
*   `viz/src/pages/DashboardPage.tsx`
