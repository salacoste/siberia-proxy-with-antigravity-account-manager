# Story-70: Z.ai Provider (Anthropic Compatible)

**Epic:** [Epic-19: Z.ai Intelligence](./epic-19-zai-intelligence.md)
**Status**: Completed
**Priority**: High
**Basis**: `proxy/providers/zai_anthropic.rs`

## Goal
Implement a dedicated provider for Z.ai that mimics the Anthropic API but forwards requests to `api.z.ai`. This enables "Hybrid Dispatching" where specific models (e.g., `claude-3-opus`) are routed to Z.ai while others go to Google.

## Tasks
- [ ] **Task 1: Provider Implementation**
    - File: `apps/backend/proxy/providers/zai.go`
    - Implement `ForwardAnthropicJSON`.
    - Handle `x-api-key` auth mapping.
- [ ] **Task 1.5: Configuration Update**
    - File: `apps/backend/config/config.go`
    - Ensure `ZaiConfig` struct has `DispatchMode` (Off, Exclusive, Pooled, Fallback) and `ModelMapping`.
    - Add `ApiKey` field if not present.
- [ ] **Task 2: Model Mapping**
    - Map `claude-3-opus` -> `metis` (or configured Z.ai model).
    - Support `zai:` prefix stripping.

## Acceptance Criteria
- [ ] Requests to `claude-3-opus` (if configured) are routed to Z.ai.
- [ ] Responses are streamed back correctly.
