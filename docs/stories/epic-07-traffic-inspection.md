# Epic-07: Traffic Inspection UI

**Goal:** Provide a real-time view of HTTP/SOCKS traffic passing through the proxy, similar to "Charles Proxy" or "Network Tab".

## Existing System Context
*   **Current Functionality:** `access_log.go` writes JSONL to disk (or stdout). No live UI.
*   **Tech Stack:** Go (Backend/Proxy), React (Frontend).
*   **Integration Points:** `siberia/proxy` middleware needs to broadcast events to Wails/Frontend.

## Enhancement Details
*   **What:**
    *   **Backend:** Implement a mechanism (Wails Events or WebSocket) to stream request metadata to the frontend.
    *   **Frontend:** A new page `/monitor` (already in routing from Story-03, but empty/placeholder) with a data grid.
*   **Performance:** High volume traffic shouldn't freeze the UI. Use windowing or capping (e.g., enable/disable toggle).

## Stories

### Story-22: Real-time Log Streaming (Backend)
*   **Goal:** Emit events to frontend when requests complete.
*   **Tasks:**
    *   Update `AccessLogMiddleware` to emit Wails Runtime Events (`runtime.EventsEmit`).
    *   Event Name: `proxy:log`
    *   Payload: `LogEntry` struct (Method, URL, Status, Duration, Size).

### Story-23: Traffic Inspector UI (Frontend)
*   **Goal:** Display streaming logs in a virtualized table.
*   **Tasks:**
    *   Update `MonitorPage.tsx`.
    *   Subscribe to `proxy:log` events using `EventsOn`.
    *   Store logs in ephemeral state (Zustand or local state), capped at 1000 items.
    *   Render Table: Timestamp, Method, URL, Status.

### Story-24: Request Details & Body Capture
*   **Goal:** Inspect headers and bodies.
*   **Tasks:**
    *   Backend: Add toggle to `ProxyService` to "Capture Body" (default off for performance/security).
    *   Frontend: Clicking a row shows a Detail View (Headers, Body preview).

## Compatibility Requirements
*   [x] Non-blocking for main proxy loop.
*   [x] Toggleable (don't stream if no one is watching? Or stream always?). -> Stream always if "Monitoring" enabled.

## Risk Mitigation
*   **Risk:** Flood of events freezing the UI.
*   **Mitigation:** Wails events are async. Frontend should throttle renders or cap list size.
