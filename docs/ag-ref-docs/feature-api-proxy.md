# Feature Deep Dive: API Proxy Module

**Scope:** Analysis of the local API Gateway that provides the OpenAI/Claude-compatible interface. Covers the Control Plane (Axum Server) and the User Interface.

## 1. System Architecture

The API Proxy is a **local reverse proxy** (Gateway) that translates standard AI protocols into Google's internal formats.

### Components
*   **Server Core (`server.rs`):** Rust/Axum HTTP server.
*   **Router:** Handles `/v1/chat/completions`, `/v1/messages`, and `/mcp/*`.
*   **State Manager (`AppState` in `server.rs`):** Holds `Arc<RwLock<...>>` for hot-reloadable configurations.
*   **Frontend UI (`ApiProxy.tsx`):** Configuration dashboard and status monitor.

---

## 2. Backend Control Plane (`server.rs`)

### A. Lifecycle Management
*   **Startup:** `AxumServer::start` binds to a TCP port (default 8045) and spawns a Tokio task.
*   **Shutdown:** Uses a `oneshot::channel` to signal graceful shutdown when the user toggles "Stop" or the app exits.

### B. Hot Reloading
The server supports dynamic updates without restarting the TCP listener. It uses `RwLock` for:
*   **Model Mappings:** `update_mapping` allows changing `gpt-4` -> `gemini-1.5-pro` on the fly.
*   **Security:** Updating API Keys or LAN access rules (`update_security`).
*   **Z.ai Config:** Toggling `exclusive` / `pooled` modes for Anthropic traffic.
*   **Observability:** Toggling Access Logs (`update_observability`).
*   **Upstream:** Changing the target Google endpoint (Prod/Daily) via `update_proxy`.

### C. Routing & Middleware
*   **Attribution:** `attribution_headers` middleware injects `X-Antigravity-Account` headers to tell the client "Who served this request?".
*   **Monitoring:** `monitor_middleware` tracks active requests (RPS) and latency for the live dashboard.
*   **Security:** `auth_middleware` enforces Bearer Token validation if configured.

---

## 3. Frontend Experience (`ApiProxy.tsx`)

### A. Configuration Interface
*   **Server Control:** "Start/Stop" button invoking `start_proxy_service` / `stop_proxy_service`.
*   **Network Config:**
    *   **Port:** Configurable (8000-65535).
    *   **Lan Access:** Toggle to allow other devices on the network to use the proxy (binds 0.0.0.0 vs 127.0.0.1).
*   **Authentication Mode:**
    *   `Off`: No auth required.
    *   `Strict`: Requires generated API Key.
    *   `Auto (LAN)`: Requires Auth only for non-localhost requests.

### B. Dynamic Code Generation
The UI provides a "Help / Connect" section that generates ready-to-run code snippets based on current settings:
*   **Supported SDKs:** Python (OpenAI/Anthropic/Google GenAI), Curl.
*   **Context Aware:** usage of `127.0.0.1:{current_port}` and `{current_api_key}`.
*   **Protocol Switcher:** User can toggle between `OpenAI` vs `Claude` protocol examples to copy the correct boilerplate.

### C. Status Monitoring
*   **Polling:** Polls `get_proxy_status` every 3s to show:
    *   Running State (Green Pulse / Grey).
    *   Active Account Count.
    *   Real-time Port binding.
*   **Z.ai Integration:**
    *   UI to fetch available models from Z.ai upstream.
    *   Custom Mapping table to override specific model routing.

---

## 4. Integration Points

### Backend -> Frontend
*   **Events:** The backend emits Tauri events (though strictly the UI currently relies on **Polling** `get_proxy_runtime_status` for the heavy metrics to avoid event bus congestion).

### Frontend -> Backend
*   **Commands:** `save_config`, `update_model_mapping` (triggers hot reload).
