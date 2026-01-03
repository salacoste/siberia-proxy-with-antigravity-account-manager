# Story-10: Implement SQLite Database Layer (Encryption)

**Epic:** [Epic-03: Accounts Manager](./epic-03-accounts.md)
**Status:** Draft

## Goal
Establish the local database layer for storing Account entities. The database must be embedded (SQLite) and support field-level encryption for sensitive data (tokens, passwords).

## Tasks
- [ ] **Task 1: Core Setup**
    -   Install `GORM` and `Pure Go SQLite` driver (`glebarez/sqlite` preferred for CGO-free builds).
    -   Create `siberia/db` package.
    -   Initialize GORM connection in `App` startup.
- [ ] **Task 2: Encryption Utility**
    -   Implement AES-GCM encryption helper (`siberia/crypto`).
    -   Key management: Generate/Load a master key stored in `config.json` (or OS keyring if possible, but config for MVP is standard per PRD).
- [ ] **Task 3: Models**
    -   Define `Account` struct:
        -   ID (uint)
        -   Email (string, unique)
        -   PasswordEncr (string)
        -   SessionTokenEncr (string)
        -   ProxyGroup (string)
        -   Stats (JSON)
        -   Timestamps
- [ ] **Task 4: Repository**
    -   Implement CRUD methods (`Create`, `List`, `Update`, `Delete`, `FindOne`).
    -   Ensure automatic encrypt/decrypt on Read/Write.

## Acceptance Criteria
- [ ] App starts with SQLite DB `siberia.db` in app/config dir.
- [ ] `Account` table is auto-migrated.
- [ ] Sensitive fields are encrypted on disk (verified by inspecting DB file).
- [ ] Application can read back decrypted values transparently.

## Technical Notes
- **Driver:** `github.com/glebarez/sqlite` to avoid CGO requirements for cross-compilation (Windows/Mac/Linux).
- **Encryption:** Use `crypto/aes` and `crypto/cipher` (GCM).

## QA Results

### Review Date: 2026-01-03

### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
Database layer implementation (`siberia/db`) uses robust encryption-at-rest via `EncryptedString` type scanner/valuer. This is a very secure pattern that makes encryption transparent to the repository layer. `gorm` + `glebarez/sqlite` correctly avoids CGO.

### Compliance Check
- Coding Standards: [✓] GORM usage.
- All ACs Met: [✓] Data persists and is encrypted.

### Gate Status
Gate: PASS → docs/qa/gates/epic-03.story-10-sqlite-layer.yml

