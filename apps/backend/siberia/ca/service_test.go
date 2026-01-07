package ca

import (
	"os"
	"testing"

	"github.com/salacoste/siberia/siberia/config"
)

func TestService_EnsureCA(t *testing.T) {
	// Setup Temp Dir
	tempDir, err := os.MkdirTemp("", "siberia-ca-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Mock Config
	cfg := &config.Manager{
		Config: config.AppConfig{
			AppDataDir: tempDir, // Use temp dir as AppDataDir
		},
	}

	svc := NewService(cfg)

	// --- 1. First Run: Generation ---
	if err := svc.EnsureCA(); err != nil {
		t.Fatalf("First EnsureCA failed: %v", err)
	}

	certPath, keyPath := svc.GetCAPath()
	if !fileExists(certPath) || !fileExists(keyPath) {
		t.Fatal("Cert or Key file NOT created")
	}

	// Verify Permissions (Unix only - check if we can skip on Windows)
	// For this environment (Mac/Linux), we expect 0600 for key.
	info, err := os.Stat(keyPath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("Key permissions mismatch. Expected 0600, got %v", info.Mode().Perm())
	}

	// Verify content
	firstCertBytes, _ := os.ReadFile(certPath)
	firstKeyBytes, _ := os.ReadFile(keyPath)
	if len(firstCertBytes) == 0 || len(firstKeyBytes) == 0 {
		t.Error("Generated files are empty")
	}

	// --- 2. Second Run: Idempotency ---
	// Mod time check
	statBefore, _ := os.Stat(path(keyPath))
	timeBefore := statBefore.ModTime()

	if err := svc.EnsureCA(); err != nil {
		t.Fatalf("Second EnsureCA failed: %v", err)
	}

	statAfter, _ := os.Stat(path(keyPath))
	if !statAfter.ModTime().Equal(timeBefore) {
		t.Error("EnsureCA overwrote existing valid CA (ModTime changed)")
	}

	// Content check
	secondKeyBytes, _ := os.ReadFile(keyPath)
	if string(firstKeyBytes) != string(secondKeyBytes) {
		t.Error("Key content changed between runs")
	}
}

func path(s string) string { return s }
