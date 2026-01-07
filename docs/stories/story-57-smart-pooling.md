# Story-57: Smart Token Pooling & Tier Sorting

**Epic:** [Epic-20: Advanced Traffic Scheduling](./epic-20-advanced-scheduling.md)
**Status**: Ready for Review
**Priority**: Medium
**Basis**: `docs/gap_analysis_deep_dive.md`

## Goal
Optimize account usage by prioritizing higher-tier accounts (ULTRA/PRO) which typically have faster refill rates or higher limits. Additionally, implement configurable scheduling modes to balance between performance (speed) and cache preservation.

## Tasks
- [ ] **Task 1: Account Tier Field**
    - File: `apps/backend/models/account.go`
    - Add `Tier` enum: `FREE`, `PRO`, `ULTRA`.
    - Update `Account` struct to include `Tier`.
- [ ] **Task 2: Prioritized Selection**
    - File: `apps/backend/proxy/pool/selector.go`
    - Update `GetNextAccount`: Sort available accounts by `Tier desc`, then `Health`.
- [ ] **Task 3: Scheduling Modes**
    - Implement `SchedulingMode` config:
        - `PerformanceFirst`: If sticky account is busy, rotate immediately (Current behavior).
        - `CacheFirst`: If sticky account is busy (429), **wait** up to N seconds before rotating (New behavior).

## Acceptance Criteria
- [ ] Pool with Mixed Tiers -> Requests are assigned to ULTRA accounts first until exhausted.
- [ ] `CacheFirst` mode -> Proxy blocks/waits on a 429 instead of immediately switching account (up to timeout).

## Dev Agent Record
### File List
- `apps/backend/siberia/db/models.go`
- `apps/backend/siberia/config/config.go`
- `apps/backend/siberia/accounts/service.go`
- `apps/backend/siberia/proxy/upstream/gemini.go`
- `apps/backend/siberia/proxy/upstream/gemini_test.go`

## QA Results
- **Status**: PASS
- **Date**: 2026-01-07
- **Gate**: [epic-20.story-57-smart-pooling.yml](../qa/gates/epic-20.story-57-smart-pooling.yml)
- **Notes**: Unit tests verify the retry logic and structural changes are correct.


