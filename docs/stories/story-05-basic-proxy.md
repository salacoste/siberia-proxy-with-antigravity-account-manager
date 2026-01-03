# Story-05: Implement Basic Request Forwarder (HTTP/SOCKS5)

**Epic:** [Epic-02: Siberia Proxy Core](./epic-02-proxy.md)
**Status:** Draft

## Goal
Implement the core Forward Proxy listener that intercepts HTTP/HTTPS requests from local tools (Cursor, VSCode). This is the entry point for all traffic.

## Tasks
- [ ] **Task 1: Proxy Module Structure**
    -   Create `siberia/proxy` package.
    -   Define `ProxyServer` struct.
- [ ] **Task 2: HTTP/HTTPS Listener**
    -   Implement listener on configurable port (default 8080).
    -   Handling CONNECT method for HTTPS tunneling.
    -   Basic request forwarding to destination (no upstream proxy yet).
- [ ] **Task 3: Integration**
    -   Start/Stop proxy server from Wails `App` struct.
    -   Log basic request info (Method, URL, Status).

## Acceptance Criteria
- [ ] Server listens on configured port (e.g., 3000).
- [ ] Can curl through it: `curl -x localhost:3000 http://example.com`.
- [ ] HTTPS tunneling works: `curl -x localhost:3000 https://google.com`.
- [ ] Basic access logs printed to console/stdout.

## Technical Notes
- **Library:** Use `elazarl/goproxy` or standard `net/http` with `httputil.ReverseProxy`. `goproxy` is recommended for ease of HTTPS CONNECT handling.
- **Axum Note:** Previous docs mentioned Axum (Rust); ignoring this as we are using Go.

## QA Results

### Review Date: 2026-01-03

### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
Core proxy implementation using `goproxy` is cleaner than standard `httputil`. `ServeHTTP` clearly separates security/routing logic from the underlying proxy handler. Integration with Wails runtime for logging is handled well via `emitFullEvent`.

### Compliance Check
- Coding Standards: [✓] Service-based architecture.
- Project Structure: [✓] `siberia/proxy` package.
- All ACs Met: [✓] Proxies HTTP/HTTPS correctly.

### Gate Status
Gate: PASS → docs/qa/gates/epic-02.story-05-basic-proxy.yml

