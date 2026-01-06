# Story-71: Robust Proxy Transformation

**Epic:** [Epic-02: Proxy Engine](./epic-02-proxy.md)
**Status**: Draft
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
- [ ] Image generation requests don't fail due to leftover tool definitions.
