# Story-47: Dashboard Layout & Widget System

**Parent**: Epic-16
**Status**: Completed

## Description
Create the foundational UI structure for the Analytics Dashboard. This involves setting up the route, the grid layout system, and the base "Widget" component container.

## Requirements
1.  **Route**: Ensure `/dashboard` is accessible (or updates the Home view).
2.  **Layout**:
    -   Use a responsive Grid layout (CSS Grid).
    -   Support "span" (widgets that take 1, 2, or 3 columns).
3.  **Components**:
    -   `DashboardGrid`: Container.
    -   `WidgetCard`: Base container with Title, Content Area, and consistent styling (Corca `card` tokens).
4.  **Skeleton**: Implement a static "Mock" version of the dashboard with placeholder widgets.

## Acceptance Criteria
- [x] `/dashboard` route renders correctly.
- [x] Grid system handles resizing gracefully.
- [x] `WidgetCard` looks consistent with the Design System (Story-45).
- [x] Placeholder widgets (RPS, Bandwidth) are visible.

## Dev Agent Record

### Completion Notes
- **Status Update**: Retroactively marked as **Completed**.
- **Implementation**: The dashboard layout, grid system, and `WidgetCard` component were implemented and verified during the development of **Story-49**.
- **Verification**: Browser verification in Story-49 confirmed the grid layout works and is responsive.

### Story DoD Checklist
- [x] Functionality verified in Story-49.
- [x] Artifacts exist (`DashboardPage.tsx`, `WidgetCard.tsx`).
- [x] Design tokens match Project.
