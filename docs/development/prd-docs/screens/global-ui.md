# Global UI (Layout + Navbar)

Source:
- Layout: `src/components/layout/Layout.tsx`
- Navbar: `src/components/layout/Navbar.tsx`
- Theme manager: `src/components/common/ThemeManager.tsx`
- Background tasks: `src/components/common/BackgroundTaskRunner.tsx`
- Toasts: `src/components/common/ToastContainer.tsx`

## Purpose

Defines app-wide UX and global controls:
- navigation across screens,
- global theme/language quick toggles,
- background automation runners,
- global toast notifications,
- window drag regions.

## Navigation

Navbar links:
- Dashboard (`/`)
- Accounts (`/accounts`)
- API Proxy (`/api-proxy`)
- Settings (`/settings`)

Routing:
- `src/App.tsx` uses `createBrowserRouter`.

## Controls (exhaustive)

### Navbar (center pill)

| Control | Type | Intent | Behavior | State source |
|---|---|---|---|---|
| Dashboard | Nav link | Go to overview | Navigates to `/` | `react-router` |
| Accounts | Nav link | Manage accounts | Navigates to `/accounts` | `react-router` |
| API Proxy | Nav link | Configure proxy | Navigates to `/api-proxy` | `react-router` |
| Settings | Nav link | App preferences | Navigates to `/settings` | `react-router` |

### Navbar (right side)

| Control | Type | Intent | Behavior | Backend |
|---|---|---|---|---|
| Theme toggle | Button | Switch theme quickly | Toggles `config.theme` (`light` ↔ `dark`) and persists | `save_config` |
| Language toggle | Button | Switch UI language quickly | Toggles `config.language` (`zh` ↔ `en`), persists, calls `i18n.changeLanguage` | `save_config` |

## Theme toggle

Navbar has a theme toggle button:
- Updates `config.theme` (`light` ↔ `dark`) and persists via `save_config`.
- If the browser supports View Transition API, animates the transition.

Backend:
- Uses `save_config` (via config store).

## Language toggle

Navbar has a language toggle button:
- Switches `config.language` (`zh` ↔ `en`)
- Persists via `save_config`
- Calls `i18n.changeLanguage` immediately.

Backend:
- Uses `save_config`.

## BackgroundTaskRunner

Always mounted at app level.
Runs periodic background tasks based on config:
- auto refresh quotas
- auto sync from DB

Details:
- `docs/development/prd-docs/screens/settings.md` (Account tab)

## Toast system

Global toast component is always mounted and used by pages for success/error notifications.

## Window drag regions

The Layout renders always-on-top “drag regions” for desktop window dragging.
- Primary drag overlay: `src/components/layout/Layout.tsx`
- Secondary overlay above navbar: `src/components/layout/Navbar.tsx`

## Diagnostic / debug UI

There is a “Network Monitor” UI (`src/components/common/NetworkMonitor.tsx`) backed by `src/stores/networkMonitorStore.ts`,
but it is not currently mounted in the Layout, so it is effectively unused in production UI.
