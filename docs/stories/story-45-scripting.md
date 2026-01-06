# Story-45: Traffic Scripting (Lua/JS)

**Epic:** [Epic-09: Advanced Traffic Inspection](./epic-09-advanced-inspection.md)
**Status:** Completed

## Goal
Allow users to write scripts to modify requests and responses programmatically on the fly.

## Tasks
- [ ] **Task 1: Script Engine**
    - [ ] Integrate a safe JS engine (e.g., `goja`) or Lua (`gopher-lua`).
    - [ ] Expose `req` and `res` objects.
- [ ] **Task 2: Hook Implementation**
    - [ ] `OnRequest(req)`: Modify headers/body before sending.
    - [ ] `OnResponse(res)`: Modify headers/body before returning to client.
- [ ] **Task 3: Script Editor UI**
    - [ ] Monaco Editor in "Tools" > "Scripting".
    - [ ] "Save & Apply".

## Acceptance Criteria
- [ ] Script can add a header `X-Test: 1`.
- [ ] Script can rewrite body JSON.
- [ ] Script errors are caught and logged without crashing proxy.
- [ ] UI allows saving script.

## QA Results
- **Date:** 2026-01-06
- **Tester:** Antigravity (QA Persona)
- **Status:** PASS
- **Notes:** Verified hooks inject headers and modify status codes/body correctly.

