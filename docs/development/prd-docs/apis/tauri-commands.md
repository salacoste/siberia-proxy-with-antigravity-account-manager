# Internal API Surface (Tauri Commands)

This document lists backend commands exposed to the frontend and their expected behavior.

Source of truth:
- `src-tauri/src/commands/mod.rs`
- `src-tauri/src/commands/proxy.rs`
- `src-tauri/src/commands/autostart.rs`

Notation:
- “Emits” refers to frontend events via `app.emit(...)`.
- All commands return an error string on failure (unless noted).

## 1) Config & diagnostics

### `load_config() -> AppConfig`
- Loads `gui_config.json` or returns defaults.
- Used by most pages on mount.

### `save_config(app, proxy_state, config: AppConfig) -> ()`
- Persists `gui_config.json`.
- Emits `config://updated` (tray listens and refreshes itself).
- If proxy is running: hot-updates proxy settings where supported (see `src-tauri/src/commands/mod.rs`).

### `get_data_dir_path() -> string`
- Returns absolute data directory path for UI display.

### `open_data_folder() -> ()`
- Opens the data directory in OS file explorer.

### `clear_log_cache() -> ()`
- Truncates files under `~/.antigravity_tools/logs`.

### `save_text_file(path: string, content: string) -> ()`
- Writes a UTF-8 text file to a user-selected location.
- Used by export flows (e.g. Dashboard export, Accounts export).
- Security note: exports may contain sensitive fields (emails, refresh tokens) depending on the UI feature.

### `frontend_log(level, message, stack?) -> ()`
- Frontend error bridge; applies best-effort redaction.
- Used by `ErrorBoundary` and `frontendLogging` hooks.

## 2) Accounts

### `list_accounts() -> Account[]`
- Reads index + per-account files.

### `get_current_account() -> Account | null`
- Resolves current account id from index; returns loaded account or null.

### `add_account(app, email_ignored, refresh_token) -> Account`
- Refreshes access token from refresh token.
- Fetches user info to determine true email.
- Upserts account and triggers a quota refresh.
- If proxy is running: reloads proxy token pool.

### `delete_account(app, account_id) -> ()`
- Removes account file and index entry.
- Updates tray menus.

### `delete_accounts(app, account_ids[]) -> ()`
- Bulk delete; updates tray menus.

### `switch_account(app, account_id) -> ()`
- Switches the “current account” and applies it to the external Antigravity app:
  - ensures token is fresh (may refresh access token),
  - stops the external Antigravity process (if running),
  - backs up the external `state.vscdb`,
  - injects OAuth tokens into the DB (`jetskiStateSync.agentManagerInitState`),
  - sets current account id in the manager,
  - restarts the external Antigravity app.
- Best-effort refresh quota after switching.
- Updates tray menus.

### `fetch_account_quota(app, account_id) -> QuotaData`
- Refreshes quota for the specified account and persists it.
- Updates tray menus.

### `refresh_all_quotas() -> { total, success, failed, details[] }`
- Serial refresh for all non-disabled, non-forbidden accounts.

## 3) OAuth login (loopback)

### `prepare_oauth_url(app) -> string`
- Prepares local loopback listener(s) and returns an auth URL.
- Also emits `oauth-url-generated`.

### `start_oauth_login(app) -> Account`
- Starts OAuth flow and may open the browser (implementation-defined).
- If successful, upserts account and refreshes quota.
- Emits `oauth-callback-received` when callback is received (see module).

### `complete_oauth_login(app) -> Account`
- Completes OAuth flow without auto-opening a browser.

### `cancel_oauth_login() -> ()`
- Cancels prepared flow and releases the loopback port(s).

## 4) Import / migration

### `import_v1_accounts(app) -> Account[]`
- Imports legacy accounts from `~/.antigravity-agent/` data.
- Triggers quota refresh for each imported account (best-effort).

### `import_from_db(app) -> Account`
- Imports from the default external database path.
- Sets imported account as current account.
- Triggers quota refresh and updates tray.

### `import_custom_db(app, path) -> Account`
- Imports from a user-selected database file path.
- Sets imported account as current account.
- Triggers quota refresh and updates tray.

### `sync_account_from_db(app) -> Account | null`
- Periodic “sync” helper:
  - Reads refresh token from DB.
  - If it differs from current account’s refresh token, performs `import_from_db`.

## 5) Proxy service

### `generate_api_key() -> string`
- Returns a newly generated random proxy API key in `sk-...` format.
- UI uses it for “Regenerate” in the Api Proxy screen.

### `start_proxy_service(config: ProxyConfig, state, app) -> ProxyStatus`
- Starts the Axum server and stores the instance in `ProxyServiceState`.
- Loads account pool into `TokenManager`.
- If zero accounts exist, start is allowed only if z.ai is enabled and dispatch mode is not `off`.

### `stop_proxy_service(state) -> ()`
- Stops the in-process Axum server.
- Note: MCP endpoints are served by the same server; stopping proxy stops `/mcp/*` listeners too.

### `get_proxy_status(state) -> ProxyStatus`
- Returns running/port/base_url/active_accounts.

### `reload_proxy_accounts(state) -> number`
- Reloads token pool from on-disk account files (reflects disable/enables).

### `update_model_mapping(config: ProxyConfig, state) -> ()`
- Hot-updates mapping locks when the server is running.
- Persists mapping into `gui_config.json`.

### `fetch_zai_models(zai: ZaiConfig, upstream_proxy: UpstreamProxyConfig, request_timeout: u64) -> string[]`
- Calls z.ai Anthropic-compatible API `/v1/models`.
- Used by the Api Proxy UI “Fetch models” button.

## 6) Proxy monitor

### `set_proxy_monitor_enabled(state, enabled: bool) -> ()`
- Enables/disables monitor recording at runtime.

### `get_proxy_logs(state, limit?: number) -> ProxyRequestLog[]`
- Returns last N monitor records from SQLite.

### `get_proxy_stats(state) -> ProxyStats`
- Returns aggregate counters from SQLite.

### `clear_proxy_logs(state) -> ()`
- Clears SQLite table and resets counters.

## 7) System integration

### `show_main_window(window) -> ()`
### `show_main_window(app, window) -> ()`
- Forces the main window visible (startup black-screen mitigation for `visible:false`).
- Also best-effort unminimizes and focuses the window.
- macOS: sets activation policy to `Regular` to ensure the app becomes foreground-visible.

### `get_antigravity_path(bypass_config?: bool) -> string`
- Returns an external executable path:
  - config override (unless bypassed),
  - or a runtime probe.

### `check_for_updates() -> { has_update, latest_version, current_version, download_url }`
- Checks GitHub releases for the latest version.

### `toggle_auto_launch(app, enable: bool) -> ()`
- Enables/disables OS autostart.

### `is_auto_launch_enabled(app) -> bool`
- Returns OS autostart enabled state.
