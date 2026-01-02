# Technology Stack Decision

**Date:** 2026-01-03
**Status:** Approved

## Core Stack

### Backend / System Layer: **Go (Golang)**
*   **Rationale:** 
    *   Ideally suited for high-performance networking and proxy management.
    *   Native concurrency (goroutines) allows handling thousands of connections efficiently.
    *   Compiles to a single, static binary for Windows, MacOS, Linux.
    *   Strong standard library and ecosystem.

### Desktop GUI Framework: **Wails (v2)**
*   **Rationale:**
    *   Provides a bridge between Go (backend) and Web technologies (frontend).
    *   Uses native system rendering engines (WebView2 on Windows, WebKit on Mac/Linux) instead of bundling Chromium (Electron).
    *   **Result:** Application size ~15-20MB (vs ~150MB+ for Electron) and significantly lower RAM usage.

### Frontend (UI): **React + TypeScript + Vite**
*   **Rationale:**
    *   Utilizes the user's existing familiarity with JavaScript/modern web stacks.
    *   Rich ecosystem of UI components (for the Account Manager interface).
    *   Fast development cycle with Vite (HMR).

## Compilation Strategy
*   **Cross-Compilation:** configurable via `GOOS` and `GOARCH` environment variables.
*   **Output:** Single executable file (possibly with an installer wrapper for Windows `.msi` or Mac `.dmg`).

## Next Steps
1.  Initialize Wails project structure.
2.  Define Functional Requirements (Proxy logic, Account Manager features).
