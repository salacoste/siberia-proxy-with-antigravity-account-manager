# Feature Analysis: OAuth & Security

**Source:** `src-tauri/src/modules/oauth.rs`, `src-tauri/src/modules/token.rs`

## 1. Product Logic
The system relies entirely on Google's OAuth 2.0 infrastructure to authenticate users and obtain the necessary credentials (`access_token`, `refresh_token`) to allow the Proxy Engine to access the internal `cloudcode-pa` API.

**Key Components:**
*   **Authentication Flow:** Standard OAuth 2.0 Authorization Code flow.
*   **Token Refresh:** Automated background refreshing of expired tokens.
*   **Scopes:** Requests broad permissions including Cloud Platform, UserInfo, and specific internal scopes (`experimentsandconfigs`, `cclog`).

## 2. Technical Implementation

### A. OAuth Configuration (`oauth.rs`)
*   **Client ID/Secret:** Embedded hardcoded credentials acting as a "Desktop App" client.
    *   `CLIENT_ID`: `1071006060591-tmhssin2h21lcre235vtolojh4g403ep...`
    *   `TOKEN_URL`: `https://oauth2.googleapis.com/token`
*   **Scopes:**
    *   `https://www.googleapis.com/auth/cloud-platform` (Core for Cloud Code)
    *   `https://www.googleapis.com/auth/userinfo.email`
    *   `https://www.googleapis.com/auth/userinfo.profile`

### B. Authorization Flow
1.  **Initiate:** `get_auth_url` generates a URL to `accounts.google.com`.
    *   `access_type=offline` (Critical for getting `refresh_token`).
    *   `prompt=consent` (Force consent screen to ensure refresh token is returned even if previously authorized).
2.  **Local Callback Handling (Robustness):**
    *   **Dual-Stack Binding:** `oauth_server.rs` attempts to bind to `[::1]:0` (IPv6) first, then gets the port to bind `127.0.0.1:{port}`. This ensures `localhost` works regardless of OS resolution preference.
    *   **Graceful Cancellation:** Supports cancellation signals if the user aborts.
3.  **Exchange:** `exchange_code` swaps the returned `code` for tokens.
    *   **Warning:** Logs a warning if `refresh_token` is missing (common issue if user didn't revoke previous access).

### C. Token Management
*   **Refresh Logic:** `ensure_fresh_token` checks local timestamp.
    *   Buffer: Refreshes if expiry is within 5 minutes (`now + 300`).
*   **Persistence:** Tokens are stored in the individual account JSON files (encrypted or plain text depending on OS/implementation, though ref code shows plain JSON serialization in `Account` struct).

## 3. Security Considerations
*   **Hardcoded Secrets:** The application masquerades as a legitimate Google Cloud Code client using extracted credentials.
*   **Scope Minimization:** While it requests `cloud-platform`, it primarily uses it for the specific internal PaLM/Gemini APIs.
*   **User Privacy:** `UserInfo` fetching tries to resolve a display name (Name -> Given+Family -> None) for UI display.
