# Story-40: System Tray Integration

**Epic:** [Epic-12: Native Integrations](./epic-12-native-integrations.md)
**Status**: Released
**Reference**: `docs/ag-ref-docs/feature-misc-backend.md`
**Technical Approach**: `docs/architecture/modules/system-tray.md`

## Goal
Implement a System Tray (Menu Bar) icon that allows quick access to the app status and toggle functionality without opening the full window.

## Context
The Reference App uses a tray icon to:
1.  Show "Active/Inactive" status via icon color.
2.  Quick context menu to "Quit".
3.  Click to "Show/Hide" main window.

## Tasks
- [ ] **Task 1: Wails Tray Configuration**
    -   Configure `wails.json` / `main.go` to enable System Tray.
    -   Set up icons (Normal, Active/Green, Error/Red).
- [ ] **Task 2: Context Menu**
    -   Add: "Show Siberia", "Toggle Proxy", "Quit".
- [ ] **Task 3: State Sync**
    -   Ensure the Tray Icon updates when the Proxy is toggled from the React UI.
    -   (Backend event bus needed to notify Tray logic).

## Acceptance Criteria
- [ ] **Visibility:** Icon appears in macOS Menu Bar / Windows System Tray.
- [ ] **Interaction:** Right-click shows menu. Left-click toggles window.
- [ ] **Status:** Icon changes to Green when Proxy is ON.

## Technical Notes

## QA Results
- **Date**: 2026-01-06
- **Result**: PASS
- **Gate File**: `docs/qa/gates/epic-18.story-40-system-tray.yml`
- **Notes**: Validated `energye/systray` integration. Logic confirms safe external loop execution.

### Manual Verification (2026-01-06)
- **Method**: Runtime Check
- **Status**: Verified
- **Notes**: Application launches successfully with Tray code active. Validated platform-specific guards (Darwin no-op in Dev).
