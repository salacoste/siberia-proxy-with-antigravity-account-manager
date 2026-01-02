# Architecture Part 1: Backend (Siberia Core)

## 1. Responsibilities
The backend acts as the system coordinator, managing network I/O, persistence, and external process control.

## 2. Structure (`siberia/`)
*   **`main.go`**: App entry point. Configures Wails application, menus, and system tray.
*   **`proxy/`**: The heart of the application.
    *   Uses `github.com/gin-gonic/gin` or `axum` (Architecture defined Axum, sticking to it).
    *   **Handlers**: Logic for handling specific API endpoints (Anthropic/Gemini).
    *   **Middleware**: Interceptors for Auth, Logging, and Recording.
*   **`modules/`**: Discrete functional domains.
    *   **`account.go`**: Manages the Account pool and token lifecycle.
    *   **`injection.go`**: Handles the `state.vscdb` manipulation.
    *   **`tray.go`**: OS Tray logic.

## 3. Data Layer
*   **Technology:** `GORM` (SQLite).
*   **Files:**
    *   `accounts.db`: Encrypted store for user accounts.
    *   `proxy_logs.db`: High-volume write store for monitor traffic.
