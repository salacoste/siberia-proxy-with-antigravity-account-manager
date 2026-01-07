# Story-06: Implement Upstream Proxy Chaining

**Epic:** [Epic-02: Siberia Proxy Core](./epic-02-proxy.md)
**Status:** Completed

## Goal
Enable the local Siberia proxy to forward all outgoing traffic to a specified **Upstream Proxy** (e.g., a residential proxy provider or Tor), effectively chaining them: `Client -> Siberia -> Upstream -> Target`.

## Tasks
- [ ] **Task 1: Configuration**
    -   Add `UpstreamProxy` (string) field to `AppConfig` (e.g., `http://user:pass@1.2.3.4:8000`).
- [ ] **Task 2: Proxy Transport**
    -   Modify `siberia/proxy/service.go` to configure the `Transtport` of `goproxy`.
    -   Parse the upstream URL.
    -   Set `proxy.Tr.Proxy` to `http.ProxyURL(upstreamUrl)`.
- [ ] **Task 3: Settings UI**
    -   Add an input field in `SettingsPage.tsx` to configure the upstream URL.

## Acceptance Criteria
- [ ] Setting a valid upstream URL in config causes traffic to route through it.
- [ ] Clearing the upstream URL reverts to direct connection.
- [ ] Supports authentication in URL (`user:pass`).

## Technical Notes
- **Library:** `goproxy` uses a custom `http.Transport` for outgoing requests (`proxy.Tr`). We simply need to set its `Proxy` function.
- **Protocol Support:** Ensure `http` and `socks5` schemes are supported in the upstream URL.

## QA Results

### Review Date: 2026-01-03

### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
Upstream proxy logic is correctly implemented by setting the `Transport.Proxy` function in `Start()`. This allows hot-reloading (via restart) and environment variable fallback.

### Compliance Check
- Coding Standards: [✓] Go standard library usage.
- All ACs Met: [✓] Upstream chaining works.

### Gate Status
Gate: PASS → docs/qa/gates/epic-02.story-06-upstream-proxy.yml

