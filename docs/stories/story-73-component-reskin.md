# Story-73: Component Reskin (Corca Adaptation)

**Epic**: [Epic-15: Visual Redesign](./epic-15-visual-redesign.md)
**Status**: Ready for Dev
**Priority**: High

## Goal
Update core UI primitives (`Button`, `Card`, `Input`, `Badge`) to strictly adhere to the "Corca" aesthetic, utilizing the design tokens established in Story-72. This involves changing specific component variants to use the new "Blue Tint" and "Soft Border" styles.

## Requirements
1.  **Buttons**:
    -   **Default/Primary**: Should use the "Blue Tint" style: Light Blue Background (`bg-blue-100` / `hsl(var(--accent))`) with Dark Blue Text (`text-blue-600` / `hsl(var(--primary))`).
    -   **Hover**: Slightly darker blue tint.
    -   **Ghost**: Standard ghost behavior but with blue text on hover.
2.  **Cards**:
    -   **Border**: Very subtle (`border-border`).
    -   **Shadow**: `shadow-sm` (gentle lift).
    -   **Background**: White (`bg-card`).
3.  **Inputs**:
    -   **Height**: Standardize to `30px` or `32px` (compact) where appropriate to match the "Density" of Corca.
    -   **Border**: `border-input`.
4.  **Badges**:
    -   Ensure Badges (used for Account Tiers, Status) match the new radius (8px usually too distinct for badges, maybe pills `rounded-full` or strictly `rounded-md`).

## Tasks
- [ ] **Task 1: Update Button Component** (`apps/frontend/src/components/ui/button.tsx`)
    -   Modify `default` variant or create a new `corca` variant.
    -   Ensure hover states are correct.
- [ ] **Task 2: Update Card Component** (`apps/frontend/src/components/ui/card.tsx`)
    -   Review border and shadow classes.
- [ ] **Task 3: Update Input Component** (`apps/frontend/src/components/ui/input.tsx`)
    -   Adjust height/padding if needed.
- [ ] **Task 4: Component Verification**
    -   Check `DashboardPage` and `AccountsPage` to ensure components render correctly.

## Acceptance Criteria
- [ ] Primary Buttons appear as Light Blue rects with Dark Blue text.
- [ ] Cards have a clean, white, minimal look on the gray background.
- [ ] Inputs feel aligned and crisp.

## QA Results
- **Status**: PASS
- **Reviewer**: qa-agent
- **Date**: 2026-01-07
- **Gate File**: `docs/qa/gates/epic-15.story-73-component-reskin.yml`
- **Summary**: Verified component reskin implementation (Buttons, Inputs, Cards) matches the Corca style guide.
