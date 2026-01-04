# Story-46: Sidebar & Layout Polish

**Parent**: Epic-15
**Status**: Pending

## Description
Update the application shell (Sidebar, Header, Main Area) to fully utilize the new design system.

## Requirements
1.  **Sidebar**:
    -   **Contrast**: Increase distinction from main body. (Option: Slightly darker gray background `bg-sidebar` or distinct right border).
    -   **Logo Text**: Ensure text next to logo is **DARK** in light mode (currently white/invisible).
2.  **Layout**:
    -   **Background**: Change global body background to a visible off-white/gray (e.g., `#F3F4F6` or `#F9FAFB`) to allow white cards to "pop".
    -   **Depth**: Ensure white content cards have appropriate `shadow-sm` or `shadow` to separate from the new gray background.
3.  **Typography**:
    -   Update Page Titles to use `GeistSans` and larger 30px+ size where appropriate.

## Acceptance Criteria
- [x] Sidebar is clearly distinct from main content.
- [x] Logo text is legible in Light Mode.
- [x] Global background is off-white; Cards are white and pop.
- [x] Page Headings use new font family.
