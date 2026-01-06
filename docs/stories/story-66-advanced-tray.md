# Story-66: Advanced Tray Telemetry

**Epic:** [Epic-20: Advanced Traffic Scheduling](./epic-20-advanced-scheduling.md)
**Status**: Completed
**Priority**: Medium
**Basis**: `docs/gap_analysis_deep_dive.md`

## Goal
Enhance the System Tray menu to display real-time telemetry about model quotas. Currently, the tray only has Start/Stop. The goal is to show the percentage of free quota or remaining limits for key models (Gemini, Claude) directly in the menu, giving users immediate visibility without opening the dashboard.

## Context
As identified in the Gap Analysis (Section 13), the reference implementation displays:
- `Gemini High: XX%`
- `Gemini Image: XX%`
- `Claude 4.5: XX%`

This story implements that visibility in the Go/Wails version using the `systray` library.

## Tasks
- [ ] **Task 1: Telemetry Data Event**
    - File: `apps/backend/siberia/proxy/handlers/telemetry.go` (or similar)
    - Enhance the `siberia/quota` service to emit a `quota:update` event with normalized percentages for tray consumption.
- [ ] **Task 2: Tray Menu Dynamic Updates**
    - File: `apps/backend/siberia/tray/manager.go`
    - Subscribe to `quota:update` events.
    - Add/Update menu items dynamically:
        - `Gemini 1.5 Pro: 80%`
        - `Claude 3.5 Sonnet: 45%`
- [ ] **Task 3: Wiring**
    - Ensure the `TrayManager` has access to the event bus or quota service.

## Acceptance Criteria
- [ ] Use `systray` to add non-clickable menu items for Quota stats.
- [ ] Stats update at least every minute or on major usage events.
- [ ] If quota is unknown/unlimited, show nothing or "Unlimited".
- [ ] Menu items should use standard formatting (e.g., "Model Name: XX%").

## QA Results
- **Status**: PASS
- **Reviewer**: default_qa_agent
- **Date**: 2026-01-06
- **Gate File**: `docs/qa/gates/epic-20.story-66-advanced-tray.yml`
- **Summary**: All acceptance criteria met. Event emission verified in `quota/service.go`. Tray updates implemented in `tray/tray_others.go`. Wiring confirmed in `app.go`. Wiring uses context injection to ensure events are emitted to Wails runtime.

