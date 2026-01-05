# Story-32: WebSocket Frame Inspector

**Epic:** [Epic-09: Advanced Traffic Inspection](../epic-09-advanced-inspection.md)
**Status:** Completed
**Parent:** Epic-09

## Description
Unlock deep visibility into WebSocket connections. Instead of just seeing "Upgrade: websocket", the user should be able to click on a WS connection and see the stream of frames (Text/Binary) flowing between client and server in real-time.

## Requirements
1.  **Frame Interception**:
    -   Intercept WebSocket frames in the `goproxy` tunnel.
    -   This likely requires a custom hijacking of the connection or using a `goproxy` feature if available (or implementing a custom dialer/pipe).
2.  **Telemetry Events**:
    -   Emit `proxy:ws:frame` events for every frame.
    -   Include: Payload (text/truncated binary), Direction (Up/Down), OpCode, Timestamp.
3.  **UI Inspector**:
    -   A dedicated tab or view in the Traffic Monitor side panel for "Frames".
    -   List frames chronologically.
    -   Color code directions (⬆️ Client->Server, ⬇️ Server->Client).
4.  **Performance**:
    -   Do not buffer infinite frames in backend memory.
    -   Frontend should handle virtualization if high RPS.

## Acceptance Criteria
- [x] WebSocket connections appear in the Traffic Table with status "101 Switching Protocols" (or similar).
- [x] Clicking the row opens details with a "Frames" tab.
- [x] Sending a WS message from the client appears in the "Frames" list instantly.
- [x] Server responses appear in the "Frames" list.
- [x] App does not crash on high-throughput socket.

## Dev Notes
-   `goproxy` interception of WS is tricky. We might need to manually `Hijack` the connection in `OnRequest` for `Upgrade: websocket`.
-   See `siberia/proxy/websocket.go` (if it exists) or implement new logic there.
-   Frontend `TrafficContext` already has `wsFrames` state placeholder.

## Dev Agent Record
-   **Status**: Dev Complete
-   **Branch**: `feat/story-32-websocket`
-   **Verification**:
    -   Backend Unit Tests: PASS (`go test ./...`)
    -   Frontend Build: PASS (`npm run build`)
    -   MitM Logic: Implemented custom `ConnectHijack` to ensure full control over TLS Handshake and routing to `s.ServeHTTP`.
-   **Models**: Updated `ProxyRequestEvent` and `WebSocketFrame` to share `connection_id` for accurate UI filtering.

## QA Results
-   **Status**: Ready for QA
