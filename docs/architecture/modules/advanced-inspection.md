# Architecture: Advanced Inspection Module

## Overview
This module handles the deep inspection, modification, and visualization of traffic beyond simple logging. It powers Regex Filtering, Breakpoints, and WebSocket Inspection.

## 1. Request Breakpoints
Located in `apps/backend/siberia/proxy/breakpoint.go`.

### Components
*   **BreakpointManager**: Global singleton embedded in the `Service`.
*   **PendingRequest**: Struct representing a paused request.
*   **Channels**: Uses Go channels (`pauseCh`, `resumeCh`) to block the Goroutine handling the request in `service.go`.

### Flow
1.  `OnRequest` handler checks `manager.Match(req)`.
2.  If match:
    *   Request body is read and stored.
    *   `PendingRequest` event emitted to Wails (`breakpoint:hit`).
    *   Handler blocks on a unique channel string (Request ID).
3.  Frontend calls `ResumeRequest(id, modifiedData)`.
4.  Manager receives data, updates the request object, and closes the channel.
5.  Handler resumes, replacing `req.Body` and `req.Header` with modified data.

## 2. WebSocket Inspection
Located in `apps/backend/siberia/proxy/websocket.go`.

### Mechanism
Does **not** use standard `httputil.ReverseProxy` for the data stream.
1.  **Detection**: Inspects `Upgrade: websocket` header in `OnRequest`.
2.  **Hijacking**: Uses `goproxy.Hijack` to take over the client TCP connection.
3.  **Tunneling**: Dials the target backend manually.
4.  **TeeReader**:
    *   Creates a pipe logic using `io.TeeReader` to copy stream data to a side-channel buffer.
    *   Parses WebSocket frames (RFC 6455) from this side-channel without blocking or modifying the main tunnel.
5.  **Events**: Emits `proxy:ws:frame` to frontend.

### Frontend
*   `WebSocketViewer.tsx`: Displays frames in an append-only log.
*   Auto-scroll logic handles high-throughput streams.

## 3. Filtering
*   **Strategy**: Client-side filtering (React).
*   **Performance**: Capable of handling ~500 visible rows.
*   **Logic**: `MonitorPage.tsx` implements a token-based parser supporting simple Regex and Key-Value pairs.
