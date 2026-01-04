# Story-35: Optimize Backend for High Concurrency

**Epic**: [Epic-11: Performance](../epics/epic-11-performance.md)
**Status**: Completed

## Goal
Optimize the proxy engine to handle high concurrency loads (1k-10k active connections) efficiently.

## Problem
Currently, we spawn a goroutine per request (standard `net/http`). Under heavy load (e.g. scraping), this might lead to excessive memory usage or GC pressure if we leak context or buffers.
We also lock the `AccessLog` slice globally, which is a contention point.

## Requirements
1.  **Benchmarking**: Establish a baseline using a tool like `go-wrk` or `k6`.
2.  **Lock Contention**: Refactor `AccessLog` and `TrafficMonitor` to use non-blocking structures (channels or buffered writes) instead of blocking Mutexes on the hot path.
3.  **Connection Pooling**: Ensure `elazarl/goproxy` is configured with optimal Transport settings (MaxIdleConns, etc.).

## Acceptance Criteria
- [x] Baseline benchmark recorded.
- [x] Hot path (Request -> Proxy -> Log) optimized (no long locks).
- [x] Benchmark post-optimization shows > 20% throughput improvement or stable memory usage under load.

## QA Results (2026-01-04)
**Agent**: Quinn (QA)
**Decision**: **PASS** (via Gate `epic-11.story-35-high-throughput.yml`)

### Verification
- **Test**: `TestTelemetryChannelSaturation` confirmed non-blocking behavior even with 2000+ pending events.
- **Architecture**: `telemetryChan` buffer (1000) and worker implemented correctly.
- **Stability**: Build successful.

## Technical Notes
-   Beware of `sync.Mutex` in `proxy.Service.Broadcast`. If the frontend is slow, we shouldn't block the proxy. Use non-blocking sends or drop frames.
