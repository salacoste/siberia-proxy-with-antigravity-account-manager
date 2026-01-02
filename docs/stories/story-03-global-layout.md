# Story-03: Create Global Layout (Navbar, Routing)

**Epic:** [Epic-01: Core Foundation & UI Shell](./epic-01-core.md)
**Status:** Draft

## Goal
Establish the core UI shell of the application, including the permanent Navigation Bar, Client-Side Routing, and the root layout structure that holds the application together.

## Tasks
- [ ] **Task 1: Routing Setup**
    -   Install `react-router-dom`.
    -   Create route constants/enums.
    -   Configure `HashRouter` (preferred for Wails/Electron-like apps to avoid history API complexity with file:// protocol).
- [ ] **Task 2: Page Placeholders**
    -   Create basic components for:
        -   `Dashboard` (`/`)
        -   `Accounts` (`/accounts`)
        -   `Proxy` (`/proxy`)
        -   `Settings` (`/settings`)
- [ ] **Task 3: Navbar Component**
    -   Implement responsive sidebar or topbar (Design decision: Sidebar for desktop app).
    -   Include Navigation Links (Dashboard, Accounts, Proxy, Settings).
    -   Include connection status indicator (from Proxy Store - future hook).
- [ ] **Task 4: Root Layout**
    -   Implement `RootLayout` component wrapping the `Outlet`.
    -   Apply global styles/theme container.

## Acceptance Criteria
- [ ] App launches to Dashboard by default.
- [ ] Clicking Nav links switches views instantly without reload.
- [ ] Active link state is visible in Navbar.
- [ ] Layout is responsive and fills the Wails window.

## Technical Notes
- **Library:** `react-router-dom` v6+.
- **Router:** Use `HashRouter` to ensure compatibility with Wails asset serving.
- **Styling:** Tailwind CSS.
