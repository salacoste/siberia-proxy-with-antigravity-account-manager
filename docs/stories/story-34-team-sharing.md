# Story-34: Implement "Share Session" Feature

**Epic:** [Epic-10: Team Collaboration & Cloud Sync](./epic-10-collaboration.md)
**Status:** Draft

## Goal
Generate a shareable public/private link for a specific captured HTTP request/response pair (HAR format).

## Tasks
- [ ] **Task 1: HAR Export**
    -   Convert internal log format to standard HAR.
- [ ] **Task 2: Upload & Link Gen**
    -   Upload JSON to temporary storage (e.g., S3/R2 presigned URL).
    -   Generate link `https://share.siberia.dev/log/<id>`.

## Acceptance Criteria
- [ ] "Share" button in Traffic Inspector.
- [ ] Link opens a web-view of that specific request.
- [ ] Option to set expiration (1h, 1d).
