# Story-33: Design & Implement Cloud Profile Sync

**Epic:** [Epic-10: Team Collaboration & Cloud Sync](./epic-10-collaboration.md)
**Status:** Draft

## Goal
Allow a user to login (e.g., via GitHub/Google) and sync their local `accounts.db` and `config.json` to the cloud.

## Tasks
- [ ] **Task 1: Auth Integration**
    -   Integrate Auth Provider (e.g., Supabase Auth).
    -   Login Screen in "Settings".
- [ ] **Task 2: Sync Engine**
    -   Detect local changes -> Push.
    -   Detect remote changes -> Pull.
    -   Conflict resolution (Last Write Wins for MVP).

## Acceptance Criteria
- [ ] User can log in.
- [ ] Adding an account on Device A appears on Device B.
- [ ] Secrets (passwords) are strictly E2EE or managed via secure vault (Key Decision required).

## PO Notes
- **Status review**: Moved to "Ready for Dev" (Architect Approved).
- **Blocker**: Resolved. See [Architecture: Cloud Sync & E2EE](../architecture/modules/cloud-sync.md).
- **Action**: Dev to implement `crypto/vault` and `sync/manager`.


## QA Results (Backend Core)

### Review Date: 2026-01-04

### Reviewed By: Quinn (Test Architect)

- **Crypto Vault**: Verified AES-256-GCM and Argon2id implementations yield expected ciphertext structure.
- **Sync Protocol**: Confirmed Last-Write-Wins (LWW) logic correctly handles Push/Pull conflicts in simulated server environment.
- **Security**: Confirmed Zero-Knowledge design (Server receives only encrypted blobs).

### Gate Status

Gate: PASS → docs/qa/gates/epic-10.story-33-cloud-profiles-backend.yml
