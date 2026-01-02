# PRD Docs (Implementation-Matched)

This folder contains a **high-fidelity PRD/spec** for another team to re-implement the current Antigravity Manager functionality in a different project.

Scope principles:
- Match the **current behavior** (including quirks) unless explicitly marked as “recommended improvement”.
- Describe **screens, actions, storage, background tasks, proxy behavior**, and how everything fits together.
- Do **not** include any real tokens, real emails, logs, or user-specific paths in examples.

## Index

### Overview
- `docs/development/prd-docs/prd.md`
- `docs/development/prd-docs/architecture.md`
- `docs/development/prd-docs/proxy.md`
- `docs/development/prd-docs/glossary.md`
- `docs/development/prd-docs/external-integration.md`

### Storage & Data
- `docs/development/prd-docs/data/storage.md`

### APIs (app internal)
- `docs/development/prd-docs/apis/tauri-commands.md` — command surface used by the UI (and additional commands available in backend).

### Screens
- `docs/development/prd-docs/screens/global-ui.md`
- `docs/development/prd-docs/screens/dashboard.md`
- `docs/development/prd-docs/screens/accounts.md`
- `docs/development/prd-docs/screens/api-proxy.md`
- `docs/development/prd-docs/screens/monitor.md`
- `docs/development/prd-docs/screens/settings.md`
- `docs/development/prd-docs/screens/tray.md`
