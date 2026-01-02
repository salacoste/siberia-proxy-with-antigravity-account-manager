# Architecture (Current Implementation)

## 1) High-level components

### Desktop app shell
- Tauri v2 application.
- Frontend: React + TypeScript + React Router.
- Backend: Rust.

### Backend subsystems
- Account storage and quota fetching: `src-tauri/src/modules/account.rs`, `src-tauri/src/modules/quota.rs`
- OAuth login flow (loopback): `src-tauri/src/modules/oauth_server.rs`, `src-tauri/src/modules/oauth.rs`
- Import/migration: `src-tauri/src/modules/migration.rs`, `src-tauri/src/modules/db.rs`
- Local proxy server (Axum/Hyper): `src-tauri/src/proxy/server.rs` and `src-tauri/src/proxy/handlers/*`
- Proxy monitor persistence (SQLite): `src-tauri/src/modules/proxy_db.rs`
- Tray menu: `src-tauri/src/modules/tray.rs`
- Logging: `src-tauri/src/modules/logger.rs`

## 2) Runtime data flow

### 2.1 Frontend ⇄ Backend control plane

Frontend calls backend via Tauri commands:
- Wrapper: `src/utils/request.ts`
- Commands: `src-tauri/src/commands/mod.rs`, `src-tauri/src/commands/proxy.rs`, `src-tauri/src/commands/autostart.rs`

Frontend errors are bridged to backend logs:
- `src/utils/frontendLogging.ts` and `src/components/common/ErrorBoundary.tsx`
- Backend sink: `frontend_log` command (`src-tauri/src/commands/mod.rs`)

### 2.2 Proxy server lifecycle

- Started/stopped via commands:
  - `start_proxy_service`, `stop_proxy_service` (`src-tauri/src/commands/proxy.rs`)
- Proxy is an in-process Axum server:
  - routes in `src-tauri/src/proxy/server.rs`
  - handlers in `src-tauri/src/proxy/handlers/*`

### 2.3 Background tasks (frontend timers)

The frontend runs periodic tasks based on persisted config:
- Auto refresh quotas (minutes) and auto sync from DB (seconds):
  - `src/components/common/BackgroundTaskRunner.tsx`

Design implication:
- These tasks do not run if the UI process is not running.

### 2.4 Tray integration

- Tray is created during app setup and stays available even when the main window is hidden.
- Tray actions emit events to the frontend:
  - `tray://account-switched`
  - `tray://refresh-current`
- Frontend listens and refreshes state in `src/App.tsx`.

## 3) Observability and logging

### 3.1 App logs
- Console + rotating file logs under the app data directory:
  - `src-tauri/src/modules/logger.rs`

### 3.2 Proxy access log (safe mode)
- Optional access log middleware logs method/path/status/latency only:
  - config: `proxy.access_log_enabled`
  - middleware: `src-tauri/src/proxy/middleware/access_log.rs`

### 3.3 Proxy monitor (payload recording)
- Optional “monitor recording” stores request/response samples and token usage:
  - config: `proxy.enable_logging`
  - middleware: `src-tauri/src/proxy/middleware/monitor.rs`
  - UI: `src/components/proxy/ProxyMonitor.tsx`
  - storage: `proxy_logs.db` (SQLite)

## 4) Cross-platform assumptions

- App data directory is under the user home folder (see `docs/development/prd-docs/data/storage.md`).
- Default database path for “import from DB” is OS-specific.
- Auto-launch uses OS integration via `tauri-plugin-autostart`.

