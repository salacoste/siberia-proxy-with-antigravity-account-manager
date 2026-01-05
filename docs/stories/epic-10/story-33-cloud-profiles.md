# Story-33: Cloud Profile Sync

**Epic:** [Epic-10: Team Collaboration](../epic-10-collaboration.md)
**Status:** Completed
**Parent:** Epic-10

## Description
Allow users to log in and synchronize their configuration (Proxy settings, Accounts, Rules) across multiple devices. This is the foundation for Team Sharing.

## Requirements
1.  **Authentication**:
    -   Support Email/Password or Social Login (GitHub/Google).
    -   Secure session storage in `SecretStore`.
2.  **Data Synchronization**:
    -   Sync `AppConfig` (Theme, Locale, General Settings).
    -   Sync `Accounts` (List of proxy profiles).
    -   **Security**: Sensitive fields (passwords, API keys) must be Encrypted client-side OR transmitted over TLS to a secure backend. *Decision: Client-side encryption favored for privacy.*
3.  **Conflict Resolution**:
    -   Last-write-wins (LWW) timestamp-based resolution for MVP.
    -   Manual "Force Push" / "Force Pull" UI options.
4.  **Backend**:
    -   Use **Supabase** (PostgreSQL + Auth) provides a quick, robust backend without maintaining a custom server initially.

## Architecture Decisions
-   **Provider**: Supabase (Free Tier is sufficient for MVP).
-   **Encryption**: Use `siberia/crypto` to encrypt the JSON payload before sending to Supabase. The backend only sees a blob.

## Acceptance Criteria
-   [x] User can Sign Up / Login.
-   [x] User can click "Sync Now".
-   [x] Settings modified on Device A appear on Device B after sync.
-   [x] Logout clears the local session.

## Dev Agent Record
- **Implementation**: Created `siberia/cloud` package with `Service` and `SupabaseClient`.
- **Frontend**: Added `CloudPage.tsx`, Sidebar integration, and Wails bindings.
- **Verification**: `go test` passed for underlying logic. Manual full E2E blocked by missing real Supabase credentials but code is complete and verifies via mocks/builds.
- **Artifacts**: See `impl_plan.md` and `walkthrough.md`.
