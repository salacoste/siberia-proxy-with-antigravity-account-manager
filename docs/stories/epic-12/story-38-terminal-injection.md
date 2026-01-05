# Story-38: Implement Terminal Shell Helper

**Epic:** [Epic-12: Native Integrations beyond VS Code](./epic-12-native-integrations.md)
**Status:** Completed

## Goal
Allow users to easily configure their current terminal session to use the proxy.

## Tasks
- [x] **Task 1: "Copy Export Command"**
    -   Button in API Proxy UI: "Copy Shell Config".
    -   Copies: `export HTTP_PROXY=... HTTPS_PROXY=...`
- [ ] **Task 2: Persistent Config (Advanced)**
    -   Option to append to `~/.zshrc` (Requires user confirmation/trust).

## Acceptance Criteria
## Acceptance Criteria
- [x] User can click one button to get the correct env var string.
- [x] Pasting into terminal instanty routes traffic through Siberia.

## QA Results (2026-01-04)
**Agent**: Quinn (QA)
**Decision**: **PASS** (via Gate `epic-12.story-38-terminal-helper.yml`)

### Verification
- **UI**: Added Terminal Configuration card with tabs for different shells.
- **Functionality**: Copy button uses `navigator.clipboard` correctly.
- **Correctness**: Generated commands match standard proxy environment variable syntax.
