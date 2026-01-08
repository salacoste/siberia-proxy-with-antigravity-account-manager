# Story-78: Legacy Codex Stream Support

**Epic:** [Epic-02: Proxy Engine](./epic-02-proxy.md)
**Status**: Done
**Priority**: Medium
**Basis**: Follow-up to Story-65
**File List**:
- `apps/backend/siberia/proxy/handlers/legacy/handler.go`
- `apps/backend/siberia/proxy/handlers/legacy/handler_test.go`


## Goal
Update the `LegacyHandler` (`/v1/completions`) to support Server-Sent Events (SSE) streaming (`stream=true`), enabling compatibility with clients that rely on real-time code completion (e.g., Copilot).

## Tasks
- [ ] **Task 1: Stream Handler**
    - Update `LegacyCompletionRequest` struct to support `stream bool`.
    - Implement `StreamGenerateContent` mapper in `handlers/legacy/handler.go`.
    - Map upstream Gemini SSE events to Legacy OpenAI SSE format (`choices[0].text`).
- [ ] **Task 2: Response Mapping**
    - Ensure `data: [DONE]` is sent at the end.
    - Validate `finish_reason` mapping.

## Acceptance Criteria
- [ ] `curl` with `stream=true` to `/v1/completions` returns `text/event-stream`.
- [ ] Chunks contain `text` deltas, not full messages.
- [ ] Client receives complete generation.

## Dev Notes
- Reuse `upstream/gemini.go`'s `StreamGenerateContent` capability.
- This is purely a translation layer task (Gemini Stream -> OpenAI Legacy Stream).
## QA Results

### Review Date: 2026-01-08
### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
Excellent implementation. The streaming logic is concise and correctly handles the translation between Gemini's modern response format and the legacy OpenAI "text_completion" SSE format. The use of channels effectively bridges the upstream client and the response writer.

### Refactoring Performed
- None required.

### Compliance Check
- Coding Standards: [✓]
- Project Structure: [✓]
- Testing Strategy: [✓] (Unit test mocks upstream stream correctly)
- All ACs Met: [✓]

### Gate Status
Gate: PASS → docs/qa/gates/epic-02.story-78-codex-stream.yml

### Recommended Status
[✓ Ready for Done]
