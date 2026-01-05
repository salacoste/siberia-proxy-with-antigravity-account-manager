# Story-46: Layout & Sidebar Polish (Corca Style)

**Epic**: [Epic-15: Visual Redesign](../epics/epic-15-visual-redesign.md)
**Status**: Completed

## Goal
Transform the application shell to match the Corca layout: a split background where the sidebar is transparent/gray (or white) and the main content area is distinct (floating or solid).

## Scope
1.  **Sidebar**:
    -   Font: `Geist Sans` (Heading font).
    -   Item Radius: `rounded-lg` or `rounded-xl`.
    -   Active State: `bg-white text-blue-600 shadow-sm` (if sidebar is gray) OR `bg-blue-50 text-blue-600` (if sidebar is white).
2.  **App Layout**:
    -   Background: `bg-gray-50` (or `bg-corca-main`) for main area.
    -   Logo: Ensure text is dark in light mode.
3.  **Typography**:
    -   Ensure Headings use `font-heading`.

## Acceptance Criteria
- [x] Sidebar uses `font-heading` for links.
- [x] Active sidebar item matches Corca style.
- [x] Main background is `#F3F4F6` (Gray-100).
- [x] Logo text is legible in Light Mode.

## Verification Plan
- Check `RootLayout.tsx` and `AppSidebar.tsx` (if exists).
- Verify "Monitor" and "Proxy" pages.
