# Story-07: Implement z.ai Provider Integration

**Epic:** [Epic-02: Siberia Proxy Core](./epic-02-proxy.md)
**Status:** Draft

## Goal
Integrate "z.ai" as a fallback or alternative provider. The proxy should be able to route requests to z.ai's API instead of the direct target (Google) based on configuration or routing rules.

## Tasks
- [ ] **Task 1: Configuration**
    -   Add `ZaiBaseURL` (string) and `ZaiApiKey` (string) to `AppConfig`.
- [ ] **Task 2: Request Inspection**
    -   Implement `OnRequest` handler in `goproxy`.
    -   Inspect request body/headers to determine if it should be routed to z.ai (e.g., if model is `gpt-4` or user explicitly selects it).
    -   *Constraint:* HTTPS traffic is encrypted (CONNECT). To inspect body, we must enable MITM (Man-In-The-Middle) which requires generating and trusting a CA certificate.
    -   *Alternative:* For now, maybe just route *everything* if a "z.ai mode" is enabled, or route based on header `X-Provider: zai`.
- [ ] **Task 3: Request Rewriting**
    -   Rewrite target URL to `ZaiBaseURL`.
    -   Inject `Authorization: Bearer <ZaiApiKey>`.

## Acceptance Criteria
- [ ] Configurable Z.ai Base URL and API Key.
- [ ] Requests can be successfully routed to z.ai.
- [ ] Authorization header is correctly injected.

## Technical Notes
- **MITM:** To inspect payloads of HTTPS requests (like "model": "gpt-4"), we MUST decrypt SSL. This requires `goproxy.MitmConnect` and a generated CA cert.
- **Scope:** For this story, let's focus on the *mechanism* to route to z.ai, perhaps triggered by a specific header or global toggle first, to avoid full MITM complexity immediately if not needed.
- **Update:** PRD mentions "Dispatch: Round-robin or Fallback".
