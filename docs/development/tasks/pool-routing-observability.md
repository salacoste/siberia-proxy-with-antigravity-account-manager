# Task: Pool Routing Attribution & Runtime Status

## Goal

Make it obvious to users (and debuggable for developers) **which upstream provider, which resolved model, and which Google account** is currently serving requests through the API Proxy—without requiring the client to explicitly choose the exact model.

This is specifically for a proxy that:

- Uses a **Google OAuth account pool** (AI Studio style), *not Vertex*.
- Optionally uses an **Anthropic-compatible upstream** (z.ai) for Claude-protocol requests.
- Supports multiple inbound protocols (Anthropic `/v1/messages`, Gemini native `/v1beta/...`, OpenAI `/v1/...`).

## Non-goals / out of scope

- Vertex AI integration.
- Forcing clients to explicitly choose exact upstream model ids (the proxy remains “mapping-driven”).

## Why this matters

Users have multiple accounts with token limits/quotas. When limits are hit and we rotate/fallback, users need to understand:

- Which account is consuming budget *now* (and in the recent history).
- Which upstream model is actually being called (after model mapping).
- Why a request rotated/fell back (429/401/403/etc).

## References (what to borrow)

### 1) `coffeegrind123/gemini-for-claude-code`

Borrow design ideas, not code:

- **Anthropic compatibility**: strict handling for `/v1/messages` and `/v1/messages/count_tokens`.
- **Streaming resilience**: retry streaming on malformed chunks, fallback to non-streaming when needed.
- **Model alias mapping**: “family mapping” (small/big) plus optional exact overrides.
- **Actionable error messages**: classify upstream errors to reduce user confusion.

### 2) `catlog22/Claude-Code-Workflow` (CCW)

Borrow dashboard/status patterns:

- **Aggregated status endpoint**: one request returns everything needed for a UI (example pattern: `ccw/src/core/routes/status-routes.ts`).
- **Provider/model configuration UX**: consistent toggles, model discovery UI, “test connection” concepts.
- **Operational focus**: clear visibility of installed components and runtime state.

## Current state (in our codebase)

### Google account selection happens in one place

- Account pool logic: `src-tauri/src/proxy/token_manager.rs` (`TokenManager::get_token`)
  - Performs 60s lock to stabilize account selection.
  - Rotates on retry attempts.
  - Refreshes OAuth tokens and disables invalid accounts (`invalid_grant`).

### Proxy logging exists but lacks attribution

- Monitor middleware persists request/response samples and token usage when available:
  - `src-tauri/src/proxy/middleware/monitor.rs`
  - `src-tauri/src/proxy/monitor.rs`
  - DB persistence: `src-tauri/src/modules/proxy_db.rs`
- Access log middleware exists (safe method/path/status/latency only):
  - `src-tauri/src/proxy/middleware/access_log.rs`

### Quota (“budget”) data already exists, but is not surfaced in Proxy UI

- Per-account quota snapshots are fetched and stored in account files (`quota.models[].percentage`, `reset_time`):
  - Fetch logic: `src-tauri/src/modules/quota.rs`
  - Data model: `src-tauri/src/models/quota.rs`
  - UI currently surfaces this on the Accounts page, not the Proxy page.

### Routing across providers exists (but attribution is not surfaced)

- Anthropic `/v1/messages` routing decision: `src-tauri/src/proxy/handlers/claude.rs`
  - Dispatch modes: `off | exclusive | pooled | fallback` via `ZaiDispatchMode` (`src-tauri/src/proxy/config.rs`).

### Important privacy gap (must be addressed while implementing attribution)

Today, some handlers emit request-derived content to runtime logs (and sometimes include full emails):
- `src-tauri/src/proxy/handlers/openai.rs` logs “Transformed Gemini Body” at `info` level.
- `src-tauri/src/proxy/handlers/claude.rs` has similar “Transformed Gemini Body” logging.
- `src-tauri/src/proxy/handlers/gemini.rs` logs “Using account: <email> …”.

This conflicts with the privacy guarantees we want for production usage. Attribution work should include tightening these logs (gate behind debug-only, or remove/replace with sanitized metadata).

## Requirements

### Functional

1. Users can see, in UI:
   - Last upstream provider used (google pool vs z.ai).
   - Last resolved upstream model.
   - Last Google account used (stable identifier, not token).
   - Whether rotation happened (and why, if known).
2. Users can see recent history (e.g. last 50–200 requests):
   - Grouped by account (requests, errors, tokens).
3. “Family mapping” remains the default:
   - Users configure only big/medium/small (or opus/sonnet/haiku), with optional exact overrides.
4. No Vertex integration and no requirement for the client to specify exact Gemini model names.

### Security / Privacy

- Never store or log secrets (API keys, OAuth access tokens/refresh tokens).
- Avoid storing full email addresses in logs/DB.
  - Prefer `account_id` (local random id) and an optional **masked email** (e.g. `jo***@g***.com`) only for display.
- Avoid emitting full emails in runtime logs (e.g. `tracing::info!("Using account: ...")`); prefer `account_id` or masked email.
- Monitor payload logging must remain optional and off by default (keep current `enable_logging` behavior).
- Attribution visibility must not require payload capture:
  - “Which account/provider served this request” should be available in a safe mode even when monitor recording is off.

## Proposed implementation (high-level)

### A) Introduce a “request attribution” struct

Add a small struct that represents **who served the request**:

- `request_id`: UUID
- `protocol`: `anthropic | gemini | openai | mcp`
- `provider`: `google_pool | zai | unknown`
- `incoming_model`: the model string in the inbound request (if any)
- `resolved_upstream_model`: the mapped model actually used upstream
- `account_id`: Google pool account id (if provider is `google_pool`)
- `quota_group`: `gemini | claude | image_gen | ...` (whatever we already use)
- `attempt`: retry attempt index (0..N)
- `rotated`: boolean
- `decision_reason`: optional short enum/string (e.g. `dispatch_mode`, `429_quota`, `invalid_grant_refresh`, `project_id_fetch_failed`)
- `quota_snapshot`: optional best-effort (remaining % + reset time) for the relevant model group at request time

Store it in **response extensions** (Axum) so middleware can read it without changing handler signatures everywhere.

We already use this pattern for provider label:

- `src-tauri/src/proxy/observability.rs` (`UpstreamRoute`)
- `src-tauri/src/proxy/handlers/claude.rs` inserts `UpstreamRoute("zai")` for z.ai passthrough responses

Extend this to cover:

- `resolved_upstream_model`
- `account_id` (google pool)
- `provider=google_pool` for all Google-backed flows (so access logs don’t stay `unknown`)

Implementation note (TokenManager API):
- Today `TokenManager::get_token(...)` returns `(access_token, project_id, email)`.
- Attribution needs the stable `account_id`, so either:
  - change the return type to include `account_id`, or
  - add a new method (e.g. `get_token_with_account_id`) and migrate callers.

### B) Extend monitor logging to capture attribution

Update the monitor middleware (`src-tauri/src/proxy/middleware/monitor.rs`) to read:

- `response.extensions().get::<...>()`

Then enrich `ProxyRequestLog` with new optional fields:

- `provider`
- `resolved_model`
- `account_id`
- `quota_group`

If DB schema is too heavy to migrate immediately, maintain:

- A small in-memory ring buffer keyed by `request_id` (but prefer DB for history).

### B2) Add a safe, always-available attribution store (separate from payload capture)

Because users need “who served my traffic” even when full recording is off:

- Maintain a lightweight in-memory ring buffer of **attribution-only** events (no headers, no query strings, no bodies).
- Optional: persist attribution-only events into SQLite in a separate table (small schema, safe by construction).
- Drive the Proxy UI “Now serving / Recent usage” from this safe store, not from the payload monitor DB.

### C) Expose an aggregated runtime status endpoint (UI-friendly)

Add a single API surface (Tauri command recommended, since UI already uses `invoke`):

- `get_proxy_runtime_status()` returns:
  - Proxy running/port (existing fields)
  - Current dispatch modes (`zai.dispatch_mode`, auth mode)
  - Google pool size / enabled accounts
  - TokenManager “lock” snapshot:
    - last used `account_id`
    - lock age / remaining seconds
  - Recent request summary:
    - last N attributions
    - per-account aggregates (requests/errors/tokens last 15m + last 24h)
    - per-provider aggregates

This mirrors the CCW “one call for the dashboard” pattern.

Include quota context in the same call (to satisfy “budget visibility”):
- For each account in the pool, include a sanitized quota summary (from on-disk account `quota` snapshots), e.g.:
  - subscription tier (`FREE/PRO/ULTRA` when available)
  - remaining % for relevant models (gemini/claude)
  - reset time
  - forbidden flag (`is_forbidden`)

### D) UI: add a “Routing / Active Usage” panel

In `src/pages/ApiProxy.tsx`:

- “Now serving” section:
  - Provider badge (google pool / z.ai)
  - Resolved model id
  - Account badge (masked email + stable id)
  - Lock indicator (e.g. “sticky for 42s”)
- “Recent usage” table:
  - Account id (and optional masked email)
  - Requests, errors
  - Input/output tokens (best-effort)
  - Last used timestamp
- Keep tooltips consistent with existing HelpTooltip behavior.

## Corner cases to handle

- Concurrent requests: there is no single “active” account globally; show “last used” and “recent N”.
- Streaming responses: token usage may only arrive at end; monitor must not block streaming.
- Missing usage metadata: still record provider/model/account; tokens optional.
- Rotations: distinguish between retry-based rotation vs dispatch-mode routing decisions.
- Internal model changes: non-trivial logic can change the effective upstream model (e.g. background task detection/downgrades); the UI must show both incoming and resolved model.
- Disabled accounts: if `invalid_grant` disables an account, reflect this in “pool size” and recent errors (without exposing email in API responses).

## Suggested acceptance criteria

1. With 2+ Google accounts, repeated requests show:
   - Account id changes when rotation happens (attempt > 0 or error conditions).
   - Account id stays stable inside the 60s lock window (current behavior).
2. For Anthropic `/v1/messages`:
   - When z.ai dispatch is `exclusive`, “provider=zai” is displayed.
   - When dispatch is `pooled`, both providers appear over time.
3. UI renders a clear “Now serving” card and a “Recent usage” table.
4. No secrets appear in logs/DB/diff (verify with repository-wide grep before merge).
5. With monitor recording disabled, the UI still shows provider/model/account attribution (safe telemetry mode).

## Implementation plan (phases)

### Phase 1 — Attribution plumbing (minimal risk)

- Add `RequestAttribution` (response extensions) and set it in:
  - Gemini handler (`src-tauri/src/proxy/handlers/gemini.rs`) after account selection + model resolution.
  - Anthropic handler (`src-tauri/src/proxy/handlers/claude.rs`) for both z.ai and Google paths.
- Also set it in OpenAI protocol handler (`src-tauri/src/proxy/handlers/openai.rs`).
- Ensure `UpstreamRoute("google")` (or equivalent) is set for Google-backed responses so access logs have a meaningful upstream label.
- Update monitor middleware to capture attribution from `response.extensions()`.
- Add a Tauri command returning a UI-friendly aggregated status snapshot (single call).
- UI: add “Now serving” and “Recent usage” panels.

### Phase 2 — Persistence + aggregation quality

- Extend `ProxyRequestLog` and DB schema to persist provider/account/resolved_model.
- Add server-side aggregation helpers (per-account, per-provider, time windows).
- Add “copy diagnostics” feature that is sanitized by construction.

### Phase 3 — Compatibility hardening (borrow from `gemini-for-claude-code`)

These are optional but strongly recommended once attribution exists (they become easier to debug):

- Implement real token counting for Anthropic-protocol clients:
  - `POST /v1/messages/count_tokens` when routed to Google pool (use Gemini countTokens upstream).
  - Ensure z.ai path stays passthrough for count_tokens.
- Add streaming recovery strategy:
  - Detect malformed SSE/JSON chunk sequences and retry streaming (bounded).
  - If still failing, fallback to non-streaming completion.
- Add clearer error classification for common upstream errors (429 quota, 401/403 auth/region, 5xx transient).

### Phase 4 — Privacy hardening (required for production readiness)

- Remove or gate request/response body debug logging in protocol handlers:
  - Keep payload capture exclusively inside the explicit “monitor recording” feature.
- Ensure any remaining “account selection” logs do not print full emails by default (use `account_id` or masked email).
