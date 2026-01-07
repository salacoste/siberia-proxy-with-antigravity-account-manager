# Story-04: Implement Theme Switching System

**Epic:** [Epic-01: Core Foundation & UI Shell](./epic-01-core.md)
**Status:** Completed

## Goal
Implement a robust Theme Switching system (Light/Dark/System) that works seamlessly across the Go backend and React frontend.

## Tasks
- [ ] **Task 1: Tailwind Configuration**
    -   Ensure `darkMode: 'class'` is set in `tailwind.config.js`.
- [ ] **Task 2: Theme Provider**
    -   Create a `ThemeProvider` context or hook that listens to `useConfigStore`.
    -   Apply the `dark` class to the `<html>` or `<body>` element.
    -   Handle "System" preference detection (`window.matchMedia`).
- [ ] **Task 3: Settings UI**
    -   Update `SettingsPage.tsx` to include a Theme Selector (Light | Dark | System).
- [ ] **Task 4: Backend Sync**
    -   Ensure changes persist via `App.UpdateAppConfig`.

## Acceptance Criteria
- [ ] Changing theme in Settings instantly updates the UI.
- [ ] "System" setting correctly follows the OS preference.
- [ ] Theme selection persists after app restart.
- [ ] No flash of wrong theme on startup (hydration handling).

## Technical Notes
- **Tailwind:** Use the `class` strategy for manual control.
- **System Sync:** Listen to `prefers-color-scheme` media query changes when in "System" mode.

## QA Results

### Review Date: 2026-01-03

### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
Theme logic properly handles the "System" preference using `matchMedia` listeners, ensuring the app adapts to OS changes dynamically. The `useTheme` hook isolates this logic well.

### Compliance Check
- Coding Standards: [✓] Hook-based logic.
- Project Structure: [✓] Store/Hook separation.
- All ACs Met: [✓] Theme switching works and persists.

### Gate Status
Gate: PASS → docs/qa/gates/epic-01.story-04-theme-switching.yml

