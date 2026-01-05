# Story-51: Security Hardening

**Epic**: [Epic-17: Launch Readiness](../epic-17-launch-readiness.md)
**Status**: Completed

## Goal
Ensure the application is secure for end-users by disabling debugging tools and preventing sensitive data leakage.

## Scope
1.  **Disable DevTools**:
    -   Ensure `wails.json` or build flags disable the "Inspect Element" context menu in production builds.
    -   (Note: `wails build -production` typically handles this, verify config).
2.  **Log Audit**:
    -   Search for `fmt.Println` or `log.Printf` that might dump raw request bodies or auth tokens.
    -   Replace with `log/slog` structured logging (if available) or remove.
3.  **App Config**:
    -   Ensure no hardcoded secrets in `config.go` or frontend code.

## Acceptance Criteria
- [x] DevTools disabled in production build.
- [x] No `fmt.Println` showing request bodies.
- [x] No hardcoded API keys in source.

## Verification
- Code Audit (grep).
- Launch production build and try to open Inspector.
