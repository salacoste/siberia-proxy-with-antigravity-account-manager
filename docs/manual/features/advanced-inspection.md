# User Guide: Advanced Traffic Inspection

## Overview
The Advanced Inspection suite allows developers to deeply analyze, filter, and modify HTTP/WebSocket traffic flowing through the Siberia Proxy.

## 1. Advanced Filtering
The Traffic Monitor search bar now supports powerful filtering syntax.

### Search Syntax
| Filter Type | Syntax | Example | Description |
|---|---|---|---|
| **Text** | `word` | `json` | Matches URL, Method, or Status containing "json". |
| **Field** | `field:value` | `method:POST` | Matches specific fields. Supported: `method`, `status`, `url`, `host`. |
| **Regex** | `/pattern/` | `/api\/v\d+/` | JavaScript-compatible Regex matching on URL. |
| **Negation** | `!term` or `-term` | `!image` | Excludes requests matching the term. |

**Examples:**
- Find all POST errors: `method:POST status:5`
- Exclude CSS/Images: `!css !image !png`
- Complex API tracking: `/api\/v1\/users/ method:GET`

## 2. Request Breakpoints
Intercept and modify requests on the fly (similar to "Charles Proxy").

### Managing Breakpoints
1.  Navigate to the **Monitor** page.
2.  Click the **Breakpoints** toggle to open the panel.
3.  Click **Add Rule**.
4.  Enter a URL pattern (e.g., `/api/auth`) and Method (e.g., `POST`);
5.  Ensure the "Enabled" checkbox is ticked.

### Handling Intercepted Requests
When a request matches a rule:
1.  The UI will show a modal "Pending Request".
2.  The request hangs at the proxy level.
3.  **Modify Headers/Body**: You can directly edit the JSON body or headers.
4.  **Actions**:
    *   **Resume**: Forward the (potentially modified) request to the server.
    *   **Drop**: Kill the connection immediately.

## 3. WebSocket Inspection
Inspect real-time WebSocket frames.

1.  Click the **WebSockets** toggle on the Monitor page.
2.  Initiate a WS connection in your app.
3.  Frames will appear in the "Live Stream" view.
    *   **Green**: Outgoing (Client -> Server)
    *   **Blue**: Incoming (Server -> Client)
4.  Payloads are displayed in text format.
