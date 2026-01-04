# Story-40: Real File Storage (S3/R2)

**Epic**: [Epic-13: Cloud Infrastructure](../epics/epic-13-cloud-infrastructure.md)
**Status**: Draft

## Goal
Replace the mock `UploadSession` in `siberia/share` with a real S3 uploader.

## Requirements
1.  **Storage Provider**: Cloudflare R2 or AWS S3 (Configurable).
2.  **Bucket Structure**: `sessions/{random_id}.json` (or `.har`).
3.  **Lifecycle**:
    -   (Optional) Auto-delete after 24h via Bucket Policy.
4.  **Public Access**:
    -   Files should be accessible via a public domain (e.g., `https://share.siberia.dev/...`).

## Technical Notes
-   Use `aws-sdk-go-v2`.
-   Credentials should be embedded (limited scope token) or proxied via a backend function (if we had one). For this Desktop App, we might need a "Presigned URL" generator or a restricted API key.
-   *Decision point for Dev*: Can we embed a "Write Only" token? Or do we need a lambda? -> *ADR Required*.

## Acceptance Criteria
-   [ ] "Share" button uploads actual JSON to the bucket.
-   [ ] The returned Link is reachable via Browser (curl).
