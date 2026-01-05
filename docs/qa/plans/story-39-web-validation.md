# QA Validation Plan: Siberia Proxy Web Interface

**Scope:**
-   **URL:** `http://localhost:7101` (Dev Server)
-   **Focus:** Story-39 features (Snippet Generator) & General Layout.

## Test Scenarios

### 1. Startup & Loading
-   **Expected Result:** "API Proxy" dashboard loads.
-   **Validation:**
    -   Page Title is present.
    -   Cards (Terminal Config, Connect to Siberia) are visible.

### 2. Snippet Generator (Story-39)
-   **Expected Result:** Code snippets update dynamically based on language selection.
-   **Steps:**
    1.  Navigate to "Connect to Siberia" card.
    2.  Click "Python" tab -> Verify Python code shown.
    3.  Click "Node.js" tab -> Verify Node.js code shown.
    4.  Click "Curl" tab -> Verify Curl command shown.
    5.  Click ".env" tab -> Verify .env format.
    6.  **Verify Port**: Code should show `http://localhost:8888/v1` (or configured port).
-   **Action:** Click "Copy" button (if verifiable via clipboard, otherwise visual check).

### 3. Terminal Configuration
-   **Expected Result:** Terminal export commands are shown.
-   **Steps:**
    1.  Check "Bash/Zsh" tab -> Verify `export HTTP_PROXY=...`.
    2.  Check "PowerShell" tab -> Verify `$env:HTTP_PROXY=...`.

## Execution strategy
1.  Launch Frontend Dev Server (`npm run dev -- --port 7101`).
2.  Use Browser Tool to visit `http://localhost:7101`.
3.  Capture Screenshots of the interface.
4.  Report findings.

> [!NOTE]
> Since this is running in a browser verify (without Wails), backend calls (`GetAppConfig`, `EventsOn`) typically fail or return mock data if implemented. We are validating the **Presentation Layer logic**.
