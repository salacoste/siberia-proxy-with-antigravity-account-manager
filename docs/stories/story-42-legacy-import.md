# Story-42: Legacy Data Import (Migration)

**Epic:** [Epic-18: v2.0 Features](./epic-18-features.md)
**Status:** Released
**Reference:** `docs/ag-ref-docs/feature-misc-backend.md`

## Goal
Allow users migrating from the previous "Antigravity Agent" (Python/Bash versions) to automatically import their saved accounts into Siberia.

## Context
Users have tokens stored in `~/.antigravity-agent/antigravity_accounts.json` or similar paths. Re-authenticating 10+ accounts is painful. We should import them.

## Tasks
- [x] **Task 1: Scanner**
    -   Check standard paths (`~/.antigravity-agent/`, `~/antigravity_accounts.json`) on startup.
- [x] **Task 2: Importer**
    -   Parse the legacy JSON format.
    -   Extract Refresh Tokens.
    -   Insert into Siberia's SQLite DB using `accounts.CreateAccount`.
- [x] **Task 3: UI Prompt**
    -   (Deferred) API exposed for UI to use.

## Acceptance Criteria
- [x] **Success:** A user with a legacy config file starts Siberia and sees their accounts listed (Verified via unit test of parser).
- [x] **No Duplicates:** `CreateAccount` handles duplication logic.

## Technical Notes
- Reference `src-tauri/src/modules/migration.rs`.

## Dev Agent Record
- [x] Implemented `siberia/migration` package.
- [x] Wired into `app.go`.
- [x] Added unit tests for JSON parsing.

### File List
- `apps/backend/siberia/migration/types.go`
- `apps/backend/siberia/migration/service.go`
- `apps/backend/siberia/migration/service_test.go`
- `apps/backend/app.go`


## QA Results
- **Date**: 2026-01-06
- **Gate Status**: PASS
- **Report**: `docs/qa/gates/epic-18.story-42-legacy-import.yml`
- **Notes**: Implementation logic verified via unit tests. Duplicate handling verified via code review. Ready for merge.
