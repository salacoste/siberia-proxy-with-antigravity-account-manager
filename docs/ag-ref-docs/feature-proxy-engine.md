# Feature Analysis: Proxy Engine

**Source:** `src-tauri/src/proxy/`
**Target Upstream:** `cloudcode-pa.googleapis.com` (Google Cloud Code "Gemini" Internal API)

## 1. Product Logic
The Proxy Engine is the "Siberia" component. It acts as a local bridge that tricks 3rd-party AI clients (like Cursor, VS Code, or general API clients) into thinking they are talking to OpenAI or Claude, while secretly routing the traffic to Google's internal Gemini models using the user's allocated quotas.

**Key Capabilities:**
*   **Protocol Translation:** Accepts OpenAI/Claude API payloads -> Converts to Google's `GenerateContent` format.
*   **Model Routing:** Maps `gpt-4`, `claude-3-opus` etc. to internal Gemini models (e.g., `gemini-ultimate`).
*   **Resilience:** Automatically balances load across multiple Google accounts and handles rate limits.

## 2. Technical Implementation

### A. HTTP Server (`server.rs`)
*   **Framework:** `Axum` (Rust).
*   **Routes:**
    *   `/v1/chat/completions` (OpenAI Chat)
    *   `/v1/images/generations` (OpenAI Image)
    *   `/v1/messages` (Claude Chat)
    *   `/v1beta/models/*` (Native Gemini Pass-through)
    *   `/mcp/*` (MCP Reverse Proxy)

### B. Protocol Translation (`handlers/openai.rs`, `handlers/claude.rs`)
*   **OpenAI Handler:**
    *   Deserializes `OpenAIRequest`.
    *   Resolves Model Mapping (e.g., `gpt-4` -> `models/gemini-1.5-pro-latest`).
    *   Transforms Messages + Tools into Gemini `Content` and `FunctionCall` objects.
    *   **Streaming:** Converts Gemini's SSE stream event format to OpenAI's `data: [DONE]` format on the fly.
*   **Upstream Logic:**
    *   Method: `generateContent` or `streamGenerateContent` (if streaming).
    *   Query Param: `?alt=sse` for streaming.

### C. Upstream Client (`upstream/client.rs`)
*   **Endpoints:**
    *   **Primary:** `https://cloudcode-pa.googleapis.com/v1internal`
    *   **Fallback:** `https://daily-cloudcode-pa.sandbox.googleapis.com/v1internal`
*   **Resilience Strategy:**
    *   **Pools:** Uses `reqwest::Client` with connection pooling (16 idle per host).
    *   **Automatic Fallback:** If Prod returns 429/5xx, seamlessly retries on Daily endpoint.
    *   **Retry Logic:**
        *   **Strategies:** Fixed Delay, Linear Backoff, and Exponential Backoff based on error type (429 vs 503 vs 500).
        *   **Retry-After:** Honors standard header and Google's custom `RetryInfo`.
        *   **Account Rotation:** If an account hits 429 (Quota Exhausted) or 401/403 (Auth Fail), it grabs a fresh token from the `TokenManager` (different account) and retries.
    *   **Connection Pooling:** High-performance pooling (16 idle/host, 60s TCP keepalive) to minimize latency.

### D. Middleware
*   **Auth:** Validates the local bearer token (optional/configurable).
*   **RateLimit:** Local rate limiting to prevent abuse.
*   **Observability:** Access Logs and Response Attribution (injects headers to tell the client which account/model was used).

## 3. Protocol Specifics
### A. OpenAI (`handlers/openai.rs`)
*   **Streaming:** Converts Gemini's SSE stream to OpenAI's `data: [DONE]` format.
*   **Tool Calls:** Maps OpenAI tools to Gemini `FunctionCall`. **Critical:** Shell command arguments are wrapped in arrays `[cmd]` to satisfy Gemini's strict validation (fixes 400 errors).
*   **Image Generation:** Maps DALL-E 3 parameters (size, quality, style) to Gemini Imagen prompts (aspect ratio 16:9, 9:16 etc.).
*   **Codex:** Supports legacy `/v1/completions` by converting prompts -> chat messages, enabling Copilot support.

### B. Claude (`handlers/claude.rs`)
*   **Thinking Blocks:**
    *   **Validation:** Strictly filters "Thinking" blocks with invalid/missing signatures to prevent upstream 400 errors (Vertex AI strictness).
    *   **Fallback:** If a "Thinking" block is invalid but contains text, it is converted to a standard Text block to preserve context.
*   **Resilience:**
    *   **Smart Background Detection:** Detects non-user-interactive tasks (VS Code agentic loops) and auto-downgrades to Flash/Lite models to save quota.
    *   **Retry:** Implements specific backoff strategies for "Thinking" signature errors (Fixed 200ms).
*   **Z.ai Integration (Advanced):**
    *   **Modes:** Support for `Exclusive` (all traffic), `Pooled` (load balanced with Google), or `Fallback` (Google first, then Z.ai).
    *   **Vision & Web Tools:** Implements native Z.ai Vision (Screenshots analysis) and Web Search/Reader tools by orchestrating calls to `api.z.ai`.
    *   **MCP Forwarding:** Acts as a reverse proxy for MCP protocols, forwarding requests to Z.ai's MCP servers (e.g., `web_search_prime`, `web_reader`), handling API key normalization and session management.

## 4. Protocol Translation & Mapping Details

### A. Shared logic & Grounding
*   **Grounding Strategy:** The engine inspects tool definitions and model suffixes (e.g., `-online`) to decide whether to use the "Agent" mode (standard) or "Web Search" mode (injected Google Search tool).
*   **Image Config:** Parses model suffixes like `-hd`, `-4k`, `-16x9` to configure Gemini Imagen parameters dynamically.

### B. OpenAI Mapper Specifics
*   **System Instructions:** Extracts messages with role `system` and converts them to Gemini's `systemInstruction` field.
*   **Tool Wrappers:**
    *   **Shell:** Automatically maps `local_shell_call` to `shell`.
    *   **Cleaning:** Removes unsupported fields from JSON Schemas (e.g., `strict`, `format`) and enforces uppercase Types (e.g., `STRING`) for Gemini compatibility.
*   **SSOP (Shell Output Parsing):**
    *   **Problem:** Some models output shell commands as plain text instead of proper tool calls.
    *   **Solution:** The streaming logic (in `streaming.rs`) scans text chunks for JSON-like patterns containing "command": "shell". If found, it synthetically injects a `local_shell_call` event into the SSE stream, allowing the client to execute the command. This is a critical resilience feature for models like Codex.

### C. Claude Mapper Specifics
*   **Thinking Mode:**
    *   **Smart Downgrade:** If the conversation history contains Tool Use but no Thinking blocks, the engine automatically disables Thinking for the next request to prevent API errors ("Assistant message must start with thinking").
    *   **Vertex Fix:** Explicitly sets `thought: true` for thinking blocks to satisfy Vertex AI strictness.
*   **Cache Control:** Proactively removes `cache_control` fields from all messages to prevent "Extra inputs" errors from upstream.
*   **Tool Priority:** If both local tools (e.g., MCP) and Google Search are requested, the engine prioritizes local tools and drops Google Search, as Gemini v1internal does not support mixed tool types in a single request.

## 5. Security & Privacy
**Source:** `proxy/privacy.rs`

The engine enforces strict PII sanitization before logging or displaying data:
*   **Email Masking:** Transforms `user.name@example.com` to `u***@example.com` to protect user identity in logs.
*   **ID Anonymization:** Truncates long IDs (e.g., `lat_123456789`) to `lat_...6789`.
*   **Stable Hashing:** Uses SHA-256 to generate consistent correlation IDs from user attributes without storing the raw attribute.

## 6. Internal Identity Handling
**Source:** `proxy/project_resolver.rs`

To authenticate with the internal `cloudcode-pa` API, the system requires a valid `cloudaicompanionProject`.
*   **Resolution Strategy:**
    1.  Call `https://cloudcode-pa.googleapis.com/v1internal:loadCodeAssist` with the user's `access_token`.
    2.  Extract `cloudaicompanionProject` from the JSON response.
*   **Fallback Logic:**
    *   If the account is ineligible (no Project ID returned), the system generates a **Mock Project ID** (e.g., `swift-core-a1b2c`).
    *   This allows the proxy to function in a "Shadow Mode" or "Tier 0" capacity depending on upstream validation rules.


## 7. Performance Optimizations (v1.2.0)
**Source:** `proxy/pool.go`, `proxy/service.go`

To handle high-concurrency enterprise workloads (10k+ RPS), several optimizations were added:
*   **Buffer Pooling:** Uses `sync.Pool` for byte slices to reduce Garbage Collection pressure during body copying/reading.
*   **Log Sampling:** Configurable `AccessLogSampleRate` (default 100) to log only 1% of requests under heavy load, preventing I/O bottlenecks while maintaining observability signals.
*   **Zero-Copy Body Peeking:** Optimized `peekAndRestore` logic to inspect bodies for PII/Logging without full memory allocation when possible.
