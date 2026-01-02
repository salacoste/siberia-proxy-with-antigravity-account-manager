# Screen: Settings

Route: `/settings`

Source:
- UI: `src/pages/Settings.tsx`
- Config store: `src/stores/useConfigStore.ts`
- Background tasks: `src/components/common/BackgroundTaskRunner.tsx`

## Purpose

Provide global app settings:
- language/theme,
- background automation (quota refresh, DB sync),
- external path preferences (export path, executable path),
- proxy upstream proxy settings (outbound),
- logs and diagnostics,
- updates and “about” metadata.

## Layout

Top:
- tab strip (General / Account / Proxy / Advanced / About)
- “Save” button (persists all settings)

Tabs:
1) General
2) Account
3) Proxy
4) Advanced
5) About

## Interactions (global)

| UI element | Intent | Behavior | Backend interaction |
|---|---|---|---|
| Save | Persist all settings | Calls store `saveConfig(formData)` | `save_config(config)` |
| Tab buttons (General/Account/Proxy/Advanced/About) | Switch settings category | Changes visible form section, does not persist by itself | None |

## Tab: General

### Language
- Select: `config.language` (`zh` or `en`)
- Note: Navbar also exposes a language toggle (quick switch) which updates the same key.

### Theme
- Select: `config.theme` (`light` / `dark` / `system`)
- Note: Navbar also exposes a theme toggle (quick switch) which updates the same key.

### Auto-launch (system startup)
- Select: enabled/disabled
- Updates:
  - OS auto-start via command
  - config key `auto_launch` for display/state

Backend:
- `toggle_auto_launch(enable: bool)`

### Controls (General tab)

| Control | Type | Config key(s) | Behavior | Backend |
|---|---|---|---|---|
| Language | Select | `language` | Updates `formData.language` | Persisted on Save |
| Theme | Select | `theme` | Updates `formData.theme` | Persisted on Save |
| Auto-launch | Select | `auto_launch` | Immediately calls OS toggle and updates `formData.auto_launch` | `toggle_auto_launch(enable)` |

## Tab: Account (automation)

### Auto refresh quotas
- Toggle: `config.auto_refresh`
- When enabled:
  - Executes immediately once.
  - Then repeats every `config.refresh_interval` minutes.

Runtime implementation:
- Frontend timer in `src/components/common/BackgroundTaskRunner.tsx`
- Uses backend `refresh_all_quotas` through the account store.

### Refresh interval
- Input/select: `config.refresh_interval` (minutes)
- Visible only when auto refresh is enabled.

### Auto sync current account from DB
- Toggle: `config.auto_sync`
- When enabled:
  - Executes immediately once.
  - Then repeats every `config.sync_interval` seconds.

Runtime implementation:
- Frontend timer in `src/components/common/BackgroundTaskRunner.tsx`
- Uses backend `sync_account_from_db`.

### Sync interval
- Input/select: `config.sync_interval` (seconds)
- Visible only when auto sync is enabled.

### Controls (Account tab)

| Control | Type | Config key(s) | Behavior | Backend |
|---|---|---|---|---|
| Auto refresh quotas | Toggle | `auto_refresh` | Updates `formData.auto_refresh` | Persisted on Save |
| Refresh interval | Number input | `refresh_interval` | Visible only when auto refresh enabled | Persisted on Save |
| Auto sync current account from DB | Toggle | `auto_sync` | Updates `formData.auto_sync` | Persisted on Save |
| Sync interval | Number input | `sync_interval` | Visible only when auto sync enabled | Persisted on Save |

## Tab: Proxy (outbound upstream proxy)

This tab configures an **outbound proxy** used by backend HTTP clients for upstream calls:
- Google calls (quota, internal calls),
- z.ai calls,
- remote MCP calls (if enabled).

Config keys:
- `proxy.upstream_proxy.enabled`
- `proxy.upstream_proxy.url`

Backend usage:
- Wired through proxy startup (`start_proxy_service`) and z.ai model fetch (`fetch_zai_models`).

### Controls (Proxy tab)

| Control | Type | Config key(s) | Behavior | Backend |
|---|---|---|---|---|
| Upstream proxy enabled | Toggle | `proxy.upstream_proxy.enabled` | Updates `formData.proxy.upstream_proxy.enabled` | Persisted on Save |
| Upstream proxy URL | Text input | `proxy.upstream_proxy.url` | Updates `formData.proxy.upstream_proxy.url` | Persisted on Save |

## Tab: Advanced

### Default export path
- Key: `default_export_path`
- Used by the Accounts screen export:
  - if set, export will write into this directory without opening a save dialog.

### Data directory path + open button
- Shows absolute data directory path from backend:
  - `get_data_dir_path()`
- Open button:
  - `open_data_folder()`

### External executable path
- Key: `antigravity_executable`
- Buttons:
  - Select file (native file picker)
  - Detect automatically (`get_antigravity_path(bypassConfig=true)`)

### Logs (app logs, not proxy monitor DB)
- Clear logs button opens confirm modal:
  - backend command: `clear_log_cache()`
- Note: proxy monitor history has its own “clear” in the Monitor screen.

### Controls (Advanced tab)

| Control | Type | Config key(s) | Behavior | Backend |
|---|---|---|---|---|
| Default export path | “Choose folder” button + display | `default_export_path` | Uses native directory picker; sets/clears value | Persisted on Save |
| Clear export path | Button | `default_export_path` | Sets to `undefined` | Persisted on Save |
| Data directory path | Read-only text | — | Loaded on mount | `get_data_dir_path()` |
| Open data folder | Button | — | Opens data directory in OS file explorer | `open_data_folder()` |
| External executable path | Text input | `antigravity_executable` | Edits path manually | Persisted on Save |
| Select executable path | Button | `antigravity_executable` | File picker → sets path | Persisted on Save |
| Detect executable path | Button | `antigravity_executable` | Auto-detect and set path | `get_antigravity_path(bypassConfig=true)` |
| Clear executable path | Button | `antigravity_executable` | Sets to `undefined` | Persisted on Save |
| Clear logs | Button + confirm modal | — | Opens confirm modal → clears app log files | `clear_log_cache()` |

## Tab: About

Contains:
- Version display (hard-coded UI string in current code)
- Update check button:
  - `check_for_updates()`
- External links (repository)

### Controls (About tab)

| Control | Type | Intent | Behavior | Backend |
|---|---|---|---|---|
| Check updates | Button | Check GitHub release | Updates in-page status; may reveal download URL | `check_for_updates()` |
| Download link (if update) | Link/button | Install newer version | Opens external URL | Browser/OS open |
| Repository link | Link/button | Open project page | Opens external URL | Browser/OS open |

## Design notes / rationale

- The Settings page edits a full in-memory copy of `AppConfig` and persists via “Save”.
- Some actions also invoke backend operations immediately (e.g. auto-launch toggle).
- Automation is implemented in the frontend (timers) rather than a backend scheduler, so it runs only when the UI is running.
