# Story-43: Proxy Integration & Wiring

**Epic:** [Epic-02: Siberia Proxy Core](./epic-02-proxy.md)
**Status:** Draft
**Reference:** `docs/ag-ref-docs/feature-proxy-engine.md`

## Goal
Wire together the isolated components (Mappers, Handlers, Upstream Client, Account Service) into a functioning Proxy Engine. This is the final step of Epic-02 to enable actual traffic flow.

## Context
We have built:
-   **Mappers:** Story-35/51
-   **Handlers:** Story-36, Story-37
-   **Upstream:** Story-38
-   **Middleware:** Story-08 (Basic Auth)

Currently, the `siberia/proxy/service.go` only starts a basic server. It needs to instantiate the new components and route traffic.

## Tasks
- [ ] **Task 1: Dependency Injection**
    -   Update `siberia/proxy/service.go` dependencies.
    -   Inject `AccountService` into Proxy Service.
    -   Initialize `upstream.NewGeminiClient(accountService)`.
- [ ] **Task 2: Handler Registration**
    -   Initialize `openai.NewHandler(client)` and `claude.NewHandler(client)`.
    -   Register routes:
        -   `POST /v1/chat/completions` -> OpenAI Handler
        -   `POST /v1/messages` -> Claude Handler
- [ ] **Task 3: Gateway Logic**
    -   Ensure standard `goproxy` listeners can seamlessly route these requests (Virtual Endpoints) OR use a dedicated internal port/mux.
    -   Verify `localhost:PORT/v1/chat/completions` works directly.

## Acceptance Criteria
- [ ] **End-to-End Flow:** A `curl` request to `localhost:3000/v1/chat/completions` returns a valid response from Gemini (via the Proxy).
- [ ] **Auth:** Request fails without authentication (if enabled).
- [ ] **Streaming:** SSE works end-to-end.

## Technical Notes
-   Consider if we run these handlers on the *same* port as the Forward Proxy (via `goproxy` hijacking or custom listeners) or a separate "Gateway" port.
-   Ref App runs everything on one port usually, intercepting specific paths. `goproxy` allows intercepting requests. Ideally we use `NonProxyHandler` or just checks in `OnRequest`.
