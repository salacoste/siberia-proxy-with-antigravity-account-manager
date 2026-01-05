# Story-29: Traffic Inspector UI (Frontend)

**Epic:** [Epic-09: Advanced Traffic Inspection](../epic-09-advanced-inspection.md)
**Status:** Ready for Dev
**Parent**: Epic-09

## Description
Develop a detailed inspection interface for the Proxy Monitor. When a user clicks on a request row, a detailed view (Sidebar/Sheet) should appear displaying the full Request and Response information, including decrypted headers and bodies.

## Requirements
1.  **Interaction**:
    -   Clicking a row in `/monitor` opens a side panel (Sheet) or Modal.
2.  **Data Display**:
    -   **General**: URL, Method, Status, Duration, Timestamp.
    -   **Request Tab**:
        -   Headers (Key-Value table).
        -   Body (Code block / JSON viewer).
    -   **Response Tab**:
        -   Headers.
        -   Body (Code block / JSON viewer).
3.  **Features**:
    -   **Pretty Print**: JSON bodies should be formatted.
    -   **Copy**: Buttons to copy Body, Headers, or "Copy as cURL".

## Acceptance Criteria
- [x] Clicking a request opens the Inspector.
- [x] Decrypted bodies (from Epic-08) are visible in the Inspector.
- [x] JSON content is syntax highlighted/formatted.
- [x] Large bodies are truncated or handled gracefully (virtualization if needed, or simple cap).
- [x] "Copy as cURL" works.

## Dev Agent Record

### Completion Notes
- Implemented `RequestDetails` using `shadcn-ui` Sheet component.
- Implemented `DataViewer` for clean JSON/Text display.
- Verified build and integration.
- Installed missing shadcn components correctly.

### Story DoD Checklist
- [x] UI matches requirements (Sheet, Tabs, JSON view).
- [x] Validated via build.
- [x] Code follows standards (shadcn).

### QA Fixes (2026-01-05)
- Implemented `generateCurlCommand` and "Copy cURL" button.
- Added "Copy Body" buttons to both Request and Response tabs.
- Fixed lint warning (unused `toast`).

## Dev Notes
-   Use `shadcn-ui` `Sheet` or `Dialog` component.
-   Use `react-syntax-highlighter` or similar for body display? Or just `<pre>`.
-   Ensure `ProxyRequestEvent` type in frontend matches the backend struct updated in Epic-08.

## QA Results
> [!CAUTION]
> **Status: FAIL**
> **Date:** 2026-01-05
> **Reviewer:** Quinn (QA Agent)
> **Gate:** [epic-09.story-29-traffic-inspector.yml](../../qa/gates/epic-09.story-29-traffic-inspector.yml)
>
> **Issues:**
> 1. Missing "Copy as cURL" implementation.
> 2. Missing "Copy Body" buttons.
>
> **Action:** Returned to Dev for fixes.
