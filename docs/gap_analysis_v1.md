# Gap Analysis: Reference vs Current Implementation (v1.2.1)

**Date**: 2026-01-06
**Scope**: Comparison of `docs/ag-manager-ref` (Legacy Codebase) against current `docs/stories` implementation.

## 1. Feature Gaps

### A. Native OAuth Integration (High Priority Gap)
*   **Reference Implementation**: `src-tauri/src/modules/oauth.rs`
    *   Spawns a local HTTP server (`127.0.0.1` + IPv6 fallback).
    *   Generates Google OAuth URL with proper scopes (`cloud-platform`).
    *   Handles callback code exchange automatically.
*   **Current Implementation**: `Story-12: Add Account`
    *   **Manual Entry**: User acts as the "human bridge", manually pasting tokens or credentials.
    *   **Impact**: High friction for adding new accounts. Missing "Sign in with Google" button.

### B. Smart Token Pooling (Performance Gap)
*   **Reference Implementation**: `token_manager.rs`
    *   **Sticky Sessions**: Routes requests from the same `session_id` to the same account to leverage upstream context caching.
    *   **Tier Prioritization**: Sorts accounts (Ultra > Pro > Free) to use faster-resetting accounts first.
    *   **Smart Wait**: Calculates exact wait time for rate limits.
*   **Current Implementation**: `Story-38: Upstream Resilience`
    *   **Logic**: Simple Account Rotation (Round-Robin) on 429/401 errors.
    *   **Impact**: Potential context thrashing and suboptimal usage of high-tier accounts.

### C. Deep State Injection (Migration Gap)
*   **Reference Implementation**: `migration.rs`
    *   Scans IntelliJ / Android Studio `jetskiStateSync` (Protobuf) binary files to extract tokens from currently authenticated IDE sessions.
*   **Current Implementation**: `Story-42: Legacy Import`
    *   Scans `~/.antigravity-agent/antigravity_accounts.json` (Text/JSON based).
    *   **Impact**: Users must have used the old CLI tool; cannot import directly from IDE state.

### D. Advanced Tray Telemetry (UX Gap)
*   **Reference Implementation**: `tray.rs`
    *   Dynamic submenu showing usage breakdown per model (e.g., "Gemini High: 50%").
*   **Current Implementation**: `Story-40: System Tray`
    *   Simple Start/Stop and Global Status.
    *   **Impact**: Less visibility into granular quota usage without opening the dashboard.

## 2. Recommendations

### For v1.3.0
1.  **Implement Sticky Sessions**: Add session affinity to the Proxy Service to improve upstream performance.
2.  **Native OAuth**: Reduce friction by implementing the Go-based OAuth Loopback (porting `oauth_server.rs` logic).

### For v2.0
1.  **Deep Import**: Add Protobuf scanning for IDEs (Complex, fragile to IDE updates).
2.  **Tray Telemetry**: Align Tray menu with new Analytics Dashboard data.
