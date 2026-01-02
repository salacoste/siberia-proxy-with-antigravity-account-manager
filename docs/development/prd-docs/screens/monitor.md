# Screen: Proxy Monitor

Routes:
- `/monitor` (full-screen monitor)
- Entry point from Api Proxy page when proxy is running.

Source:
- `src/pages/Monitor.tsx`
- `src/components/proxy/ProxyMonitor.tsx`
- Backend:
  - monitor: `src-tauri/src/proxy/monitor.rs`
  - middleware: `src-tauri/src/proxy/middleware/monitor.rs`
  - DB: `src-tauri/src/modules/proxy_db.rs`
  - commands: `src-tauri/src/commands/proxy.rs`

## Purpose

Provide **local traffic inspection** for debugging proxy usage:
- last N requests (real time + historical),
- quick filters,
- per-request details (including optional payload capture),
- and a “recording” toggle.

## Core concepts

### “Recording” is separate from safe access logs

- Recording toggle (`proxy.enable_logging`) controls whether the proxy captures payloads and persists to SQLite.
- Safe access logs (`proxy.access_log_enabled`) are configured elsewhere and do not store payloads.

## Layout

1) Header
- Back button (to `/api-proxy`)
- Title and subtitle

2) Monitor dashboard
- Stats summary (total/success/error)
- Recording toggle
- Search filter input
- Quick filters
- Table of recent requests
- Request details modal
- Clear logs (with confirmation)

## Interactions (buttons & controls)

| UI element | Intent | Backend behavior |
|---|---|---|
| Back button | Return to Api Proxy | Client-side navigation |
| Recording toggle | Enable/disable payload monitoring | Loads config, sets `proxy.enable_logging`, calls `save_config`, calls `set_proxy_monitor_enabled(enabled)` |
| Search input | Filter logs client-side | No backend calls |
| Quick filter chips | Set filter string | No backend calls |
| Reset filter | Clear current filter | Sets filter to empty string | No backend calls |
| Click log row | View details | Client-side modal |
| Clear logs | Reset history | Calls `clear_proxy_logs()`, clears UI state |

## Details modal (per-request)

Opened by clicking a table row.

| Control | Intent | Behavior |
|---|---|---|
| Click outside modal | Close | Closes modal (does not change filter/history) |
| “X” close button | Close | Closes modal |
| Request payload section | Inspect | Renders JSON prettified if parsable; otherwise raw text |
| Response payload section | Inspect | Renders JSON prettified if parsable; otherwise raw text |

## Real-time updates

When recording is enabled:
- Backend emits `proxy://request` events to the frontend for each stored log record.
- UI listens via `@tauri-apps/api/event` and prepends the new record.

Event emission:
- `src-tauri/src/proxy/monitor.rs` (`app.emit("proxy://request", &log)`)

## What is captured (current implementation)

Stored fields include:
- method, url, status, duration
- model (best-effort)
- request_body (up to 1MB; only for POST)
- response_body (up to 512KB; for JSON/text responses; streams store `"[Stream Data]"`)
- input_tokens/output_tokens (best-effort; attempts to parse from SSE tail or JSON `usage`)

Implementation details:
- `src-tauri/src/proxy/middleware/monitor.rs`
- `docs/proxy-monitor-technical.md`

Security implications:
- Payload capture can include sensitive content. It must be opt-in and off by default.
