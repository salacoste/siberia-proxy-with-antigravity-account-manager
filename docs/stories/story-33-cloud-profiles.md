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
