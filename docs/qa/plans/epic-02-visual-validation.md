# Epic-02 Visual Validation Plan

**Goal:** Verify the visual integrity and basic interactivity of the Siberia Web Interface (Frontend).
**Environment:** `npm run dev` (Port 7101).
**Note:** Backend calls will use Mock Data (`useConfigStore` mock).

## Scope

### 1. Dashboard Page (`/`)
-   **Check:** Layout, "Overview" cards, Recent Activity placeholder.
-   **Action:** Verify navigation links are present.

### 2. Accounts Page (`/accounts`)
-   **Check:** Account List layout.
-   **Action:**
    -   Verify "Add Account" button exists.
    -   Click "Add Account" to open Modal (Visual check).
    -   Verify "Activate" toggle/button state.

### 3. Monitor Page (`/monitor`)
-   **Check:** Traffic Inspector layout.
-   **Action:** Verify table headers and "Listening..." state (visual).

### 4. Proxy Page (`/proxy`) (Re-verify)
-   **Check:** Snippet Tabs and Terminal Config.

### 5. Settings Page (`/settings`)
-   **Check:** General settings layout (Theme toggle, etc.).

## Execution Steps
1.  Start Frontend Server (`npm run dev -- --port 7101`).
2.  Browser Agent navigates to each page.
3.  Capture Screenshot of each view.
4.  Compile findings into Walkthrough.
