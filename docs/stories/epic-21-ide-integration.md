# Epic-21: Zero-Config Onboarding (Deep Integration)

**Status**: Completed
**Priority**: Medium (Tier 3 - UX)
**Basis**: `docs/gap_analysis_deep_dive.md` (Section 8)

## Executive Summary
Reduce user friction by implementing "Deep State Token Injection". This feature allows the application to detect existing installations of JetBrains IDEs (IntelliJ, Android Studio), read their local SQLite configuration databases, and extracting authenticated sessions (Refresh Tokens) directly. This enables users to start using the proxy immediately without manual credential entry.

## Key Features
1.  **IDE Database Discovery**:
    -   Scan standard paths (check both macOS and Linux/Windows conventions) for `~/.config/Google/AndroidStudio*` or `~/Library/Application Support/JetBrains/*/options/other.xml` (or specific SQLite DBs).
2.  **Protobuf Decoding**:
    -   Read the binary blob from `jetskiStateSync`.
    -   Decode nested Protobuf fields (Field 6 -> Field 3) to extract the refresh token string.
3.  **Migration Wizard**:
    -   UI to show found accounts and offer "One-Click Import".

## Success Metrics
-   [ ] Application detects installed Android Studio session on startup.
-   [ ] Correctly extracts valid Refresh Token from the binary definition.
-   [ ] Extracted token successfully validates and refreshes quota.
