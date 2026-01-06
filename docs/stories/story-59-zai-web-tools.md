# Story-59: Z.ai Web Search & Reader Tools

**Epic:** [Epic-19: Z.ai Intelligence](./epic-19-zai-intelligence.md)
**Status**: Completed
**Priority**: High
**Basis**: `docs/gap_analysis_deep_dive.md`

## Goal
Port the web intelligence tools from the reference project's Rust implementation to Go. This includes a privacy-focused web search integration and a robust web reader that converts HTML to clean Markdown for LLM consumption.

## Tasks
- [x] **Task 1: Web Search Tool**
    - File: `apps/backend/zai/tools/search.go`
    - Port `call_web_search_prime`.
    - Inputs: `query`, `recency` (e.g., "past_24h"), `domain_filter`.
- [x] **Task 2: Web Reader Tool**
    - File: `apps/backend/zai/tools/reader.go`
    - Port `call_web_reader`.
    - Use a Go HTML-to-Markdown library (e.g., `go-readability` or similar).
    - Implement URL normalization: Strip `utm_`, `gclid`, `fbclid` parameters (Reference: `normalize_web_reader_url`).
- [x] **Task 3: Register with MCP**
    - Register these two functions into the Tool Registry created in Story-58.

## Acceptance Criteria
- [x] MCP Call to `search_web` returns valid JSON search results.
- [x] MCP Call to `read_web_page` returns clean Markdown content.
- [x] Privacy check: Tracking parameters are removed from URLs before fetching.

## QA Results
**Date**: 2026-01-06
**Reviewer**: Quinn (QA)
**Status**: PASS

### Verification Summary
- **Code Audit**: Validated strict privacy controls in `reader.go` and secure configuration usage in `search.go`.
- **Test Results**: All unit tests passed (Tools & MCP Server).
- **Compliance**: Meets all acceptance criteria.

**Gate Decision**: [PASS](../qa/gates/epic-19.story-59-zai-web-tools.yml) - Ready for Merge.
