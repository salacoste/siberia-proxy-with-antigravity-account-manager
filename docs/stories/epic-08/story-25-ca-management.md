# Story-25: Certificate Authority Management (Backend)

**Epic:** [Epic-08: HTTPS Decryption (MitM)](../epic-08-mitm.md)
**Status**: Completed
**Feature Branch**: `antigravity/feat/mitm-ca`
**Parent**: Epic-08

## Description
Implement the backend capability to generate, store, and retrieve the Root Certificate Authority (CA) used for HTTPS decryption (MitM). If a CA does not exist, one should be generated securely.

## Requirements
1.  **Module**: `apps/backend/siberia/ca` (or `mitm` package).
2.  **Storage**:
    -   Store `ca.pem` (Cert) and `ca.key` (Private Key) in the application data directory (`UserDataDir/certificates`).
    -   Secure storage permissions (readable only by owner).
3.  **Generation**:
    -   Use `crypto/x509` to generate a Root Certificate.
    -   **Algorithm**: RSA 2048 or ECDSA P256 (P256 preferred for performance).
    -   Subject Name: "Siberia Proxy CA".
    -   Validity: 10 years.
    -   Key Usage: `CertSign`, `CRLSign`.
4.  **Validation**:
    -   **Permissions**: `ca.key` MUST be set to `0600` (Read/Write by Owner only) on Unix-like systems.
    -   **Idempotency**: `EnsureCA()` must NEVER overwrite an existing valid CA.
5.  **API**:
- [x] `GetCAPath()`: Returns full paths to cert and key.
    - [x] `EnsureCA()`: Checks existence, generates if missing.

## Acceptance Criteria
- [x] `EnsureCA()` creates `ca.pem` and `ca.key` if they don't exist.
- [x] Generated certificate has the correct Subject Name ("Siberia Proxy CA").
- [x] Generated certificate is valid for ~10 years.
- [x] Private key is stored securely (file permissions).
- [x] Subsequent calls to `EnsureCA()` do not overwrite existing CA.
- [x] Unit tests verify generation and parsing.

## Dev Agent Record

### Completion Notes
- Implemented `siberia/ca` module using `crypto/x509` (RSA 2048).
- Enforced strict `0600` permissions on `ca.key`.
- Verified idempotency via unit tests.
- Wired `EnsureCA()` into `app.go` startup.

### Files Modified
- `apps/backend/siberia/ca/cert.go`
- `apps/backend/siberia/ca/service.go`
- `apps/backend/siberia/ca/install_darwin.go`
- `apps/backend/app.go`



## QA Results

### Review Date: 2026-01-07

### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment

The implementation is robust and follows security best practices. Separating the CA logic into its own module (`siberia/ca`) ensures clean separation of concerns.

-   **Security**: Key generation uses standard `crypto/rsa` (2048 bit). Private key storage strictly enforces `0600` permissions, ensuring only the owner can read/write.
-   **Idempotency**: `EnsureCA` correctly handles existing files, preventing accidental overwrite of the Root CA (which would break trust on all client devices).
-   **Testing**: Unit tests cover generation parameters (Subject, Validity) and storage mechanism (Permissions, file existence).

### Refactoring Performed

None required. The initial implementation by Dev was clean.

### Compliance Check

-   Coding Standards: [✓]
-   Project Structure: [✓] (`siberia/ca` matches module pattern)
-   Testing Strategy: [✓] (Unit tests in place)
-   All ACs Met: [✓]

### Improvements Checklist

-   [ ] Consider ECDSA P256 support in future for performance (currently RSA 2048).
-   [ ] Add `install_linux.go` and `install_windows.go` for cross-platform support (currently `install_darwin.go` only).

### Security Review

-   **Critical**: Private key (`ca.key`) permissions are enforced to `0600`.
-   **Generation**: Self-signed root CA is valid for 10 years, minimizing rotation needs.

### Files Modified During Review

None.

### Gate Status

Gate: PASS → docs/qa/gates/epic-08.story-25-ca-management.yml

### Recommended Status

[✓ Ready for Done]
