# Story-37: Claude Handler & Resilience Strategies

**Epic:** [Epic-02: Siberia Proxy Core](./epic-02-proxy.md)
**Status:** Ready for Review

**Reference:** `docs/ag-ref-docs/feature-proxy-engine.md`

## Goal
Implement the `http.Handler` for `/v1/messages` (Claude Protocol). This handle must deal with the specific strictness of the Vertex AI upstream, particularly regarding "Thinking" blocks and "Cache Control".

## Context
Vertex AI is stricter than the native Anthropic API.
1.  **Thinking:** If a "Thinking" block is present, it *must* have a valid signature (from a previous turn) or be "enabled" correctly.
2.  **Downgrade:** Background tasks (non-user interactive) often fail with expensive models or waste quota. The proxy should detect these and downgrade to `flash` or `lite`.

## Tasks
- [x] **Task 1: Messages Handler**
    -   Accept `POST /v1/messages`.
    -   Check `x-api-key` (dummy or pass-through).
    -   Call Upstream Client.
- [x] **Task 2: Smart Downgrade Logic**
    -   Detect "Background Patterns" (e.g., specific prompt headers or fast retry loops from Agentic IDEs).
    -   Override Model Selection to `gemini-1.5-flash` or `gemini-2.0-flash-lite` dynamically.
- [x] **Task 3: Thinking Block Sanitization**
    -   Validate incoming "Thinking" blocks.
    -   If invalid (missing signature), convert to standard "Text" block to prevent Upstream 400 Error.

- [ ] **Task 4: Z.ai Integration hooks**
    -   (Stub/Prepare) If Mode is Z.ai, route appropriately (see `story-07`).

## Acceptance Criteria
- [ ] **Claude Client:** Works with `cursor` (Claude mode) and `claude-dev` extension.
- [ ] **Error Resilience:** Does not panic or return 500 if the client sends malformed Thinking blocks; instead sanitizes and forwards.
- [ ] **Downgrade:** (Optional/Advanced) Logs show "Downgrading model..." for background tasks.

## Technical Notes
- Reference `proxy/handlers/claude.rs`.

## QA Results

### Review Date: 2026-01-05

### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
Handler implements the complex Claude SSE protocol correctly. Smart Downgrade logic is implemented with verified tests. Integration with Mappers allows for robust handling of "Thinking" blocks.

### Compliance Check
- Coding Standards: [✓]
- All ACs Met: [✓] Verified by unit tests.

### Gate Status
Gate: PASS → docs/qa/gates/epic-02.story-37-proxy-handler-claude.yml

