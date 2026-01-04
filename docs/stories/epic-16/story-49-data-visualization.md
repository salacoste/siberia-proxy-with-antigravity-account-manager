# Story-49: Data Visualization (Recharts)

## Description
Implement rich data visualization for the Analytics Dashboard using `recharts`. The dashboard should display real-time traffic statistics, including request rates, status code distribution, and top domains.

## Requirements
-   **Library**: Use `recharts` for charting.
-   **Components**:
    -   `TrafficChart`: Area chart showing Requests per Second (RPS) over time.
    -   `StatusChart`: Donut/Pie chart showing distribution of HTTP status codes (2xx, 4xx, 5xx).
    -   `TopDomains`: Bar chart or List showing most accessed domains.
-   **Data Source**: Fetch data from `AnalyticsService.GetStats` via a React Context poller.
-   **Design**: Matches "Saint Vandal" aesthetic (dark mode, sleek, animations).

## Acceptance Criteria
-   [ ] `recharts` is installed.
-   [ ] `AnalyticsContext` polls `GetStats` every 1-2 seconds.
-   [ ] Dashboard displays a live-updating Area Chart for traffic volume.
-   [ ] Dashboard displays a Pie Chart for status codes.
-   [ ] Dashboard displays Top Domains.
-   [ ] UI is responsive and fits the grid layout.
