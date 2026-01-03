# Story-18: Update UI & Notifications

**Epic:** [Epic-05: Distribution](./epic-05-distribution.md)
**Status:** Draft
**Feature Branch:** `antigravity/feat/distribution`

## Goal
Frontend UI to expose update status.

## Tasks
- [ ] **Task 1: Settings UI**
    - [ ] Add "Check for Updates" button in `SettingsPage`.
    - [ ] Display Current Version (from Wails `GetAppConfig` or runtime).
- [ ] **Task 2: Update Notification**
    - [ ] If update available: Show Toast or Badge.
    - [ ] Click -> Triggers `Download` or `OpenBrowser`.

## Acceptance Criteria
- [ ] User can manually check for updates.
- [ ] UI reflects "Up to Date" or "Update Available".

## QA Results

### Review Date: 2026-01-04

### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
Frontend `SettingsPage.tsx` correctly calls `CheckForUpdates` and handles the response. UI Feedback (Toasts) is appropriate. Wails Model synchronization update was verified.

### Compliance Check
- Coding Standards: [✓] React Hooks / Async handling.
- All ACs Met: [✓] UI Update flow verified.

### Gate Status
Gate: PASS → docs/qa/gates/epic-05.story-18-update-ui.yml

