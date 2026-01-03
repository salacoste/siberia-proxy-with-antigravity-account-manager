# Story-22: Real-time Log Streaming (Backend)

**Epic:** [Epic-07: Traffic Inspection UI](./epic-07-traffic-inspection.md)
**Status:** Completed

## Goal
Emit events to frontend when requests complete.

## Tasks
- [x] Update `AccessLogMiddleware` to emit Wails Runtime Events.
- [x] Event Name: `proxy:log`.
- [x] Payload: `LogEntry` struct (Method, URL, Status, Duration, Size).

## Acceptance Criteria
- [x] Wails events fire on every request.
- [x] Context is handled correctly to avoid potential nil panics.

## QA Results

### Review Date: 2026-01-04

### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
Backend implementation in `service.go` (`emitFullEvent`) correctly broadcasts `proxy:log`. The event payload includes all necessary metadata (Headers, Body, Timings).

### Compliance Check
- Coding Standards: [✓] Wails Events best practices.
- All ACs Met: [✓] Streaming verified.

### Gate Status
Gate: PASS → docs/qa/gates/epic-07.story-22-log-streaming.yml
