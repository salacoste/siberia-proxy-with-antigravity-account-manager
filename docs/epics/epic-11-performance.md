# Epic-11: Performance & Scalability

**Status**: Planned

## Goal
Ensure the proxy engine can handle high-throughput scenarios (e.g., scraping, load testing) without crashing or degrading UI performance.

## Stories

### Story-35: High Concurrency Backend
-   **Goal**: Optimize `proxy.Service` to use worker pools or non-blocking IO where possible to support 10k+ concurrent connections.
-   **Acceptance**: `go-wrk` benchmark shows < 50ms latency at 1k concurrent reqs.

### Story-36: Virtualized Traffic List & Memory Safeguards
-   **Goal**: Prevent frontend crashing when thousands of logs arrive.
-   **Tech**: Implement "Windowing" / "Virtualization" in React (TanStack Virtual).
-   **Backend**: Implement log rotation or capping (e.g., keep last 5000 requests only).
