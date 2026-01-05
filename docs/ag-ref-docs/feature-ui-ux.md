# Feature Analysis: UI/UX & Frontend Architecture

**Source:** `src/App.tsx`, `src/pages/`, `src/services/`, `src/stores/`
**Framework:** React + Vite + TailwindCSS + Zustand + Tauri API

## 1. Product Logic
The frontend serves as the control center for the Accounts Manager. It provides a clean, modern interface for users to:
*   **Visualize Accounts:** Grid/List views of all Google accounts.
*   **Key Actions:** Switch Identity (Process Injection), Refresh Quota, Toggle Proxy participation.
*   **Interactive Feedback:** Toast notifications for success/failure, Modal dialogs for confirmations.

## 2. Technical Implementation

### A. Routing (`App.tsx`)
*   **Router:** `react-router-dom` (BrowserRouter).
*   **Structure:**
    *   `/` -> `Layout` (Sidebar + Header + Content Area).
    *   `index` -> `Dashboard` (Overview).
    *   `accounts` -> `Accounts` (Main management page).
    *   `api-proxy` -> `ApiProxy` (Proxy configuration & logs).
    *   `monitor` -> `Monitor` (Real-time charts).
    *   `settings` -> `Settings` (App preferences).
*   **Events:** Listens for backend events (e.g., `tray://account-switched`) to auto-refresh data.

### B. State Management (Zustand)
*   **Store:** `stores/useAccountStore.ts`
*   **State:**
    *   `accounts`: Array of all accounts.
    *   `currentAccount`: The currently active/injected account.
    *   `loading`/`error`: UI states.
*   **Actions:**
    *   `fetchAccounts()`: Calls `accountService.listAccounts()`.
    *   `switchAccount(id)`: Invokes backend switch logic, then optimistic UI updates.
    *   `refreshQuota(id)`: Force refreshes an account's quota.
*   **Optimistic UI:** Used in `reorderAccounts` to respond instantly to drag-and-drop before the backend confirms.

### C. Services Layer (`services/`)
*   **Strict Bridge:** `services/accountService.ts` wraps `invoke` calls.
*   **Tauri Command Mapping:**
    *   `listAccounts` -> `invoke('list_accounts')`
    *   `switchAccount` -> `invoke('switch_account', { accountId })`
    *   `refreshQuota` -> `invoke('fetch_account_quota', { accountId })`
*   **Environment Check:** `ensureTauriEnvironment()` guards against running in a standard browser without Tauri injection.

### D. Views & Components
*   **`Accounts.tsx`:**
    *   **Dual View:** List Table vs. Card Grid.
    *   **Filtering:** Filter by Subscription Tier (Free/Pro/Ultra).
    *   **Pagination:** Dynamic client-side pagination based on window height/width.
    *   **Bulk Actions:** Batch Delete, Batch Enable/Disable Proxy.
*   **Modals:** Used for "Add Account", "Delete Confirmation", "Account Details".

## 3. Specifics & Edge Cases
*   **Responsive Pagination:** `ITEMS_PER_PAGE` is memoized and calculates available screen real estate (ResizeObserver) to prevent scrolling within the view if possible.
*   **I18n:** Fully integrated `react-i18next` for English/Chinese/Russian support.
*   **Tray Integration:** Bi-directional sync. If user switches account via System Tray, the UI listens for event and re-fetches state.
