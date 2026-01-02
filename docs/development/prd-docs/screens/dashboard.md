# Screen: Dashboard

Route: `/`

Source:
- UI: `src/pages/Dashboard.tsx`
- Components:
  - `src/components/dashboard/CurrentAccount.tsx`
  - `src/components/dashboard/BestAccounts.tsx`
  - `src/components/accounts/AddAccountDialog.tsx`

## Purpose

Provide a fast “control plane” overview:
- how many accounts exist,
- average quota health for key models,
- quick actions: add account, refresh current quota, export, switch to best account.

## Data sources

- Zustand stores:
  - `useAccountStore` (`src/stores/useAccountStore.ts`)
- Backend commands:
  - `list_accounts`, `get_current_account`, `fetch_account_quota`, `save_text_file`

## Layout & sections

1) Header / greeting
- Uses current account name/email prefix when available.

2) Primary actions (right side)
- Add account (opens a multi-tab dialog).
- Refresh current account quota.

3) Stats cards (5)
- Total accounts.
- Avg quota % for Gemini Pro.
- Avg quota % for Gemini Image.
- Avg quota % for Claude.
- Count of “low quota” accounts (threshold-based).

4) Two panels
- Current account summary (quota + subscription tier).
- Best accounts suggestion (one for Gemini, one for Claude).

5) Quick links
- Navigate to Accounts screen.
- Export all accounts to JSON.

## Interactions (buttons)

| UI element | Intent | Backend/API behavior | Notes |
|---|---|---|---|
| “Add account” (dialog open) | Onboard new accounts | UI-only; actions inside dialog call account commands | See `docs/development/prd-docs/screens/accounts.md` “AddAccountDialog” section. |
| “Refresh quota” (current) | Update current account quota snapshot | Calls `fetch_account_quota(account_id)` and then refreshes current account state | UI disables button if no current account. |
| “Switch account” (inside CurrentAccount panel) | Go to Accounts list | UI navigation to `/accounts` | This button does not switch account directly. |
| “Switch to best” (inside BestAccounts panel) | Switch to recommended account | Calls `switch_account(account_id)` | Selection heuristic prefers higher combined Gemini score vs Claude score. |
| “View all accounts” | Jump to Accounts page | UI navigation to `/accounts` | — |
| “Export data” | Export all accounts to JSON | Uses native save dialog; writes via `save_text_file(path, content)` | Export includes emails + refresh tokens (sensitive). |

## Business logic notes

### “Avg quota” computation

Dashboard computes averages only across accounts where the model quota is present and > 0.
Key model names are compared case-insensitively in some places; exact keys differ between components:
- Dashboard uses a strict lookup for several hardcoded model ids (see `src/pages/Dashboard.tsx`).
- Tray uses strict lowercased matching (see `src-tauri/src/modules/tray.rs`).

Rewrite guidance:
- Keep a central “well-known model ids” mapping, because many UI elements assume exact ids.
