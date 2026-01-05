# Story-44: Design System Foundation (Tokens)

**Parent**: Epic-15
**Status**: Ready

## Description
Implement the core design tokens derived from the Corca.app analysis. This involves updating Tailwind configuration, CSS variables, and font assets.

## Requirements
1.  **Typography**:
    -   Install/Add `GeistSans` (or use a compatible web font host/local file).
    -   Update `tailwind.config.js` to set `font-sans` to `Inter` (Body) and `font-heading` to `GeistSans`.
2.  **Colors**:
    -   Update `primary` to `#DBEAFE` (Light Blue) and `primary-foreground` to `#3B82F6` (Dark Blue) for that specific "Corca Button" look.
    -   Update `background` to `#FFFFFF` and `muted/secondary` to `#F3F4F6`.
3.  **Radius**:
    -   Update global radius to `0.6rem` (approx 10px) to match the "friendly" aesthetic.

## Acceptance Criteria
- [x] Tailwind Config uses `Inter` and `GeistSans`.
- [x] CSS Variables reflect the Corca palette (Light Blue actions).
- [x] Border Radius is updated globally.
