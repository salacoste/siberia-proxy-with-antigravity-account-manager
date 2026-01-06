# Epic-13: Cloud Infrastructure

**Goal:** Provide robust backend infrastructure for data syncing and file storage, moving beyond mock implementations to real, scalable services (Supabase & MinIO).

## Scope
*   **Authentication & Relational Data:** Integrate with Supabase (PostgreSQL + Auth) for user profiles and settings sync.
*   **File Storage:** Integrate with S3-compatible storage (MinIO) for persisting large artifacts like exported sessions (HAR files).

## Stories
*   [x] **Story-39:** Real Cloud Backend (Supabase Sync)
    *   *Status: Completed*
*   [x] **Story-40:** Real File Storage (Docker/MinIO)
    *   *Status: Completed*

## Status
**Status:** Completed

### Manual Verification (2026-01-06)
- **Method**: Browser Subagent
- **Steps**: Checked "Cloud Sync" section in Settings and `/cloud` page.
- **Result**: Login UI elements present.
- **Evidence**: `advanced_settings_validation_1767697033575.webp`
