# Story-17: Implement Auto-Update Logic

**Epic:** [Epic-05: Distribution](./epic-05-distribution.md)
**Status:** Draft
**Feature Branch:** `antigravity/feat/distribution`

## Goal
Implement the backend service to check for updates and download the new version.

## Tasks
- [ ] **Task 1: Update Service**
    - [ ] Create `siberia/updater/service.go`.
    - [ ] Method `CheckForUpdates()` -> Returns `UpdateInfo` (Version, Notes, URL).
    - [ ] Call GitHub API: `GET /repos/{owner}/{repo}/releases/latest`.
- [ ] **Task 2: Version Comparison**
    - [ ] Parse `CurrentVersion` vs `LatestVersion` (use `Masterminds/semver` or simple string compare).
    - [ ] Return `UpdateAvailable: true/false`.
- [ ] **Task 3: Download & Apply**
    - [ ] Method `DownloadUpdate(url)`.
    - [ ] Use `minio/selfupdate` or `equinox-io/selfupdate` (or generic replacement logic).
    - [ ] **Simplification:** For MVP, just download the file and prompt user to replace, OR open the release page in browser if self-update is too risky for now.
    -   *Decision:* Open Browser URL is safer for MVP v1.

## Acceptance Criteria
- [ ] Backend can fetch latest release info from GitHub.
- [ ] Correctly identifies if remote > local.

## QA Results

### Review Date: 2026-01-04

### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
Backend `updater/service.go` correctly implements `CheckForUpdates`. Uses GitHub API to fetch latest release. Comparison logic (string based for MVP) is acceptable.

### Compliance Check
- Coding Standards: [✓] Safe HTTP Client usage (timeouts).
- All ACs Met: [✓] Version check logic verified.

### Gate Status
Gate: PASS → docs/qa/gates/epic-05.story-17-auto-update.yml

