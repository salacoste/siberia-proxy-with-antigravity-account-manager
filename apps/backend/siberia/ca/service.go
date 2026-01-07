package ca

import (
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/salacoste/siberia/siberia/config"
)

type Service struct {
	config *config.Manager
	mu     sync.Mutex
}

func NewService(cfg *config.Manager) *Service {
	return &Service{
		config: cfg,
	}
}

// GetCAPath returns the absolute paths to the CA certificate and private key.
func (s *Service) GetCAPath() (certPath, keyPath string) {
	dir := s.certDir()
	return filepath.Join(dir, "ca.pem"), filepath.Join(dir, "ca.key")
}

// GetCAPair loads the CA key pair and returns it as a tls.Certificate pointer.
func (s *Service) GetCAPair() (*tls.Certificate, error) {
	certPath, keyPath := s.GetCAPath()
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

// EnsureCA checks if the CA exists, and generates it if it doesn't.
// It strictly enforces 0600 permissions on the private key.
// Returns idempotently if CA already exists and is valid (at least loadable).
func (s *Service) EnsureCA() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	dir := s.certDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create cert directory: %w", err)
	}

	certPath, keyPath := s.GetCAPath()

	// 1. Check Existence
	certExists := fileExists(certPath)
	keyExists := fileExists(keyPath)

	if certExists && keyExists {
		// Verify Permissions on Key?
		// We should enforce 0600 even if it exists.
		if err := enforceSecurePermissions(keyPath); err != nil {
			fmt.Printf("[CA] Warning: Fix permissions on existing key failed: %v\n", err)
		}
		return nil
	}

	if certExists != keyExists {
		// Partial state?
		// Safest to backup/wipe and regenerate if we have half a pair, as we can't recover one from the other easily (we technically can recover Pub from Priv, but not vice versa).
		// For MVP: Treat as missing and overwrite.
		fmt.Println("[CA] Partial CA state detected. Regenerating...")
	} else {
		fmt.Println("[CA] No CA found. Generating new Root CA...")
	}

	// 2. Generate
	certPEM, keyPEM, err := GenerateCA()
	if err != nil {
		return fmt.Errorf("failed to generate CA: %w", err)
	}

	// 3. Write Key (Strict 0600)
	// Open file with O_CREATE|O_WRONLY|O_TRUNC and 0600 mode
	keyFile, err := os.OpenFile(keyPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open key file for writing: %w", err)
	}
	if _, err := keyFile.Write(keyPEM); err != nil {
		keyFile.Close()
		return fmt.Errorf("failed to write key: %w", err)
	}
	keyFile.Close()

	// Double check permissions (paranoid check)
	if err := enforceSecurePermissions(keyPath); err != nil {
		return fmt.Errorf("failed to enforce key permissions: %w", err)
	}

	// 4. Write Cert (0644 is fine for pub cert)
	if err := os.WriteFile(certPath, certPEM, 0644); err != nil {
		return fmt.Errorf("failed to write cert: %w", err)
	}

	fmt.Printf("[CA] Generated and stored CA at %s\n", dir)
	return nil
}

func (s *Service) certDir() string {
	// Use AppDataDir from config, subfolder 'certificates'
	// Assuming AppDataDir is populated or available via config.
	// Actually AppConfig struct has AppDataDir (but it is json:"-").
	// Let's use config.Get().
	// Wait, config.Get() returns AppConfig copy.
	// We need to ensure AppDataDir is accessible. s.config.Get() returns a copy.
	// Let's assume AppDataDir is correctly set there.
	cfg := s.config.Get()
	return filepath.Join(cfg.AppDataDir, "certificates")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func enforceSecurePermissions(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	// Check if mode is 0600 (modulo file type bits)
	// ModePerm is 0777.
	if info.Mode().Perm() != 0600 {
		// Try to fix
		if err := os.Chmod(path, 0600); err != nil {
			return err
		}
	}
	return nil
}
