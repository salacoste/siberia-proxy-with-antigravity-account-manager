# Story-69: Native Gemini Protocol Support

**Epic:** [Epic-02: Proxy Engine](./epic-02-proxy.md)
**Status**: Completed
**Priority**: High
**Basis**: `proxy/handlers/gemini.rs`

## Goal
Support standard Google Gemini API requests (`POST /v1beta/models/...:generateContent`) directly, handling them with the same sticky-session and rotation logic as OpenAI requests.

## Tasks
- [x] **Task 1: Handler Implementation**
    - File: `apps/backend/proxy/handlers/gemini.go`
    - Handle `generateContent`, `streamGenerateContent`, `countTokens`.
    - Support path parameters like `models/gemini-1.5-pro:generateContent`.
- [x] **Task 2: Request/Response Wrapping**
    - Reuse wrapping logic to inject `project_id`.
    - Handle SSE format (which differs slightly effectively, just passing through but wrapping/unwrapping).

## Acceptance Criteria
- [x] `curl ...:generateContent` works against the proxy.
- [x] Sticky sessions work for Gemini protocol requests.
- [x] Quota errors trigger rotation.

## QA Results
- **Date**: 2026-01-06
- **Reviewer**: Quinn (QA)
- **Status**: **PASS**
- **Notes**:
    - Native handler implementation verified.
    - Model ID parsing logic is robust.
    - Unit tests confirm request/response mapping.
    - Ready for merge.


