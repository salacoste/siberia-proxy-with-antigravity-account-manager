# Story-39: Real Cloud Backend (Supabase Sync)

**Epic**: [Epic-13: Cloud Infrastructure](../epics/epic-13-cloud-infrastructure.md)
**Status**: Ready for Dev
**Feature Branch**: `antigravity/feat/cloud-backend`

## Goal
Replace the in-memory `MockServer` in `siberia/modules/sync` with a real client that talks to Supabase.

## Requirements
1.  **Auth**:
    -   User can sign up/sign in via the Settings UI.
    -   Backend stores the Session/JWT securel.
2.  **Storage (DB)**:
    -   Table: `profiles`
    -   Columns: `user_id` (UUID), `data_blob` (ByteA/Text), `updated_at` (Timestamp).
    -   RLS (Row Level Security): Users can only read/write their own rows.
3.  **Sync Logic**:
    -   Update `SyncManager` to use the Supabase client for `Push` and `Pull`.

## Technical Notes
-   Use `env` variables or build flags for Supabase URL/Key.
-   **Supabase Go-Lang Strategy**:
    -   There is no official Go SDK.
    -   Use `resty` or standard `net/http` to hit the REST API (`/auth/v1` for GoTrue, `/rest/v1` for PostgREST).
    -   **Security**: The "Master Password" encryption (Story-33) remains unchanged. We just store the *result* in Supabase.

## Acceptance Criteria
-   [ ] Login with real credentials works.
-   [ ] "Push" persists data across app restarts (checked via Supabase Dashboard or second device).
-   [ ] "Pull" retrieves the correct encrypted blob.


