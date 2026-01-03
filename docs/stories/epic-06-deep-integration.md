# Epic-06: Deep Integration - Real Process Injection

**Goal:** Replace the mocked process management and injection logic with real OS-specific implementations to enable functional account switching in VS Code.

## Existing System Context
*   **Current Functionality:** The `process` and `injection` modules currently run in `DryRun` mode, logging actions to stdout without affecting the system.
*   **Tech Stack:** Go (Backend), SQLite (Target App Data).
*   **Integration Points:** `apps/backend/siberia/accounts/service.go` calls `process.Kill`, `injector.Inject`, and `process.Start`.

## Enhancement Details
*   **What:**
    *   Implement `process.Service.Kill()` using `os/exec` or `syscall` to terminate "Code Helper" or "Electron" processes.
    *   Implement `injector.Service.Inject()` to write tokens into VS Code's `state.vscdb` (SQLite).
*   **Success Criteria:**
    *   Clicking "Switch" actually closes VS Code.
    *   Clicking "Switch" modifies the target SQLite DB file.
    *   VS Code restarts with the new state.

## Stories

### Story-19: Implement Real Process Termination
*   **Goal:** Ability to find and kill a process by name (e.g., "Visual Studio Code").
*   **Tasks:**
    *   Refactor `modules/process` to use `gopsutil` or `pgrep/pkill`.
    *   Handle permission errors gracefully.
    *   **Compatibility:** Ensure it works on macOS (primary target).

### Story-20: Implement SQLite Token Injection
*   **Goal:** Write connection tokens into VS Code's `globalStorage/state.vscdb`.
*   **Tasks:**
    *   Locate default VS Code `state.vscdb` path on macOS.
    *   Open target DB in `modules/injection`.
    *   Update specific keys (e.g., `github.auth`).
    *   **Risk:** Corrupting VS Code state. **Mitigation:** Backup file before writing.

### Story-21: Integration & Verification
*   **Goal:** Wire it all together and verify end-to-end.
*   **Tasks:**
    *   Update `App` checks to verify paths exist.
    *   End-to-End test: Login -> Switch -> Verify VS Code changed.

## Compatibility Requirements
*   [x] Existing APIs remain unchanged (Service interfaces are stable).
*   [x] Database schema changes are N/A (Internal Siberia DB).
*   [ ] VS Code version compatibility (Target external app may change schema).

## Risk Mitigation
*   **Primary Risk:** Corrupting user's VS Code installation or data.
*   **Mitigation:** Always create a `.bak` copy of `state.vscdb` before writing.
*   **Rollback:** Button to "Restore Backup" if VS Code fails to launch.
