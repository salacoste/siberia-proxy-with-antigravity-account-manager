# Screen: API Proxy

Route: `/api-proxy`

Source:
- UI: `src/pages/ApiProxy.tsx`
- Proxy implementation (backend):
  - `src-tauri/src/commands/proxy.rs`
  - `src-tauri/src/proxy/server.rs`
  - handlers: `src-tauri/src/proxy/handlers/*`

## Purpose

Configure and operate the local proxy server:
- Start/stop the proxy.
- Configure binding/timeout/auth.
- Provide user-friendly model mapping controls.
- Provide “how to call the proxy” code snippets for different protocols.
- Configure optional external providers (z.ai) and MCP servers.
- Configure safe access logging.

## Data sources

- App config is loaded/saved as a single object:
  - `load_config`, `save_config`
- Runtime proxy status:
  - `get_proxy_status`
- Proxy lifecycle:
  - `start_proxy_service`, `stop_proxy_service`
- Hot updates:
  - `update_model_mapping`
- z.ai model discovery:
  - `fetch_zai_models`
- API key generation:
  - `generate_api_key`

## Sections & interactions

### 1) Proxy status + primary controls

UI shows:
- Running status indicator + active Google accounts count.
- Buttons:
  - “Open Monitor” (visible only when running)
  - “Start/Stop”

Behavior:
- Start: calls `start_proxy_service({ config: appConfig.proxy })`.
- Stop: calls `stop_proxy_service()`.
- After either action, UI refreshes status via `get_proxy_status`.

Controls table:

| Control | Intent | Backend / state changes |
|---|---|---|
| Open Monitor | Inspect proxy traffic | Navigate to `/monitor` |
| Start | Start server | `start_proxy_service(proxy)` |
| Stop | Stop server | `stop_proxy_service()` |

### 2) Core proxy configuration

All config edits on this page are “live persisted”:
- The page performs optimistic UI updates and then calls `save_config` after each change.
- On save failure, it shows an alert and reverts to the previous config snapshot (best-effort).

#### Port
- Input: `proxy.port`
- Disabled while running.
- Expected range: `8000..65535`.

#### Request timeout
- Input: `proxy.request_timeout`
- Disabled while running.
- Clamped in UI to `30..600` seconds.

#### Auto start
- Toggle: `proxy.auto_start`
- Used by backend startup to auto-start the proxy.

#### Allow LAN access
- Toggle: `proxy.allow_lan_access`
- Disabled while running.
- Changes bind address from loopback-only to `0.0.0.0`.
- Security note: when enabled, auth defaults become stricter when `auth_mode=auto`.

#### Authorization

UI includes:
- “Enabled” switch (maps to `proxy.auth_mode != off`)
- Mode selector: `off | strict | all_except_health | auto`

Behavior:
- Enabled switch sets mode to `all_except_health` when turned on; sets `off` when turned off.
- “Effective mode” is computed on the UI side for user clarity.

Request contract (when auth required):
- Clients must send `Authorization: Bearer <proxy.api_key>`

References:
- `docs/proxy/auth.md`

#### Proxy API Key

UI includes:
- Read-only field displaying `proxy.api_key`
- Buttons:
  - Regenerate (calls `generate_api_key` and persists via `save_config`)
  - Copy to clipboard

Note:
- Regeneration should require confirmation (UI confirm dialog).

#### Access logging (safe)

UI toggle:
- `proxy.access_log_enabled`

Behavior:
- Logs method/path/status/latency only (no headers, query, or bodies).
- Exists to debug routing without leaking secrets.

Reference:
- `docs/proxy/logging.md`

Controls table (core config):

| Control | Config key(s) | Backend side effects |
|---|---|---|
| Port input | `proxy.port` | Persisted; requires restart to take effect |
| Request timeout input | `proxy.request_timeout` | Persisted; requires restart to take effect |
| Auto start toggle | `proxy.auto_start` | Persisted; used on app start |
| Allow LAN access toggle | `proxy.allow_lan_access` | Persisted; requires restart to take effect |
| Auth enabled toggle | `proxy.auth_mode` (off ↔ all_except_health) | Persisted; hot-applied if proxy supports security hot update |
| Auth mode select | `proxy.auth_mode` | Persisted; hot-applied if supported |
| Regenerate API key | `proxy.api_key` | Calls `generate_api_key`, then `save_config` |
| Copy API key | — | Clipboard only |
| Access log toggle | `proxy.access_log_enabled` | Hot-updates access log middleware |

### 3) External providers: z.ai (Anthropic-compatible)

Rendered as a collapsible card with enable toggle.

Config keys:
- `proxy.zai.enabled`
- `proxy.zai.base_url`
- `proxy.zai.api_key`
- `proxy.zai.dispatch_mode`: `off | exclusive | pooled | fallback`

Dispatch rules:
- `exclusive`: all Anthropic `/v1/messages` → z.ai
- `pooled`: z.ai participates as one “slot” among `(google_accounts + 1)` (no priority guarantee)
- `fallback`: z.ai only when Google pool has zero available accounts

Model mapping UI:
- “Fetch models” button calls `fetch_zai_models(zai, upstream_proxy, request_timeout)`
- User configures:
  - family defaults: `proxy.zai.models.{opus,sonnet,haiku}`
  - exact overrides: `proxy.zai.model_mapping[from] = to`
- The UI supports:
  - dropdown selection from fetched models
  - a “custom” option that reveals a text input

References:
- `docs/zai/provider.md`
- `docs/proxy/routing.md` (provider selection section)

Controls table (z.ai provider):

| Control | Config key(s) | Backend side effects |
|---|---|---|
| Enabled toggle | `proxy.zai.enabled` | Changes routing eligibility for `/v1/messages` |
| Dispatch mode select | `proxy.zai.dispatch_mode` | Alters provider selection behavior |
| Base URL input | `proxy.zai.base_url` | Used for upstream calls |
| API key input | `proxy.zai.api_key` | Used for upstream calls |
| Fetch models | — | Calls `fetch_zai_models(zai, upstream_proxy, request_timeout)` |
| Family defaults | `proxy.zai.models.*` | Affects how `claude-*` is mapped when routed to z.ai |
| Add/Update exact override | `proxy.zai.model_mapping[from]=to` | Applied only on z.ai routed requests |
| Remove exact override | delete `proxy.zai.model_mapping[from]` | — |

#### Exact override controls (UI details)

The “Advanced overrides” section is editable:
- `from` field: text input (incoming model string to match)
- `to` field: dropdown of fetched z.ai model ids, with a “custom” option that reveals a manual text input
- Add/Update button: upserts `proxy.zai.model_mapping[from]=to`
- Delete button per row: removes one mapping

### 4) External providers: z.ai MCP

Rendered as a collapsible card with enable toggle.

Master switch:
- `proxy.zai.mcp.enabled`

Per-server toggles:
- `web_search_enabled`
- `web_reader_enabled`
- `zread_enabled`
- `vision_enabled` (local, in-process server)

Optional upstream key override (MCP only):
- `proxy.zai.mcp.api_key_override`

Web Reader URL normalization:
- `proxy.zai.mcp.web_reader_url_normalization`:
  - `off | strip_tracking_query | strip_query`

References:
- `docs/zai/mcp.md`

Controls table (z.ai MCP):

| Control | Config key(s) | Backend side effects |
|---|---|---|
| Master enabled toggle | `proxy.zai.mcp.enabled` | When off, `/mcp/*` returns 404 |
| API key override | `proxy.zai.mcp.api_key_override` | Used only for remote MCP upstream calls |
| Web Search toggle | `proxy.zai.mcp.web_search_enabled` | Enables `/mcp/web_search_prime/mcp` |
| Web Reader toggle | `proxy.zai.mcp.web_reader_enabled` | Enables `/mcp/web_reader/mcp` |
| zread toggle | `proxy.zai.mcp.zread_enabled` | Enables `/mcp/zread/mcp` |
| Vision toggle | `proxy.zai.mcp.vision_enabled` | Enables `/mcp/zai-mcp-server/mcp` (local) |
| Web Reader URL normalization | `proxy.zai.mcp.web_reader_url_normalization` | Normalizes URLs before forwarding tool calls |

#### MCP endpoints (UI helper panel)

When MCP is enabled and at least one server is toggled on, the UI displays local endpoint URLs for copy/paste:
- `http://127.0.0.1:<port>/mcp/web_search_prime/mcp`
- `http://127.0.0.1:<port>/mcp/web_reader/mcp`
- `http://127.0.0.1:<port>/mcp/zread/mcp`
- `http://127.0.0.1:<port>/mcp/zai-mcp-server/mcp` (local vision server)

### 5) Router: model mapping controls

UI exposes “group mapping” and “custom mapping” as user-friendly controls over:
- `proxy.anthropic_mapping`
- `proxy.openai_mapping`
- `proxy.custom_mapping`

Key operations:
- Update one mapping entry:
  - Persist to config via `save_config`
  - Hot update backend via `update_model_mapping`
- Reset mapping:
  - Restores default mappings (UI-side defaults), then calls the same persist/hot-update path.

Routing implications:
- These mappings affect the Google-backed flows (Gemini/OpenAI and the non-z.ai Anthropic path).
- When requests are routed to z.ai, z.ai mapping logic is separate and only applied in that path.

Reference:
- `docs/proxy/routing.md`

#### Group mappings (predefined)

The UI provides “series” selectors (dropdowns) that update either:
- `proxy.anthropic_mapping` (Claude series keys), or
- `proxy.openai_mapping` (GPT series keys)

Each dropdown includes curated Gemini targets (with optgroups in the UI).

#### Exact mappings (custom)

| Control | Intent | Config key(s) | Behavior |
|---|---|---|---|
| “Original” input | Define match key | (input field) | User types an incoming model id string (e.g. `gpt-4`) |
| “Target” input | Define route target | (input field) | User types a local target model id (e.g. `gemini-2.5-pro`) |
| Add button | Persist custom mapping | `proxy.custom_mapping[from]=to` | Calls `update_model_mapping` (hot) and persists |
| Table row delete | Remove mapping | delete `proxy.custom_mapping[from]` | Calls `update_model_mapping` (hot) and persists |
| Reset mapping | Restore defaults | clears `anthropic_mapping/openai_mapping/custom_mapping` | Calls `update_model_mapping` (hot) and persists |

### 6) Protocol usage examples

UI offers protocol selection and a curated model list, then renders code samples for:
- Anthropic-compatible SDK usage
- Gemini-native SDK usage
- OpenAI-compatible SDK usage

Copy-to-clipboard support is provided for quick setup.

Notes:
- Examples prefer `127.0.0.1` over `localhost` to avoid IPv6 resolution delays in some environments.
- For OpenAI-compatible calls, base URL is `http://127.0.0.1:<port>/v1`.

#### Controls (usage examples)

| Control | Intent | Behavior |
|---|---|---|
| Protocol cards (OpenAI / Anthropic / Gemini) | Pick calling convention | Sets selected protocol; changes example code and “copy endpoints” buttons |
| Copy base URL | Quick configuration | Copies base URL for the selected protocol |
| Copy endpoint path buttons | Quick configuration | Copies full endpoint URL (`/v1/...`, `/v1/messages`, `/v1beta/...`) |
| Model table rows | Pick model | Sets `selectedModelId`; highlights row |
| Copy model id | Quick configuration | Copies `modelId` string |
| Example code panel | Show snippet | Displays Python snippet for selected protocol + model |
| Copy example code | Quick configuration | Copies rendered example snippet |

## “Open Monitor” button behavior

Visible only when proxy is running.
Navigates to `/monitor`, which hosts the Proxy Monitor UI.
