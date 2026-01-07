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
    - **Specific Payload Schema**:
      - `input`: Array of strings (code lines). Join with newline.
      - `instructions`: String (System Message).
      - `stop`: Array of strings (Forward to upstream).
    - Handle specific Codex tool types if present (Mapping `tool_calls` to OpenAI standard if schema differs, otherwise pass-through).
- [ ] **Task 3: Response Mapping**
    - Convert Chat Response (`choices[0].message.content`) back to Legacy Text Response (`choices[0].text`).
    - ensure `object`: "text_completion" in response JSON.

## Acceptance Criteria
- [ ] `curl /v1/completions` with `model: gpt-3.5-turbo-instruct` (or mapped model) works.
- [ ] "Codex" style request with `instructions` is correctly processed as a System Prompt.
