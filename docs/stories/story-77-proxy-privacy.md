# Story-77: Proxy Privacy & Identity Logic

**Epic:** [Epic-22: Enterprise Compliance & Legacy Support](./epic-22.md)
**Status**: Done
**Priority**: High
**Basis**: `docs/ag-ref-docs/feature-proxy-engine.md` (Sections 5 & 6)

## Story
**As a** System Administrator,
**I want** user emails to be masked in logs and internal identity to be correctly resolved,
**so that** I can comply with privacy regulations (GDPR/PII) and ensure reliable access to internal tools.

## Acceptance Criteria
- [ ] **Email Masking**: Emails logged in access logs are consistently masked (e.g., `u***@host.com`).
- [ ] **Project ID Resolution**: Upstream client can resolve a valid `cloudaicompanionProject` ID given a token.
- [ ] **Mock Fallback**: Fallback to Mock Project ID occurs if upstream resolution fails (or 403).
- [ ] **No ID Leakage**: No raw user IDs are logged (must be anonymized or hashed).

## Tasks / Subtasks
- [ ] **Task 1: Privacy Package** (AC: Email Masking, No ID Leakage)
    - [ ] Create `apps/backend/siberia/proxy/privacy/masking.go`.
    - [ ] Implement `MaskEmail(email string) string`.
    - [ ] Implement `AnonymizeID(id string) string`.
    - [ ] Implement `StableHash(value string) string`.
    - [ ] Add unit tests in `masking_test.go`.
- [ ] **Task 2: Identity Package** (AC: Project ID Resolution, Mock Fallback)
    - [ ] Create `apps/backend/siberia/proxy/identity/resolver.go`.
    - [ ] Implement `FetchProjectID(accessToken string)`.
    - [ ] Implement `GenerateMockProjectID()`.
    - [ ] Integration: Call upstream `loadCodeAssist` endpoint.
- [ ] **Task 3: Integration**
    - [ ] Usage in `middleware/logging.go` (Mask emails).
    - [ ] Usage in `upstream/client.go` (Inject Project ID header if needed, or just resolve it).

## Dev Notes
- **Reference**: See `docs/ag-ref-docs/feature-proxy-engine.md`.
- **Upstream Endpoint**: `https://cloudcode-pa.googleapis.com/v1internal:loadCodeAssist`
- **Mock Format**: `{adj}-{noun}-{5chars}` (e.g., `swift-core-x9y8z`).
- **Files to Create**:
    - `apps/backend/siberia/proxy/privacy/masking.go`
    - `apps/backend/siberia/proxy/identity/resolver.go`

### Testing
- **Location**: `apps/backend/siberia/proxy/privacy/masking_test.go`, `apps/backend/siberia/proxy/identity/resolver_test.go`
- **Standards**: Standard Go unit tests.
- **Requirements**: Ensure 100% coverage of masking logic.

## Change Log
| Date | Version | Description | Author |
|---|---|---|---|
| 2026-01-08 | 1.0 | Initial Draft | PO |

## Dev Agent Record
### Agent Model Used
_TBD_

### Debug Log References
_TBD_

### Completion Notes List
_TBD_

### File List
_TBD_

## QA Results

### Review Date: 2026-01-08
### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
The implementation is robust and follows the "Secure by Design" principle. By moving masking logic to a dedicated `privacy` package and enforcing it via centralized middleware (`attribution.go`), the risk of accidental PII leakage in individual handlers is minimized. The code is clean, well-tested, and adheres to Go standards.

### Refactoring Performed
- **Handlers**: Removed ad-hoc manual masking from `openai`, `claude`, `gemini`, and `legacy` handlers to use the centralized `privacy` package.
  - **Why**: Reduces code duplication and potential for inconsistent masking rules.
  - **How**: Passed raw identity to `attribution.go` where `privacy.MaskEmail` is applied.

### Compliance Check
- Coding Standards: [✓]
- Project Structure: [✓]
- Testing Strategy: [✓] (Unit + Integration)
- All ACs Met: [✓]

### Security Review
- **PII Protection**: Excellent. Centralized masking ensures no email leaks in logs.
- **Identity Resolution**: Securely fetches Project ID. Fallback mechanism is safe for non-prod but should be monitored.

### Gate Status
Gate: PASS → docs/qa/gates/epic-22.story-77-proxy-privacy.yml

### Recommended Status
[✓ Ready for Done]
