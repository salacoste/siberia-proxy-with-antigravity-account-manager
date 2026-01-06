# Story 53: Display Account Tiers in UI

**Epic**: Epic-03 Accounts Manager

## Goal
As a user, I want to see visual indicators (badges) in the Accounts list that show whether an account is on a "Free" or "Paid" tier, so I can easily distinguish capabilities.

## Context
The Backend `QuotaService` correctly calculates tiers, but the Frontend `AccountsPage` currently hides this information.

## Tasks
- [ ] Update `AccountDTO` in backend to explicitly include `Tier` or `BadgeLabel`.
- [ ] Add "Tier" column to the Accounts Table in `AccountsPage.tsx`.
- [ ] Add visual badges (e.g., Gray for Free, Gold/Green for Paid).

## Acceptance Criteria
- [x] "Paid" accounts show a "Pro" or "Paid" badge.

## QA Results

### Manual Verification (2026-01-06)
- **Method**: Browser Subagent + Code Injection
- **Steps**:
    1. Verified default "Free" (available in list).
    2. Injected `Tier="Pro"` in backend.
    3. Verified "Pro" badge (Amber/Gold).
- **Result**: PASS
- **Evidence**: `accounts_tier_badges_1767698898447.png`
