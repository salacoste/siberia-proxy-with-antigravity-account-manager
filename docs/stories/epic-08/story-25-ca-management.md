# Story-25: Certificate Authority Management (Backend)

**Epic:** [Epic-08: HTTPS Decryption (MitM)](../epic-08-mitm.md)
**Status:** In Progress
**Parent**: Epic-08

## Description
Implement the backend capability to generate, store, and retrieve the Root Certificate Authority (CA) used for HTTPS decryption (MitM). If a CA does not exist, one should be generated securely.

## Requirements
1.  **Module**: `apps/backend/siberia/ca` (or `mitm` package).
2.  **Storage**:
    -   Store `ca.pem` (Cert) and `ca.key` (Private Key) in the application data directory (`UserDataDir/certificates`).
    -   Secure storage permissions (readable only by owner).
3.  **Generation**:
    -   Use `crypto/x509` to generate an RSA or ECDSA root certificate.
    -   Subject Name: "Siberia Proxy CA".
    -   Validity: 10 years.
    -   Key Usage: `CertSign`, `CRLSign`.
4.  **API**:
    -   `GetCAPath()`: Returns full paths to cert and key.
    -   `EnsureCA()`: Checks existence, generates if missing.

## Acceptance Criteria
- [x] `EnsureCA()` creates `ca.pem` and `ca.key` if they don't exist.
- [x] Generated certificate has the correct Subject Name ("Siberia Proxy CA").
- [x] Generated certificate is valid for ~10 years.
- [x] Private key is stored securely (file permissions).
- [x] Subsequent calls to `EnsureCA()` do not overwrite existing CA.
- [x] Unit tests verify generation and parsing.

## Dev Agent Record

### Completion Notes
- Implemented `siberia/ca` service with RSA 2048 key generation.
- Enforced `0600` permissions for `ca.key`.
- Integrated `EnsureCA` into `app.go`.
- Verified via `service_test.go` confirming idempotency and regeneration logic.

### Story DoD Checklist
- [x] All functional requirements met.
- [x] Tests passed.
- [x] Code follows patterns.

## QA Results

### Review Date: 2026-01-05

### Reviewed By: Quinn (Test Architect)

- **Audit**: Code reviewed. `EnsuraCA` logic is sound. Tests cover idempotency.
- **Verdict**: Approved.

### Gate Status

Gate: PASS → docs/qa/gates/epic-08.story-25-ca-management.yml
