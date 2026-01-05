# Feature Analysis: Process Injection & Switching

**Source:** `src-tauri/src/modules/db.rs`, `src-tauri/src/modules/process.rs`, `src-tauri/src/utils/protobuf.rs`

## 1. Product Logic
The "Switch Account" feature creates a seamless transition between identities in the target external application (e.g., VS Code or Antigravity/Cursor).
**Goal:** Click "Switch" -> Target App restarts logged in as the selected user.

## 2. Technical Implementation

### A. Process Management (`process.rs`)
*   **Running Detection (`is_antigravity_running`):**
    *   Iterates system processes via `sysinfo`.
    *   **Filtering Strategy:**
        *   Excludes itself (the Manager).
        *   Excludes "Helpers" (Renderer, GPU, Utility, Crashpad, etc.) by checking process Name and Args (`--type=`).
        *   **macOS:** Checks for `.app` bundle paths.
        *   **Windows:** Checks for `antigravity.exe`.
        *   **Linux:** Process tree analysis (Ancestors/Descendants) to avoid killing the launcher or shell.
*   **Graceful Shutdown (`close_antigravity`):**
    1.  **Identify Main Process:** Tries to find the browser "Main" process (no `--type` arg).
    2.  **Phase 1 (SIGTERM):** Sends `kill -15` (or `taskkill`) to the Main process.
    3.  **Wait:** Loops for 70% of timeout duration checking if processes took the hint.
    4.  **Phase 2 (SIGKILL):** If still running, sends `kill -9` to *all* remaining related PIDs (cleanup).
*   **Startup (`start_antigravity`):**
    *   Support manual executable path configuration.
    *   **macOS:** Uses `open -a "/path/to.app"`.
    *   **Others:** Spawns binary directly.

### B. Database Injection (`db.rs`)
*   **Target:** `state.vscdb` (SQLite database used by VS Code/Cursor logic).
*   **Path Resolution:**
    *   Heuristic checks: standard config paths (`~/Library/...`, `%APPDATA%`, `~/.config`).
    *   Support for Portable Mode (relative `data/user-data` paths).
    *   User-defined custom paths.
*   **Injection Logic (`inject_token`):**
    1.  **Connect:** Open SQLite connection to `state.vscdb`.
    2.  **Read:** key `jetskiStateSync.agentManagerInitState` from `ItemTable`.
    3.  **Decode:** Base64 decode the blob.
    4.  **Protobuf Manipulation:**
        *   **Remove:** Field `6` (Existing Auth Data).
        *   **Create:** New Field `6` containing:
            *   `access_token` (Field 1, String)
            *   `token_type` = "Bearer" (Field 2, String)
            *   `refresh_token` (Field 3, String)
            *   `expiry` (Field 4, Timestamp Message)
        *   **Merge:** Append new Field 6 to the blob.
    5.  **Encode:** Base64 encode.
    6.  **Write:** Update `ItemTable`.
    7.  **Flag:** Force Onboarding flag (`antigravityOnboarding = true`).

### C. Protobuf Engine (`protobuf.rs`)
*   **Custom Implementation:** Does not use a full Protobuf library. Implements raw Varint reading/writing and Field skipping/replacement.
*   **Why?** To preserve unknown fields in the binary blob without needing the full `.proto` definition of the target app's private schema.

## 3. Specifics & Edge Cases
*   **macOS Path Fix:** Auto-corrects if user points to a Helper inside the `.app` bundle.
*   **Windows PID Kill:** Uses `taskkill /PID` for precision.
*   **Backup:** Creates `state.vscdb.backup` before injecting.
