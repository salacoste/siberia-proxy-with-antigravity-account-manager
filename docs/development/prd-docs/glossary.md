# Glossary

## Account (managed)
An on-disk record under `accounts/<id>.json` containing identity, OAuth tokens, and quota snapshots.

## Account pool
The set of enabled, non-forbidden accounts loaded into the proxy’s `TokenManager` for request execution.

## Current account
The account id stored in `accounts.json` as the “selected” account for UI and tray display. This is not the same as the proxy pool rotation.

## Disabled account
An account marked `disabled: true` on disk. Disabled accounts are skipped by the proxy pool loader.
Typical reason: OAuth refresh returned `invalid_grant`.

## Forbidden account
An account whose quota fetch returned `403 Forbidden` and is recorded as `quota.is_forbidden=true`. These accounts are typically skipped by quota refresh automation and should not be used for proxy traffic.

## Quota snapshot
Per-account data describing remaining percentage and reset time for multiple models. Stored in the account file (`quota.models[]`).

## Protocol surface
A group of HTTP endpoints representing one client-facing protocol:
- Anthropic-compatible: `/v1/messages`
- Gemini-native: `/v1beta/...`
- OpenAI-compatible: `/v1/...`
- MCP: `/mcp/...`

## Mapping (model mapping)
A set of config dictionaries that map inbound model names to upstream model ids for Google-backed flows:
- `proxy.anthropic_mapping`
- `proxy.openai_mapping`
- `proxy.custom_mapping`

## z.ai dispatch mode
A rule controlling whether Anthropic-compatible requests go to the Google-backed flow or to z.ai passthrough:
`off | exclusive | pooled | fallback`.

## Monitor recording
Payload-level proxy monitoring controlled by `proxy.enable_logging` and stored in `proxy_logs.db`.

## Access log
A safe-by-default per-request log controlled by `proxy.access_log_enabled` (method/path/status/latency only).

