# Story-68: Native OAuth Loopback

**Epic:** [Epic-22: Native Authentication](./epic-22-native-auth.md) (New Epic)
**Status**: Completed
**Priority**: Medium
**Basis**: `modules/oauth_server.rs`

## Goal
Implement a proper OAuth 2.0 Authorization Code flow with a local loopback server, replacing the manual "copy-paste token" method.

## Tasks
- [x] **Task 1: Loopback Server**
    - File: `apps/backend/oauth/server.go`
    - Listen on ephemeral port (check IPv4/IPv6).
    - Serve Success/Fail HTML pages.
- [x] **Task 2: Flow Orchestration**
    - Generate Auth URL with `redirect_uri=http://localhost:PORT/callback`.
    - Handle `code` exchange for tokens.
    - Emit events: `oauth-url-generated`, `oauth-callback-received`.

## Acceptance Criteria
- [x] User clicks "Add Account" -> Browser opens.
- [x] User authorizes -> Browser redirects to localhost -> "Success" page.
- [x] App automatically receives token and adds account.

## Dev Agent Record
### File List
- `apps/backend/siberia/oauth/server.go` [NEW]
- `apps/backend/siberia/oauth/server_test.go` [NEW]
- `apps/backend/app.go` [MODIFY]
- `apps/frontend/src/components/accounts/AddAccountDialog.tsx` [MODIFY]

### Change Log
1.  **Backend**: Implemented `oauth.Server` with channel-based synchronization and state validation.
2.  **App Wiring**: Added `LoginWithOAuth` to `App` struct to orchestrate the flow (Start Server -> Open Browser -> Wait).
3.  **Frontend**: Added "Add with Google" button to `AddAccountDialog`.
4.  **Tech Debt**: Fixed lint error in `mcp/server.go` (unused ctx).
### QA Results
- **Date**: 2026-01-07
- **Evaluator**: Quinn (QA)
- **Status**: PASS
- **Gate Artifact**: `docs/qa/gates/epic-22.story-68-native-oauth.yml`
- **Notes**: Excellent improvement to UX and Security. Loopback server implementation is clean and follows best practices for native apps (RFC 8252).
