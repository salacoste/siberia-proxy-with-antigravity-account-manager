# Architecture Decision: System Tray Strategy

**Status:** Proposed
**Date:** 2026-01-06
**Author:** Winston (Architect)

## 1. Problem
The application requires a persistent background presence (System Tray/Menu Bar) to allow users to toggle the proxy without keeping the main window open. Wails v2.11.0 lacks robust, unified APIs for this.

## 2. Options Analysis

### Option A: Wails v2 + `energye/systray` (Recommended)
Use an external Cgo-based library to manage the tray.
- **Pros:** Works with current Wails v2 tech stack. `energye` fork offers better event handling (Right Click, etc.) and `RunWithExternalLoop` to avoid fighting Wails' main loop.
- **Cons:** Requires Cgo. Potential stability issues on macOS if not handled correctly (NSApplication conflicts).
- **Risk:** Medium.

### Option B: Wails v3 Alpha
Upgrade the entire stack to Wails v3.
- **Pros:** Native, beautiful API. No Cgo hacks needed.
- **Cons:** Alpha software. API flux. Potential breaking changes.
- **Risk:** High.

### Option C: Native Bridge
Write custom Objective-C/Swift bridge for macOS and DLL calls for Windows.
- **Pros:** Maximum control. Zero external dependencies.
- **Cons:** High maintenance burden. Requires specialized platform knowledge.
- **Risk:** High (Cost).

## 3. Decision
**Proceed with MVP using `energye/systray`**.
We will utilize the `RunWithExternalLoop` pattern to integrate with the Wails application lifecycle.

## 4. Implementation Plan (Story-40)

1.  **Dependency**: Add `github.com/energye/systray`.
2.  **Lifecycle**:
    -   Initialize Tray in `main.go`.
    -   Run Tray logic in a separate goroutine or via external loop hook.
3.  **Events**:
    -   Bind `OnReady` to setup menu.
    -   Bind `OnExit` to cleanup.
4.  **Integration**:
    -   Expose `UpdateTrayState(bool)` via `App` struct to change icon (Green/Red).

## 5. Fallback
If `energye/systray` proves unstable on macOS, we will pivot to a simple "Dock Menu" (macOS only) via `wails.json` config if available, or defer to Wails v3.
