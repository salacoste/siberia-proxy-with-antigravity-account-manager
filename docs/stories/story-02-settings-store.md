# Story-02: Implement Settings Store & Persistence

**Epic:** [Epic-01: Core Foundation & UI Shell](./epic-01-core.md)
**Status:** Completed

## Goal
Implement the backend configuration system and the frontend state store to manage application settings. This is the "brain" of the app's preferences.

## Tasks
- [ ] **Task 1: Backend Config Module**
    -   Define `Config` struct (Theme, Language, RefreshInterval, etc.).
    -   Implement `Load()` and `Save()` logic using `config.json` in the app data directory.
    -   Create `ConfigService` Wails struct.
- [ ] **Task 2: Wails Bindings**
    -   Expose `GetAppConfig() Config`
    -   Expose `UpdateAppConfig(Config)`
    -   Expose `GetAppVersion() string`
- [ ] **Task 3: Frontend Store**
    -   Create `useConfigStore` using Zustand.
    -   Initialize store by calling `GetAppConfig` on hydration.

## Acceptance Criteria
- [ ] `config.json` is created in the correct OS-specific data path on first launch.
- [ ] Frontend can read the configuration.
- [ ] Frontend can update a setting (e.g., Theme), and it persists to disk after restart.

## Technical Notes
- **Library:** Use standard `encoding/json` for simplicity, or `viper` if complexity requires. (Stick to `encoding/json` for now unless validation is heavy).
- **Paths:** Use `os.UserConfigDir()` to determine the correct location for `siberia/config.json`.
- **Concurrency:** Ensure `Save()` is thread-safe if multiple writes happen (use `sync.RWMutex`).

## QA Results

### Review Date: 2026-01-03

### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
Configuration management uses a robust `Manager` struct with `sync.RWMutex` for thread safety, meeting technical requirements. Encryption for `MasterKey` is handled correctly using `crypto` package.

### Compliance Check
- Coding Standards: [✓] Thread-safe access.
- Project Structure: [✓] `siberia/config` package used.
- All ACs Met: [✓] Config persists to disk.

### Gate Status
Gate: PASS → docs/qa/gates/epic-01.story-02-settings-store.yml

