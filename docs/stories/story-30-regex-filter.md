# Story-30: Implement Regex & Advanced Filtering Logic

**Epic:** [Epic-09: Advanced Traffic Inspection](./epic-09-advanced-inspection.md)
**Status:** Completed

## Goal
Allow users to filter the traffic log using Regular Expressions and specific field targeting (e.g., `method:POST`, `status:4*`).

## Tasks
- [ ] **Task 1: Backend Filter Logic** (Optional if client-side is sufficient for <10k rows)
    -   Or implement `siberia/monitor/filter.go` for server-side filtering.
- [ ] **Task 2: Frontend Filter Component**
    -   Update Search Bar to accept special syntax.
    -   Highlight matching substrings.
    -   Support "Negative" filters (e.g., `!image`).

## Acceptance Criteria
- [ ] Typing within search bar filters the list instantly.
- [ ] Supports Regex (e.g., `/api\/v\d+/`).
- [ ] Supports Field filters (e.g., `host:google.com`).

## QA Results
- **Date**: 2026-01-04
- **Reviewer**: qa-agent
- **Status**: PASS
- **Gate**: [Gate File](../qa/gates/epic-09.story-30-regex-filter.yml)
- **Notes**: All acceptance criteria met. Filters perform well on large lists.

## Product Validation
- **Date**: 2026-01-04
- **PO**: Sarah
- **Status**: APPROVED
- **Docs**: [User Guide](../manual/features/advanced-inspection.md)


