# Story-30: Regex & Advanced Filtering Logic

**Epic:** [Epic-09: Advanced Traffic Inspection](../epic-09-advanced-inspection.md)
**Status:** Completed
**Parent:** Epic-09

## Description
Implement advanced filtering capabilities for the Proxy Monitor. Currently, we likely have simple text search. We need to support Regex filtering and potentially specific field filtering (e.g., `method:POST`, `status:500`) to help developers find specific traffic events in high-volume streams.

## Requirements
1.  **Regex Support**: Allow users to toggle "Regex" mode for the search bar.
2.  **Field Filtering**: Supports structured queries (conceptual or explicit UI):
    -   Filter by Method (GET, POST, etc.)
    -   Filter by Status Code (200, 404, 500, etc.)
    -   Filter by URL path pattern.
3.  **UI Updates**:
    -   Update Search Bar to include a "Regex" toggle button.
    -   (Optional) Add "Filter Pill" UI or dropdowns for common filters (Method/Status).
4.  **Backend vs Frontend**:
    -   Decide if filtering happens in Frontend (on the received list) or Backend (before sending over WS).
    -   *Decision*: Frontend filtering is usually sufficient for < 10k items and provides instant feedback. Backend filtering requires changing the WebSocket subscription model. Let's stick to **Frontend Filtering** for now unless performance dictates otherwise.

## Acceptance Criteria
- [ ] Users can type a Regex string into the search bar (e.g., `api\/v1\/(users|products)`) and filter the list.
- [ ] Invalid Regex patterns do not crash the app (show red border or error).
- [ ] Users can filter by HTTP Method (e.g., via a dropdown or special syntax like `method:POST`).
- [ ] Filtering is case-insensitive by default (unless Regex flag says otherwise).
- [ ] Filter applies to both URL and potentially Body content (bonus). *MVP: URL only first.*

## Dev Notes
-   Update `useMonitorStore` to handle filtering logic.
-   Use `new RegExp(pattern, 'i')` for implementation.
-   Consider adding a `FilterBar` component above the table.

## Dev Agent Record

### Completion Notes
- Implemented `filterUtils.ts` with comprehensive unit tests (8/8 passed).
- Created `FilterBar.tsx` using `shadcn` components (`Select`, `Toggle`).
- Integrated `FilterBar` into `MonitorPage.tsx`, replacing the legacy Input.
- Verified compilation and test suite.

### Story DoD Checklist
- [x] UI matches requirements (Regex toggle, Dropdowns).
- [x] Unit tests passed.
- [x] Code follows standards.

## QA Results
> [!TIP]
> **Status: PASS**
> **Date:** 2026-01-05
> **Reviewer:** Quinn (QA Agent)
> **Gate:** [epic-09.story-30-filter-logic.yml](../../qa/gates/epic-09.story-30-filter-logic.yml)
>
> **Action:** Approved for Merge.
