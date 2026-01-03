# Story-37: Add Support for Cursor and Windsurf IDEs

**Epic:** [Epic-12: Native Integrations beyond VS Code](./epic-12-native-integrations.md)
**Status:** Draft

## Goal
Extend the `ProcessManager` and `Injector` to handle Cursor and Windsurf, which are forks of VS Code but use different paths/DB names.

## Tasks
- [ ] **Task 1: Reconnaissance**
    -   Identify Process Names (e.g., `Cursor`, `Windsurf`).
    -   Identify DB Paths (`~/Library/Application Support/Cursor/...`).
- [ ] **Task 2: Implementation**
    -   Refactor `siberia/modules/injection` to accept TargetType enum.
    -   Add logic for new targets.

## Acceptance Criteria
- [ ] Can kill/restart Cursor.
- [ ] Can inject auth tokens into Cursor's state DB.
