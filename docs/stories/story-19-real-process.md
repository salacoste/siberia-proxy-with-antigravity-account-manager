# Story-19: Implement Real Process Termination

**Epic:** [Epic-06: Deep Integration](./epic-06-deep-integration.md)
**Status:** Completed

## Goal
Ability to find and kill a process by name (e.g., "Visual Studio Code") on macOS.

## Tasks
- [x] Refactor `modules/process` to use `gopsutil`.
- [x] Handle permission errors gracefully.

## Acceptance Criteria
- [x] Can kill "Code Helper" or generic processes.
- [x] DryRun mode validated (via interfaces).

## QA Results

### Review Date: 2026-01-04

### Reviewed By: Quinn (Test Architect)

### Code Quality Assessment
`process.go` correctly uses `shirou/gopsutil` to enumerate and terminate processes matching the name. Includes DryRun safety flag.

### Compliance Check
- Coding Standards: [✓] Interface-based design (`Manager`).
- All ACs Met: [✓] Implementation verified.

### Gate Status
Gate: PASS → docs/qa/gates/epic-06.story-19-real-process.yml

### Manual Verification (2026-01-06)
- **Method**: Indirect (via Story-13)
- **Status**: Implicitly Verified
- **Notes**: Account Activation (Story-13) successfully stopped processes before injection, proving this module works.
