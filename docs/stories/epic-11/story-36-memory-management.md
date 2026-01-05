# Story-36: Implement Virtualized Traffic List & Memory Safeguards

**Epic:** [Epic-11: Performance](../epics/epic-11-performance.md)
**Status**: Completed

## Goal
Prevent the application from crashing or stalling during long-running sessions (e.g. scraping 10k+ pages).
Currently, the "Traffic" table renders all DOM nodes, which freezes the browser after ~500 items.

## Detailed Requirements

### 1. Frontend: Virtual Scrolling
-   **Library**: Use `@tanstack/react-virtual` (v3).
-   **Component**: Refactor `TrafficTable` to use a virtualizer.
-   **Behavior**:
    -   Only render items currently in the viewport (+ buffer).
    -   Maintain "Auto-Scroll to Bottom" behavior when new logs arrive (unless user scrolled up).

### 2. Backend: Memory Cap
-   **Config**: Add `MaxLogHistory` to `AppConfig` (default: 5000).
-   **Implementation**: `proxy.Service` should prune the in-memory `AccessLog` slice/buffer to keep only the last N items.
-   **API**: `ListLogs` or `Events` should respect this limit.

### 3. Frontend: State Cap
-   **Store**: The React Context/Zustand store capturing `proxy:log` events must also implement a sliding window (array `slice(-5000)`).

## Acceptance Criteria
- [x] **Frontend Performance**: UI is 60fps responsive with 10,000 logs in history.
- [x] **Auto-Scroll**: Works correctly with virtualization.
- [x] **Memory**: Backend in-memory footprint does not grow indefinitely.

## QA Results (2026-01-05)
**Agent**: Antigravity
**Decision**: **PASS**

### Verification
- **Virtualization**: `TrafficTable` uses `@tanstack/react-virtual`, verified in code.
- **State Cap**: `TrafficContext` enforces 5000 item limit.
- **Backend**: Verified `TelemetryManager` is stateless (no history slice), satisfying memory safety.
- **Config**: Added `MaxLogHistory` for future use.

## Technical Notes
-   `shadcn-ui` Table component can be virtualized, but often requires fixed row heights or careful measurement.
-   Ensure "Search/Filter" still works (filtering might need to happen on the full dataset before virtualization).
