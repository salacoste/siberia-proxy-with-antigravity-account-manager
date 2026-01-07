# Story-56: Session Fingerprinting & Sticky Sessions

**Epic:** [Epic-20: Advanced Traffic Scheduling](./epic-20-advanced-scheduling.md)
**Status**: Done
**Priority**: Medium
**Basis**: `docs/gap_analysis_deep_dive.md`

## Goal
Implement "Sticky Sessions" to improve upstream cache hit rates (KV Caching). By fingerprinting incoming requests based on their model and initial prompt, the proxy can consistently route related requests to the same upstream account, mimicking a stateful session even for stateless clients.

## Tasks
- [ ] **Task 1: Session Fingerprinting**
    - File: `apps/backend/siberia/proxy/session/fingerprint.go`
    - Implement `GenerateFingerprint(model string, messages []Message) string`.
    - Logic: Hash `model_name` + `first_user_message_content`.
- [ ] **Task 2: Sticky Routing Logic**
    - File: `apps/backend/siberia/proxy/loadbalancer.go`
    - Update routing logic: Check if a fingerprint points to an active, valid account. If yes, prioritize that account.
- [ ] **Task 3: Session Header Support**
    - Respect specific headers (e.g., `x-siberia-session-id`) if provided by the client, overriding the implicit fingerprint.

## Acceptance Criteria
- [ ] Two requests with the exact same first prompt are routed to the same Account ID (assuming the account is healthy).
- [ ] Requests with explicit `x-siberia-session-id` always map to the same account.

## Dev Agent Record
- **Agent**: James (Dev)
- **Model**: Gemini 2.0 Flash
- **Date**: 2026-01-07

### Modified Files
- `apps/backend/siberia/proxy/session/fingerprint.go` (New)
- `apps/backend/siberia/proxy/session/fingerprint_test.go` (New)
- `apps/backend/siberia/accounts/service.go`
- `apps/backend/siberia/proxy/upstream/gemini.go`
- `apps/backend/siberia/proxy/upstream/gemini_test.go`
- `.github/workflows/build.yml` (CI Fix)

### Completion Notes
- Implemented deterministic Sticky Sessions based on Prompt Hash.
- Refactored `GetRotatingToken` to accept `fingerprint`.
- Added comprehensive unit tests.
- Fixed CI/CD `PKG_CONFIG_PATH` issue.
- **[Post-QA Fix]**: Implemented `x-siberia-session-id` header support (AC2).
   - Updated `handlers/gemini/handler.go` to inject header into Context.
   - Updated `upstream/gemini.go` to prioritize context session ID.
   - Added `TestGeminiClient_SessionHeader_Override`.

### Change Log
- **Feature**: Added `GenerateFingerprint`.
- **Refactor**: Updated `AccountProvider` interface.
- **Test**: Added `TestGeminiClient_StickySession`.

### Debug Log
- Encountered `go.mod` missing deps -> Ran `go mod tidy`.
- Encountered Lint Error in tests -> Updated Mock in `gemini_test.go`.

## QA Results

### Review Date: 2026-01-07
### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
Implementation of fingerprinting logic is clean and well-tested (`session/fingerprint.go`). However, the implementation is incomplete against the Acceptance Criteria.

### Compliance Check
- Coding Standards: [✓]
- Project Structure: [✓]
- Testing Strategy: [✓] Unit tests present.
- All ACs Met: [✗] **AC2 Failed**

### Review Date: 2026-01-07 (Re-review)
### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
Dev has addressed the missing Acceptance Criteria (AC2). The implementation correctly prioritizes the `x-siberia-session-id` header over the content fingerprint.

### Compliance Check
- Coding Standards: [✓]
- Project Structure: [✓]
- Testing Strategy: [✓] New test added `TestGeminiClient_SessionHeader_Override`.
- All ACs Met: [✓] **AC2 Passed**

### Gate Status
Gate: PASS → docs/qa/gates/epic-20.story-56-sticky-sessions.yml

### Recommended Status
[✓ Ready for Done]
