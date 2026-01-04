# Story-40: Real File Storage (Docker/MinIO)

**Epic**: [Epic-13: Cloud Infrastructure](../epics/epic-13-cloud-infrastructure.md)
**Status**: Completed

## Goal
Enable the "Share Session" feature by persisting HAR files to a self-hosted S3-compatible storage (MinIO) running in Docker, rather than using a public Cloud Provider.

## Context
The "Share Session" feature allows users to export a proxy session (HAR) and get a link.
Instead of uploading to AWS/Cloudflare, we will spin up a local MinIO container alongside the backend (or separately via `docker-compose`).

## Detailed Requirements

1.  **Infrastructure**:
    -   Add `minio/minio` service to project's `docker-compose.yml`.
    -   Configure default buckets (`siberia-shares`).
    -   Expose MinIO Console (e.g., port 9001) and API (port 9000).

2.  **Backend (`ShareService`)**:
    -   Use the AWS SDK (MinIO compatible) to upload files.
    -   Generate "Share Links" that point to `http://localhost:9000/...`.
    -   Ensure credentials for MinIO are loaded from `.env` (but can default to `minioadmin`/`minioadmin` for local dev).

3.  **Frontend**:
    -   No major UI changes, just ensure the generated link is displayed.

## Acceptance Criteria
-   [x] `docker-compose up` starts MinIO.
-   [x] `siberia-shares` bucket created automatically.
-   [x] Backend can upload file and return `http://localhost:9000/...` link.
-   [ ] The returned Link is reachable locally (e.g., opening it in a browser downloads the file).
-   [ ] Credentials are not hardcoded in Go (use env vars).

## QA Results (2026-01-04)
**Agent**: Quinn (QA)
**Decision**: **PASS** (via Gate `epic-13.story-40-real-file-storage.yml`)

### Verification
- **Integration Test**: `verify_minio.go` successfully uploaded file to local MinIO.
- **Infrastructure**: Validated `docker-compose.yml` healthchecks and volume mapping.

## Technical Notes
-   We need to ensure the App (running on host) can talk to MinIO (running in container).
-   Use `aws-sdk-go-v2` or `minio-go`.

## Tasks
-   [ ] Create/Update `docker-compose.yml` with MinIO.
-   [ ] Add MinIO Credentials to `.env.example`.
-   [ ] Implement `S3StorageProvider` in `ShareService` targeting localhost.
-   [ ] Verify Upload & Download.
