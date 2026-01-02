# Storage & Local State (Current Implementation)

This document describes what is persisted locally and why.

## 1) App data directory

Root directory:
- `~/.antigravity_tools/`

Implementation:
- `src-tauri/src/modules/account.rs` (`get_data_dir`)

## 2) Primary config file

- `~/.antigravity_tools/gui_config.json`

Contains:
- UI preferences (language, theme).
- Background task settings (auto refresh/sync).
- Proxy configuration (`proxy.*`).
- Optional convenience paths (default export path, executable path).

Implementation:
- `src-tauri/src/modules/config.rs`
- Type shapes:
  - frontend: `src/types/config.ts`
  - backend: `src-tauri/src/models/config.rs`

## 3) Account storage

### 3.1 Account index
- `~/.antigravity_tools/accounts.json`

Stores:
- list of account ids
- current selected account id

### 3.2 Per-account files
- `~/.antigravity_tools/accounts/<account_id>.json`

Contains:
- user identity (email, optional display name)
- OAuth tokens (access + refresh + expiry)
- quota snapshot (models, remaining %, reset time, subscription tier)
- disable flags:
  - `disabled`, `disabled_at`, `disabled_reason`

Implementation:
- `src-tauri/src/modules/account.rs`
- Disable-on-invalid-grant:
  - `src-tauri/src/proxy/token_manager.rs` (proxy-side refresh)
  - `src-tauri/src/modules/account.rs` (quota refresh path)

Security note:
- This design intentionally stores credentials on disk (local-first).
- A rewrite should define equivalent storage but must ensure **no secrets are committed or logged**.

## 4) Proxy monitor database

- `~/.antigravity_tools/proxy_logs.db`

Purpose:
- Persist proxy monitor records (method/url/status/duration/model + optional payload samples and token usage).

Implementation:
- `src-tauri/src/modules/proxy_db.rs`
- UI:
  - `src/components/proxy/ProxyMonitor.tsx`

## 5) App logs

- `~/.antigravity_tools/logs/`

Implementation:
- `src-tauri/src/modules/logger.rs`

## 6) Default external database path (import/sync)

Used by “Import from DB” and “Sync from DB” features.

Default OS-specific path:
- macOS: `~/Library/Application Support/Antigravity/User/globalStorage/state.vscdb`
- Windows: `%APPDATA%\\Antigravity\\User\\globalStorage\\state.vscdb`
- Linux: `~/.config/Antigravity/User/globalStorage/state.vscdb`

Implementation:
- `src-tauri/src/modules/db.rs`

## 7) Exported files

User-chosen output path or `default_export_path`:
- Export format is JSON:
  - `[ { "email": "...", "refresh_token": "..." }, ... ]`

Backend write bypass is used to avoid sandbox scope constraints:
- command: `save_text_file`
- implementation: `src-tauri/src/commands/mod.rs`

