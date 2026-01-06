# Manual Verification Plan (v1.2.0)

This plan outlines the manual validation steps for all completed features in the current release.

## 1. Core Proxy & Configuration
| Feature | User Action | Expected Interface Result | Expected Backend Log |
| :--- | :--- | :--- | :--- |
| **Proxy Startup** | Launch App | Dashboard shows "System Status: Online" | `[Proxy] Starting on :7100` |
| **Terminal Setup** | Run `siberia set-env` | Output export commands | N/A (CLI) |
| **Legacy Import** | Start app with `~/antigravity_accounts.json` present | Accounts tab shows inherited accounts | `Imported legacy account: ...` |

## 2. Accounts Management
| Feature | User Action | Expected Interface Result | Expected Backend Log |
| :--- | :--- | :--- | :--- |
| **Add Account** | Use "Add Account" Dialog | New row appears in table | `[Accounts] Created account ...` |
| **List Accounts** | View Accounts Tab | Table populated with DB data | `SELECT * FROM accounts` |
| **Switch Account** | Click "Activate" on row | Status changes to "Active", others "Inactive" | `[Injector] Successfully injected token` |
| **Delete Account** | Click Trash Icon | Row disappears | `UPDATE accounts SET deleted_at ...` |
| **Quota Display** | View Account Details | Badges show "Free" or "Paid" tier | `[Quota] Calculated tier for ...` |

## 3. Traffic Inspection
| Feature | User Action | Expected Interface Result | Expected Backend Log |
| :--- | :--- | :--- | :--- |
| **Live Monitor** | Browse web via proxy (`:7100`) | New rows appear in Traffic Table | `[Proxy] Req: GET ...` |
| **Filtering** | Type "google" in filter bar | Only matching rows displayed | N/A (Frontend filter) |

## 4. Advanced Features
| Feature | User Action | Expected Interface Result | Expected Backend Log |
| :--- | :--- | :--- | :--- |
| **Map Local** | Add Rule `.*test.com` -> `local.json` | Request to test.com returns JSON content | `[MapLocal] Serving ...` |
| **Scripting** | Add Script `return {status: 418}` | Request returns "I'm a teapot" | `[Scripting] Executing hook ...` |

## 5. System Integration
| Feature | User Action | Expected Interface Result | Expected Backend Log |
| :--- | :--- | :--- | :--- |
| **Auto-Update** | Go to Settings -> Check Update | "You are up to date" (or update avail) | `[Updater] Checking for updates...` |
| **System Tray** | (Linux/Win) Right-click Tray Icon | Menu appears (Start/Stop) | `[Tray] OnReady` |

## Execution Status
- [ ] Core Proxy
- [ ] Accounts Management
- [ ] Traffic Inspection
- [ ] Advanced Features
- [ ] System Integration
