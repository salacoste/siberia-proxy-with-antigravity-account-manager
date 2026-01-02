# PRD Part 2: Accounts Manager

## 1. Data Model (`Account`)
*   **Fields:** ID (UUID), Platform (Google/etc), Email, Password, 2FA Secret, Status (Health), Tier (Free/Pro/etc), User Notes.
*   **Security:** Refresh Tokens and Authentication Cookies must be stored securely.
*   **Relations:** Assigned Proxy ID (optional).

## 2. Views
*   **List View:** Sortable table, bulk selection checkboxes.
*   **Grid View:** Card-based layout with visual health indicators.
*   **Pagination:** Local page size override.

## 3. Account Actions
*   **Import / Onboarding:**
    *   **OAuth Loopback:** Webserver on localhost to capture OAuth callbacks.
    *   **Refresh Token Paste:** Batch import via JSON array or raw list.
    *   **Legacy Import:** Read from `~/.antigravity/v1/accounts`.
*   **Management:**
    *   **Edit:** Update credentials, notes, proxy assignment.
    *   **Delete:** Soft/Hard delete confirmation.
    *   **Export:** Dump selected accounts to JSON (with warnings for sensitive data).
*   **Operation:**
    *   **Refresh Quota:** Fetch latest usage limits from upstream servce.
    *   **Switch Account (Critical Integration):**
        1.  **Process Control:** Soft-terminate the external target process (e.g., VS Code or Antigravity Browser) using OS logic (`modules/process.rs`).
        2.  **State Backup:** Back up the external DB (`state.vscdb`) to `state.vscdb.backup`.
        3.  **Token Injection:**
            *   Read `state.vscdb` (SQLite).
            *   Locate target row.
            *   **Inject Protobuf Blob:** Decode existing base64, replace Auth Token fields (Access/Refresh/Expiry), re-encode to base64, update DB.
        4.  **Restart:** Relaunch the external process.

## 4. Search & Filtering
*   **Search:** Client-side fuzzy search on Email/ID.
*   **Quick Filters:** Chips for "Alla", "Available", "Low Health", "Pro Tier".
