# Product Requirements Document (PRD)
**Project:** Siberia Proxy with Antigravity Account Manager
**Version:** 1.0
**Source of Truth:** Legacy Docs Analysis (docs/development/prd-docs)

## 1. Executive Summary

A unified, cross-platform desktop application (Go backend + Wails/React frontend) that integrates two core functions:
1.  **Siberia Proxy:** A local proxy service that routes traffic to upstream proxies, supports provider integrations (z.ai, MCP), and manages network identity.
2.  **Antigravity Account Manager:** An identity management system for secure storage and automated handling of accounts (Google, etc.), including OAuth flows and browser isolation.

## 2. Technical Architecture (Confirmed)
*   **Backend:** Go (Golang)
    *   Handles: Network I/O, Proxy Server, SQLite DB, File System, Sub-process management (Browsers).
*   **Frontend:** React + TypeScript + Vite
    *   Handles: UI, State Management (Zustand), optimistic updates, localized interactions.
*   **Bridge:** Wails v2
    *   Compiles to single binary (Windows/macOS/Linux).

## 3. Detailed Functional Requirements

### 3.1 Global UI & Navigation
*   **Layout:**
    *   Permanent **Navbar** with navigation links: Dashboard, Accounts, API Proxy, Settings.
    *   **Theme Toggle:** Fast switch (Light/Dark) with persistence.
    *   **Language Toggle:** Fast switch (EN/ZH).
    *   **Status Indicators:** Network connectivity, background task status (quota refresh).
*   **System Integration:**
    *   Native window drag regions.
    *   Toast notification system for global alerts (Success/Error).

### 3.2 Screen: Accounts Manager (`/accounts`)
The central hub for managing digital identities.
*   **Data Structure:**
    *   Account ID, Email/Username, Password, 2FA Secret.
    *   **State:** Health (Available/Low/Pro/Ultra/Free), Quota details.
    *   **Tokens:** Refresh tokens (securely stored), Cookies.
*   **Views:** Switchable **Grid (Card)** and **List (Table)** views.
*   **Filtering:**
    *   Client-side search (Email/Name).
    *   Quick Filters: By Health Status (Available, Pro, etc.).
*   **Actions:**
    *   **Import:**
        *   OAuth (Local loopback flow for acquiring tokens).
        *   Manual Entry (Refresh Token list paste).
        *   Legacy DB Import (`v1`).
    *   **Batch Operations:** Select multiple -> Refresh Quota, Delete, Export.
    *   **External App Switch:** "Switch Account" button effectively restarts the external target application with the selected account's tokens injected.

### 3.3 Screen: API Proxy (`/api-proxy`)
Configuration interface for the local proxy server.
*   **Core Control:**
    *   Start/Stop button for the `Siberia` proxy service.
    *   Port selection (default: 8000-65535).
    *   **Auth Mode:** `Off`, `Strict` (Bearer token required), `Auto`.
    *   **Allow LAN:** Toggle to bind 0.0.0.0 vs 127.0.0.1.
*   **Upstream Providers:**
    *   **z.ai Integration:**
        *   Enable/Disable.
        *   Dispatch Modes: `Exclusive`, `Pooled`, `Fallback` (Google pool backup).
        *   Model Mapping: "Exact" (e.g., `gpt-4` -> `gemini-pro`) and "Family" defaults.
    *   **MCP (Model Context Protocol):**
        *   Enable/Disable MCP server endpoints.
        *   Features: Web Search, Web Reader, Vision.
*   **Developer Tools:**
    *   **Code Snippets:** Auto-generated Python/cURL examples for selected protocols (OpenAI/Anthropic/Gemini).
    *   **Monitor Button:** Quick jump to traffic inspector.

### 3.4 Screen: Proxy Monitor (`/monitor`)
A debugging tool for network traffic.
*   **Recorder:** Toggle `Enable Logging`.
    *   *Security Note:* Payloads are only captured if explicitly enabled by user.
*   **Inspector:**
    *   Table view of recent requests (Method, URL, Status, Duration).
    *   Details Modal: Request/Response Headers and Bodies (JSON pretty-print).
*   **Filtering:** Real-time search/filter logs.

### 3.5 Screen: Settings (`/settings`)
Global configuration persistence.
*   **General:** Auto-launch on startup, Language, Theme.
*   **Account:** Authorization to automate background tasks (Auto-refresh quotas interval, Sync interval).
*   **Proxy:** "Upstream Proxy" configuration (for the app's own outbound traffic, unrelated to the forwarded traffic).
*   **Advanced:** Data directory paths, External executable paths (for the "Switch Account" feature).

## 4. Key Workflows (The "Switch Account" Logic)
A critical legacy feature documented is the "External App Integration":
1.  User clicks "Switch" on an Account.
2.  **Siberia Backend:**
    *   Identifies running external process (e.g., VS Code or Antigravity specialized browser).
    *   Terminates the process.
    *   Injects the selected account's authentication state (tokens/db) into the external app's storage.
    *   Restarts the external process.

## 5. Security Requirements
*   **Credentials:** Refresh tokens and passwords must be stored encrypted at rest (e.g., utilizing OS keychain or encrypted SQLite).
*   **Logs:** Access logs must NOT contain request bodies by default. Reviewing payloads requires explicit user opt-in (Monitor).
*   **Updates:** Built-in update checker (GitHub Releases).
