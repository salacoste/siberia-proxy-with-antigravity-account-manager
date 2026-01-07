# Story-55: Intelligent Rate Limit Parsing

**Epic:** [Epic-20: Advanced Traffic Scheduling](./epic-20-advanced-scheduling.md)
**Status**: Completed
**Priority**: High
**Basis**: `docs/gap_analysis_deep_dive.md`

## Goal
Enhance the reliability of the proxy by implementing intelligent parsing of rate limit errors. Instead of relying solely on HTTP `Retry-After` headers, the system should scrape error bodies for natural language wait times and distinguish between short-term rate limits and long-term quota exhaustion.

## Tasks
- [ ] **Task 1: Regex Parser Implementation**
    - File: `apps/backend/proxy/ratelimit/parser.go`
    - Implement regex patterns to extract durations from error bodies:
        - "Try again in 2m 30s"
        - "quotaResetDelay": "42s"
        - "rate limit exceeded, retry after X"
- [ ] **Task 2: Error Classification**
    - Implement distinct error types:
        - `RateLimitExceeded`: Soft penalty (e.g., 30s-2m).
        - `QuotaExhausted`: Hard penalty (e.g., 1h default if not specified).
- [ ] **Task 3: Integration with Token Manager**
    - Update `TokenManager.ReportError` to accept these parsed durations and apply them to the account's cooldown state.

## Acceptance Criteria
- [ ] Unit Test: Error body "Try again in 5m" -> Returns `300s` duration.
- [ ] Unit Test: Unknown 429 error -> Returns default fallback duration.
- [ ] `QuotaExhausted` error triggers a longer backoff than `RateLimitExceeded`.

## QA Results
- **Status**: PASS
- **Reviewer**: default_qa_agent
- **Date**: 2026-01-06
- **Gate File**: `docs/qa/gates/epic-20.story-55-rate-limit-parsing.yml`
- **Summary**: Rate Limit Parser successfully implemented and tested. Unit tests cover NLP, JSON, and raw formats. Integration into Gemini client is complete.

