# Deep Gap Analysis: Reference Source Code Audit
**Date**: 2026-01-06
**Method**: manual inspection of `src-tauri` Rust files vs `apps/backend` Go files.

## 1. Z.ai Intelligence Layer (Missing Epic)
**Source**: `src-tauri/src/proxy/zai_web_tools.rs`
**Status**: 🔴 Completely Missing in v1.2

The reference app implements a full **MCP (Model Context Protocol)** bridge for Z.ai:
-   **Web Search Prime**: `call_web_search_prime`
    -   Supports `search_domain_filter`, `search_recency_filter`.
-   **Web Reader**: `call_web_reader`
    -   Scrapes URLs to Markdown.
    -   **Advanced Feature**: `normalize_web_reader_url` with `StripTrackingQuery` (removes `utm_`, `gclid`).
-   **Impact**: Our current implementation cannot browse the web or perform RAG-like activities via Z.ai Tools, restricting our "Agentic" capabilities.

## 2. Session Intelligence (The "Sticky" Brain)
**Source**: `src-tauri/src/proxy/session_manager.rs` & `token_manager.rs`
**Status**: 🔴 Missing

The reference app manages upstream context caching much more aggressively:
-   **Session Fingerprinting**:
    -   Generates a stable `sid` by hashing `model_name + first_user_message`.
    -   Allows "Implicit Sticky Sessions" even for stateless clients (like `curl`).
-   **Scheduling Modes**:
    -   `CacheFirst`: If the sticky account is rate-limited, it **waits** (up to N seconds) rather than rotating, to preserve the KV Cache.
    -   `PerformanceFirst`: Rotates immediately on 429.
-   **Tier Sort**: Explicitly sorts accounts `ULTRA > PRO > FREE` to burn fast-resetting quotas first.

## 3. Process Management Robustness
**Source**: `src-tauri/src/modules/process.rs` (1000+ lines)
**Status**: 🟡 Partial / Simplified in v1.2

-   **Family Killing**: On Linux, it kills the entire process tree (`get_self_family_pids`) to avoid zombie processes.
-   **Heuristic Detection**: Distinguishes "Main" process from "Helper/Renderer" processes (`--type=` arg checking).
-   **Clean Shutdown**: Tries `SIGTERM` -> Waits -> `SIGKILL`. Our Go implementation currently just calls `process.Kill()`.

## 4. Privacy & Utils
**Source**: `src-tauri/src/proxy/privacy.rs`
**Status**: 🟡 Partial
-   Standardized `mask_email` (e.g., `exa***@gm***.com`) and `anonymize_id`. We have this scattered in line code.

---

## Strategic Recommendations

### Create **Epic-19: Z.ai Intelligence**
-   Port `zai_web_tools.rs` to a new Go module `siberia/zai/tools`.
-   Expose these as MCP Tools to the client.

### Create **Epic-20: Advanced Traffic Scheduling**
-   Implement `SessionManager.Fingerprint()`.
-   Refactor `TokenManager` to support `SchedulingMode` (Sticky vs Rotate).
-   Add Tier-based sorting.

### Update **Epic-06: Process Management**
-   Backport the Linux Family Kill logic to `siberia/modules/process`.

## 5. Claude Protocol Hardening (Stability Gap)
**Source**: `src-tauri/src/proxy/handlers/claude.rs`
**Status**: 🔴 Missing / Critical
-   **Thinking Block Sanitization**: The ref app aggressively filters/sanitizes "Thinking" blocks.
    -   Validates signatures (min length 10).
    -   Converts invalid blocks to text to prevent Vertex AI 400 errors.
    -   Removes trailing unsigned thinking blocks.
    -   **Impact**: without this, Claude 3.7 requests via Vertex AI may fail randomly.
-   **Retry Strategy**: Implements specific backoff for "Invalid signature" (400) errors.

## 6. Internal MCP Server (Feature Gap)
**Source**: `src-tauri/src/proxy/handlers/mcp.rs`
**Status**: 🔴 Missing
-   The Reference App acts as a local **MCP Server** (JSON-RPC over HTTP).
-   It exposes `tools/list` and `tools/call` endpoints.
-   This allows the IDE (Cursor/Windsurf) to connect to Siberia not just as an OpenAI proxy, but as a Tool Provider (for Web Search, etc.).

## 7. Hybrid Dispatching
**Source**: `src-tauri/src/proxy/handlers/claude.rs`
**Status**: 🟡 Simplified in v1.2
-   Ref app supports `Pooled` dispatch mode: treating Z.ai as just "Account N+1" in the load balancing pool.
-   v1.2 only supports explicit Switch (Exclusive/Fallback).

## 8. Deep State Token Injection (Onboarding Gap)
**Source**: `src-tauri/src/modules/migration.rs`
**Status**: 🔴 Missing / High Value
-   **Mechanism**: The reference app connects to local SQLite databases of **IntelliJ / Android Studio** (`ItemTable`).
-   **Protobuf Logic**: It decodes `jetskiStateSync.agentManagerInitState` and extracts the Refresh Token from nested fields (Field 6 -> Field 3).
-   **Impact**: Allows users to "Import from Android Studio" instantly without manual login.

## 9. Intelligent Rate Limit Parsing
**Source**: `src-tauri/src/proxy/rate_limit.rs`
**Status**: 🟡 Simplified in v1.2
-   **Body Scraping**: The ref app uses Regex to parse wait times from error bodies when headers are missing (e.g., "Try again in 2m 30s", "quotaResetDelay": "42s").
-   **Error Classification**: Distinguishes between `QuotaExhausted` (force 1h wait) and `RateLimitExceeded` (30s wait).
-   **Gap**: v1.2 relies primarily on standard `Retry-After` headers.

## 10. Z.ai Vision & Multimodal (Missing Capability)
**Source**: `src-tauri/src/proxy/zai_vision_tools.rs`
**Status**: 🔴 Completely Missing
-   Reference app supports full multimodal interactions:
    -   `ui_to_artifact`: Screenshot to Code.
    -   `diagnose_error_screenshot`: Visual debugging.
-   Uses `glm-4.6v` with specialized prompt templates.

## 11. IDE State Write-Back (Injector)
**Source**: `src-tauri/src/modules/db.rs` & `account.rs`
**Status**: 🔴 Missing
-   Not only does it *read* from the IDE (Story-60), it *writes back* to it.
-   It terminates the IDE process, injects the new token via Protobuf manipulation, and restarts the IDE.
-   This enables the "Switch Account" feature to propagate to the external environment.

## 12. Legacy Codex Support
**Source**: `src-tauri/src/proxy/handlers/openai.rs`
**Status**: 🟡 Missing
-   Explicit handling for `/v1/completions` with "Codex" style bodies (`input` + `instructions`).
-   Maps `local_shell_call` and `web_search_call` from legacy formats to modern Tool Calls.

## 13. Advanced Tray Telemetry
**Source**: `src-tauri/src/modules/tray.rs`
**Status**: 🟡 Partial / Simplified
-   Ref app displays **3 specific model quotas** directly in the tray menu:
    -   `Gemini High: XX%`
    -   `Gemini Image: XX%`
    -   `Claude 4.5: XX%`
-   Updates dynamically on refresh events.

## 14. Response Attribution & Access Logs
**Source**: `src-tauri/src/proxy/middleware/attribution_headers.rs`
**Status**: 🔴 Missing
-   Injects `x-antigravity-provider`, `x-antigravity-model`, `x-antigravity-account` into responses.
-   Critical for verifying "Sticky Sessions" (to see if you actually hit the same account).

## 15. Native Authentication (OAuth Loopback)
**Source**: `src-tauri/src/modules/oauth_server.rs`
**Status**: 🔴 Missing
-   Ref app uses a local loopback server (port 0) for "Authorize" flow.
-   Current app supports only Manual Token Entry.

## 16. Native Gemini Protocol & Z.ai Provider
**Source**: `src-tauri/src/proxy/handlers/gemini.rs` & `providers/zai_anthropic.rs`
**Status**: 🔴 Missing
-   Ref app natively handles `POST /v1beta/models/...:generateContent`.
-   Ref app has specific Z.ai provider for "Hybrid Dispatching".

## 17. Proxy Robustness & Sanitization
**Source**: `src-tauri/src/proxy/mappers/gemini/wrapper.rs`
**Status**: 🟡 Missing Details
-   **Deep Cleaning**: Removes `[undefined]` strings (Cherry Studio bug).
-   **Schema Strictness**: Enforces `type: "OBJECT"` (uppercase) and removes unsupported JSON Schema fields (`format`, `strict`) recursively.



