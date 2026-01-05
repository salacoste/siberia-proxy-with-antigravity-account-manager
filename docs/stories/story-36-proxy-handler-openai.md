# Story-36: OpenAI Handler & SSOP

**Epic:** [Epic-02: Siberia Proxy Core](./epic-02-proxy.md)
**Status:** Ready
**Reference:** `docs/ag-ref-docs/feature-proxy-engine.md`

## Goal
Implement the `http.Handler` for `/v1/chat/completions` that orchestrates the request workflow: Auth -> Mapper -> Upstream -> Mapper -> Response. Crucially, implement the SSOP (Shell Output Parsing) logic to fix broken shell commands from models.

## Context
Some models output shell commands as plain text (Markdown code blocks) instead of proper Tool Calls. The Reference App solves this by scanning the *stream* of text and synthetically injecting a Tool Call event if it detects a shell command pattern.

## Tasks
- [ ] **Task 1: Chat Completion Handler**
    -   Accept `POST /v1/chat/completions`.
    -   Perform Model Routing (using `story-35` logic).
    -   Call Upstream Client (Google API).
    -   Handle Non-Streaming Responses.
- [ ] **Task 2: Streaming Response Support**
    -   Implement SSE (Server-Sent Events) forwarding.
    -   Convert Google's SSE format (`data: {...}`) to OpenAI's SSE format on the fly.
- [ ] **Task 3: SSOP (Shell Output Parsing)**
    -   **Logic:** Implement a Stream Scanner that looks for `json` code blocks with `command="shell"`.
    -   **Action:** If found, suppress the text output and instead emit a `tool_calls` event chunk.
    -   **Reference:** See `proxy/streaming.rs` in Reference App.

## Acceptance Criteria
- [ ] **Standard Chat:** Works with `curl` and standard OpenAI libs.
- [ ] **Streaming:** Works seamlessly (tokens appear in real-time).
- [ ] **SSOP:** A model outputting "```json ... shell ...```" triggers a real tool execution in the client (Cursor/VSCode), not just text display.

## Technical Notes
- This handler sits *behind* the main Proxy mux.
- It is a "Virtual Endpoint" - the client thinks it's hitting OpenAI, but we handle it entirely in code.
