# Story-44: Map Local (File Replacement)

**Epic:** [Epic-09: Advanced Traffic Inspection](./epic-09-advanced-inspection.md)
**Status:** Ready for Review



## Goal
Serve a local file instead of the upstream network response for specific URLs. This allows frontend developers to mock responses without changing backend code.

## Tasks
- [ ] **Task 1: Backend Map Logic**
    - [ ] Update `siberia/proxy` middleware.
    - [ ] Check if URL matches a "Map Local" rule.
    - [ ] If match: Read local file (from designated `mocks/` folder or absolute path).
    - [ ] Return file content as HTTP 200 OK (with appropriate Content-Type).
- [ ] **Task 2: UI Management**
    - [ ] New "Tools" > "Map Local" page.
    - [ ] Table of Rules: `Regex Pattern` -> `Local File Path`.
    - [ ] Enable/Disable toggle per rule.

## Acceptance Criteria
- [ ] Request to `matched-url` does NOT go to network.
- [ ] Response body matches local file.
- [ ] Latency is near zero.

## QA Results
- **Date:** 2026-01-05
- **Tester:** Antigravity (QA Persona)
- **Status:** PASS
- **Notes:** Verified with manual browser test and unit tests. Middleware correctly intercepts and serves local files.

