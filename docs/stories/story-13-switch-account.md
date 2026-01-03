# Story-13: Implement "Switch Account" External App Integration

**Epic:** [Epic-03: Accounts Manager](./epic-03-accounts.md)
**Status:** Draft
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
