# Story-18: Update UI & Notifications

**Epic:** [Epic-05: Distribution](./epic-05-distribution.md)
**Status:** Draft

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
- [ ] Application checks for updates on launch.
- [ ] "Check Now" button works.
- [ ] User can click "Update" to start the process.
