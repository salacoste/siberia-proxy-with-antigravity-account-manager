# Epic-15: Visual Redesign (Corca Adaptation)

**Status**: Completed
**Parent**: Project "Siberia"

## Executive Summary
Adapt the Siberia Proxy UI to match the visual language of `https://corca.app/`.
The goal is to transition from a generic "Shadcn/Slate" look to a distinct "Academic/Industrial Tech" aesthetic characterized by soft contrast, large border radii, and a refined blue/monochrome palette.

## Design Reference
-   **Source**: `https://corca.app/`
-   **Data Specification**: `docs/ui/design-system.md`

## Key Design Traits (The "Corca Vibe")
1.  **Palette**:
    -   **Backgrounds**: `#FFFFFF` (Card) on `#F3F4F6` (Main) or `#FCFCFD`.
    -   **Primary Action**: `#DBEAFE` (Light Blue BG) with `#3B82F6` (Blue Text).
    -   **Text**: `#000000` (Headings) and `#6B7280` (Muted).
2.  **Typography**:
    -   **Headings**: `GeistSans` (or Inter tight tracking).
    -   **Body**: `Inter`.
    -   **Scale**: Large, confident headings (H1: 80px, H2: 30px).
3.  **Geometry**:
    -   **Radius**: `6px` - `10px` (refining previous assumption of 16px).
    -   **Spacing**: Base unit `4px`. Loose, airy layouts.

## Scope (Child Stories)

### 1. Story-44: Design System Foundation (Tokens)
-   **Goal**: Update `tailwind.config.js` and `App.css` (variables).
-   **Tasks**:
    -   Import colors from `design-system.md` (Primary `#DBEAFE`, Accent `#EFF6FF`).
    -   Add `GeistSans` (or configure Font Family stack).
    -   Update `radius` variable to `0.5rem` (8px) or `0.625rem` (10px).

### 2. Story-45: Component Reskin
-   **Goal**: Update core UI primitives to match Corca style.
-   **Components**:
    -   **Buttons**: `bg-blue-100 text-blue-600` style.
    -   **Cards**: Minimal borders, subtle shadows.
    -   **Inputs**: align with 30px height patterns if applicable.

### 3. Story-46: Layout & Sidebar Polish
-   **Goal**: Apply the new aesthetic to the main shell.
-   **Tasks**:
    -   **Sidebar**: Use `GeistSans` for navigation items.
    -   **Backgrounds**: Implement the split light-gray/white hierarchy.
