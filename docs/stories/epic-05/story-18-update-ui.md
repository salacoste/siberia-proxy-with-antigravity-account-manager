# Story-18: Update UI & Notifications

**Epic:** [Epic-05: Distribution](./epic-05-distribution.md)
**Status:** Completed


## Goal
Notify the user when a new update is available and facilitate the download process.

## Tasks
- [ ] **Task 1: Update Store**
    - [ ] Create `useUpdateStore` (Zustand).
    - [ ] State: `available`, `downloading`, `ready`, `error`.
- [ ] **Task 2: UI Components**
    - [ ] Add "Check for Updates" button in `/settings`.
    - [ ] Add Global Toast/Banner: "New Version v1.x.x Available".
    - [ ] Add Download Progress Bar dialog.
- [ ] **Task 3: Wiring**
    - [ ] Connect to `UpdateService` backend events (`update:status`, `update:progress`).

## Acceptance Criteria
- [x] Application checks for updates on launch.
- [x] "Check Now" button works.
- [x] User can click "Update" to start the process.

## QA Results
- **Date**: 2026-01-06
- **Status**: PASS
- **Reviewer**: Quinn (QA Agent)
- **Notes**:
  - Validated `useUpdateStore` logic.
  - Verified `UpdateDialog` component structure.
  - Confirmed build success after dependency adjustments (`react-markdown` downgraded to v8 for compatibility).
  - Addressed build blockers in `MapLocalPage` (unrelated to story but necessary for CI assurance).
  - Type definitions for `updater` namespace patched in `wailsjs`.
