# Story-67: Attribution Headers & Access Log Middleware

**Epic:** [Epic-04: Monitoring & Logs](./epic-04-monitor.md)
**Status**: Draft
**Priority**: Low
**Basis**: `src-tauri/src/proxy/middleware/attribution_headers.rs`

## Goal
Implement missing middleware that injects attribution headers into HTTP responses. This allows downstream clients (agents, scripts) to know *which* account and model actually served their request, which is vital for debugging sticky sessions and routing logic.

## Tasks
- [ ] **Task 1: Attribution Middleware**
    - File: `apps/backend/proxy/middleware/attribution.go`
    - Check config `response_attribution_headers`.
    - Inject headers:
        - `x-antigravity-provider`: (google/zai)
        - `x-antigravity-model`: (resolved model name)
        - `x-antigravity-account`: (masked account ID or email hash)
- [ ] **Task 2: Access Log Configuration**
    - File: `apps/backend/proxy/middleware/logger.go`
    - Ensure `access_log_enabled` config toggles a concise one-line log: `[Access] POST /v1/chat/completions 200 450ms`.

## Acceptance Criteria
- [ ] `curl -v ...` shows `x-antigravity-account: ***` in response headers.
- [ ] Logs show clean one-liners when enabled.
