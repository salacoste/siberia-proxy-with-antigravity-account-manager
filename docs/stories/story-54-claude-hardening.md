# Story-54: Claude Protocol Hardening

**Epic:** [Epic-20: Advanced Traffic Scheduling](./epic-20-advanced-scheduling.md)
**Status**: Completed
**Priority**: Critical
**Basis**: `docs/gap_analysis_deep_dive.md`

## Goal
Implement robust sanitization and handling for Claude 3.7 "Thinking" blocks to prevent `400 Invalid Argument` errors from upstream providers (Vertex AI/Google). This ensures stability for reasoning models.

## Tasks
- [ ] **Task 1: Thinking Block Validator**
    - File: `apps/backend/proxy/handlers/claude.go`
    - Implement regex/logic to validate "Thinking" block signatures (e.g., minimum length, valid start/end tags).
    - Logic: If a block is partial or invalid, convert it to a standard text block or drop it to prevent upstream rejection.
- [ ] **Task 2: Trailing Block Cleanup**
    - Logic: Detect and remove trailing, unsigned, or empty thinking blocks often accumulated at the end of context.
- [ ] **Task 3: Vertex AI Error Handler**
    - File: `apps/backend/proxy/retry.go` (or similar)
    - Implement specific retry logic for `400 Invalid Argument` that identifies signature issues and retries *without* the thinking blocks (fallback strategy).

## Acceptance Criteria
- [ ] Unit Test: Invalid thinking block input -> Sanitized/Converted output.
- [ ] Unit Test: Context with trailing partial thinking block -> Trailing block removed.
- [ ] Integration: Claude 3.7 requests with "Thinking" enabled succeed without 400 errors.

## QA Results
- **Date**: 2026-01-06
- **Reviewer**: Quinn (QA)
- **Status**: **PASS**
- **Notes**:
    - Robustness logic verified via `sanitize_test.go`.
    - Fallback mechanism for 400 errors is a strong safety net.
    - Code quality adheres to standards.

