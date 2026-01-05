# Epic-10: Team Collaboration & Cloud Sync

**Goal:** Enable users to synchronize their settings across devices and share proxy configurations/sessions with teammates.
**Status:** Completed

## Scope
*   **Cloud Profiles:** Encrypted sync of `AppConfig` and `Accounts` (excluding sensitive secrets if preferred, or using E2EE).
*   **Team Sharing:** "Share Link" functionality for sending a captured session or a proxy group config.
*   **Backend:** Introduction of a lightweight remote backend (e.g., Firebase, Supabase, or custom Go server).

## Stories
*   [x] **Story-33:** Design & Implement Cloud Profile Sync
*   [x] **Story-34:** Implement "Share Session" Feature
