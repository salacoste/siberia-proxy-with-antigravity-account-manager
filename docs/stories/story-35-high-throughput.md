# Story-35: Optimize Backend for High Concurrency

**Epic:** [Epic-11: Performance Tuning](./epic-11-performance.md)
**Status:** Draft

## Goal
Ensure the proxy adds minimal latency (<5ms) and can handle stress tests of 10k RPS.

## Tasks
- [ ] **Task 1: Benchmark Suite**
    -   Create `wrk` or `k6` test suite.
    -   Establish baseline.
- [ ] **Task 2: Bottleneck Removal**
    -   Profile CPU/Memory (pprof).
    -   Optimize Logger (switch to ring buffer + batch flush).
    -   Remove reflection usage in `goproxy` if any.

## Acceptance Criteria
- [ ] P99 Latency < 10ms under 1k RPS load.
- [ ] No goroutine leaks during extended runs.
