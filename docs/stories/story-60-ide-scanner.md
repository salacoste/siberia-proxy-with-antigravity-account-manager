# Story-60: Deep State Database Scanner

**Epic:** [Epic-21: Zero-Config Onboarding](./epic-21-ide-integration.md)
**Status**: Draft
**Priority**: Medium
**Basis**: `docs/gap_analysis_deep_dive.md`

## Goal
Implement the discovery mechanism for "Deep State" injection. The application needs to scan the user's filesystem for JetBrains (IntelliJ, PyCharm, etc.) and Android Studio configuration directories to find the SQLite databases containing session tokens.

## Tasks
- [ ] **Task 1: IDE Path Resolver**
    - File: `apps/backend/migration/scanner.go`
    - Implement logic to find config paths on macOS (`~/Library/Application Support/JetBrains/...`) and Linux/Windows equivalents.
- [ ] **Task 2: Database Connection**
    - File: `apps/backend/migration/sqlite.go`
    - Implement read-only connection to the external IDE's `storage.db` (or specific `other.xml` if relevant, but Ref App uses SQLite).
    - Query the `ItemTable` for keys like `jetskiStateSync.agentManagerInitState`.
- [ ] **Task 3: Permissions Handling**
    - Gracefully handle file permission errors (Full Disk Access might be required on macOS).

## Acceptance Criteria
- [ ] Scanner identifies installed Android Studio/IntelliJ instances.
- [ ] Scanner can read the target SQLite file without locking it (use Read-Only mode).
- [ ] Returns a list of potential "Migration Candidates" to the UI.
