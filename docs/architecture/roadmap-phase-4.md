# Phase 4 Roadmap: Deep Integration & Intelligence (v2.0)

**Status:** Draft
**Owner:** Sarah (PO)
**Date:** 2026-01-06

## 1. Executive Summary
Building on the v1.2.0 release, Phase 4 aims to make Siberia a true "Set and Forget" system utility. This involves solving the System Tray challenge for persistent operation and adding "Intelligence" to the account management (Quota/Tier detection).

## 2. Strategic Objectives

### 2.1 Persistent Presence (System Tray)
- **Problem:** Users currently have to keep the window open.
- **Solution:** Implement robustness System Tray support (Story-40).
- **Technical Challenge:** Wails v2 limitations.
- **Strategy:** Explore native Cgo bindings or upgrading to Wails v3 Alpha (Risk Assessment required).

### 2.2 Account Intelligence (Quota)
- **Problem:** Users don't know if their accounts are "Free" or "Pro" without checking manually.
- **Solution:** Automated Tier Detection (Story-41).
- **Impact:** Enables "Smart Pooling" (Use Pro accounts first).

### 2.3 User Migration
- **Problem:** Existing Python-agent users have data trapped in JSON files.
- **Solution:** Legacy Import (Story-42).

## 3. Sprint Plan

1.  **Story-40: System Tray 2.0**
    -   *Risk:* High (Dependencies).
    -   *Strategy:* Use `energye/systray` with `RunWithExternalLoop`.
    -   *Risk:* Medium (Cgo/Main Loop).
    -   *Assignee:* Dev (w/ Architect oversight).

2.  **Story-41: Quota Management**
    -   *Risk:* Medium (Reverse Engineering).
    -   *Assignee:* Dev.

3.  **Story-42: Legacy Import**
    -   *Risk:* Low.
    -   *Assignee:* Dev.

## 4. Success Metrics
-   Can close window -> App stays in Tray.
-   Account Table shows "Pro" badge automatically.

## 5. Handoff
-   **Recommended Command:** `/architect` to research System Tray solutions or design Quota Service.
