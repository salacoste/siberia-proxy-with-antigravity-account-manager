# Story-20: Implement SQLite Token Injection

**Epic:** [Epic-06: Deep Integration](./epic-06-deep-integration.md)
**Status:** Completed

## Goal
Write connection tokens into VS Code's `globalStorage/state.vscdb`.

## Tasks
- [x] Locate default VS Code `state.vscdb` path.
- [x] Open target DB in `modules/injection`.
- [x] Update specific keys (e.g., `siberia.auth_token`).
- [x] Create backup `.bak` before writing.

## Acceptance Criteria
- [x] DB is modified with new JSON token.
- [x] Backup file is created.

## QA Results

### Review Date: 2026-01-04

### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
`injector.go` uses `database/sql` + `go-sqlite` for robust DB handling. Backup logic (`copyFile`) acts as a critical safety net before mutation.

### Compliance Check
- Coding Standards: [✓] Error handling and resource cleanup (`defer`).
- All ACs Met: [✓] Injection logic verified.

### Gate Status
Gate: PASS → docs/qa/gates/epic-06.story-20-file-injection.yml
