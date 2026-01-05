# Full Regression Test Plan v1.2.0

**Date**: 2026-01-05
**Version**: 1.2.0-RC1
**Tester**: Antigravity (AI Agent)

## 1. Installation & Startup
- [ ] **Install**: Install the DMG/Exe on a fresh VM (Simulated).
- [ ] **Launch**: Application opens without crash.
- [ ] **First Run**: User sees "Welcome" or default Dashboard.
- [ ] **Version Check**: Sidebar shows "v1.2.0".

## 2. Proxy Core
- [ ] **Start Proxy**: Click "Start Proxy" -> Status indicator turns Green (Port 8080).
- [ ] **HTTP Traffic**: Browser configured to 127.0.0.1:8080 loads non-HTTPS sites (e.g. `http://example.com`) -> Request appears in Monitor.
- [ ] **HTTPS Traffic (No Cert)**: HTTPS site fails (CORRECT BEHAVIOR).
- [ ] **Install Cert**: Settings -> Install CA Cert -> System Prompt appears.
- [ ] **HTTPS Traffic (With Cert)**: HTTPS site loads -> Decrypted request appears in Monitor.

## 3. Account Management
- [ ] **Create Profile**: Accounts -> New Profile -> Enter Name/Color -> Save.
- [ ] **Switch Profile**: Click new profile -> Active profile changes.
- [ ] **Persistence**: Restart App -> Last used profile is active.

## 4. Cloud Features
- [ ] **Sign Up**: Cloud -> Sign Up -> Form validation works.
- [ ] **Login**: Login with valid credentials (mocked) -> specific mock behavior (`test@antigravity.dev`).
- [ ] **Sync**: Click "Sync" -> Spinner appears -> Success toast.

## 5. UI/UX
- [ ] **Theme Config**: Settings -> Dark Mode -> UI updates instantly.
- [ ] **Responsive**: Resize window -> Sidebar and Tables adjust.
- [ ] **Corca Design**: Verify "Soft Blue" buttons and "Split Background".

## 6. Security
- [ ] **DevTools**: Right click -> "Inspect" is NOT available (Production build).
- [ ] **Logs**: No PII in `~/Library/Logs/siberia/`.

## 7. Performance
- [ ] **Long Run**: Leave app running for 1 hour with traffic -> Memory usage stable < 200MB.
