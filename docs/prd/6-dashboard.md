# PRD Part 6: Dashboard

## 1. Overview
The Dashboard (`/`) is the landing screen providing a "Control Plane" view of the system's health.

## 2. Stats & Metrics
*   **Cards:**
    *   **Total Accounts:** Count of managed identities.
    *   **Quota Health:** Average remaining quota % for key models (Gemini Pro, Gemini Image, Claude).
    *   **Low Quota Alert:** Count of accounts below defined threshold.

## 3. Suggestions & Actions
*   **Current Account Panel:** Displays currently active identity, Tier, and Quota.
*   **Best Accounts:** Smart suggestion for "Best Available" account for Gemini and Claude (based on quota health).
    *   **Action:** "Switch to Best" button (Switches external app to this account).
*   **Quick Actions:**
    *   Add Account (Dialog).
    *   Refresh Current Quota.
    *   Export Data.

## 4. Business Logic
*   **Averages:** Computed only across accounts where quota > 0.
*   **Model IDs:** Dashboard uses specific hardcoded IDs for "Key Models" (needs config mapping).
