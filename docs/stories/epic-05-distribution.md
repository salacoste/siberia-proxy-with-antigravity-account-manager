# Epic-05: Distribution & Auto-Updater

**Goal:** Enable the application to automatically check for updates, download them, and facilitate a smooth upgrade process for the user. Also includes setting up the CI pipeline to publish releases to GitHub.

## Scope
*   **CI/CD:** Automate GitHub Release creation with binaries attached when a tag is pushed.
*   **Backend:** Logic to query GitHub API for latest release, compare versions, and download assets.
*   **Frontend:** UI for "Check for Updates", "Update Available", and progress bars.
*   **Signing:** (Optional/Setup) Configuration for signing binaries if certificates are provided.

## Stories
- [ ] **Story-16: Automate GitHub Releases**
    -   Update `build.yml` to trigger on tags.
    -   Add step to upload build artifacts (binaries/installers) to the GitHub Release.
- [ ] **Story-17: Implement Auto-Update Logic**
    -   Add `UpdateService` in Go.
    -   Use `go-github` or raw HTTP to fetch `latest` release.
    -   Compare `CurrentVersion` vs `RemoteVersion` (Semver).
    -   Asset download logic (self-update lifting).
- [ ] **Story-18: Update UI & Notifications**
    -   Add "Check for Updates" button in Settings.
    -   Toast/Modal when update is found.
    -   Progress UI during download.
