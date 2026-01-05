# Story-48: Analytics Aggregation Service

**Parent**: Epic-16
**Status**: Completed

## Description
Implement a backend service that aggregates traffic data in real-time to power the Analytics Dashboard. It should process `ProxyRequestEvent`s and maintain counters for key metrics.

## Requirements
1.  **Analytics Engine**:
    -   In-memory aggregation (reset on restart is acceptable for v1).
    -   Thread-safe counters.
2.  **Metrics to Track**:
    -   **Request Rate (RPS)**: Current requests per second.
    -   **Bandwidth**: Total bytes transferred (In/Out) and current throughput.
    -   **Active Connections**: Current open connections gauge.
    -   **Top Domains**: Map of Host -> Count (approximated or limited top-k).
    -   **Response Codes**: Distribution (2xx, 4xx, 5xx).
3.  **Wails API**:
    -   Expose `GetAnalyticsSnapshot()` to frontend.
    -   (Optional) Emit `analytics:update` event every 1s.

## Acceptance Criteria
- [x] `AnalyticsEngine` struct implemented.
- [x] `Track(event)` updates internal counters safely.
- [x] `GetAnalyticsSnapshot()` returns JSON-friendly stats.
- [x] Integrated into `TelemetryManager.worker()`.
- [x] Unit tests for aggregation logic.

## Implementation Notes (2026-01-04)

### Implemented Features
- Real-time aggregation of RPS, Bandwidth, and Status Codes.
- `AnalyticsService` exposed via Wails.
- `scripts/traffic_gen.go` created for validation.

### Challenges & Solutions
1.  **Wails Binding Failure**:
    -   *Issue*: Frontend threw `TypeError: Cannot read properties of undefined (reading 'analytics')`. The `window.go` object was missing the `analytics` namespace.
    -   *Root Cause*: `wails.json` had an incorrect `wailsjs:dir` pointing to `src/wailsjs`, but imports expected `wailsjs/`.
    -   *Fix*: Updated `wails.json` path to `../frontend/wailsjs` and restarted `wails dev`.

2.  **DevServer Injection**:
    -   *Issue*: Even after config fix, `window.go` was undefined on the main dev port.
    -   *Workaround*: Manually injected `<script src="/wails/ipc.js"></script>` and `<script src="/wails/runtime.js"></script>` into `index.html` for development mode.

3.  **Validation Tooling**:
    -   *Issue*: `traffic_gen.go` was targeting port `9090` (old config) instead of `7100` (actual proxy).

## QA Results

### Review Date: 2026-01-04

### Reviewed By: Quinn (Test Architect)

### Gate Status

Gate: PASS → docs/qa/gates/epic-16.story-48-analytics-backend.yml
