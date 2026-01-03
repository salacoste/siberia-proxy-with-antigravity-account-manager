# 🗺️ Project Roadmap: Siberia Proxy (Post-v1.0)

## 🌟 Current Status: v1.0.0 (Ready for Release)
- **Core Features:** Complete (Proxy, Accounts, Monitor, MitM, Auto-Update).
- **Stability:** QA Validated (Epics 01-08).
- **CI/CD:** Automated builds and release generation fully functional.

## 🚀 Phase 2: Refinement & Scale (Q1 2026)

### Epic-09: Advanced Traffic Inspection
- [ ] **Regex Filtering:** Allow users to filter traffic logs by complex patterns.
- [ ] **Breakpoint / Rewrite:** Ability to pause requests and edit body/headers on the fly (like Charles Proxy).
- [ ] **WebSocket Inspection:** Decode and display WS frames.

### Epic-10: Team Collaboration & Cloud Sync
- [ ] **Cloud Profiles:** Sync accounts/proxy settings across devices (requires cloud backend).
- [ ] **Team Sharing:** Share "Proxy Groups" or "Captured Sessions" with teammates via link.

### Epic-11: Performance Tuning
- [ ] **High Throughput Mode:** Optimize `goproxy` for >10k req/s (currently optimized for dev/debug).
- [ ] **Memory Management:** Better distinct handling of large traffic logs in the frontend (virtualization optimization).

### Epic-12: Native Integrations beyond VS Code
- [ ] **Cursor / Windsurf Support:** Verify injection paths for other AI IDEs.
- [ ] **Terminal Injection:** `export ALL_PROXY=...` helper for shell sessions.
