# Story-27: MitM Proxy Logic (Backend)

**Epic:** [Epic-08: HTTPS Decryption (MitM)](../epic-08-mitm.md)
**Status**: Completed
**Parent**: Epic-08

## Description
Configure the `proxy` service to use the Root CA managed by `siberia/ca` to intercept HTTPS CONNECT requests. This enables the proxy to decrypt traffic, inspect headers/bodies, and re-encrypt it for the client.

## Requirements
1.  **CA Integration**: Use `ca.Service.GetCAPair()` to retrieve the signing certificate.
2.  **MitmConnect**: Configure `goproxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)` (or conditional based on config).
3.  **Toggle**:
    -   Add `MitmEnabled` boolean to `AccessLogMiddleware` (or Proxy Service state).
    -   If `MitmEnabled` is false, fallback to `goproxy.AlwaysReject` (for MITM) or just pass through via `goproxy.AlwaysMitm` but without inspection? No, if MITM is disabled, we should just tunnel (which `goproxy` does by default if we don't say `HandleConnect(Mitm)`).
    -   **Correction**: `goproxy` default for CONNECT is blindly tunnel. To MITM, we MUST call `HandleConnect(goproxy.AlwaysMitm)`.
    -   So, if `MitmEnabled` is false -> Do `HandleConnect(goproxy.AlwaysReject)` (don't mitm) or simply *don't* add the `AlwaysMitm` handler.
4.  **Safety**:
    -   Handle CA load failures gracefully (fallback to blind tunnel).

## Acceptance Criteria
- [x] Proxy service loads CA pair on startup.
- [x] If `MitmEnabled` (config) is true, `goproxy` intercepts HTTPS.
- [x] If `MitmEnabled` is false, `goproxy` tunnels HTTPS blindly (Status 200/0, no body).
- [x] Decrypted traffic logs headers and body to `AccessLogMiddleware`.

## Dev Agent Record

### Completion Notes
- Updated `AppConfig` to include `MitmEnabled`.
- Verified `NewService` dynamically checks `MitmEnabled` in the `HandleConnect` handler.
- Verified `goproxy.GoproxyCa` global is set via `mitm_test.go`.
- Code supports dynamic toggling if `svc.config` is updated in memory (pointer usage).

### Story DoD Checklist
- [x] AppConfig updated.
- [x] Service Logic Verified.
- [x] Test coverage for CA loading.

## Dev Notes
-   `goproxy` requires setting `proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)`
-   `proxy.MitmChunkedResponse = true` might be needed for streaming?

## QA Results

### Review Date: 2026-01-05

### Reviewed By: Quinn (Test Architect)

- **Audit**: `MitmEnabled` config integration verified. `service.go` logic handles dynamic switching.
- **Verdict**: Approved.

### Gate Status

Gate: PASS → docs/qa/gates/epic-08.story-27-mitm-logic.yml
