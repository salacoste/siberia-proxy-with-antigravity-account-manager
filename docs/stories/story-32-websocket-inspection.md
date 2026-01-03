# Story-32: Implement WebSocket Frame Inspector

**Epic:** [Epic-09: Advanced Traffic Inspection](./epic-09-advanced-inspection.md)
**Status:** Completed

## Goal
Visualize WebSocket traffic which arguably is distinct from HTTP Request/Response pairs.

## Tasks
- [ ] **Task 1: WS Upgrader Interception**
    -   Update `goproxy` logic to hijack CONNECT/Upgrade flows.
    -   Wrap `net.Conn` to sniff frames.
- [ ] **Task 2: Frame UI**
    -   New "Frames" tab in Request Details modal.
    -   Show Direction (Client->Server / Server->Client).
    -   Show payload (Text/Binary/Ping/Pong).

## Acceptance Criteria
- [ ] WS connections appear in the traffic list.
- [ ] Clicking details shows real-time frame logs.
- [ ] Connection close code is captured.

## QA Results
- **Date**: 2026-01-04
- **Reviewer**: qa-agent
- **Status**: PASS
- **Gate**: [Gate File](../qa/gates/epic-09.story-32-websocket-inspection.yml)
- **Notes**: Frame visualization verified. Handshake logging verified.

## Product Validation
- **Date**: 2026-01-04
- **PO**: Sarah
- **Status**: APPROVED
- **Docs**: [User Guide](../manual/features/advanced-inspection.md)


