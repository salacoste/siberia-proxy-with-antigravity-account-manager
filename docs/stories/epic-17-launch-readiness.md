# Epic-17: Launch Readiness

**Status**: Completed
**Parent**: Project "Siberia"

## Executive Summary
Prepare the application for public release. This involves ensuring the application is securely packaged, signed (where applicable/possible), stripped of debug capabilities in production, and fully verified against the critical user journeys.

## Goals
1.  **Professional Packaging**: Produce clean, versioned artifacts for macOS (DMG/App) and Windows (NSIS/Exe).
2.  **Security Assurance**: Ensure no sensitive data is logged, debug tools are disabled in prod, and headers are secure.
3.  **Final Confidence**: Run a full regression suite (Manual + Auto) to "Green Light" the release.

## Stories

### Story-50: Final Packaging & CI Workflow
- **Goal**: Polish the build pipeline and manual build scripts.
- **Tasks**:
    -   Update `wails.json` metadata (Copyright 2026, version).
    -   Verify `wails build -platform windows/amd64,darwin/universal`.
    -   Create `release-build.sh` helper script.
    -   (Supersedes `story-99-packaging.md`).

### Story-51: Security Hardening
- **Goal**: Lock down the application for end-users.
- **Tasks**:
    -   Disable `devtools` in production Wails config.
    -   Review `fmt.Println` usage -> Change to structural logging or remove.
    -   Audit `TrafficController` for PII leakage (ensure Body capture is opt-in or safe).

### Story-52: Final Verification (Green Light)
- **Goal**: End-to-end verification.
- **Tasks**:
    -   Execute `docs/qa/test-plans/full-regression.md` (will create).
    -   Produce `RELEASE_CANDIDATE.md`.
