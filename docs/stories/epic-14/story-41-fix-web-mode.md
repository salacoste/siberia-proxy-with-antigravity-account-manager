# Story-41: Fix Web Mode Crashes

**Epic**: [Epic-14: Release Polish](../epics/epic-14-release-polish.md)
**Status**: Planning

## Goal
Ensure the application runs gracefully in a standard web browser ("Web Mode") where the Wails runtime (`window.runtime`, `window.go`) is not available.

## Problem
Currently, the `SettingsPage` and potentially other pages assume `window.go` exists. Accessing it directly may cause "Cannot read properties of undefined" or blank screens in a normal browser, blocking development/demo flows.

## Detailed Requirements

1.  **Settings Page**:
    -   Guard all Wails calls (`GetVersion`, `CheckCertTrust`, `CheckForUpdates`).
    -   Show "Mock" or "Web Mode" indicators instead of crashing or showing valid data.
    -   Disable/Hide controls that strictly require backend (e.g., "Install Certificate").

2.  **Traffic Monitor**:
    -   Ensure `TrafficProvider` doesn't crash if `window.runtime.EventsOn` is missing.
    -   Provide Mock Data generation if in Web Mode (optional, but helpful).

3.  **Global Error Boundary**:
    -   Wrap the main App router in a React Error Boundary to catch unanticipated crashes and show a friendly UI.

## Acceptance Criteria
- [x] `npm run dev` (Web Mode) loads `http://localhost:5173/settings` without white screen or console errors.
- [x] "Install Certificate" button is disabled/hidden in Web Mode with a tooltip explanation.
- [x] Version displays "Web Mode" or "Mock".

## QA Results (2026-01-05)
**Agent**: Antigravity
**Decision**: **PASS**

### Verification
- **Code Audit**: Validated that `SettingsPage`, `SyncSettings`, `CloudPage`, and `App` guard all `window.go` calls.
- **Safety**: Added Global `ErrorBoundary` to catch leftover crashes.
- **Mocking**: Confirmed `useConfigStore` provides mock data when Wails runtime is missing.cSettings`, and `App` guard all `window.go` calls.
- **Safety**: Added Global `ErrorBoundary` to catch leftover crashes.
- **Mocking**: Confirmed `useConfigStore` provides mock data when Wails runtime is missing.
