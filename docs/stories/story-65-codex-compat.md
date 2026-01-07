# Story-65: Legacy Codex Compatibility

**Epic:** [Epic-02: Proxy Engine](./epic-02-proxy.md)
**Status**: Completed
**Feature Branch**: `antigravity/feat/codex-compat`
**Priority**: Low
**Basis**: `docs/gap_analysis_deep_dive.md`

## Goal
Support legacy OpenAI completions endpoints (`/v1/completions`) and the specific "Codex" style payload (`input` + `instructions`) used by some older agentic tools. The proxy must map these to standard Chat Completion messages.

## Tasks
- [x] **Task 1: Handler Implementation**
    - File: `apps/backend/siberia/proxy/handlers/legacy/handler.go`
    - Handle `/v1/completions` and `/v1/responses`.
- [x] **Task 2: Payload Transformation**
    - Map `input` (code blocks) + `instructions` (system prompt) -> `messages` array.
    - **Specific Payload Schema**:
      - `input`: Array of strings (code lines). Join with newline.
      - `instructions`: String (System Message).
      - `stop`: Array of strings (Forward to upstream).
    - Handle specific Codex tool types if present (Mapping `tool_calls` to OpenAI standard if schema differs, otherwise pass-through).
- [x] **Task 3: Response Mapping**
    - Convert Chat Response (`choices[0].message.content`) back to Legacy Text Response (`choices[0].text`).
    - ensure `object`: "text_completion" in response JSON.

## Acceptance Criteria
- [x] `curl /v1/completions` with `model: gpt-3.5-turbo-instruct` (or mapped model) works.
- [x] "Codex" style request with `instructions` is correctly processed as a System Prompt.

## Dev Agent Record
### Files
- `apps/backend/siberia/proxy/handlers/legacy/handler.go` (NEW: Legacy Handler)
- `apps/backend/siberia/proxy/handlers/legacy/handler_test.go` (NEW: Unit Tests)
- `apps/backend/siberia/proxy/service.go` (MOD: Registered legacy handler)

### Completion Notes
- Implemented `LegacyCompletionRequest` to `ChatCompletionRequest` transformation.
- Mapped Legacy Codex "Instructions" to System Prompt.
- Mapped "Input" strings array to User Prompt.
- Added full unit test coverage for logic.

### Debug Log
- Fixed nil pointer dereference in `Usage` struct mapping when Upstream returns nil usage metadata.

## Change Log
## QA Results

### Review Date: 2026-01-07

### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment

The implementation of the `LegacyHandler` is clean, modular, and effectively bridges the gap between the legacy Codex format and the modern Chat Completion format.
- **Robustness**: The handler correctly manages potential nil references in the `Usage` struct (a common issue with mocked upstream responses).
- **Architecture**: The decision to map `Legacy -> Chat -> Gemini` reuses the existing robust `openai` request mappers, avoiding code duplication at the cost of slight object overhead.
- **Testing**: The unit tests verify the core transformation logic (Input array joining, Instructions to System prompt) and the response mapping.

### Refactoring Performed

- None required during QA. The developer proactively fixed a nil pointer panic during their verification phase.

### Compliance Check

- Coding Standards: [✓] Adheres to Go idioms and project structure.
- Project Structure: [✓] Placed correctly in `handlers/legacy`.
- Testing Strategy: [✓] Unit tests present and passing.
- All ACs Met: [✓]

### Improvements Checklist

- [ ] Future: Add support for `LegacyCompletionRequest.Stream` if older clients require it (currently unary only).

### Security Review

- No new security risks introduced. The handler delegates to the existing `UpstreamClient` which handles authentication and attribution.

### Files Modified During Review

- None.

### Gate Status

Gate: PASS → docs/qa/gates/epic-02.story-65-codex-compat.yml

### Recommended Status

[✓ Ready for Done]
