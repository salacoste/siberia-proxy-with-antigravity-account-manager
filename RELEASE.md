# Siberia Release Guide

This document outlines the build and release process for Siberia Proxy v1.2.0+.

## Prerequisites

- **Go**: v1.21+
- **Node.js**: v18+
- **Wails**: v2.6+ (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

## Project Structure

The project uses a monorepo-style layout:
- `apps/backend`: Go Backend & Wails Configuration (Run Wails commands here)
- `apps/frontend`: React/Vite Frontend

## Building the Application

**IMPORTANT**: All `wails` commands must be run from `apps/backend` directory.

### macOS (Universal)

To build a universal binary for Intel and Apple Silicon:

```bash
cd apps/backend
wails build -platform darwin/universal
```

The output will be located at: `apps/backend/build/bin/Siberia Proxy.app`

### Windows

```bash
cd apps/backend
wails build -platform windows/amd64
```

### Linux

Ensure `libgtk-3-dev` and `libwebkit2gtk-4.0-dev` are installed.

```bash
cd apps/backend
wails build -platform linux/amd64
```

## CI/CD Pipeline

Builds are automatically triggered via GitHub Actions (`.github/workflows/build.yml`) on:
- Push to `main`
- Tags start with `v*`

Artifacts are uploaded to GitHub Releases.
