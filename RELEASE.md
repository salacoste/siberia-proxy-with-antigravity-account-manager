# Release Instructions

## Building Siberia

### Prerequisites
*   Go 1.21+
*   Node.js 20+
*   Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Build Command
Run the following from the root directory:

```bash
cd apps/backend
wails build
```

The output binary will be located in `apps/backend/build/bin/`.

### Cross-Compilation (CI/CD)
The project is configured with GitHub Actions to build for macOS, Windows, and Linux automatically on push to `main`.

### Versioning
Update `wails.json` version field before tagging a release.

## Creating a Release (CI/CD)
The project uses GitHub Actions to automate releases.

1.  **Bump Version:**
    *   Update `version` in `apps/backend/wails.json` (e.g., `1.0.1`).
    *   Update `GetVersion()` in `apps/backend/app.go`.
    *   Commit changes: `git commit -m "chore: bump version to v1.0.1"`
2.  **Tag & Push:**
    *   `git tag v1.0.1`
    *   `git push origin v1.0.1`
3.  **Automation:**
    *   GitHub Actions (`release.yml`) will trigger.
    *   Binaries for Linux, Windows, and macOS will be built.
    *   A GitHub Release will be created/updated with these assets.

## Auto-Updater
The application includes a built-in auto-updater.
*   It queries the GitHub Releases API for the `latest` tag.
*   If `latest > current`, it prompts the user to update.
*   Users can manually check in **Settings > Application Info**.
