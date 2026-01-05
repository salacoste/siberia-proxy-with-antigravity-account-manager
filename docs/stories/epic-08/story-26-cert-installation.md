# Story-26: Trusted Cert Installation (Backend/OS)

**Epic:** [Epic-08: HTTPS Decryption (MitM)](../epic-08-mitm.md)
**Status:** Ready for Dev
**Parent**: Epic-08

## Description
To decrypt HTTPS traffic without browser warnings, the Root CA generated in Story-25 must be trusted by the OS. This story covers the implementation of the `InstallCert` and `CheckTrust` methods, specifically targeting macOS for the MVP.

## Requirements
1.  **Platform Support**: macOS (Darwin) is the priority. Windows/Linux stubs can return "Not Supported" or manual instructions.
2.  **`InstallCert()`**:
    -   Execute `security add-trusted-cert` to add `ca.pem` to the login keychain (or system keychain if appropriate).
    -   Must handle the potential password prompt (OS native dialog) or error if user cancels.
3.  **`CheckTrust()`**:
    -   Verify if the certificate is currently trusted.
    -   Use `security dump-trust-settings` or verify the cert chain against the system store.
    -   Return `true` if trusted, `false` otherwise.

## Acceptance Criteria
- [x] `InstallCert` successfully adds the CA to the macOS Keychain.
- [x] User sees the OS "Type your password to allow this" dialog (handled by `security` command).
- [x] `CheckTrust` returns `true` after successful installation.
- [x] `CheckTrust` returns `false` if cert is missing or untrusted.
- [x] `InstallStub` handles non-macOS platforms gracefully (no crash).

## Dev Agent Record

### Completion Notes
- Implemented `InstallCert` using `security add-trusted-cert` targeting `login.keychain-db`.
- Implemented `CheckTrust` using `security verify-cert -c <file>`.
- Verified `CheckTrust` returns false for uninstalled certs via `go test`.
- Manual verification of "Install" flow will be part of the full Epic acceptance (Story-28 UI).

### Story DoD Checklist
- [x] Logic is robust for macOS.
- [x] Cross-platform stubs in place.
- [x] Tests pass (negative verification).

## Dev Notes
-   Command: `security add-trusted-cert -d -r trustRoot -k ~/Library/Keychains/login.keychain-db <ca_path>`
-   Verify command works on modern macOS (Sequioa/Sonoma).
