# Feature Deep Dive: Observability & Monitoring

**Scope:** Analysis of how the system tracks, logs, and visualizes traffic data, ensuring user privacy while providing actionable insights.

## 1. Data Architecure (`proxy_db.rs`)

### A. Storage
*   **Engine:** SQLite
*   **Location:** `proxy_logs.db` in the user's data directory.
*   **Schema:** `proxy_request_logs` table.
*   **Retention:** Configurable `MAX_LOGS` (default 500) kept in memory ring buffer for instant UI access, while SQLite stores history.

### B. Fields Tracked
*   `access_log_enabled`: Feature toggle.
*   **Metrics:**
    *   Timestamp, Duration (latency).
    *   Input/Output Tokens (cost tracking).
    *   HTTP Status Code.
*   **Context:**
    *   `model`: The requested model (e.g. `gpt-4`).
    *   `resolved_model`: The actual upstream model used (e.g. `gemini-1.5-pro`).
    *   `provider`: Who served it (Google / Z.ai).
*   **Privacy:**
    *   `account_email_masked`: Stores `s***@gmail.com` to allow debugging specific accounts without storing full PII.
    *   `request_body` / `response_body`: Optional full logging (defaults to disabled for privacy).

---

## 2. Runtime Monitoring (`monitor.rs`)

### A. In-Memory Buffer
*   **Structure:** `VecDeque` (Ring Buffer).
*   **Purpose:** Powers the "Live Monitor" dashboard graph without querying disk on every frame.
*   **Stats:** Maintains running counters for `total_requests`, `success_count`, `error_count` for O(1) status checks.

### B. Event Emitters
*   **Channel:** `proxy://request`
*   **Data:** Emits `ProxyRequestLog` struct for every completed request.
*   **Frontend Consumer:** The Dashboard listens to this event to render the "traffic light" dots and real-time RPS graph.

---

## 3. Frontend Visualization

### A. Dashboard Components
*   **Traffic Graph:** Uses `recharts` (implied from React structure) to plot RPS/Latency over time.
*   **Request Table:**
    *   Columns: Time, Method, Path, Status, Latency, Model.
    *   **Attribution:** Visual indicators showing which specific Google Account handled the request (using the masked email).
*   **Health Check:** Visual distinction between "success" (2xx - Green), "client error" (4xx - Yellow), and "server error" (5xx - Red).
