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
