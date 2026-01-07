# Story-15: Implement Access Log Middleware

**Epic:** [Epic-04: Proxy Monitor](./epic-04-monitor.md)
**Status:** Completed
**Feature Branch:** `antigravity/feat/access-logs`

## Goal
Implement a persistent access log (standard Apache/Nginx format or JSON line) for compliance and debugging, separate from the real-time monitor.

## Tasks
- [ ] **Task 1: Log Rotation & Setup**
    -   Use `lumberjack` or similar for log rotation in `siberia/logger`.
    -   Configure log path in `siberia/config`.
- [ ] **Task 2: Middleware**
    -   Update `siberia/proxy/service.go`.
    -   Write to log file asynchronously to avoid blocking implementation.
    -   Log: Time, Source IP, Method, URL, Status, Size, Duration, User-Agent.
    -   **Privacy:** Mask emails in the request/response context (e.g. `s***@gmail.com`) to prevent PII leakage.


## Acceptance Criteria
- [ ] Requests are written to `access.log` in the config directory.
- [ ] Logs rotate (e.g., at 10MB).
- [ ] Log format is parseable (JSON or CLF). (Preferred: JSONL for modern tools).
- [ ] **Privacy:** Email addresses are masked in the logs.


## QA Results

### Review Date: 2026-01-03

### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
Access logging is implemented using `lumberjack` for rotation and `encoding/json` for structured output, meeting modern observability standards. The implementation in `siberia/logger` is thread-safe and asynchronous (via goroutine) to prevent blocking the proxy.

### Compliance Check
- Coding Standards: [✓] Async logging pattern.
- All ACs Met: [✓] JSON Logs written to `access.log`.

### Gate Status
Gate: PASS → docs/qa/gates/epic-04.story-15-access-log.yml

