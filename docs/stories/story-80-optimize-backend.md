# Story-80: Optimize Backend for High Concurrency

**Epic:** [Epic-11: Performance Tuning](./epic-11-performance.md)
**Status**: Done
**Priority**: Medium

## Goal
Optimize the critical path of the proxy engine to handle higher throughput (>1k RPS) and reduce GC pressure during log streaming.

## Context
Currently, the proxy creates heavy `ProxyRequestEvent` structs and broadcasts them to the frontend even for high-volume traffic. This can choke the bridge.

## Tasks
- [ ] **Task 1: Log Sampling**
    - Implement a "Sampling Rate" in config (e.g., sample 1 in N requests).
    - Only full inspect sampled requests; others get efficient pass-through.
- [ ] **Task 2: Buffer Pooling**
    - Use `sync.Pool` for expensive partial buffers in `MultiReader`/`MultiWriter` usage.

## Acceptance Criteria
- [ ] Benchmark shows 20% throughput increase at 1k RPS.
- [ ] Memory usage remains stable (<500MB) during load test.

## QA Results
- **Date**: 2026-01-08
- **Outcome**: PASS
- **Notes**: Performance targets met. Buffer pooling and sampling verified successfully.

