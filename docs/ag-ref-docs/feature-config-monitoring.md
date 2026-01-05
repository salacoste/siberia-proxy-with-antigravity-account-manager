# Feature Analysis: Configuration & Monitoring

**Source:** `src-tauri/src/modules/config.rs`, `src-tauri/src/proxy/monitor.rs`

## 1. Product Logic
*   **Configuration:** Simple, persistent settings for the application (e.g., UI language, page size).
*   **Monitoring:** "Black box" flight recorder for the Proxy. Captures every request/response for debugging, auditing, and "Real-time" charts in the UI.

## 2. Technical Implementation

### A. Configuration System (`config.rs`)
*   **Storage:** `gui_config.json` in the app data directory.
*   **Structure (`AppConfig`):**
    *   `language`: UI Language (en, zh-CN, ru).
    *   `accounts_page_size`: User preference for list density.
    *   `default_export_path`: Path for JSON exports.
    *   `antigravity_executable`: Manual override for the target app path (for Injection).
*   **Logic:** Simple Load/Save functions using `serde_json`.

### B. Monitoring System (`monitor.rs`)
*   **Dual Storage Strategy:**
    1.  **Memory (Ring Buffer):** `VecDeque` with a fixed capacity (`max_logs`) for instant UI access.
    2.  **Database (SQLite via `proxy_db`):** *Implied* persistent storage for historical analysis (referenced as `crate::modules::proxy_db`).
*   **Data Captured (`ProxyRequestLog`):**
    *   Method, URL, Status, Duration.
    *   Model, Provider, Resolved Model (Best-effort parsing).
    *   **Attribution:** Account ID/Email (masked) to track which account handled the request.
    *   Token Usage: Input/Output tokens (if available in response).
*   **Real-time Updates:**
    *   Uses `tauri::AppHandle::emit("proxy://request", ...)` to push new logs to the Frontend immediately.
*   **Stats:**
    *   `total_requests`, `success_count`, `error_count`.
    *   Thread-safe counters via `RwLock`.

## 3. Specifics & Edge Cases
*   **Performance:** Memory logs are a Ring Buffer (`pop_back` when full) to prevent leaks.
*   **Privacy:** Account emails are masked in logs (e.g., `s***@gmail.com`) before storage/emission.
*   **Global Toggle:** `enabled` atomic flag allows completely disabling monitoring to save resources.
