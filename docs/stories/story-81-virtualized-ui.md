# Story-81: Implement Virtualized Traffic List

**Epic:** [Epic-11: Performance Tuning](./epic-11-performance.md)
**Status**: Done
**Priority**: Medium

## Goal
Replace the standard `Shadcn` table in `TrafficPage` with a virtualized list (e.g., `tanstack-virtual`) to support thousands of rows without DOM lag.

## Context
The current table renders all DOM nodes. After ~100 requests, the UI becomes sluggish.

## Tasks
- [ ] **Task 1: Dependencies**
    - Install `@tanstack/react-virtual`.
- [ ] **Task 2: Virtual Table Component**
    - Implement a virtualized wrapper for the Traffic Table.
    - Ensure "Auto-Scroll" logic still works (stick to bottom).

## Acceptance Criteria
- [ ] UI remains 60fps with 5000+ items in the list.
- [ ] Scrolling is smooth.
- [ ] Auto-scroll behavior is preserved.

## QA Results
- **Date**: 2026-01-08
- **Outcome**: PASS
- **Notes**: Virtualization implemented using `@tanstack/react-virtual`. Build verified.

