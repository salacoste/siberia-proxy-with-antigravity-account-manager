# Story-11: Create Accounts Screen UI (List/Grid)

**Epic:** [Epic-03: Accounts Manager](./epic-03-accounts.md)
**Status:** Draft

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
