# Story-61: Protobuf Token Decoder

**Epic:** [Epic-21: Zero-Config Onboarding](./epic-21-ide-integration.md)
**Status**: Draft
**Priority**: Medium
**Basis**: `docs/gap_analysis_deep_dive.md`

## Goal
The tokens in the IDE databases are stored as binary Protobuf blobs. This story covers the logic to decode these blobs and extract the actual `refresh_token` string, completing the "Zero Config" onboarding flow.

## Tasks
- [ ] **Task 1: Protobuf Definition / Decoding**
    - File: `apps/backend/migration/decoder.go`
    - Since we might not have the `.proto` file, implement a "Raw Field Extraction" logic similar to the Ref App (Field 6 -> Field 3).
    - Alternately, define a minimal Go struct with Protobuf tags if the schema is known.
- [ ] **Task 2: Token Validation**
    - Verify the extracted string looks like a JWT or Opaque Token.
- [ ] **Task 3: Import Logic**
    - Save the extracted token as a new Account in the local `accounts.db`.

## Acceptance Criteria
- [ ] Input: Binary blob from Android Studio DB.
- [ ] Output: Valid Refresh Token string.
- [ ] Unit Test: Verify against a known sample blob (if available/mocked).
