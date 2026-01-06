# Story-13: Implement "Switch Account" External App Integration

**Epic:** [Epic-03: Accounts Manager](./epic-03-accounts.md)
**Status:** Completed
**Feature Branch:** `antigravity/feat/switch-account`

## Goal
Implement the "Activate" / "Switch Account" functionality that orchestrates the external environment switch.

## Tasks
- [ ] **Task 1: Process Module (`siberia/modules/process`)**
    *   Create helper to Find/Kill/Start external processes.
    *   Implement `ProcessManager` interface.
- [ ] **Task 2: Injection Module (`siberia/modules/injection`)**
    *   Create `Injector` interface.
    *   Implement `MockInjector` (logs changes, validates DB paths).
    *   *Note:* Actual Protobuf injection requires `.proto` files which are not present. We will scaffold the architecture so the real implementation can be swapped in later.
- [ ] **Task 3: Backend Orchestration**
    *   Update `siberia/accounts/service.go` with `ActivateAccount(id)`.
    *   Flow:
        1.  Fetch Account (decrypt tokens).
        2.  `ProcessManager.Kill("TargetApp")`.
        3.  `Injector.Inject(tokens)`.
        4.  `ProcessManager.Start("TargetApp")`.
- [ ] **Task 4: Frontend Integration**
    *   Add "Switch" button to `AccountsPage`.
    *   Display "Activating..." spinner.
    *   Handle errors (e.g., App not found).

## Acceptance Criteria
- [ ] Clicking "Activate" on an account triggers the backend flow.
- [ ] Backend logs the orchestration steps (Kill -> Inject -> Start).
- [ ] Frontend receives success/failure feedback.

## QA Results

### Review Date: 2026-01-03

### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
Backend orchestration in `service.go` correctly implements the Kill -> Inject -> Start flow. Hardcoded paths are acceptable for MVP status as noted in code comments. Error propagation to frontend via Promises/Toast is handled well.

### Compliance Check
- Coding Standards: [✓] Error wrapping pattern.
- All ACs Met: [✓] External app orchestration works.

### Gate Status
Gate: PASS → docs/qa/gates/epic-03.story-13-switch-account.yml

### Manual Verification (2026-01-06)
- **Method**: Browser Subagent
- **Steps**: Used "Activate" button on account list.
- **Result**: Success Toast received. Visual status confirmed.
- **Evidence**: `switch_account_proxy_validation_1767696947035.webp`

