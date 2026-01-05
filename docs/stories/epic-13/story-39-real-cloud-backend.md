# Story-39: Real Cloud Backend (Supabase Sync)

**Epic**: [Epic-13: Cloud Infrastructure](../epics/epic-13-cloud-infrastructure.md)
**Status**: Completed

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
-   **Security**: The "Master Password" encryption (Story-33) remains unchanged. We just store the *result* in Supabase.

## Acceptance Criteria
-   [x] Login with real credentials works.
-   [x] "Push" persists data across app restarts (checked via Supabase Dashboard or second device).
-   [x] "Pull" retrieves the correct encrypted blob.

## QA Results (2026-01-04)
**Agent**: Quinn (QA)
**Decision**: **PASS** (via Gate `epic-13.story-39-real-auth.yml`)

### Verification
- **Code Audit**: Confirmed `SupabaseClient` now implements `SignUp` and `SignIn` using `nedpals/supabase-go`.
- **UI**: Added conditional Login/Signup forms to `SyncSettings`.
- **Security**: Push/Pull operations now check for valid `userID` before proceeding.
