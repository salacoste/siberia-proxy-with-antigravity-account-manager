# Story-62: Migration Wizard UI

**Epic:** [Epic-21: Zero-Config Onboarding](./epic-21-ide-integration.md)
**Status**: Completed
**Priority**: High
**Basis**: `docs/gap_analysis_deep_dive.md`

## Goal
Implement a frontend "Migration Wizard" that utilizes the backend Scanner (Story-60) and Decoder (Story-61) to present discovered accounts to the user and allow them to one-click import them into Siberia.

## Tasks
- [ ] **Task 1: Backend Endpoint**
    - File: `apps/backend/siberia/handler/migration.go` (New)
    - Implement `ScanAccounts` endpoint -> Calls `migration.ScanIDEProfiles` -> `migration.ReadStateToken` -> `migration.ExtractRefreshToken`.
    - Returns JSON list of discovered accounts (with masked tokens).
- [ ] **Task 2: Frontend Wizard Component**
    - File: `apps/frontend/src/components/migration/MigrationWizard.tsx`
    - UI: Modal or Card showing "Found X Accounts".
    - Action: "Import Selected" which calls backend `POST /api/migration/import` (or reuses `CreateAccount`).
- [ ] **Task 3: Integration**
    - Show this wizard automatically on first run if accounts are empty (Onboarding flow).

## Acceptance Criteria
- [ ] Backend exposes scanning endpoint.
- [ ] UI lists discovered JetBrains/Android Studio accounts.
- [ ] "Import" button adds them to the local `siberia.db` as valid accounts.

## Dev Agent Record
### File List
- `apps/backend/siberia/migration/types.go`
- `apps/backend/siberia/migration/service.go`
- `apps/backend/app.go`
- `apps/frontend/src/components/migration/MigrationWizard.tsx`
- `apps/frontend/src/App.tsx`

## QA Results
- **Status**: PASS
- **Date**: 2026-01-07
- **Gate**: [epic-21.story-62-migration-wizard.yml](../qa/gates/epic-21.story-62-migration-wizard.yml)
- **Notes**: Fixed backend types and logger issues during QA. Frontend state cleanup applied.


