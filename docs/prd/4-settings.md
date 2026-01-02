# PRD Part 4: Settings & Configuration

## 1. Config Persistence
*   **Storage:** JSON file (`config.json`), managed by Backend, synced to Frontend store.
*   **Optimistic UI:** UI updates immediately; reverts on save failure.

## 2. General Settings
*   **Auto-Launch:** Toggle OS startup registration.
*   **Language/Theme:** (See Global UI).

## 3. Automation Settings
*   **Account Quota Refresh:**
    *   Enable/Disable auto-refresh.
    *   Interval (minutes).
    *   Logic: Frontend timer triggers backend `refresh_all_quotas`.
*   **DB Sync:**
    *   Enable/Disable.
    *   Interval (seconds).
    *   Logic: Syncs local state with external DB changes.

## 4. Advanced Settings
*   **Exports:** Default export path selection.
*   **Paths:**
    *   Data Directory (Read-only view + "Open in Explorer" button).
    *   **External Executable:** Path to the target application (VS Code / Browser) for "Switch Account" functionality. Auto-detect capability.
*   **Logs:**
    *   Clear App Logs (Cache).

## 5. Updates
*   **Check for Updates:** GitHub API integration to check latest release.
*   **Download:** Open browser to download page.
