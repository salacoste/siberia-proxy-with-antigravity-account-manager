# Epic-13: Cloud Infrastructure

**Status**: Planned
**Phase**: 13
**Feature Branch**: `feat/epic-13-cloud-infra`

## Description
Replace the "Mock" implementations of Cloud Sync (Epic-10) with real, production-ready infrastructure.
This Epic ensures that:
1.  **Profiles** are actually synced to a persistent database (Supabase).
2.  **Shared Sessions** are actually uploaded to object storage (S3/R2) and accessible via public links.

## Goals
-   [ ] **Zero-Knowledge Persistence**: Store encrypted blobs in Supabase without the server knowing the keys.
-   [ ] **Real Authentication**: Use Supabase Auth (Email/Pass or OAuth) instead of "Mock Login".
-   [ ] **Public Sharing**: Generate valid, expirable URLs for shared HAR sessions.

## Stories
-   **Story-39**: Real Cloud Backend (Supabase Sync).
-   **Story-40**: Real File Storage (S3/R2).

## Architecture Reference
-   Authentication: Supabase Go Client (`github.com/supabase-community/supabase-go` or raw REST).
-   Storage: AWS SDK for Go v2 (compatible with R2/S3).
