# Epic-08: HTTPS Decryption (MitM)

**Goal:** Enable inspection of encrypted HTTPS traffic by acting as a Man-in-the-Middle (MitM) proxy with a locally trusted Certificate Authority.

## Context
*   **Current State:** The proxy handles HTTPS via `CONNECT`, but acts as a blind tunnel. We cannot see headers or bodies for HTTPS requests (Status 0).
*   **Requirement:** To inspect traffic for `z.ai` or other targets, we need to decrypt it.
*   **Mechanism:** `goproxy` supports MitM. We need to generate a Root CA, trust it on the OS, and use it to sign dynamic certificates for visited domains.

## Stories

### Story-25: Certificate Authority Management (Backend)
*   **Goal:** Generate and manage a persistent Root CA (cert + key).
*   **Tasks:**
    *   Create `siberia/ca` module.
    *   Check for existing `ca.pem` / `ca.key` in app data.
    *   If missing, generate a new RSA/ECDSA Root CA with a friendly name (e.g., "Siberia Proxy CA").
    *   Expose method to get CA path.

### Story-26: Trusted Cert Installation (Backend/OS)
*   **Goal:** Automate the installation of the Root CA into the OS Trust Store.
*   **Tasks:**
    *   Implement `InstallCert()` for macOS (using `security add-trusted-cert`).
    *   (Optional MVP) Linux/Windows support (start with manual instructions or `certutil`).
    *   Implement `CheckTrust()` to verify if already installed.

### Story-27: MitM Proxy Logic (Backend)
*   **Goal:** Enable MitM in `proxy` service.
*   **Tasks:**
    *   Configure `goproxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)`.
    *   Set `goproxy.MitmConnect` handler to sign certs using our CA.
    *   Add a config toggle `MitmEnabled` (default false for safety/performance).

### Story-28: Trust UI (Frontend)
*   **Goal:** Allow user to install/verify the certificate from Settings.
*   **Tasks:**
    *   Add "HTTPS Inspection" section to `SettingsPage`.
    *   Show status: "Certificate Not Installed" / "Trusted".
    *   Button: "Install Certificate" (triggers backend `InstallCert`, might prompt generic OS password dialog).
    *   Toggle: "Enable Decryption".

## Risks
*   **Security:** CA Key must be readable only by the app.
*   **Privacy:** Decrypting banking/auth traffic is dangerous. We should support an "Ignore / Pass-through" list (e.g., `*bank*`, `*google*`).
*   **Performance:** Signing certs incurs CPU overhead.

## Mitigation
*   **MVP Scope:** Decrypt *everything* when enabled, but default to DISABLED.
*   **Future:** Whitelist/Blacklist domains.
