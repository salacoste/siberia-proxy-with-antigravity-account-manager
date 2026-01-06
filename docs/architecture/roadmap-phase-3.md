# Phase 3 Roadmap: The Last Mile

**Status:** Released
**Owner:** Sarah (PO)
**Date:** 2026-01-06

## 1. Executive Summary
The project has achieved significant maturity with Core Proxy, Accounts, Cloud Sync, and Visual Redesign completed. Phase 3 focuses on **Deep Integration** (making the app invisible/native) and **Launch Readiness** (hardening for public release).

## 2. Strategic Objectives

### 2.1 Deep Integration (Native UX)
Move the application from "just a window" to a persistent system utility.
-   **System Tray (Story-40):** [BLOCKED] Persistent background presence (Deferred to v2).
-   **Terminal Helper (Story-38):** [DONE] Shell integration (`siberia set-env`) for CLI tools.

### 2.2 Documentation Hygiene
Resolve accumulated technical debt in the project management artifacts.
-   **ID Collision Resolution:** Resolve conflicts between Epic-04/05/09 story IDs (e.g., Story-16, Story-44 duplication).
-   **Epic Closure:** Formally archive completed epics and updated their statuses from "Draft" to "Released".

### 2.3 Launch Readiness (Epic-17)
Prepare the binary for the wild.
-   **Security Audit:** [DONE] Final check of PII logging (Story-51).
-   **Signed Builds:** Release automation with signing (Story-50).
-   **Regression Testing:** [DONE] Full manual pass.

## 3. Immediate Action Plan (Sprint 6)

1.  **Prioritized Feature:** `Story-40: System Tray Integration`.
    -   *Why:* Critical for "always-on" proxy usage pattern.
    -   *Assignee:* Dev Agent.

2.  **Secondary Feature:** `Story-38: Terminal Helper`.
    -   *Why:* Completes the "Developer Experience" aspect.

3.  **Ops Task:** Documentation Refactor.
    -   *Why:* Ensures future maintainability.
    -   *Assignee:* PM/PO Agent.

## 4. Architect's Notes
-   The "Intelligent Proxy" layer (Mappers, Handlers) appears complete (Stories 35-39). Verify integration tests for these during Release Polish.
-   **Risk:** `MapLocalPage` had build issues (patched in QA). This technical debt needs a permanent fix before v1.0.

## 5. Handoff
-   **Recommended Next Command:** `/dev` to start **Story-40 (System Tray)**.
