//go:build darwin

package ca

import (
	"fmt"
	"os"
	"os/exec"
)

// InstallCert installs the root CA into the macOS Keychain
func (s *Service) InstallCert() error {
	certPath, _ := s.GetCAPath()
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return fmt.Errorf("CA certificate not found at %s", certPath)
	}

	// Dynamic keychain detection
	keychainPath := os.Getenv("HOME") + "/Library/Keychains/login.keychain-db"
	if _, err := os.Stat(keychainPath); os.IsNotExist(err) {
		// Fallback for older macOS
		keychainPath = os.Getenv("HOME") + "/Library/Keychains/login.keychain"
	}

	// Command: security add-trusted-cert -d -r trustRoot -k <keychain> <certPath>
	// -d: Add to admin cert store (requires auth, but applies to system/admin level)
	// -r trustRoot: Trust as a Root CA
	// -k: Specific keychain

	cmd := exec.Command("security", "add-trusted-cert", "-d", "-r", "trustRoot", "-k", keychainPath, certPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install cert: %v, output: %s", err, string(output))
	}

	return nil
}

// CheckTrust checks if the cert is trusted using 'security verify-cert'
func (s *Service) CheckTrust() bool {
	certPath, _ := s.GetCAPath()
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return false
	}

	// -c: Cert file
	// -l: Leaf verification (sanity check)
	// -L: Local cert checking (don't network check revocation, though for root CA irrelevant)
	// Note: verifying a Root CA with verify-cert asks if it's trusted.
	cmd := exec.Command("security", "verify-cert", "-c", certPath, "-l", "-L")

	// We only care if the command succeeds (exit code 0)
	if err := cmd.Run(); err != nil {
		// e.g. "Cert Verify Failed: CSSMERR_TP_NOT_TRUSTED"
		return false
	}

	return true
}
