# Story-75: Z.ai Vision & Multimodal Tools

**Epic:** [Epic-19: Z.ai Intelligence](./epic-19-zai-intelligence.md)
**Status**: Completed
**Priority**: High
**Basis**: `docs/gap_analysis_deep_dive.md`

## Goal
Implement the multimodal capabilities found in the reference `zai_vision_tools.rs`. This allows the application to process images and videos using Z.ai's `glm-4.6v` model, enabling features like "UI to Code", "Error Diagnosis", and "Diagram Understanding".

## Tasks
- [ ] **Task 1: Vision Client**
    - File: `apps/backend/zai/vision.go`
    - Implement `vision_chat_completion` that supports `thinking: {type: "enabled"}` schema.
    - Implement `image_source_to_content` handling Base64 encoding and size checks.
- [ ] **Task 2: Tool Implementation**
    - Port the following tools:
        - `ui_to_artifact`: Screenshot -> Code/Spec.
        - `extract_text_from_screenshot`: OCR.
        - `diagnose_error_screenshot`: Debugging.
        - `understand_technical_diagram`.
        - `ui_diff_check`.
- [ ] **Task 3: MCP Registration**
    - Register these as tools in the Internal MCP Server (Story-58).

## Acceptance Criteria
- [ ] Can send an image path to `ui_to_artifact` and receive generated frontend code.
- [ ] Large images (>5MB) are rejected or compressed (logic to be decided, ref app rejects).
- [ ] Supports both HTTP URLs and local paths.
- [ ] Supports both HTTP URLs and local paths.

## Dev Agent Record
### File List
- `apps/backend/siberia/zai/vision.go` (New)
- `apps/backend/siberia/zai/tools/vision_tools.go` (New)
- `apps/backend/siberia/zai/vision_test.go` (New)
- `apps/backend/siberia/mcp/registry.go` (Modified)
- `apps/backend/siberia/config/config.go` (Checked)

## QA Results
- **Status**: PASS
- **Date**: 2026-01-07
- **Gate**: [epic-19.story-75-zai-vision.yml](../qa/gates/epic-19.story-75-zai-vision.yml)
- **Notes**: All unit tests passed. Integration confirmed via `registry` tests.
