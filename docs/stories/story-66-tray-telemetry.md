# Story-66: Advanced Tray Telemetry

**Epic:** [Epic-20: Advanced Traffic Scheduling](./epic-20-advanced-scheduling.md) (or new UX Epic)
**Status**: Draft
**Priority**: Medium
**Basis**: `docs/ag-manager-ref/src-tauri/src/modules/tray.rs`

## Goal
Enhance the System Tray menu to display real-time, granular quota information across different model families without opening the main application window.

## Tasks
- [ ] **Task 1: Model-Specific Quota Extraction**
    - File: `apps/backend/tray/menu.go`
    - Update the logic to fetch specific model quotas: `gemini-3-pro-high`, `gemini-3-pro-image`, `claude-sonnet-4-5`.
    - Match logic in `tray.rs:188` (strict string matching).
- [ ] **Task 2: Dynamic Menu Repainting**
    - File: `apps/backend/tray/tray.go` (Wails tray logic)
    - Implement `updateTrayMenus()` that redraws the menu items with percentages (e.g., "Gemini High: 85%").
    - Ensure it runs on `config://updated`, `tray://refresh-current` events.

## Acceptance Criteria
- [ ] Tray menu displays 3 distinct lines for different model quotas.
- [ ] Clicking "Refresh" in the tray updates these numbers immediately.
- [ ] Displays "🚫 Forbidden" if the account is 403.
