# Architecture: Siberia Proxy & Antigravity Account Manager

**Version:** 4.0
**Status:** Approved
**Stack:** Go (Backend) + Wails v2 (Bridge) + React/Vite (Frontend)

## 1. High-Level Overview

The system is a unified desktop application where:
1.  **Siberia (Backend):** A high-performance Go process hosting the `Wails` application, a `net/http` Proxy Server (Axum), and the SQLite database layer. It runs the system tray and manages external process injection.
2.  **Antigravity (Frontend):** A React SPA running inside the system WebView, providing the control plane for accounts, settings, and proxy monitoring.

## 2. Directory Structure

```
.
├── apps/
│   ├── frontend/            # React Application (Antigravity)
│   │   ├── src/
│   │   └── ...
│   │
│   └── backend/             # Go Application (Siberia)
│       ├── siberia/         # Go Modules
│       ├── main.go          # Entry point
│       ├── wails.json       # Wails Config (points to ../frontend)
│       └── go.mod
│
└── pkg/                     # (Optional) Shared libraries
```

## 3. Backend Architecture (Siberia)

### 3.1 Proxy Engine (`siberia/proxy`)
*   **Framework:** `Axum` (High performance, ergonomic).
*   **Listener:** Dynamic port binding (default 8000).
*   **Middleware Chain:**
    1.  `Recovery`: Panic safety.
    2.  `Auth`: Bearer token validation (Configurable).
    3.  `AccessLog`: Info-level logging (Method/URL/Duration).
    4.  `Recorder`: Payload capture -> SQLite (Async Buffer).
*   **Routing:**
    *   `/v1/messages` -> `handlers.HandleAnthropic`
    *   `/v1beta/*` -> `handlers.HandleGemini`
    *   `/mcp/*` -> `handlers.HandleMCP`

### 3.2 Account Manager (`siberia/modules/accounts`)
*   **DB:** `sqlite3` with `SQLCipher` (optional) or standard encryption-at-rest.
*   **Token Rotation:** Background goroutine checks token expiry in `accounts.db` and refreshes via Google OAuth loopback.

### 3.3 External Injection (`siberia/modules/injection`)
*   **Logic:**
    1.  **Identify:** Find process ID of target app (e.g., `code`, `chrome`).
    2.  **Terminate:** Send `SIGTERM`/`WM_CLOSE`.
    3.  **Manipulate:** Open target SQLite DB (`state.vscdb`), decode Protobuf blob, replace `accessToken`, `refreshToken`, `expiry`.
    4.  **Respawn:** Relaunch application.

## 4. Frontend Architecture (Antigravity)

### 4.1 State Management
*   **Zustand Stores:**
    *   `useAppStore`: Config, Theme, Language.
    *   `useAccountStore`: List of Accounts, Current Selection, computed "Best".
    *   `useProxyStore`: Server Status (On/Off), Metrics stream.

### 4.2 Wails Bridge
*   Frontend invokes Backend via `window.go.main.App.*`.
*   Backend pushes events via `runtime.EventsEmit`.
    *   `proxy:request`: Real-time traffic log.
    *   `tray:switch`: External tray action occurred.

## 5. Deployment
*   **Build:** `wails build` produces single binary.
*   **CI:** GitHub Actions for multi-platform build (Window/Mac/Linux).
