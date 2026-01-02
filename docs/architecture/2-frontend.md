# Architecture Part 2: Frontend (Antigravity Shell)

## 1. Stack
*   **Core:** React 18+, TypeScript, Vite.
*   **Styling:** Tailwind CSS + `shadcn/ui` (or similar component library based on Radix).
*   **State:** `Zustand` for global manageable state.

## 2. Structure (`frontend/`)
*   **`src/lib/wails/`**: Dedicated bridge layer.
    *   `events.ts`: Strongly-typed event listeners wrapper.
    *   `commands.ts`: Wrappers for `window.go.main.App` calls.
*   **`src/stores/`**:
    *   `useAppStore`: Theme, Language, Init status.
    *   `useAuthStore`: Current User Token (if app had its own auth, but mostly config here).
*   **`src/components/`**:
    *   Atomic Design: `atoms`, `molecules`, `organisms`.
    *   **Layouts**: `RootLayout` (Navbar + Outlet).

## 3. Communication Pattern
*   **Command (Request/Response):** Frontend calls Backend -> Await Promise -> Result.
    *   Example: `AddAccount(details)` -> `AccountID`.
*   **Event (Push):** Backend emits -> Frontend subscribes.
    *   Example: `proxy:log` -> Update Monitor Table.
