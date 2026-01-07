# Story-74: Layout & Sidebar Polish (Corca Adaptation)

**Epic**: [Epic-15: Visual Redesign](./epic-15-visual-redesign.md)
**Status**: Ready for Dev
**Priority**: High

## Goal
Apply the final "Corca" aesthetic polish to the application shell. This focuses on the Sidebar typography and the main content area's background hierarchy to ensure the app feels like a cohesive professional tool.

## Requirements
1.  **Sidebar**:
    -   **Font**: Update navigation items to use `GeistSans` (or `Inter tight`) for a technical look.
    -   **Active State**: Ensure the active item uses the "Blue Tint" style (already partially verified in Story-73, but needs explicit layout confirmation).
    -   **Background**: Ensure sidebar contrasts correctly with the main content (Light Gray vs White).
2.  **Main Layout**:
    -   **Background Split**: The main "App Shell" should be Light Gray (`#F3F4F6` or `bg-muted/30`).
    -   **Content Area**: The actual page content (dashboard widgets, tables) should sit on this gray background, effectively "floating" or being distinct.
3.  **Spacing**:
    -   Verify global padding matches the "Airy" Corca vibe (looser spacing than standard Shadcn).

## Tasks
- [ ] **Task 1: Update Sidebar Typography** (`apps/frontend/src/components/layout/Sidebar.tsx`)
    -   Apply font-family utility.
- [ ] **Task 2: Global Layout Background** (`apps/frontend/src/components/layout/RootLayout.tsx`)
    -   Verify `bg-background` application produces the desired gray canvas.
- [ ] **Task 3: Spacing Review**
    -   Adjust main content padding if too tight.

## Acceptance Criteria
- [ ] Sidebar navigation uses the correct technical font.
- [ ] Application has a distinct "Page on Canvas" feel (Gray background, White Cards).
- [ ] No visual clash between the new Blue components and the layout.

## QA Results
- **Status**: PASS
- **Reviewer**: qa-agent
- **Date**: 2026-01-07
- **Gate File**: `docs/qa/gates/epic-15.story-74-layout-sidebar.yml`
- **Summary**: Verified layout polish. Sidebar uses correct technical typography, and the gray application background provides excellent separation for content cards.
