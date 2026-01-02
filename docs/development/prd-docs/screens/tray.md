# Feature: System Tray

Source:
- `src-tauri/src/modules/tray.rs`
- Frontend listeners:
  - `src/App.tsx`

## Purpose

Provide a minimal always-available control surface:
- show the main window,
- switch to the next account quickly,
- refresh quota for the current account,
- and display a compact quota snapshot.

## Tray menu structure (conceptual)

1) Current account label (read-only)
2) Quota summary lines (read-only)
3) Separator
4) Actions
   - Switch next
   - Refresh current
5) Separator
6) System
   - Show window
   - Quit

Texts are localized using a dedicated tray i18n map:
- `src-tauri/src/modules/i18n.rs`

## Interactions (menu items)

| Menu item | Intent | Backend behavior | Frontend behavior |
|---|---|---|---|
| Show window | Bring UI to front | Shows main window and focuses it | None (backend-only) |
| Quit | Exit app | `app.exit(0)` | None |
| Refresh current | Refresh quota for current account | Loads current account, calls `fetch_quota_with_retry`, persists quota | Emits `tray://refresh-current` for UI to refresh list/current |
| Switch next | Cycle current account | Loads account list, chooses next by index, calls `switch_account` (includes external DB token injection + external app restart) | Emits `tray://account-switched` for UI to refresh list/current |

## Quota summary logic

The tray shows a fixed set of model lines (strict matching):
- Gemini High (e.g. `gemini-3-pro-high`)
- Gemini Image (e.g. `gemini-3-pro-image`)
- Claude (e.g. `claude-sonnet-4-5`)

If account is forbidden:
- shows a “forbidden” indicator

If quota is missing:
- shows an “unknown quota” placeholder

Implementation:
- `update_tray_menus(...)` in `src-tauri/src/modules/tray.rs`

## Event contract

Backend emits:
- `tray://account-switched`
- `tray://refresh-current`

Frontend listens and triggers:
- `fetchCurrentAccount()`
- `fetchAccounts()`

Implementation:
- `src/App.tsx`
