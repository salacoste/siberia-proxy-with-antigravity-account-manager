# Story-36: Implement Virtualized Traffic List & Memory Safeguards

**Epic:** [Epic-11: Performance Tuning](./epic-11-performance.md)
**Status:** Draft

## Goal
Prevent the Frontend from crashing when the traffic log contains 100,000+ items.

## Tasks
- [ ] **Task 1: Virtualization**
    -   Replace standard Table with `tanstack/react-virtual` or similar.
    -   Only render visible rows.
- [ ] **Task 2: Data Cap**
    -   Backend enforces hard limit (e.g., keep last 5000 requests in memory).
    -   Frontend drops old items from state store.

## Acceptance Criteria
- [ ] UI remains responsive with 50,000 logs in state.
- [ ] Memory usage remains stable (flat line) over time.
