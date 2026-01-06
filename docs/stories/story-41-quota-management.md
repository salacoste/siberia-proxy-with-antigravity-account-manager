# Story-41: Quota Management & Paid Tier Detection

**Epic:** [Epic-03: Accounts Manager](./epic-03-accounts.md)
**Status:** Ready for Dev
**Reference:** `docs/ag-ref-docs/feature-misc-backend.md`
**Technical Approach:** `docs/architecture/modules/quota-service.md`

## Goal
Implement the logic to fetch quota usage from Google's internal APIs and correct detection of "Paid Tier" status (Free vs Pro vs Ultra), essential for the "Smart Pooling" logic.

## Context
The Reference App (`modules/quota.rs`) calls `v1internal:loadCodeAssist` to get the `project_id` and billing status, then `v1internal:fetchAvailableModels`.
It distinguishes between:
-   **Free:** Standard limits.
-   **Pro/Ultra:** Higher limits.
-   **Forbidden (403):** Account cannot be used.

## Tasks
- [ ] **Task 1: Quota Service**
    -   Create `siberia/quota` package.
    -   Implement `FetchQuota(token)` method.
    -   Map the complex nested JSON response from Google.
- [ ] **Task 2: Tier Detection**
    -   Implement logic: if `paid_tier` is true -> Pro/Ultra.
    -   Store this metadata in the Account record (`stats` JSON field).
- [ ] **Task 3: Integration**
    -   Call this service periodically (e.g. every 30 mins) or on "Refresh" click.

## Acceptance Criteria
- [ ] **Accuracy:** Correctly identifies a "Gemini Advanced" subscription vs a Free one.
- [ ] **Resilience:** Handles 403 Forbidden without crashing, marking account as invalid.

## Technical Notes
- **API:** `https://cloudcode-pa.googleapis.com/v1internal/projects/...`
- **Fields:** Look for `account_state` and `reason` in the response.
