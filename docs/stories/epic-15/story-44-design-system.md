# Story-44: Design System Foundation (Corca)

**Epic**: [Epic-15: Visual Redesign](../epics/epic-15-visual-redesign.md)
**Status**: Completed

## Goal
Establish the visual foundation (colors, typography, radius) inspired by the Corca Design System.

## Requirements
1.  **Colors (Tailwind)**:
    -   Define `bg-corca-main`: `#F3F4F6` (Gray-100/50 mix).
    -   Define `bg-corca-action`: `#DBEAFE` (Blue-100).
    -   Define `text-corca-action`: `#2563EB` (Blue-600).
2.  **Typography**:
    -   Adopt `GeistSans` (or `Inter` via Google Fonts) as primary font.
3.  **Radius**:
    -   Increase Global Radius `--radius` to `0.5rem` or `0.6rem`.
4.  **Dark Mode**:
    -   Ensure the new palette has sensible Dark Mode mappings (e.g., `#0F172A` for Background).

## Acceptance Criteria
- [ ] `App.css` variables updated.
- [ ] `tailwind.config.js` updated with custom colors/fonts.
- [ ] Sample UI element (e.g., Button) reflects new radius and color.

## Verification Plan
- Inspect `npm run dev` in browser.
- Verify Dark/Light mode switching.
