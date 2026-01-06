# Story-64: IDE State Injector (Write)

**Epic:** [Epic-21: Zero-Config Onboarding](./epic-21-ide-integration.md)
**Status**: Draft
**Priority**: High (Risky)
**Basis**: `docs/gap_analysis_deep_dive.md`

## Goal
The counterpart to "Deep State Reading". This story enables the application to *write* the currently selected account's session *back* into the external IDE (Antigravity/VSCode) database. This allows the user to management accounts in Siberia and "Push" them to the IDE.

## Tasks
- [ ] **Task 1: Process Management**
    - File: `apps/backend/sys/process.go`
    - Implement `CloseAntigravity()`: finding and terminating the IDE process.
    - Implement `StartAntigravity()`: relaunching it.
- [ ] **Task 2: Protobuf Injection**
    - File: `apps/backend/migration/injector.go`
    - Logic: Read DB -> Decode Base64 -> Remove Field 6 -> Insert New Field 6 (with new Token) -> Encode -> Write DB.
    - Set `antigravityOnboarding = true`.
- [ ] **Task 3: Backup Safety**
    - Always backup `state.vscdb` to `state.vscdb.backup` before writing.

## Acceptance Criteria
- [ ] User clicks "Switch to IDE".
- [ ] IDE quits (if running).
- [ ] Token is injected.
- [ ] IDE relaunches and is logged in as the selected user.
