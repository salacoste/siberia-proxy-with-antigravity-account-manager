# PRD Part 7: System Tray

## 1. Overview
A persistent OS-level tray icon providing minimal control without opening the main window.

## 2. Menu Structure
1.  **Current Account Label:** (Read-only) e.g., "demo@gmail.com".
2.  **Quota Summary:** Condensed lines showing "High/Image/Claude" remaining %.
    *   *Logic:* Strict model ID matching. Shows "Forbidden" or "Unknown" placeholders if data missing.
3.  **Separator**
4.  **Actions:**
    *   **Switch Next:** Cycle to the next available account in the list.
    *   **Refresh Current:** Trigger quota refresh for active account.
5.  **Separator**
6.  **System:**
    *   **Show Window:** Bring main app to front.
    *   **Quit:** Terminate app (and background proxy).

## 3. Localization
*   Tray menu items must support dynamic localization (i18n) based on global app usage.
