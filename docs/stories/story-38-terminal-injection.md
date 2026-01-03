# Story-38: Implement Terminal Shell Helper

**Epic:** [Epic-12: Native Integrations beyond VS Code](./epic-12-native-integrations.md)
**Status:** Draft

## Goal
Allow users to easily configure their current terminal session to use the proxy.

## Tasks
- [ ] **Task 1: "Copy Export Command"**
    -   Button in API Proxy UI: "Copy Shell Config".
    -   Copies: `export HTTP_PROXY=... HTTPS_PROXY=...`
- [ ] **Task 2: Persistent Config (Advanced)**
    -   Option to append to `~/.zshrc` (Requires user confirmation/trust).

## Acceptance Criteria
- [ ] User can click one button to get the correct env var string.
- [ ] Pasting into terminal instanty routes traffic through Siberia.
