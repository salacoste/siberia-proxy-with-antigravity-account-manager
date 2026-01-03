# Story-15: Implement Access Log Middleware

**Epic:** [Epic-04: Proxy Monitor](./epic-04-monitor.md)
**Status:** Draft
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

## Acceptance Criteria
- [ ] Requests are written to `access.log` in the config directory.
- [ ] Logs rotate (e.g., at 10MB).
- [ ] Log format is parseable (JSON or CLF). (Preferred: JSONL for modern tools).
