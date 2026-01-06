# Story-65: Legacy Codex Compatibility

**Epic:** [Epic-02: Proxy Engine](./epic-02-proxy.md)
**Status**: Draft
**Priority**: Low
**Basis**: `docs/gap_analysis_deep_dive.md`

## Goal
Support legacy OpenAI completions endpoints (`/v1/completions`) and the specific "Codex" style payload (`input` + `instructions`) used by some older agentic tools. The proxy must map these to standard Chat Completion messages.

## Tasks
- [ ] **Task 1: Handler Implementation**
    - File: `apps/backend/proxy/handlers/legacy.go`
    - Handle `/v1/completions` and `/v1/responses`.
- [ ] **Task 2: Payload Transformation**
    - Map `input` (code blocks) + `instructions` (system prompt) -> `messages` array.
    - Handle specific Codex tool types (`local_shell_call`, `web_search_call`) if they appear in the input array.
- [ ] **Task 3: Response Mapping**
    - Convert the Chat Response back to a "Text Completion" response format (`choices[0].text`).

## Acceptance Criteria
- [ ] `curl /v1/completions` with `model: gpt-3.5-turbo-instruct` (or mapped model) works.
- [ ] "Codex" style request with `instructions` is correctly processed as a System Prompt.
