# External app integration (DB injection + process control)

This feature is critical to understanding “why switching accounts matters”.

## Goal

When the user switches the current account in the manager, the external desktop app (“Antigravity”) should:
- start using the newly selected account credentials,
- without requiring manual relogin inside that app.

## How it works (current implementation)

Switch sequence (simplified):
1) Load target account by id (manager account file).
2) Ensure access token is fresh:
   - refresh if expiring.
3) If the external app is running:
   - close it gracefully (with timeout, platform-specific).
4) Locate the external app’s local state database (`state.vscdb`).
5) Back up the DB to `state.vscdb.backup`.
6) Inject new OAuth fields into the DB row:
   - table: `ItemTable`
   - key: `jetskiStateSync.agentManagerInitState`
   - value: base64-encoded protobuf blob
7) Persist “current account id” in the manager index.
8) Restart the external app.

Primary entry point:
- `src-tauri/src/modules/account.rs` (`switch_account`)

## DB format details

The DB entry is a base64-encoded protobuf blob.
Injection logic:
- decode base64
- remove old Field 6 (oauth token info)
- create a new Field 6 containing:
  - access token
  - refresh token
  - expiry timestamp
- re-encode base64 and write back

Implementation:
- DB path: `src-tauri/src/modules/db.rs` (`get_db_path`)
- Extraction/import: `src-tauri/src/modules/migration.rs` (`extract_refresh_token_from_file`)
- Injection: `src-tauri/src/modules/db.rs` (`inject_token`)
- Protobuf helpers: `src-tauri/src/utils/protobuf.rs`

## Process control details

The manager tries hard to:
- correctly detect whether the external app is running,
- avoid killing its own process,
- avoid killing helper/renderer processes incorrectly.

Implementation:
- `src-tauri/src/modules/process.rs`

User configuration:
- `antigravity_executable` can be set in Settings to improve process matching and launching.

## User-facing behavior

From the user’s perspective:
- Switching account in Accounts screen or tray:
  - updates current account label,
  - updates external app credentials,
  - restarts external app if needed.

## Risks / caveats (for rewrite)

- DB schema/protobuf fields are implementation details of the external app and may change.
- Injection is inherently sensitive:
  - requires careful backups,
  - must avoid logging secrets,
  - should provide clear error messaging when DB format/path is unexpected.

