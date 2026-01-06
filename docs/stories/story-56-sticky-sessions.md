# Story-56: Session Fingerprinting & Sticky Sessions

**Epic:** [Epic-20: Advanced Traffic Scheduling](./epic-20-advanced-scheduling.md)
**Status**: Draft
**Priority**: Medium
**Basis**: `docs/gap_analysis_deep_dive.md`

## Goal
Implement "Sticky Sessions" to improve upstream cache hit rates (KV Caching). By fingerprinting incoming requests based on their model and initial prompt, the proxy can consistently route related requests to the same upstream account, mimicking a stateful session even for stateless clients.

## Tasks
- [ ] **Task 1: Session Fingerprinting**
    - File: `apps/backend/proxy/session/fingerprint.go`
    - Implement `GenerateFingerprint(model string, messages []Message) string`.
    - Logic: Hash `model_name` + `first_user_message_content`.
- [ ] **Task 2: Sticky Routing Logic**
    - File: `apps/backend/proxy/loadbalancer.go`
    - Update routing logic: Check if a fingerprint points to an active, valid account. If yes, prioritize that account.
- [ ] **Task 3: Session Header Support**
    - Respect specific headers (e.g., `x-siberia-session-id`) if provided by the client, overriding the implicit fingerprint.

## Acceptance Criteria
- [ ] Two requests with the exact same first prompt are routed to the same Account ID (assuming the account is healthy).
- [ ] Requests with explicit `x-siberia-session-id` always map to the same account.
