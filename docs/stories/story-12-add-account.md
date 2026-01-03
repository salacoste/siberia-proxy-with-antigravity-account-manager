# Story-12: Implement Add Account Logic (OAuth Loopback)

**Epic:** [Epic-03: Accounts Manager](./epic-03-accounts.md)
**Status:** Draft
**Feature Branch:** `antigravity/feat/add-account-logic`

## Goal
Implement the ability to add new Google Accounts to the system.
*Note:* The title mentions OAuth Loopback, but for this initial iteration (MVP), we will implement a **Manual Addition** flow (Email + Password + Recovery Email) to populate the DB, as full OAuth integration requires a more complex browser automation setup which might be a subsequent story.
*Correction:* PRD mentions "Simple Add" vs "OAuth". We will implement the **Simple Add Dialog** first to ensure CRUD works.

## Tasks
- [ ] **Task 1: Backend API**
    -   Update `siberia/accounts/service.go`.
    -   Implement `CreateAccount(email, password, recovery, proxyGroup)`.
    -   Encryption is handled automatically by `EncryptedString`.
- [ ] **Task 2: Frontend UI**
    -   Create `components/AddAccountDialog.tsx` (Shadcn Dialog).
    -   Form with validation (clean email, non-empty password).
    -   Integrate into `AccountsPage.tsx`.

## Acceptance Criteria
- [ ] User can click "Add Account".
- [ ] Dialog appears.
- [ ] User enters credentials.
- [ ] Account is saved to SQLite (encrypted).
- [ ] List refreshes automatically.

## QA Results

### Review Date: 2026-01-03

### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
`AddAccountDialog.tsx` correctly implements form validation and uses the Wails binding `CreateAccount` to securely transmit credentials. Error handling is present.

### Compliance Check
- Coding Standards: [✓] Controlled components used.
- All ACs Met: [✓] Account creation verified.

### Gate Status
Gate: PASS → docs/qa/gates/epic-03.story-12-add-account.yml

