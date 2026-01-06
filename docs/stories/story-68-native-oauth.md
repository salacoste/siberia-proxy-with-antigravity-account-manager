# Story-68: Native OAuth Loopback

**Epic:** [Epic-22: Native Authentication](./epic-22-native-auth.md) (New Epic)
**Status**: Draft
**Priority**: Medium
**Basis**: `modules/oauth_server.rs`

## Goal
Implement a proper OAuth 2.0 Authorization Code flow with a local loopback server, replacing the manual "copy-paste token" method.

## Tasks
- [ ] **Task 1: Loopback Server**
    - File: `apps/backend/oauth/server.go`
    - Listen on ephemeral port (check IPv4/IPv6).
    - Serve Success/Fail HTML pages.
- [ ] **Task 2: Flow Orchestration**
    - Generate Auth URL with `redirect_uri=http://localhost:PORT/callback`.
    - Handle `code` exchange for tokens.
    - Emit events: `oauth-url-generated`, `oauth-callback-received`.

## Acceptance Criteria
- [ ] User clicks "Add Account" -> Browser opens.
- [ ] User authorizes -> Browser redirects to localhost -> "Success" page.
- [ ] App automatically receives token and adds account.
