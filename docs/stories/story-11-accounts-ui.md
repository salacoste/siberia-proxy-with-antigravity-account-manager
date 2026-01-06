# Story-11: Create Accounts Screen UI (List/Grid)

**Epic:** [Epic-02: Siberia Proxy Core](./epic-02-proxy.md)
**Status:** Completed
**Feature Branch:** `antigravity/feat/upstream-proxy`
**Feature Branch:** `antigravity/feat/sqlite-layer`

## Goal
Create the frontend interface to view and manage Google Accounts stored in the local SQLite database.

## Tasks
- [ ] **Task 1: Backend API**
    -   Create `siberia/accounts/service.go`.
    -   Implement methods: `ListAccounts()`, `DeleteAccount(id)`, `GetAccount(id)`.
    -   Expose these via `App` struct (Wails binding).
- [ ] **Task 2: Frontend UI**
    -   Create `src/pages/AccountsPage.tsx`.
    -   Use `shadcn/ui` Table component.
    -   Columns: Email, Proxy Group, Status (Active/Inactive), Last Used, Actions (Edit/Delete).
    -   Add "Add Account" button (placeholder for now, links to Story-12).

## Acceptance Criteria
- [ ] `AccountsPage` displays a list of accounts fetched from the backend.
- [ ] Sensitive data (password/token) is NOT sent to the frontend (backend DTOs should omit them).
- [ ] "No Accounts" state is handled gracefully.

## Technical Notes
- **DTOs:** Create a specific struct for frontend display (e.g., `AccountDTO`) that excludes `Password` and `SessionToken`.
- **UI:** Use `Table` from shadcn-ui.

## QA Results

### Review Date: 2026-01-03

### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
Accounts UI (`AccountsPage.tsx`) correctly implements the ShadCN `Table` component. Data fetching logic handles the Wails boundary appropriately with mock fallbacks for dev environments.

### Compliance Check
- Coding Standards: [✓] React best practices.
- All ACs Met: [✓] List view works.

### Gate Status
Gate: PASS → docs/qa/gates/epic-03.story-11-accounts-ui.yml

### Manual Verification (2026-01-06)
- **Method**: Browser Subagent
- **Status**: Verified
- **Notes**: Account List renders correctly. Active status displayed.

