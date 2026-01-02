# PRD Part 5: Proxy Monitor

## 1. Monitor UI (`/monitor`)
*   **Access:** Full-screen interface accessible from `/api-proxy` when proxy is running.
*   **Recorder Toggle:** "Enable Logging" switch.
    *   **Off (Default):** No payloads stored (only metadata).
    *   **On:** Payloads (Body) captured to local SQLite monitoring DB.

## 2. Traffic Inspector
*   **Live Table:** Real-time stream of requests.
    *   Columns: Status Code, Method, URL path, Duration (ms), Model ID (if detected).
*   **Details Modal:**
    *   **Metadata:** Full Timestamp, Headers.
    *   **Payloads:** Request/Response Body with JSON syntax highlighting.
    *   **Metrics:** Calculated Input/Output tokens (if parseable).

## 3. Data Management
*   **Persistence:** SQLite database (separate from Accounts DB).
*   **Retention:** Manual "Clear Logs" button. Auto-pruning (optional future feature).
*   **Privacy:** Sensitive headers (Authorization) must be masked in the UI or DB unless explicitly debugged.
