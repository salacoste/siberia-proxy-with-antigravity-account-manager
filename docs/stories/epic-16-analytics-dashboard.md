# Epic-16: Analytics Dashboard

**Status**: In Progress
**Goal**: Transform the simple "proxy monitor" into a comprehensive "traffic intelligence dashboard".
**Value**: Users can see traffic patterns, errors, and bandwidth usage at a glance.

## Scope
1.  **Dashboard Screen**: A new home screen (or updated Dashboard tab) with widgets.
2.  **Metrics**:
    -   Requests per Second (RPS) - Line Chart.
    -   Response Codes (2xx, 4xx, 5xx) - Pie/Bar Chart.
    -   Bandwidth (Up/Down) - Area Chart.
    -   Top Domains - List.
3.  **Backend Aggregation**: Efficiently calculate stats from the existing access logs (SQLite).

## Stories
-   **Story-47**: Dashboard Layout & Widget System (Frontend).
-   **Story-48**: Analytics Aggregation Service (Backend).
-   **Story-49**: Data Visualization Implementation (Recharts).
