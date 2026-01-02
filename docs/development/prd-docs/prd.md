# Product Requirements Document (PRD) — Antigravity Manager (Current Implementation)

## 1) Summary

Antigravity Manager is a **local desktop control plane** for:
- managing a pool of Google OAuth accounts (stored locally),
- monitoring per-account quotas for key models,
- and running a local **multi-protocol API Proxy** (Gemini-native, OpenAI-compatible, Anthropic-compatible),
- plus optional z.ai passthrough for Anthropic-compatible requests and MCP endpoints.

Primary goal for users:
- Keep an always-available local setup by **rotating across accounts** and surfacing quota/health status clearly.

## 2) Target users

- Power users who maintain **multiple accounts** and want:
  - quick switching,
  - quota visibility,
  - and a single local proxy endpoint that multiple tools can share.

## 3) Core capabilities (feature list)

### Account lifecycle
- Add account from a refresh token (single or batch).
- Add account via loopback OAuth login (with “missing refresh token” education).
- Import account from:
  - legacy v1 data folder (`~/.antigravity-agent/...`),
  - a default local database path,
  - or a user-picked database file.
- “Current account” concept (selected account id persisted locally).
- Switching the current account is integrated with an external desktop app:
  - switch injects OAuth tokens into a local SQLite DB and restarts the external app so it picks up the new credentials.
- Delete single / delete bulk.
- Auto-disable accounts when OAuth refresh fails with `invalid_grant`.

### Quota monitoring
- Fetch per-account quota snapshot via Google endpoints.
- Store quota snapshot into the account file.
- Surface quota in:
  - Dashboard (summary + current account),
  - Accounts list (filters),
  - Tray menu (compact view).
- Background auto-refresh (interval, minutes).

### Local API Proxy (multi-protocol)
- Start/stop proxy server in-process.
- Proxy routes:
  - Anthropic-compatible (`POST /v1/messages`, `POST /v1/messages/count_tokens`),
  - Gemini-native (`/v1beta/...`),
  - OpenAI-compatible (`/v1/...`),
  - MCP endpoints (`/mcp/*`).
- Global proxy authorization modes (off/strict/all_except_health/auto).
- Safe access logging toggle (method/path/status/latency only).
- Payload-level monitor/recording toggle (stores request/response samples locally).
- Model mapping UI (group mappings + custom mapping overrides).

### z.ai integration (optional)
- Anthropic-compatible upstream passthrough for `/v1/messages` and `/v1/messages/count_tokens`.
- Dispatch modes: off / exclusive / pooled / fallback.
- Model mapping and per-family defaults (opus/sonnet/haiku) with exact overrides.
- Optional MCP reverse-proxy endpoints to upstream z.ai MCP services plus one local “Vision MCP”.

### System integration
- Tray icon with:
  - show window,
  - switch next account,
  - refresh current quota,
  - and a compact quota summary.
- Auto-launch on system startup (OS integration).
- Localization (EN/ZH) and theme (light/dark/system).

## 4) Key product decisions (why it works this way)

### Local-first storage
- Credentials and configs are stored under the app data directory to enable:
  - offline availability,
  - reproducible state,
  - and deterministic proxy behavior.

### “Pool” strategy and resilience
- Proxy account selection uses a rotating pool with a short stability window, and aggressively disables invalid accounts to avoid repeated failures.

### Multi-protocol proxy
- Different clients speak different protocols. The proxy supports multiple “protocol surfaces” by routing on HTTP path, allowing concurrent usage.

### Optional external provider (z.ai) for Anthropic protocol
- When the upstream is Anthropic-compatible, passthrough preserves maximum compatibility and avoids mixing Google/Gemini-specific transforms.

### External app integration via local DB injection
- The “current account” is not just a UI preference. Switching it updates the external app’s credential state stored in a local SQLite DB and restarts that app.
- This approach keeps the external app aligned with the manager without requiring official APIs.

## 5) Non-goals (as implemented)

- No cloud-hosted service.
- No Vertex AI integration.
- No multi-user tenancy.

## 6) Risks & known quirks (important for the rewrite)

These are present in current code; the rewrite should either match them or intentionally fix them (explicitly decided):
- Some backend logs print full account emails and/or request-derived payloads at runtime log level (privacy concern).
- The Accounts page “Refresh one account” currently calls quota refresh **three times** in a row (likely a reliability hack or an accidental duplication).

## 7) Implementation references

See:
- Proxy routing docs: `docs/proxy/routing.md`
- Proxy config keys: `docs/proxy/config.md`
- Account pool behavior: `docs/proxy/accounts.md`
- z.ai details: `docs/zai/implementation.md`
- PRD screens and actions: `docs/development/prd-docs/screens/*`
- Internal API surface: `docs/development/prd-docs/apis/tauri-commands.md`
