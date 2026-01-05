# Story-38: Upstream Resilience (Pools & Rotation)

**Epic:** [Epic-02: Siberia Proxy Core](./epic-02-proxy.md)
**Status:** Ready
**Reference:** `docs/ag-ref-docs/feature-proxy-engine.md`

## Goal
Implement a robust HTTP Client for talking to Google's internal APIs (`cloudcode-pa`). It must handle high concurrency, quota exhaustion (429), and service outages (503) by rotating accounts or switching environments.

## Context
The Reference App uses a sophisticated strategy:
-   **Host Pooling:** 16 idle connections per host.
-   **Environment Fallback:** Tries `prod` -> fails -> tries `daily` (sandbox).
-   **Account Rotation:** If an account hits 429 (Quota), it grabs a *different* token from the `AccountManager` and retries automatically.

## Tasks
- [ ] **Task 1: Upstream Client**
    -   Configure `http.Client` with `MaxIdleConnsPerHost`, `IdleConnTimeout`.
    -   Implement `Do()` method with Retry Policy.
- [ ] **Task 2: Account Rotation Logic**
    -   Inject `AccountService` dependency.
    -   On 429/401:
        -   Mark current token as "Exhausted" (temporary).
        -   Fetch next available active account token.
        -   Retry request.
- [ ] **Task 3: Environment Fallback**
    -   Define endpoints: Primary (`cloudcode-pa`), Fallback (`daily-cloudcode-pa`).
    -   On 5xx: Retry on Fallback endpoint.
- [ ] **Task 4: Retry Strategy**
    -   Implement "Linear Backoff" for 429.
    -   Implement "Exponential Backoff" for 5xx.

## Acceptance Criteria
- [ ] **Resilience:** A 429 error from Google is *never* seen by the user; the proxy masks it by retrying with a new account.
- [ ] **Performance:** High throughput (keep-alive connections leveraged).
- [ ] **Availability:** Works even if the Primary endpoint is flaky (fails over to Daily).

## Technical Notes
- Use a library like `hashicorp/go-retryablehttp` or implement custom middleware for the `http.Client`.
- Reference `proxy/upstream/client.rs`.
