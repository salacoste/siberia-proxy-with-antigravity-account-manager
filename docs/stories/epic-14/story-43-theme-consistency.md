# Story-43: Fix Theme Inconsistency (Light Mode)

**Parent Epic**: Epic-14 (Release Polish)

## Problem
When switching to Light Mode, the application retained dark backgrounds for the Sidebar and Main Content area, causing severe visual mismatch (White Cards on Black Background).
This was caused by hardcoded Tailwind classes (`bg-slate-900`, `bg-slate-950`) that ignored the theme context.

## Solution
1.  Introduced Semantic CSS Variables for Sidebar (`--sidebar-background`, etc.) in `App.css`.
2.  Updated `tailwind.config.js` to map `bg-sidebar` to these variables.
3.  Refactored `RootLayout`, `Sidebar`, and `WebSocketViewer` to use:
    -   `bg-background` instead of `bg-slate-950`
    -   `bg-sidebar` instead of `bg-slate-900`
    -   `bg-card` / `bg-muted` for containers.

## Acceptance Criteria
- [x] Light Mode: Sidebar is light (off-white).
- [x] Light Mode: Main Background is light (white).
- [x] Dark Mode: Retains original dark aesthetic.
- [x] Text is readable in both modes.

## QA Results
- **Verified By**: Antigravity Agent (Browser Subagent).
- **Status**: PASS. Verified at `http://localhost:7101/settings`.
