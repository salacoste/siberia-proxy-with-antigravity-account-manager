# Epic-20: Advanced Traffic Scheduling & Resilience

**Status**: Planning
**Priority**: Critical (Tier 1 - Stability / Tier 3 - Optimization)
**Basis**: `docs/gap_analysis_deep_dive.md` (Section 2, 5, 9)

## Executive Summary
Upgrade the core proxy engine to support "Sticky Sessions" for context caching, intelligent rate limit handling for non-standard errors, and protocol hardening for Claude 3.7. This epic moves the system from simple round-robin rotation to a sophisticated scheduling engine that maximizes upstream performance and stability.

## Key Features
1.  **Claude Protocol Hardening (Critical)**:
    -   Implement sanitization for "Thinking" blocks (validate signatures, remove partials).
    -   Prevent `400 Invalid Argument` errors on Vertex AI.
2.  **Intelligent Rate Limit Parsing**:
    -   Regex-based parsing of error bodies for wait times when headers are missing.
    -   Distinction between short `RateLimit` (30s) and long `QuotaExhausted` (1h) penalties.
3.  **Session Fingerprinting (Sticky Sessions)**:
    -   Generate deterministic Session IDs from `model + initial_prompt_hash`.
    -   Route same-context requests to the same upstream account to leverage KV Caching.
4.  **Smart Token Pooling**:
    -   **Tier Sort**: Prioritize `ULTRA` > `PRO` > `FREE` accounts.
    -   **Scheduling Modes**: Implement `CacheFirst` (Wait if busy) vs `PerformanceFirst` (Rotate immediately).

## Success Metrics
-   [ ] Claude 3.7 "Thinking" responses flow without 400 errors.
-   [ ] "Try again in 42s" body message triggers exact 42s backoff, not default.
-   [ ] Repeated requests with same prompt hit the same upstream account (verified via logs).
