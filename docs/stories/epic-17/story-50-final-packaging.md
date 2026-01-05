# Story-50: Final Packaging & CI Workflow

**Epic**: [Epic-17: Launch Readiness](../epic-17-launch-readiness.md)
**Status**: Completed

## Goal
Establish a reliable, repeatable build process for producing production-ready artifacts for macOS and Windows.

## Scope
1.  **Metadata**:
    -   Ensure `wails.json` has correct version, product name, and description.
2.  **Build Automation**:
    -   Create `scripts/build_release.sh` to automate the build command with correct flags (`-clean`, `-production`).
3.  **Cross-Platform Check**:
    -   Verify the build works for macOS (Intel/M1) and Windows (AMD64).

## Acceptance Criteria
- [x] `wails.json` metadata is accurate.
- [x] `scripts/build_release.sh` exists and runs successfully.
- [x] `build/bin` contains the expected artifacts.
- [x] Application icons are correctly applied (visual check).

## Verification
- Run `scripts/build_release.sh`.
- Launch the resulting app.
