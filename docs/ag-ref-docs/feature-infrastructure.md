# Feature Analysis: Infrastructure & System Integration

**Source:** `src-tauri/src/modules/proxy_db.rs`, `src-tauri/src/modules/tray.rs`, `src-tauri/src/modules/oauth_server.rs`

## 1. Persistence Layer (`proxy_db.rs`)
*   **Technology:** SQLite (via `rusqlite`).
*   **Location:** `proxy_logs.db` in the application data directory.
*   **Schema (`request_logs`):**
    *   Stores comprehensive request metrics: `id`, `timestamp`, `method`, `url`, `status`, `duration`.
    *   AI-specific fields: `model`, `provider`, `input_tokens`, `output_tokens`.
    *   Attribution: `account_id`, `account_email_masked`.
    *   Payloads: `request_body`, `response_body` (text).
*   **Evolution:** Uses "Best Effort" migration strategy (running `ALTER TABLE` statements on init and ignoring "duplicate column" errors).

## 2. System Tray (`tray.rs`)
The system tray acts as a mini-dashboard and quick-control center.

*   **Dynamic Menu:** The menu is rebuilt dynamically based on application state.
    *   **User Info:** Shows current active account email.
    *   **Quota Display:** Hardcoded logic to extract and display percentages for specific "Hot" models:
        *   `Gemini High` (gemini-3-pro-high)
        *   `Gemini Image` (gemini-3-pro-image)
        *   `Claude 4.5` (claude-sonnet-4-5)
    *   **Actions:**
        *   `Switch Next`: Rotates to the next available account in the list.
        *   `Refresh Current`: Triggers a quota refresh for the active account.
*   **Events:** Listens for `config://updated` and internal state changes to refresh the menu UI.
*   **Platform Specifics:** macOS icon template support.

## 3. Local OAuth Server (`oauth_server.rs`)
Handles the OAuth 2.0 Authorization Code callback loop.

*   **Architecture:** Spawns a temporary generic TCP listener (`TcpListener`).
*   **Dual-Stack Binding:**
    *   Attempts to bind ephemeral ports on both IPv6 (`[::1]`) and IPv4 (`127.0.0.1`) to the *same port*.
    *   This mitigates issues where browsers/OS resolvers inconsistently prefer IPv6 or IPv4 for `localhost`.
*   **Flow:**
    1.  Starts listener.
    2.  Generates Authorization URL with `redirect_uri` pointing to `http://localhost:<port>/oauth-callback`.
    3.  Opens system default browser (via `tauri_plugin_opener`).
    4.  Waits for HTTP request containing `code` query parameter.
    5.  Responds with a static HTML "Success/Failure" page.
    6.  Shuts down listener and returns the code to the main thread via `tokio::sync::oneshot`.
