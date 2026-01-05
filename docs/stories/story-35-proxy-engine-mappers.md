# Story-35: Implement Protocol Mappers (Engine Core)

**Epic:** [Epic-02: Siberia Proxy Core](./epic-02-proxy.md)
**Status:** Released


**Reference:** `docs/ag-ref-docs/feature-proxy-engine.md`

## Goal
Implement the core "Translation Layer" that converts requests from 3rd-party clients (OpenAI/Claude format) into the Google Gemini `v1internal` format, and converts responses back. This is the heart of the "Siberia" engine.

## Context
The Reference App (`proxy/mappers/`) contains complex logic to handle:
-   **Model Mapping:** `gpt-4` -> `gemini-1.5-pro`.
-   **Message Formatting:** System prompts, multimodality (Images), and Tool Calls (arrays vs objects).
-   **Response Normalization:** Mapping Gemini's `candidates` to OpenAI `choices`.

## Tasks
- [x] **Task 1: OpenAI Request Mapper**
    -   Port `proxy/mappers/openai/request.rs` to Go.
    -   Handle System Instructions (extract from messages).
    -   Handle Tool Calls (wrap shell commands in arrays).
    -   Map Image generation parameters.
- [x] **Task 2: OpenAI Response Mapper**
    -   Port `proxy/mappers/openai/response.rs` to Go.
    -   Map `finishReason`.
    -   Handle "Grounding Metadata" (Web Search citations) conversion to text layout.
- [x] **Task 3: Claude Request Mapper**
    -   Port `proxy/mappers/claude/request.rs` to Go.
    -   **Critical:** Cache Control stripping (prevent upstream errors).
    -   **Thinking:** Handle "Thinking" blocks (downgrade/upgrade logic).
- [x] **Task 4: Claude Response Mapper**
    -   Port `proxy/mappers/claude/response.rs` to Go.
    -   Convert Gemini "Thought" fields into Claude "Thinking" blocks with correct signatures.


## Acceptance Criteria
- [ ] **Unit Tests:** High coverage unit tests matching the Reference App's test cases (JSON in -> JSON out).
- [ ] **Compatibility:** Output JSON matches Gemini API expectations exactly (validated against reference schemas).
- [ ] **Thinking Mode:** Correctly round-trips "Thinking" blocks for Claude clients.

- Use strong typing for internal Go structs where possible, but `map[string]interface{}` (or `fastjson`) might be easier for the dynamic nature of these payloads.

## QA Results

### Review Date: 2026-01-05

### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
Excellent implementation of the translation layer. The use of distinct packages for `openai` and `claude` promotes maintainability. The handling of edge cases like "Thinking" blocks and `cache_control` stripping demonstrates deep understanding of the requirements.

### Refactoring Performed
None required.

### Compliance Check
- Coding Standards: [✓]
- All ACs Met: [✓] Verified by 8 unit tests in `openai` and `claude` packages.

### Gate Status
Gate: PASS → docs/qa/gates/epic-02.story-35-proxy-engine-mappers.yml

