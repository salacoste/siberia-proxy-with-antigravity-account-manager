# Story-45: Component Reskin (Corca Style)

**Epic**: [Epic-15: Visual Redesign](../epics/epic-15-visual-redesign.md)
**Status**: Completed

## Goal
Update the core Shadcn UI components to align with the soft, airy "Corca" aesthetic, favoring subtle backgrounds and large radii over heavy borders.

## Scope
1.  **Button**:
    -   Primary/Default: `bg-primary text-primary-foreground hover:bg-primary/90`. (Blue-600)
    -   Secondary/Outline: `bg-accent text-accent-foreground border-transparent hover:bg-accent/80`. (Blue-100)
2.  **Card**:
    -   Remove heavy borders (`border-border`).
    -   Add subtle shadow (`shadow-sm` or `shadow-md`).
    -   Background `bg-white` (Light) or `bg-zinc-950` (Dark).
3.  **Input**:
    -   Softer border (`border-gray-200`).
    -   Focus ring (`ring-blue-500`).
    -   Height `h-9` or `h-10`.

## Acceptance Criteria
- [x] `components/ui/button.tsx` updated.
- [x] `components/ui/card.tsx` updated.
- [x] `components/ui/input.tsx` updated.
- [x] Visual check: Buttons feel "clickable" and "soft".

## Verification Plan
- Check `SettingsPage` (uses all 3 components).
- Ensure no layout shift or breaking changes.
