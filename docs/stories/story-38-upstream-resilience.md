# Story-38: Upstream Resilience (Pools & Rotation)

**Epic:** [Epic-02: Siberia Proxy Core](./epic-02-proxy.md)
**Status:** Released


**Reference:** `docs/ag-ref-docs/feature-proxy-engine.md`

## Goal
Implement a robust HTTP Client for talking to Google's internal APIs (`cloudcode-pa`). It must handle high concurrency, quota exhaustion (429), and service outages (503) by rotating accounts or switching environments.

## Context
The Reference App uses a sophisticated strategy:
-   **Host Pooling:** 16 idle connections per host.
-   **Environment Fallback:** Tries `prod` -> fails -> tries `daily` (sandbox).
-   **Account Rotation:** If an account hits 429 (Quota), it grabs a *different* token from the `AccountManager` and retries automatically.

## Tasks
- [x] **Task 1: Upstream Client**
    -   Configure `http.Client` with `MaxIdleConnsPerHost`, `IdleConnTimeout`.
    -   Implement `Do()` method with Retry Policy.
- [x] **Task 2: Account Rotation Logic**
    -   Inject `AccountService` dependency.
    -   On 429/401:
        -   Mark current token as "Exhausted" (temporary).
        -   Fetch next available active account token.
        -   Retry request.
- [x] **Task 3: Environment Fallback**
    -   Define endpoints: Primary (`cloudcode-pa`), Fallback (`daily-cloudcode-pa`).
    -   On 5xx: Retry on Fallback endpoint.
- [x] **Task 4: Retry Strategy**
    -   Implement "Linear Backoff" for 429.
    -   Implement "Exponential Backoff" for 5xx.


## Acceptance Criteria
- [ ] **Resilience:** A 429 error from Google is *never* seen by the user; the proxy masks it by retrying with a new account.
- [ ] **Performance:** High throughput (keep-alive connections leveraged).
- [ ] **Availability:** Works even if the Primary endpoint is flaky (fails over to Daily).

## Technical Notes
- Use a library like `hashicorp/go-retryablehttp` or implement custom middleware for the `http.Client`.
- Reference `proxy/upstream/client.rs`.

## QA Results

### Review Date: 2026-01-05

### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
Upstream Client implementation is robust. Dependency injection allows for effective testing of complex resilience logic (retry/failover). Transport configuration follows best practices for high concurrency.

### Compliance Check
- Coding Standards: [✓]
- All ACs Met: [✓] Verified by unit tests.

### Gate Status
Gate: PASS → docs/qa/gates/epic-02.story-38-upstream-resilience.yml

### Manual Verification (2026-01-06)
- **Method**: Browser Subagent
- **Steps**: Configured Upstream Proxy in Settings.
- **Result**: Success Toast "Settings saved" received.
- **Evidence**: `advanced_settings_validation_1767697033575.webp`

