# Story 4.2: Trend Chart Component (uPlot)

Status: done

## Story

**As a** Engineer,
**I want** high-performance charts,
**So that** I can zoom into millisecond-level data without the browser freezing.

## Acceptance Criteria

1.  ✅ **Given** a chart displaying 100,000 points
2.  ✅ **When** I zoom or pan
3.  ✅ **Then** the interaction is smooth (60 FPS)
4.  ✅ **And** new real-time points append to the right side instantly
5.  ✅ **And** the chart handles window resizing correctly

## Tasks / Subtasks

- [x] **Setup uPlot Integration**
  - [x] Install `uplot` and `@types/uplot`
  - [x] Create `viz/src/features/dashboard/components/TrendChart.tsx`
  - [x] Create `viz/src/features/dashboard/components/uPlotWrapper.tsx` (Generic wrapper)
  - [x] Implement `useEffect` to initialize `uPlot` instance on mount
  - [x] Implement cleanup to `destroy()` instance on unmount

- [x] **Real-time Data Handling**
  - [x] Connect `TrendChart` to `useDashboardStore` (or a dedicated `useDataStream` store if created)
  - [x] **Critical Optimization:** Do NOT use React State for the data array. Use a `useRef` to hold the data and call `uplot.setData()` directly.
  - [x] Implement a "Streaming Mode" where the x-axis window shifts automatically
  - [x] Implement "Pause/Freeze" on hover or click

- [x] **Zoom & Pan Configuration**
  - [x] Configure `uPlot` options for `scales.x.auto: false` (manual control for streaming)
  - [x] Enable `cursor.drag` for zooming
  - [x] Add "Reset Zoom" button/double-click handler

- [x] **Performance Optimization**
  - [x] Implement throttling for `setData` calls (e.g., max 30fps updates)
  - [x] Ensure `ResizeObserver` handles container resizing efficiently

- [x] **Integration**
  - [x] Replace `WidgetPlaceholder` in `DashboardWidget` with `TrendChart`
  - [x] Verify 100k point performance with mock data generator

## Dev Notes

### Technical Requirements
*   **Library:** `uplot` v1.6+
*   **Wrapper Pattern:**
    ```typescript
    const chartRef = useRef<uPlot | null>(null);
    const containerRef = useRef<HTMLDivElement>(null);
    
    useEffect(() => {
      const u = new uPlot(opts, data, containerRef.current);
      chartRef.current = u;
      return () => u.destroy();
    }, []);
    ```

### Architecture Compliance
*   **State Management:** The chart component should subscribe to the NATS stream (via Zustand or direct WebSocket) but manage its own high-frequency render loop.
*   **Styling:** Use `uPlot`'s default CSS but override colors to match the application theme (Tailwind colors).

### Performance Guardrails
*   **NO React Renders on Data:** The component receiving data must NOT trigger a React re-render for every point. It should append to a mutable buffer and call `uplot.setData()`.
*   **Throttling:** Use `requestAnimationFrame` or a throttle function to limit `setData` calls.

### Project Structure
```
viz/src/features/dashboard/
├── components/
│   ├── TrendChart.tsx      # Business logic (NATS subscription)
│   ├── uPlotWrapper.tsx    # Generic uPlot wrapper
│   └── ChartControls.tsx   # Chart control buttons
├── hooks/
│   └── useTrendData.ts     # Data management hook
└── config/
    └── uPlotConfig.ts      # Chart configuration
```

### References
*   [Epics: Epic 4](./../epics.md#epic-4-real-time-visualization-viz)
*   [uPlot Demos](https://leeoniya.github.io/uPlot/demos/realtime.html)

## Dev Agent Record

### Context Reference
*   **Architecture:** `docs/architecture.md`
*   **Epics:** `docs/epics.md`
*   **Previous Story:** `docs/sprint-artifacts/4-1-real-time-dashboard-framework.md`

### Agent Model Used
Antigravity (Google Deepmind)

### Completion Notes List
*   [x] Confirmed `uplot` installation
*   [x] Verified 60fps with 100k points
*   [x] Verified memory cleanup on unmount
*   [x] Implemented proper useRef pattern for data
*   [x] Added chart controls (pause, reset, export)

### File List
*   `viz/package.json`
*   `viz/src/features/dashboard/components/TrendChart.tsx`
*   `viz/src/features/dashboard/components/uPlotWrapper.tsx`
*   `viz/src/features/dashboard/components/ChartControls.tsx`
*   `viz/src/features/dashboard/hooks/useTrendData.ts`
*   `viz/src/features/dashboard/config/uPlotConfig.ts`
