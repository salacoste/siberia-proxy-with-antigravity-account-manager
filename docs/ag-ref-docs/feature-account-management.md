# Feature Deep Dive: Account Management

**Scope:** End-to-end analysis of how user accounts are stored, authenticated, managed, and injected into target IDEs.

## 1. Architecture Overview

The Account Management system is designed for **high availability** and **atomicity**. It avoids a single large database file, preferring a distributed file structure to prevent corruption and allow easy manual intervention.

### Components
*   **Backend Storage (`account.rs`):** JSON-based file system storage.
*   **Authentication (`oauth.rs`):** Google OAuth 2.0 implementation with automatic token refreshing.
*   **Frontend State (`useAccountStore.ts`):** Zustand-based store with optimistic updates for UI responsiveness.
*   **Injection Logic:** Direct manipulation of IDE (VS Code) SQLite databases to inject credentials.

---

## 2. Backend Implementation (`modules/account.rs`)

### A. Storage Strategy
Instead of a single `db.sqlite` or `accounts.json`, the system uses a **Split-Index** approach:

1.  **Index File (`accounts.json`):**
    *   Stores a lightweight list of `AccountSummary` (ID, Email, Name, Order).
    *   Maintains the `current_account_id`.
    *   **Purpose:** Fast listing and sorting without loading full token data.
    *   **Safety:** Uses an auxiliary `.tmp` file + atomic `fs::rename` to prevent corruption during write.

2.  **Data Files (`accounts/{uuid}.json`):**
    *   Stores full `Account` details including sensitive `TokenData` (Access/Refresh Tokens), Quota cache, and Project ID.
    *   **Purpose:** Atomicity. Corruption of one account file does not affect others.

### B. Key Operations
*   **Upsert (Add/Update):** standard read-modify-write pattern with a global `Mutex` lock (`ACCOUNT_INDEX_LOCK`) to ensure thread safety during index updates.
*   **Switch Account:**
    1.  Validates existence.
    2.  **Auto-Refresh:** Checks `ensure_fresh_token` to guarantee the token is valid before injection.
    3.  **Process Control:** Gracefully terminates any running `antigravity` process (IDE agent).
    4.  **DB Injection:** Connects to the IDE's state DB (e.g., `state.vscdb`), backs it up, and injects the new OAuth token.
    5.  **Restart:** Relaunches the agent process.

### C. Injection Mechanics (`modules/db.rs`)
*   **Path Resolution:**
    *   **Priority 1:** `--user-data-dir` (Active Process args).
    *   **Priority 2:** Portable Mode path relative to executable (`data/user-data/...`).
    *   **Priority 3:** OS Default (MacOS `Library/...`, Win `AppData`, Linux `.config`).
*   **Protobuf Hack:**
    *   Target Key: `jetskiStateSync.agentManagerInitState` in `state.vscdb`.
    *   **Operation:** Read Blob -> Base64 Decode -> **Remove Field 6** (Old Token) -> **Inject New Field 6** (New Token + Expiry) -> Base64 Encode -> Write Back.
    *   **Why Field 6?** Reverse engineering revealed Field 6 holds the OAuth structure for the Agent Manager state.
    *   **Onboarding:** Also sets `antigravityOnboarding = true` to skip the "Welcome" wizard.

*   **Reorder:** Updates the order in `accounts.json` to match user preference (Drag & Drop support).

---

## 3. Authentication Flow (`modules/oauth.rs`)

### A. Google OAuth 2.0
*   **Client:** Uses a hardcoded `CLIENT_ID` (Antigravity/CloudCode app).
*   **Scopes:**
    *   `cloud-platform` (Project access)
    *   `userinfo.email` / `profile`
    *   `cclog`, `experimentsandconfigs` (Internal APIs)

### B. Token Lifecycle
*   **Exchange:** Swaps Authorization Code for Access + Refresh Token.
*   **Refresh Strategy (`ensure_fresh_token`):**
    *   Buffer: 5 minutes. If `expiry - now < 300s`, triggers refresh.
    *   **Persistence:** If a refresh yields a new `access_token`, it is immediately saved to the individual account file.
    *   **Error Handling:** If `invalid_grant` (revoked), the account is marked as `disabled` to prompt user re-login.

---

## 4. Frontend Integration (`useAccountStore.ts`)

### A. State Management
*   **Framework:** Zustand.
*   **State:** `accounts` (List), `currentAccount` (Detail), `loading`, `error`.

### B. Optimistic UI
For **Reordering** (Drag and Drop):
1.  **Immediate:** The Store updates the local `accounts` array state immediately.
2.  **Async:** Sends the `reorderAccounts` command to the backend.
3.  **Rollback:** If the backend fails, reverts the local state to the previous order.

### C. Service Layer (`accountService.ts`)
Acts as the IPC bridge. Uses `invoke` to call Rust commands:
*   `get_accounts`
*   `switch_account`
*   `refresh_token_force`
*   `delete_account`

---

## 5. Security & Isolation
*   **Token Storage:** Stored in plain text JSON in the user's local `AppData` directory. Relies on OS-level file permissions.
*   **Privacy:** `proxy_db` (request logs) masks user emails (`s***@gmail.com`) to prevent accidental leaks in screenshots/logs.
