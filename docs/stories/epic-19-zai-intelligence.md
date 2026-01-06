# Epic-19: Z.ai Intelligence & Agentic Capabilities

**Status**: Planning
**Priority**: High (Tier 2 - Intelligence)
**Basis**: `docs/gap_analysis_deep_dive.md` (Section 1 & 6)

## Executive Summary
Transform the application from a passive proxy into an active "Tool Provider" for AI agents. This Epic implements an internal MCP (Model Context Protocol) server and ports the Z.ai web intelligence tools from the reference project. This enables connected IDEs (Cursor/Windsurf) and the app itself to perform Reference Augmented Generation (RAG) and live web retrieval.

## Key Features
1.  **Internal MCP Server**:
    -   Host a local bridge `http://localhost:port/mcp` answering JSON-RPC 2.0 requests.
    -   Expose `tools/list` and `tools/call` endpoints.
2.  **Z.ai Web Search Tool**:
    -   Port `call_web_search_prime` functionality.
    -   Support domain filtering and recency filtering.
3.  **Z.ai Web Reader Tool**:
    -   Port `call_web_reader` functionality.
    -   Implement advanced URL normalization (stripping tracking parameters/UTM tags) to ensure clean context for LLMs.

## Success Metrics
-   [ ] `curl` request to local MCP endpoint lists available tools.
-   [ ] Agent can successfully perform a web search via the proxy.
-   [ ] Agent can scrape a documentation page to Markdown via the proxy.
