# PRD Part 1: Global UI & Core App Structure

## 1. Core Application
*   **Platform:** Windows 10/11, MacOS (Intel/Silicon), Linux (Debian/Ubuntu).
*   **Technology:** Go (Backend) + Wails (Bridge) + React/Vite (Frontend).
*   **Single Binary:** Application must ship as a single executable (plus OS-specific bundle/installer).

## 2. Global Layout (Shell)
*   **Window Frame:** Custom title bar with drag regions on all platforms.
*   **Navbar (Permanent):**
    *   **Dashboard:** Home screen.
    *   **Accounts:** Identity management.
    *   **API Proxy:** Local proxy server config.
    *   **Settings:** Global config.
*   **Quick Actions (Navbar):**
    *   **Theme Toggle:** Light/Dark mode.
    *   **Language Toggle:** English / Chinese (Simplified).

## 3. System Integration
*   **Toast Notifications:** Global system for "Success", "Error", "Info" messages (e.g., "Proxy Started", "Account Saved").
*   **Background Tasks:**
    *   Quota Refresh Runner (Interval based).
    *   DB Sync Runner (Interval based).
