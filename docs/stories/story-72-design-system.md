# Story-72: Design System Foundation (Corca Adaptation)

**Epic**: [Epic-15: Visual Redesign](./epic-15-visual-redesign.md)
**Status**: Ready for Dev
**Priority**: High

## Goal
Establish the foundational design tokens to align the Siberia Proxy UI with the "Corca" aesthetic. This involves updating the Tailwind configuration and global CSS variables to use specific font stacks, colors, and radii.

## Requirements
1.  **Typography**:
    -   Primary Font: `Inter` (Body).
    -   Heading Font: `Space Grotesk` (or `GeistSans` if available/preferred, sticking to current `Space Grotesk` if closest match to "Tech"). *Correction*: Epic doc mentions GeistSans. If we don't have Geist, we stick to `Inter` with tight tracking or `Space Grotesk`. Let's configure `Inter` as default sans.
2.  **Colors** (Light Mode Focus for strict Corca match, but must support Dark Mode):
    -   **Primary**: Light Blue Background (`#DBEAFE` / `bg-blue-100`) with Dark Blue Text (`#3B82F6` / `text-blue-600`) for buttons/accents.
    -   **Backgrounds**:
        -   `--background`: `#F3F4F6` (Gray-100) -> The app shell background.
        -   `--card`: `#FFFFFF` -> Cards float on top.
        -   `--sidebar`: `#FCFCFD` (Very light gray).
3.  **Geometry**:
    -   **Radius**: Update global radius to `0.5rem` (8px) or `0.625rem` (10px).

## Tasks
- [ ] **Task 1: Update Tailwind Config** (`apps/frontend/tailwind.config.js`)
    -   Refine color palette (ensure `primary`, `background`, `muted` map to new values).
    -   Set `borderRadius` defaults.
- [ ] **Task 2: Update CSS Variables** (`apps/frontend/src/index.css` or `App.css`)
    -   Update HSL values for `:root` and `.dark`.
- [ ] **Task 3: Font Configuration**
    -   Ensure `Inter` is loaded (already appears to be).
    -   Verify `font-heading` usage.

## Acceptance Criteria
- [ ] App background is light gray (`#F3F4F6`) in light mode.
- [ ] Cards are white with subtle borders.
- [ ] Primary buttons use the "BlueTint" style (Light blue BG, Dark blue text).
- [ ] Border radius is consistent (8px-10px).

## QA Results
- **Status**: PASS
- **Reviewer**: qa-agent
- **Date**: 2026-01-07
- **Gate File**: `docs/qa/gates/epic-15.story-72-design-system.yml`
- **Summary**: Verified Tailwind configuration and App.css variables match the Corca design specifications for both Light and Dark modes.
