# Screen: Accounts

Route: `/accounts`

Source:
- UI: `src/pages/Accounts.tsx`
- Components:
  - `src/components/accounts/AccountTable.tsx`
  - `src/components/accounts/AccountGrid.tsx`
  - `src/components/accounts/AccountRow.tsx`
  - `src/components/accounts/AccountCard.tsx`
  - `src/components/accounts/AccountDetailsDialog.tsx`
  - `src/components/accounts/AddAccountDialog.tsx`

## Purpose

This screen is the primary management UI for the local account pool:
- add/import accounts,
- switch current account,
- refresh quota snapshots,
- export refresh tokens,
- delete accounts,
- view per-model quota details,
- filter/sort by health and subscription tier.

## Data sources

- `useAccountStore` (`src/stores/useAccountStore.ts`)
- `useConfigStore` for `default_export_path` and page size (`src/stores/useConfigStore.ts`)

Backend commands used via services:
- `list_accounts`, `get_current_account`
- `add_account`, `delete_account`, `delete_accounts`, `switch_account`
- `fetch_account_quota`, `refresh_all_quotas`
- OAuth flow: `prepare_oauth_url`, `start_oauth_login`, `complete_oauth_login`, `cancel_oauth_login`
- Imports: `import_v1_accounts`, `import_from_db`, `import_custom_db`, `sync_account_from_db`
- Export write: `save_text_file`

## Layout & sections

1) Top bar
- Search input
- View switch (list / grid)
- Filter chips:
  - All / Available / Low / Pro / Ultra / Free
- Action buttons:
  - Add account (dialog)
  - Delete selected (only if selection exists)
  - Refresh (all or selected; confirmation dialog)
  - Export (all or selected)

2) Accounts list
- List mode: table view with bulk checkboxes.
- Grid mode: card view.

3) Pagination (bottom)
- Auto page size based on container height, overrideable by user (persists to config key).

4) Modals
- Account details dialog (quota per model)
- Confirm delete (single/batch)
- Confirm refresh (all/selected)

## Interactions (top-level buttons)

| UI element | Intent | Behavior | Backend interaction |
|---|---|---|---|
| Search input | Narrow results | Filters by email/name/id (client-side) | None |
| View: list | Table view | Changes presentation, preserves filters | None |
| View: grid | Card view | Changes presentation, preserves filters | None |
| Filter chips | Focus by health/tier | Client-side filter, resets selection & page | None |
| “Add account” | Onboard | Opens AddAccountDialog | See below |
| “Delete selected” | Remove selected accounts | Opens confirmation modal, then bulk delete | `delete_accounts(accountIds)` |
| “Refresh” | Refresh quotas | Opens confirmation modal, then refresh selected or all | Selected: `fetch_account_quota` per id; All: `refresh_all_quotas()` |
| “Export” | Export tokens | Exports selected or all | `save_text_file(path, content)` |

## Controls (exhaustive)

### Top bar (left)

| Control | Type | State | Behavior | Backend |
|---|---|---|---|---|
| Search input | Text input | `searchQuery` | Client-side filter by email substring; resets selection & page when changed | None |
| View mode: list | Icon button | `viewMode = "list"` | Switches to table view; resets to page 1 | None |
| View mode: grid | Icon button | `viewMode = "grid"` | Switches to card view; resets to page 1 | None |

### Filter chips (center)

| Control | Type | State | Behavior | Backend |
|---|---|---|---|---|
| All / Available / Low / Pro / Ultra / Free | Buttons | `filter` | Client-side filter over already-searched results; clears selection & resets page | None |
| Count badges | Labels | derived | Shows counts per filter based on current search | None |

### Top bar (right actions)

| Control | Type | Visibility | Behavior | Backend |
|---|---|---|---|---|
| Add account | Button → dialog | Always | Opens `AddAccountDialog` | See dialog section |
| Delete selected | Button | Only when `selectedIds.size > 0` | Opens confirm modal then bulk delete | `delete_accounts(accountIds)` |
| Refresh | Button | Always | Opens confirm modal then refreshes selected (parallel) or all (serial) | Selected: `fetch_account_quota` per id; All: `refresh_all_quotas()` |
| Export | Button | Always | Exports selected or all as JSON | `save_text_file(path, content)` |

### List mode table controls

| Control | Intent | Behavior | Notes |
|---|---|---|---|
| Select all (header checkbox) | Batch actions | Selects/clears all rows on the current page | Implemented in `AccountTable` header |
| Row select checkbox | Batch actions | Adds/removes row id in `selectedIds` | Click does not trigger row actions |
| Row actions: Details | Inspect quota | Opens `AccountDetailsDialog` | UI-only |
| Row actions: Switch | Make current account | Calls switch; disabled when switching or account disabled | `switch_account(accountId)` |
| Row actions: Refresh | Update quota | Calls quota refresh; disabled while refreshing or account disabled | `fetch_account_quota(accountId)` (note: current UI calls it 3x; see quirks) |
| Row actions: Export | Export one | Exports this account only | `save_text_file(path, content)` |
| Row actions: Delete | Remove account | Opens confirm modal then deletes | `delete_account(accountId)` |

### Grid mode card controls

Card controls mirror list mode:
- checkbox select
- details / switch / refresh / export / delete buttons
- disabled/forbidden/subscription badges are visual-only

### Pagination controls

| Control | Type | Behavior | State source |
|---|---|---|---|
| Page navigation | Buttons | Changes `currentPage` | UI state |
| Page size selector | Select | Overrides computed page size for this session | `localPageSize` (does not persist) |
| Total items / page size display | Labels | Derived | From `filteredAccounts.length` and `ITEMS_PER_PAGE` | Derived |

## AddAccountDialog — controls (explicit)

The dialog is the primary onboarding UI and has multiple buttons per tab.

### OAuth tab

| Control | Intent | Behavior | Backend / events |
|---|---|---|---|
| Dialog open | Prepare OAuth URL early | Pre-calls `prepare_oauth_url` to show URL before starting | Emits `oauth-url-generated` (also UI sets URL directly) |
| Copy URL | Easy paste into browser | Copies URL to clipboard | None |
| Open URL (if present) | Start auth in browser | Opens URL externally | OS-level browser open (frontend plugin) |
| Complete (auto) | Save account after callback | Triggered by event; UI auto-calls completion | Event: `oauth-callback-received` → calls `complete_oauth_login` |
| Complete (manual) | Save account after user authorizes | User-initiated completion | `complete_oauth_login` |
| Cancel | Release local callback port | Cancels in-progress/prepared flow | `cancel_oauth_login` |
| Switch away from OAuth tab | Prevent port leak | Auto-cancel prepared flow | `cancel_oauth_login` |

### Refresh Token tab

| Control | Intent | Behavior | Backend |
|---|---|---|---|
| Submit / Add | Add accounts from input | Parses input (JSON array or regex), dedupes, adds sequentially | Calls `add_account(refresh_token)` for each |

### Import tab

| Control | Intent | Backend | Notes |
|---|---|---|---|
| Import from DB | Import current external account | `import_from_db` | Also sets imported account as current account. |
| Import from custom DB | Import from user-selected DB file | `import_custom_db(path)` | Uses file picker to choose DB path. |
| Import v1 accounts | Import legacy accounts | `import_v1_accounts` | Reads legacy directory under home. |

## Row/Card actions

Each account row/card offers:
- Select checkbox (for batch actions).
- Switch current account.
- Refresh quota (per account).
- Export single account.
- View details.
- Delete account.

Back-end calls:
- Switch: `switch_account(account_id)` (also injects tokens into external DB and restarts external app)
- Refresh: `fetch_account_quota(account_id)`
- Delete: `delete_account(account_id)`
- Export: `save_text_file(...)`

## AccountDetailsDialog

Displays per-model quota entries:
- model name
- remaining percentage
- reset time

Implementation:
- `src/components/accounts/AccountDetailsDialog.tsx`

## AddAccountDialog (multi-tab onboarding)

Source:
- `src/components/accounts/AddAccountDialog.tsx`

Tabs:
1) OAuth
2) Refresh Token
3) Import

### OAuth tab

Goal: obtain a refresh token via loopback OAuth.

Key mechanics:
- When dialog opens on OAuth tab, it pre-generates an auth URL by calling `prepare_oauth_url`.
  - The backend starts local listeners immediately to avoid “user authorized before clicking start” races.
- The UI listens for:
  - `oauth-url-generated` (receives URL)
  - `oauth-callback-received` (callback arrived; UI auto-completes if dialog is open)

Primary actions:
- Copy URL
- Open URL in browser (optional in UI)
- Complete OAuth (manual)
- Cancel OAuth (auto-cancel when switching tabs)

Backend commands:
- `prepare_oauth_url`, `complete_oauth_login`, `cancel_oauth_login`
- (Also `start_oauth_login` exists and is used depending on UI path)

Important edge case:
- Google may not return a refresh token if the user previously authorized the app. The backend returns a long error message instructing the user to revoke permissions and retry.

### Refresh Token tab (manual / batch)

Goal: add accounts from refresh tokens (single or many).

Input parsing behavior:
- If input is a JSON array like `[{"refresh_token":"1//..."}, ...]`, tokens are extracted.
- Otherwise, a regex extracts tokens matching `1//...` patterns.
- Tokens are deduplicated and added sequentially, with a small delay.

Backend command:
- `add_account(refresh_token=...)`

Security note:
- This feature handles secrets. The rewrite must ensure tokens never appear in logs, analytics, or crash reports.

### Import tab

Goal: import accounts from existing local sources:
- “Import from DB” (default external DB path)
- “Import from custom DB” (file picker)
- “Import v1 accounts” (legacy directory)

Backend commands:
- `import_from_db`
- `import_custom_db(path)`
- `import_v1_accounts`

## Known quirks (current behavior)

- Single-account refresh on the Accounts screen calls quota refresh three times in a row (`src/pages/Accounts.tsx`), likely unintended.
- Export includes refresh tokens in plaintext JSON (intended for power users; high risk if mishandled).

## External app integration (important)

Switching the current account is not a UI-only concept. The backend performs an “apply to external app” sequence:
- close external Antigravity process (if running),
- back up the external DB,
- inject tokens into `state.vscdb`,
- restart the external app.

Implementation:
- `src-tauri/src/modules/account.rs` (`switch_account`)
- DB injection: `src-tauri/src/modules/db.rs`
- Process control + detection: `src-tauri/src/modules/process.rs`
