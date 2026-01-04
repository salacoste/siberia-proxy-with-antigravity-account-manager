# Story-41: Fix Settings Page Crash in Web Mode

**Parent Epic**: Epic-14

## Problem
The Settings Page crashes (blank white screen) when accessed in a standard web browser (e.g., during `wails dev` or PO validation) because the Wails runtime (`window.go`, `window.runtime`) is missing.
The dashboard loads fine, but `useConfigStore` or `SettingsPage` calls Wails backend methods that fail fatally.

## Requirements
1.  **Graceful Degridation**: The app MUST NOT crash if Wails runtime is missing.
2.  **Mock Data**: In Web Mode (determined by missing Wails runtime), `GetAppConfig` should return a default/mock configuration.
3.  **UI Feedback**: The Settings Page should render successfully in a browser, possibly with a "Web Mode / Mock Data" indicator.
4.  **Error Handling**: Wrap Wails calls in try/catch or checks to prevent unhandled promise rejections crashing the React tree.

## Acceptance Criteria
- [x] Navigate to `/settings` in Chrome (http://localhost:7101) DO NOT crash.
- [x] Settings Page renders usage form (even if data is mock).
- [x] Console logs do not show "Uncaught Error".

## Completion Notes
- Implemented `window.go` and `window.runtime` checks in `useConfigStore.ts` and `SettingsPage.tsx`.
- Added Mock Configuration implementation when runtime is missing.
- Added "WEB MODE (Mock Data)" badge to Settings header.
- Fixed infinite render loop in `ThemeToggle.tsx`.

## QA Results (2026-01-04)
- **Status**: PASS
- **Verified By**: QA Agent
- **Notes**: Regression test suite passed in Web Mode (http://localhost:7101). No crashes observed. "WEB MODE" badge verified.

