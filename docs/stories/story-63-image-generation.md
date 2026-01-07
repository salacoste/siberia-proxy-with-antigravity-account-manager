# Story-63: Image Generation Handler

**Epic:** [Epic-19: Z.ai Intelligence](./epic-19-zai-intelligence.md) (or Core Proxy)
**Status**: Completed
**Priority**: Medium
**Basis**: `docs/gap_analysis_deep_dive.md`

## Goal
Support the OpenAI Image Generation API (`/v1/images/generations`) by translating requests to the upstream provider (Gemini/Z.ai). The reference app maps these to `gemini-3-pro-image` or similar models, handling aspect ratio mapping and prompt enhancement.

## Tasks
- [ ] **Task 1: Handler Implementation**
    - File: `apps/backend/proxy/handlers/openai_images.go`
    - Implement `POST /v1/images/generations`.
    - Map OpenAI `size` (e.g., "1024x1024") to aspect ratios if needed.
    - Implement Prompt Enhancement ("vivid", "hd" styles appending to prompt).
- [ ] **Task 2: Upstream Dispatch**
    - Route to `image_gen` account pool.
    - Handle concurrency for `n > 1` (parallel requests).

## Acceptance Criteria
- [ ] Standard OpenAI Image Request works.
- [ ] `size: 1920x1080` is correctly handled (or mapped).
- [ ] `style: vivid` modifies the prompt sent to upstream.

## Dev Agent Record

### File List
#### [NEW] [apps/backend/siberia/proxy/mappers/image.go](file:///Users/r2d2/Documents/Code_Projects/siberia-proxy-with-antigravity-account-manager/apps/backend/siberia/proxy/mappers/image.go)
#### [NEW] [apps/backend/siberia/zai/image.go](file:///Users/r2d2/Documents/Code_Projects/siberia-proxy-with-antigravity-account-manager/apps/backend/siberia/zai/image.go)
#### [NEW] [apps/backend/siberia/proxy/handlers/openai/images.go](file:///Users/r2d2/Documents/Code_Projects/siberia-proxy-with-antigravity-account-manager/apps/backend/siberia/proxy/handlers/openai/images.go)
#### [MODIFY] [apps/backend/siberia/proxy/upstream/client.go](file:///Users/r2d2/Documents/Code_Projects/siberia-proxy-with-antigravity-account-manager/apps/backend/siberia/proxy/upstream/client.go)
#### [MODIFY] [apps/backend/siberia/proxy/upstream/gemini.go](file:///Users/r2d2/Documents/Code_Projects/siberia-proxy-with-antigravity-account-manager/apps/backend/siberia/proxy/upstream/gemini.go)
#### [MODIFY] [apps/backend/siberia/proxy/service.go](file:///Users/r2d2/Documents/Code_Projects/siberia-proxy-with-antigravity-account-manager/apps/backend/siberia/proxy/service.go)
#### [MODIFY] [apps/backend/siberia/proxy/handlers/openai/handler_test.go](file:///Users/r2d2/Documents/Code_Projects/siberia-proxy-with-antigravity-account-manager/apps/backend/siberia/proxy/handlers/openai/handler_test.go)

### Debug Log
- Tests Passed: `go test ./apps/backend/siberia/proxy/... ./apps/backend/siberia/zai/...`

