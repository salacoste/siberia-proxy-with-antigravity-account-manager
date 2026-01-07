# Story-01: Initialize Project Structure

**Epic:** [Epic-01: Core Foundation & UI Shell](./epic-01-core.md)
**Status:** Completed

## Goal
Initialize the repository with the correct project structure for a Go + Wails + React + Vite application.

## Tasks
- [ ] **Task 1:** Initialize Wails Project
    - Command: `mkdir -p apps/backend && cd apps/backend && wails init -n siberia -t react-ts`
    - Output: `apps/backend/frontend/`, `apps/backend/main.go`, `apps/backend/wails.json`.
- [ ] **Task 2:** Restructure Directories
    - Command: `mv apps/backend/frontend apps/frontend`.
    - Action: Edit `apps/backend/wails.json` to set `"frontend:dir": "../frontend"`.
    - Action: Edit `apps/backend/wails.json` to set `"wailsjs:dir": "../frontend/src/wailsjs"`.
- [ ] **Task 2:** Setup Go Module
    - Ensure `go.mod` is named `github.com/salacoste/siberia`.
    - Configure `build/` directory for keeping binaries.
- [ ] **Task 3:** Configure Frontend
    - Install `lucide-react`, `clsx`, `tailwind-merge` (UI utils).
    - Setup standard folder structure: `src/{components,pages,stores,hooks,lib}`.
- [ ] **Task 4:** CI/CD Prep
    - Create `.github/workflows/build.yml` for automated cross-platform builds.

## Acceptance Criteria
- [ ] `wails dev` starts the application without errors.
- [ ] Frontend displays "Hello World" (or basic template).
- [ ] Directory structure matches the desired architecture (`/frontend`, `/backend` separation).

## Technical Notes
- Use Wails v2.
- Ensure `go` version >= 1.21.
- Ensure `node` version >= 18.

## QA Results

### Review Date: 2026-01-03

### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
The project initialization strictly follows the Wails+React+Go architecture. The directory structure is correct (`apps/backend`, `apps/frontend`). The CI/CD pipeline setup (`.github/workflows/build.yml`) was not explicitly verified but the local build structure is valid.

### Compliance Check
- Coding Standards: [✓] Standard Go/TS structure.
- Project Structure: [✓] Matches Monorepo design.
- All ACs Met: [✓] App launches, Hello World (Dashboard) visible.

### Gate Status
Gate: PASS → docs/qa/gates/epic-01.story-01-initialize-project.yml

