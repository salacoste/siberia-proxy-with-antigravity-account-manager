# Story-31: Request Breakpoint & Rewrite System

**Epic:** [Epic-09: Advanced Traffic Inspection](../epic-09-advanced-inspection.md)
**Status:** Completed
**Parent:** Epic-09

## Description
Implement a breakpoint system that allows users to pause specific HTTP requests or responses matching defined criteria. When paused, the user should be able to inspect and modify the content (Headers/Body) before resuming (Action: Allow/Inject).

## Requirements
1.  **Breakpoint Rules**:
    -   Define rules by URL pattern (regex/glob) and Method.
    -   Toggle enabled/disabled state for rules.
2.  **Interception State**:
    -   When a request matches an active rule, the proxy (Go backend) must hold the connection.
    -   Frontend receives a "Pending Request" event.
    -   Global "Pause" mode (already exists) should be integrated or distinct? *Decision: Breakpoints are selective pauses.*
3.  **UI Interception Panel**:
    -   Alert/Modal when a request is intercepted.
    -   Editable Headers (Key-Value editor).
    -   Editable Body (Text area).
    -   Actions: "Forward", "Drop" (optional), "Reply with..." (optional).
4.  **Backend Logic**:
    -   Ensure `goproxy` hooks wait for a signal from the frontend (channels or maps).
    -   Timeout logic (don't hang forever if frontend disconnects).

## Acceptance Criteria
- [ ] Users can add a breakpoint rule (e.g., `*api/checkout*`).
- [ ] Requests matching the rule are paused; frontend shows "Pending Interception".
- [ ] User can edit the request body (e.g., change `{ "price": 100 }` to `{ "price": 1 }`) and headers.
- [ ] Clicking "Forward" sends the modified request to the destination.
- [ ] (Bonus) Intercept Responses as well.

## Dev Notes
-   Refine `BreakpointPanel` component.
-   Backend needs a map of `pendingID -> channel` to hold requests.
-   Check `PendingRequestDialog` implementation status.
## Dev Agent Record
-   **Status**: Dev Complete
-   **Verification**: Unit Tests Passed, Build Passed.
-   **Artifacts**: [Walkthrough](../../../../../.gemini/antigravity/brain/be569adf-0306-4d1c-af4e-4dccfcc5e976/walkthrough.md)

### Changes
1.  Backend: `BreakpointManager` implemented in `breakpoint.go`.
2.  Backend: `GetBreakpointRules` exposed in `App`.
3.  Frontend: `MonitorPage` loads rules on mount.
4.  Frontend: `PendingRequestDialog` wired to `ResumeRequest`.
### QA Results
- **Decision**: PASS
- **Confidence**: HIGH 🟢
- **Gate File**: [Gate Decision](../../qa/gates/epic-09.story-31-breakpoints.yml)
- **Notes**: Verified implementation against logic requirements. Unit tests confirm robustness of BreakpointManager. Frontend components are wired correctly.
