# Story-90: Add Support for Cursor and Windsurf IDEs

**Epic:** [Epic-12: Native Integrations beyond VS Code](./epic-12-native-integrations.md)
**Status:** Completed

## Goal
Extend the `ProcessManager` and `Injector` to support **Cursor**, a popular AI code editor fork of VS Code.
Also enable support for **Windsurf** (experimental).

## Detailed Requirements

### 1. Refactor: Polymorphic Injector
-   **Current**: `Injector` assumes VS Code paths.
-   **New**: Introduce `TargetIDE` interface or Enum (`VSCode`, `Cursor`, `Windsurf`).
-   **Config**: Allow user to select target IDE in `siberia` config or auto-detect.

### 2. Cursor Support (Prioritized)
- [x] **Task 1: Reconnaissance**
    -   Identify Process Names (e.g., `Cursor`, `Windsurf`).
    -   Identify DB Paths (`~/Library/Application Support/Cursor/...`).
- [x] **Task 2: Implementation**
    -   Refactor `siberia/modules/injection` to accept TargetType enum.
    -   Add logic for new targets.

### 3. Windsurf Support (Experimental)
-   **Process Name**: `Windsurf` (TBD)
-   **DB Path**: TBD. Add configuration option `CustomVSCodePath` to allow users to manually point to other forks.

## Implementation Plan
1.  **Refactor**: Modify `siberia/modules/injection` to accept a `Target` struct containing `Name`, `ProcessName`, `DBPath`.
2.  **Detect**: Check for existence of paths on startup.
3.  **UI**: Add "Target IDE" dropdown in `AccountsPage` or `SettingsPage`.

## Acceptance Criteria
-   [x] `siberia` can detect running Cursor process.
-   [x] `siberia` can inject auth token into Cursor's `state.vscdb`.
-   [x] "Kill/Restart" command works for Cursor.

## QA Results (2026-01-04)
**Agent**: Quinn (QA)
**Decision**: **PASS** (via Gate `epic-12.story-37-ide-support.yml`)

### Verification
- **Registry**: Validated `ide/profiles.go` contains correct paths for VS Code, Cursor, and Windsurf.
- **Dynamic Injection**: Verified `AccountService` resolves target IDE from config.
- **UI**: Settings page includes IDE selector.
