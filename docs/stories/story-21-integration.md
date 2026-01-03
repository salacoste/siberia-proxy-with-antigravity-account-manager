# Story-21: Integration & Verification

**Epic:** [Epic-06: Deep Integration](./epic-06-deep-integration.md)
**Status:** Completed

## Goal
Wire Process Termination and DB Injection into the Backend Service flow.

## Tasks
- [x] Update `accounts.Service` to use `process` and `injector` modules.
- [x] Define orchestration: Kill -> Inject -> Start.

## Acceptance Criteria
- [x] `ActivateAccount` calls steps in correct order.
- [x] Errors in any step abort the flow.

## QA Results

### Review Date: 2026-01-04

### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
Integration logic in `accounts/service.go` correctly sequences the operations. Hardcoded paths are acceptable for the current MVP stage but should be parameterized later.

### Compliance Check
- Coding Standards: [✓] Modular dependency injection.
- All ACs Met: [✓] Orchestration logic verified.

### Gate Status
Gate: PASS → docs/qa/gates/epic-06.story-21-integration.yml
