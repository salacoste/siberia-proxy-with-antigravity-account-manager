# Story-16: Automate GitHub Releases

**Epic:** [Epic-05: Distribution](./epic-05-distribution.md)
**Status:** Completed


**Feature Branch:** `antigravity/feat/distribution`

## Goal
Automate the creation of GitHub Releases and asset uploads when a tag is pushed.

## Tasks
- [ ] **Task 1: Update Workflow**
    - [ ] Modify `.github/workflows/build.yml`.
    - [ ] Trigger on `push: tags: 'v*'`.
    - [ ] Build for Matrix (Ubuntu, macOS, Windows).
- [ ] **Task 2: Release Creation**
    - [ ] Use `softprops/action-gh-release` or similar.
    - [ ] Upload artifacts:
        - `siberia-darwin-universal`
        - `siberia-windows-amd64.exe`
        - `siberia-linux-amd64`
    - [ ] Generate Checksums (optional but recommended).

## Acceptance Criteria
- [ ] Pushing `v0.0.1` creates a Draft Release on GitHub.
- [ ] Release contains binary assets.

## QA Results


### Review Date: 2026-01-06 (Final)

### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
Workflow logic significantly improved since initial draft. Asset renaming script handles case-sensitivity and OS naming (`darwin` vs `macos`) correctly. Upgrade to `v2` of release action ensures future compatibility.

### Compliance Check
- Coding Standards: [✓] YAML syntax valid.
- All ACs Met: [✓] Release automation configured.

### Gate Status
Gate: PASS → docs/qa/gates/epic-05.story-16-github-release.yml

### Manual Verification (2026-01-06)
- **Method**: Code Audit
- **Status**: Verified
- **Notes**: Reviewed `.github/workflows/build.yml`. Confirmed `v*` tag trigger and matrix build strategy.

