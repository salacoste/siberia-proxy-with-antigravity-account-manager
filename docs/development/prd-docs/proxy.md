# Proxy subsystem (behavioral spec)

This is a PRD-oriented description of the in-app proxy (what it serves, how it routes, and how it uses accounts).

For deeper implementation notes, see:
- `docs/proxy/routing.md`
- `docs/proxy/config.md`
- `docs/proxy/auth.md`
- `docs/proxy/logging.md`
- `docs/zai/provider.md`
- `docs/zai/mcp.md`

## 1) Goals

- Provide a **single local endpoint** that multiple tools can use concurrently.
- Support multiple protocols by routing on **HTTP path**, not by client identity.
- Use a local **Google OAuth account pool** to distribute load and survive quota limits.
- Offer optional external upstream (z.ai) for Anthropic-compatible requests.

## 2) Network surface (endpoints)

### Health
- `GET /healthz` → `{"status":"ok"}`

### Anthropic-compatible (“Claude protocol”)
- `POST /v1/messages`
- `POST /v1/messages/count_tokens`
- `GET /v1/models/claude` (stub list)

### Gemini-native
- `GET /v1beta/models`
- `GET /v1beta/models/:model`
- `POST /v1beta/models/:model` (generateContent / streamGenerateContent via `:method` suffix)
- `POST /v1beta/models/:model/countTokens` (currently returns stub tokens)

### OpenAI-compatible
- `POST /v1/chat/completions`
- `POST /v1/completions`
- `POST /v1/responses` (alias)
- `POST /v1/images/generations`
- `POST /v1/images/edits`

### MCP
- `ANY /mcp/web_search_prime/mcp`
- `ANY /mcp/web_reader/mcp`
- `ANY /mcp/zread/mcp`
- `ANY /mcp/zai-mcp-server/mcp` (local Vision MCP)

## 3) Account pool behavior (Google-backed flows)

### Selection & stability

When the proxy needs a Google account, it requests a token from `TokenManager::get_token(quota_group, force_rotate)`.

Key behaviors:
- short “stickiness” window (60s) to reduce thrash across accounts for bursts,
- rotates to another account on retry attempts,
- refreshes expiring tokens automatically,
- disables accounts on `invalid_grant` to prevent retry storms.

Implementation:
- `src-tauri/src/proxy/token_manager.rs`

### Quota groups

Requests carry a request type / quota group (e.g. `gemini`, `claude`, `image_gen`), used to decide:
- which account to pick/lock (some groups may bypass the lock),
- which upstream behavior is required.

Implementation:
- `src-tauri/src/proxy/mappers/common_utils.rs` (request type inference)
- Handlers pass `request_type` into `get_token(...)`

## 4) Provider routing (z.ai vs Google pool)

Only Anthropic-compatible routes can be redirected to z.ai.

Decision inputs:
- `proxy.zai.enabled`
- `proxy.zai.api_key` present
- `proxy.zai.dispatch_mode`:
  - `off`: never use z.ai
  - `exclusive`: always use z.ai for `/v1/messages`
  - `pooled`: z.ai participates as one slot with Google accounts
  - `fallback`: use z.ai only when Google pool has 0 accounts

Implementation:
- `src-tauri/src/proxy/handlers/claude.rs`
- `src-tauri/src/proxy/providers/zai_anthropic.rs`

## 5) Authorization model

Authorization is global (applies to all proxy routes) and controlled by `proxy.auth_mode`:
- `off`: open access
- `strict`: require auth for all routes
- `all_except_health`: require auth for all routes except `/healthz`
- `auto`: derived default (based on `allow_lan_access`)

Auth header:
- `Authorization: Bearer <proxy.api_key>`

Implementation:
- `src-tauri/src/proxy/middleware/auth.rs`
- `src-tauri/src/proxy/security.rs`

## 6) Observability

Two different toggles:

### Safe access log
- `proxy.access_log_enabled`
- Logs only: method, path (no query), status, latency, upstream label.

### Proxy monitor (payload capture)
- `proxy.enable_logging`
- Stores request/response payload samples + token usage into `proxy_logs.db`.

## 7) Shutdown semantics (proxy + MCP)

MCP endpoints are served by the same in-process Axum server instance as all other routes.
When the proxy service is stopped:
- listener stops accepting new connections,
- `/mcp/*` endpoints become unavailable immediately,
- in-flight requests may complete (normal async behavior).

Implementation:
- stop command: `src-tauri/src/commands/proxy.rs` (`stop_proxy_service`)
- server shutdown: `src-tauri/src/proxy/server.rs` (`AxumServer::stop`)

