# PRD Part 3: API Proxy & Network Handlers

## 1. Local Proxy Server (`Siberia Core`)
*   **Protocol:** HTTP & SOCKS5 support.
*   **Binding:** Port 8000-65535, localhost vs 0.0.0.0 toggle.
*   **Authorization:** Off / Strict / Auto modes.

## 2. Account Pool & Routing
*   **Goal:** Distribute load across a local pool of Google accounts.
*   **Token Manager:**
    *   **Stickiness:** Short window (e.g. 60s) to keep a client on one account for bursts.
    *   **Rotation:** Automatic failover on `429` or error.
    *   **Refresh:** Auto-refresh tokens on `invalid_grant` or expiry.
    *   **Locking:** Prevent concurrent usage of same account for exclusive operations (like Image Gen) where applicable.
*   **Request Type Inference:** Proxy must detect payload type (Gemini vs Claude vs Image) to select appropriate quota group.

## 3. Upstream Management
*   **Providers:**
    *   **z.ai Integration:**
        *   **Dispatch:** Round-robin or Fallback logic between z.ai and Google Account pool.
        *   **Model Mapping:** "Exact" and "Family" based routing.
*   **MCP Support:** Embedded servers for Web/Reader/Vision.

## 4. Developer Tools
*   **Code Generators:** Python/cURL/Node snippets.
*   **Monitor Button:** Quick jump to traffic inspector.
