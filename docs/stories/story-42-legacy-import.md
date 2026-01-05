# Story-42: Legacy Data Import (Migration)

**Epic:** [Epic-05: Distribution](./epic-05-distribution.md)
**Status:** Ready
**Reference:** `docs/ag-ref-docs/feature-misc-backend.md`

## Goal
Allow users migrating from the previous "Antigravity Agent" (Python/Bash versions) to automatically import their saved accounts into Siberia.

## Context
Users have tokens stored in `~/.antigravity-agent/antigravity_accounts.json` or similar paths. Re-authenticating 10+ accounts is painful. We should import them.

## Tasks
- [ ] **Task 1: Scanner**
    -   Check standard paths (`~/.antigravity-agent/`, `~/antigravity_accounts.json`) on startup.
- [ ] **Task 2: Importer**
    -   Parse the legacy JSON format.
    -   Extract Refresh Tokens.
    -   Insert into Siberia's SQLite DB using `accounts.CreateAccount`.
- [ ] **Task 3: UI Prompt**
    -   (Optional) Ask user "Found X legacy accounts. Import?" or just do it silently and notify.

## Acceptance Criteria
- [ ] **Success:** A user with a legacy config file starts Siberia and sees their accounts listed.
- [ ] **No Duplicates:** Does not re-import if email already exists.

## Technical Notes
- Reference `src-tauri/src/modules/migration.rs`.
