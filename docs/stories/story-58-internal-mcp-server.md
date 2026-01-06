# Story-58: Internal MCP Server

**Epic:** [Epic-19: Z.ai Intelligence](./epic-19-zai-intelligence.md)
**Status**: Completed
**Priority**: High
**Basis**: `docs/gap_analysis_deep_dive.md`

## Goal
Implement a lightweight JSON-RPC 2.0 server within the Go backend that adheres to the Model Context Protocol (MCP). This allows the proxy to expose its own tools (and bridged tools) to local MCP clients like Cursor, Windsurf, or the proxy's own agentic loop.

## Tasks
- [ ] **Task 1: MCP Protocol Handler**
    - File: `apps/backend/mcp/server.go`
    - Implement `ServeHTTP` handler for JSON-RPC.
    - Support standard methods: `tools/list`, `tools/call`, `resources/list`.
- [ ] **Task 1.5: Configuration Update**
    - File: `apps/backend/config/config.go`
    - Add `ZaiMcpConfig` struct (Enabled, WebSearchEnabled, Port).
    - Ensure these load from `settings.json`.
- [ ] **Task 2: Tool Registry**
    - File: `apps/backend/mcp/registry.go`
    - Create a registry map to dynamically register Go functions as tools.
- [ ] **Task 3: Server Exposure**
    - Expose on a dedicated port or sub-path (e.g., `http://localhost:6200/mcp`).

## Acceptance Criteria
- [ ] `curl -X POST ... {"method": "tools/list"}` returns a valid JSON-RPC response.
- [ ] External MCP Client (e.g., Inspector) can connect to the local server.

## QA Results
- **Date**: 2026-01-06
- **Reviewer**: Quinn (QA)
- **Status**: **PASS**
- **Notes**:
    - JSON-RPC 2.0 implementation is verified.
    - Registry is thread-safe and functional.
    - Configuration wiring is correct.
    - Ready for tool implementation.

