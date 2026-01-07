# Story-14: Implement Proxy Traffic Monitor

**Epic:** [Epic-04: Proxy Monitor](./epic-04-monitor.md)
**Status:** Completed
**Feature Branch:** `antigravity/feat/traffic-monitor`

## Goal
Provide a real-time (or near real-time) view of HTTP requests passing through the proxy.

## Tasks
- [ ] **Task 1: Backend Middleware**
    -   Update `siberia/proxy/service.go`.
    -   Capture Request (Method, URL, Host) and Response (Status, Size, Duration).
    -   Emit Wails Event `proxy:request` with this payload.
- [ ] **Task 2: Frontend UI**
    -   Update `src/pages/MonitorPage.tsx`.
    -   Listen for `proxy:request` events.
    -   Display in a scrolling Log/Table (Shadcn Table or ScrollArea).
    -   Show: Method (Badge), URL, Status (Color-coded), Duration.

## Acceptance Criteria
- [ ] When I make a request through the proxy, it appears on the Monitor page.
- [ ] List auto-scrolls or new items appear at top.
- [ ] Valid Status codes (200) are green, Errors (4xx/5xx) are red.

## Technical Notes
- **Event-Driven:** Use `runtime.EventsEmit` for live updates.
- **Performance:** Frontend should limit list size (e.g., last 50 requests) to prevent memory leaks during long sessions.

## QA Results

### Review Date: 2026-01-03

### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
Real-time monitoring is implemented via Wails Events (`proxy:log`). The frontend `MonitorPage.tsx` correctly subscribes to these events and displays them in a `Table`. Payload capture (Body) fits the "Traffic Recorder" requirement.

### Compliance Check
- Coding Standards: [✓] React/Wails Events integration.
- All ACs Met: [✓] Live traffic visible, Status codes colored.

### Gate Status
Gate: PASS → docs/qa/gates/epic-04.story-14-traffic-monitor.yml

