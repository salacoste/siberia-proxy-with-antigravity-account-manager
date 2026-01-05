# Feature Analysis: Config, Tray, and Quota

**Source:** 
*   `src-tauri/src/modules/config.rs`
*   `src-tauri/src/modules/tray.rs`
*   `src-tauri/src/modules/quota.rs`

## 1. Helper Modules Overview

These modules provide essential infrastructure for the "invisible helper" nature of the application.

### A. Configuration (`config.rs`)
*   **Storage:** `gui_config.json` in the OS-specific data directory (e.g., `~/Library/Application Support/...`).
*   **Format:** Standard JSON struct (`AppConfig`).
*   **Capabilities:** Simple Load/Save operations.

### B. System Tray (`tray.rs`)
*   **Framework:** Tauri Tray capabilities.
*   **Behavior:**
    *   **Left Click:** Shows/Focuses the main window.
    *   **Menu Dynamic Updates:** The menu is NOT static. It rebuilds dynamically based on specific events (`tray://refresh-current`, `tray://account-switched`).
*   **Key Menu Items:**
    *   **Status Info:** Displays current email and quota summary.
    *   **Quota Breakdown:** Dynamic list showing specific model usage (Gemini High, Image, Claude 4.5).
    *   **Quick Actions:** "Switch Next" (Instant account cycling) and "Refresh" (Force quota check).
*   **Implementation Note:** Uses `tauri::async_runtime::spawn` to perform heavy lifting (network requests) without freezing the UI thread, then signals back via Events.

### C. Quota System (`quota.rs`)
*   **Target:** Google Cloud Code API (`cloudcode-pa.googleapis.com`).
*   **Workflow:**
    1.  **Auth Check (`fetch_project_id`):** Calls `v1internal:loadCodeAssist` to get the `project_id` and, crucially, the **Paid Tier** status.
    2.  **Quota Fetch:** Calls `v1internal:fetchAvailableModels` on the project.
*   **Tier Detection Logic:** Prioritizes `paid_tier` over `current_tier` to correctly identify Pro/Ultra/Advanced subscriptions.
*   **Handling 403:** Explicitly handles `403 Forbidden` (often happens with "consumer" types or invalid regions) by marking the account as `is_forbidden`.
*   **Mapping:** Filters raw API results to only track relevant models ("gemini", "claude") for display in the Tray/UI.

## 2. Data Persistence & Migration

### A. Data Migration (`migration.rs`)
*   **V1 Import:** capable of scanning `~/.antigravity-agent/` and `antigravity_accounts.json` to import legacy accounts.
*   **Protobuf Decoding:** Uses `prost` and raw field scanning to extract Refresh Tokens from binary `jetskiStateSync` blobs (found in IntelliJ/Android Studio local DBs). This allows "stealing" the token from an authenticated IDE instance.

### B. Proxy Database (`proxy_db.rs`)
*   **Technology:** SQLite (`proxy_logs.db`).
*   **Schema:** Tracks every request's `timestamp`, `model`, `status`, `duration`, and importantly `input_tokens` / `output_tokens`.
*   **Privacy:** Stores `account_email_masked` instead of full emails.
*   **Purpose:** Provides the raw data for the "Monitor" page graphs and statistics.

### C. Smart Token Pooling (`token_manager.rs`)
*   **Sticky Sessions:**
    *   **Cache First:** Tries to route requests from the same `session_id` to the same account to maximize context caching.
    *   **Smart Wait:** If the sticky account is rate-limited, it calculates the exact wait time. If < 5s, it waits; otherwise, it breaks the stickiness and rotates.
*   **Account Tiers:** Sorts accounts by subscription (Ultra > Pro > Free) to prioritize faster-resetting accounts.
*   **Global Lock:** Implements a 60s global lock for non-session requests to prevent "thundering herd" rotation on a single conversation.

