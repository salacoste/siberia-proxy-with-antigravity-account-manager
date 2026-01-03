# Story-08: Implement Proxy Authentication (Bearer Token)

**Epic:** [Epic-02: Siberia Proxy Core](./epic-02-proxy.md)
**Status:** Draft

## Goal
Secure the local proxy server so that only authorized clients (configured with a specific token) can use it. This prevents unauthorized usage if the proxy is exposed on a public interface (e.g., `0.0.0.0`).

## Tasks
- [ ] **Task 1: Configuration**
    -   Add `AuthEnabled` (bool) and `AuthToken` (string) to `AppConfig`.
- [ ] **Task 2: Middleware Implementation**
    -   Forward Proxy: Check `Proxy-Authorization: Bearer <token>`.
    -   Reverse Proxy (Gateway): Check `Authorization: Bearer <token>`.
    -   Reject with `407 Proxy Authentication Required` or `401 Unauthorized`.
- [ ] **Task 3: Settings UI**
    -   Add Toggle and Token Input in `SettingsPage`.

## Acceptance Criteria
- [ ] When enabled, requests without valid token are rejected.
- [ ] Works for standard `http_proxy` usage (using `Proxy-Authorization`).
- [ ] Works for "API Gateway" usage (using `Authorization`).

## Technical Notes
- **Goproxy:** We can attach a handler to `proxy.OnRequest()`.
- **Security:** Token storage in `config.json` is acceptable for now (local user tool).

## QA Results

### Review Date: 2026-01-03

### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
Security middleware acts as a gatekeeper at the very top of `ServeHTTP`. It correctly handles both `Proxy-Authorization` (forward) and `Authorization` (gateway) headers, and importantly strips them before forwarding to prevent upstream leaks.

### Compliance Check
- Coding Standards: [✓] Early exit pattern used.
- All ACs Met: [✓] Unauthorized requests rejected (401/407).

### Gate Status
Gate: PASS → docs/qa/gates/epic-02.story-08-proxy-auth.yml

