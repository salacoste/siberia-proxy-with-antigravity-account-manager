# Feature Analysis: Reference Frontend

**Source:** `src/pages/`, `src/stores/`, `src/services/`

## 1. Overview
The frontend is built with React, TypeScript, Vite, and Tailwind CSS. It uses `zustand` for state management and `react-router-dom` for navigation. It interacts with the Rust backend via Tauri's IPC `invoke`.

## 2. Key Pages

### A. Accounts Page (`pages/Accounts.tsx`)
A comprehensive dashboard for managing Google accounts.
*   **Views:** supports both List (`AccountTable`) and Grid (`AccountGrid`) views.
*   **Filtering:** Filter by Subscription Tier (Pro, Ultra, Free) or Search by Email.
*   **Pagination:** Custom client-side pagination with dynamic page size calculation based on container height/width.
*   **Actions:**
    *   **Add Account:** Triggers OAuth flow.
    *   **Switch:** Activates an account for the "Switch Account" feature.
    *   **Refresh:** Updates quota information.
    *   **Toggle Proxy:** Enables/Disables an account for the Proxy Engine.
    *   **Export:** Exports account credentials to JSON.
    *   **Batch Operations:** Delete, Toggle Proxy, Refresh selected accounts.
    *   **Z.ai Integration:** Z.ai model configuration and custom mapping UI.

### B. API Proxy Page (`pages/ApiProxy.tsx`)
The control center for the local Proxy Engine.
*   **Status Indicators:** Shows Running/Stopped status, Port, and Active Account count.
*   **Configuration:**
    *   Port, Request Timeout, Auto-start.
    *   **Auth Mode:** Off, Strict, All Except Health, Auto (Lan Access).
*   **Code Generator:** Dynamically generates Python code snippets for:
    *   **OpenAI SDK:** Examples for Chat and Image Generation (DALL-E 3 compatible).
    *   **Anthropic SDK:** Example for using Claude models.
    *   **Google GenAI SDK:** Example for native Gemini calls.
*   **Z.ai Configuration:** UI for managing Z.ai model mappings and overrides.

### C. Monitoring Page (`pages/Monitor.tsx`)
A real-time dashboard for traffic analysis.
*   **Live Metrics:** Displays Request Per Second (RPS) and Latency graphs.
*   **Request Log:** Detailed table of recent requests with Status, Model, and Latency.
*   **Account Attribution:** Shows which specific account (masked) handled each request, helping debug quota/tier issues.


## 3. State Management (`stores/useAccountStore.ts`)
Uses `zustand` to manage global account state.
*   **Actions:** Wraps `accountService` calls.
*   **Optimistic Updates:** `reorderAccounts` implements optimistic UI updates for drag-and-drop reordering.
*   **OAuth Flow:** Manages the multi-step OAuth process (`startOAuthLogin` -> `completeOAuthLogin`).

## 4. Services (`services/accountService.ts`)
A clean abstraction layer over Tauri's `invoke` API.
*   **Error Handling:** Enhances raw error messages (especially for OAuth failures).
*   **Type Safety:** Typed return values matching the backend `Account` and `QuotaData` structs.
