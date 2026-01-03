//go:build darwin

package ca

import (
	"fmt"
	"os"
	"os/exec"
)

// InstallCert installs the root CA into the macOS Keychain
func (s *Service) InstallCert() error {
	certPath := s.GetCAPath()
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return fmt.Errorf("CA certificate not found at %s", certPath)
	}

	// Command: security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain <certPath>
	// or user keychain. "login.keychain" is standard for user.
	// Using /Library/Keychains/System.keychain requires sudo and makes it system-wide.
	// Using login.keychain is easier but might not cover all use cases (System services).
	// For a dev tool, login keychain is often sufficient.
	// However, browsers usually respect the login keychain.

	// Let's try attempting to add to the login keychain first.
	// Note: 'security' command might interactively prompt for password.

	cmd := exec.Command("security", "add-trusted-cert", "-d", "-r", "trustRoot", "-k", os.Getenv("HOME")+"/Library/Keychains/login.keychain-db", certPath)
	// Fallback to older path if -db doesn't exist? modern macOS uses .keychain-db.
	// Actually, just omitting -k might default correctly or prompt.
	// Let's try specific user keychain to be safe.

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install cert: %v, output: %s", err, string(output))
	}

	return nil
}

// CheckTrust checks if the cert is trusted (simplified check)
func (s *Service) CheckTrust() bool {
	// TODO: implement 'security dump-trust-settings' parsing or verify cert
	// For now, return false to force "Install" button availability
	return false
}
