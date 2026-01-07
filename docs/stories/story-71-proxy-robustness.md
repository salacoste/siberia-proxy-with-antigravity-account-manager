# Story-71: Robust Proxy Transformation

**Epic:** [Epic-02: Proxy Engine](./epic-02-proxy.md)
**Status**: Completed
**Priority**: Medium
**Basis**: `proxy/mappers/gemini/wrapper.rs` & `openai/request.rs`

## Goal
Implement advanced request cleaning and transformation logic found in the reference app to ensure compatibility with various clients (Cherry Studio, etc.) and strict upstream requirements.

## Tasks
- [ ] **Task 1: Deep Cleaning**
    - Implement `deep_clean_undefined` to recursively remove strings like "[undefined]".
- [ ] **Task 2: Schema Sanitization**
    - Implement `clean_json_schema`: remove `format`, `strict`, `additionalProperties` recursively.
    - Implement `enforce_uppercase_types`: ensure `type: "OBJECT"` (uppercase) for Gemini.
- [ ] **Task 3: Tool Fixes**
    - Filter out `web_search` if re-injecting.
    - Handle `imageConfig` injection vs `tools` conflict.

## Acceptance Criteria
- [ ] Requests with `[undefined]` strings are cleaned.
- [ ] Function calling schemas with `type: "object"` (lowercase) are converted to uppercase.
- [x] Image generation requests don't fail due to leftover tool definitions.

## Dev Agent Record
### File List
- `apps/backend/siberia/proxy/mappers/robustness.go` [NEW]
- `apps/backend/siberia/proxy/mappers/robustness_test.go` [NEW]
- `apps/backend/siberia/proxy/providers/zai.go` [MODIFY]
- `apps/backend/siberia/proxy/upstream/gemini.go` [MODIFY]

### Change Log
1.  **Robustness Logic**: Created `mappers/robustness.go` with `DeepClean`, `SanitizeSchema`, and `FilterWebSearch`.
2.  **Z.ai Provider**: Integrated `DeepClean` and `SanitizeSchema` into `ForwardAnthropicJSON` to fix compatibility with Z.ai/Anthropic upstreams.
3.  **Gemini Client**: Added `cleanAndMarshal` helper to `GeminiClient` to ensure all native Gemini requests (and OpenAI-shimmed ones) are sanitized before transmission.
4.  **Verification**: Added comprehensive unit tests in `robustness_test.go`. Verified integration with existing proxy tests.
### QA Results
- **Date**: 2026-01-07
- **Evaluator**: Quinn (QA)
- **Status**: PASS
- **Gate Artifact**: `docs/qa/gates/epic-02.story-71-proxy-robustness.yml`
- **Notes**: Implementation is solid. Unit coverage is excellent for the sanitization usage logic. Integration into both Z.ai and Gemini pathways ensures broad coverage.
