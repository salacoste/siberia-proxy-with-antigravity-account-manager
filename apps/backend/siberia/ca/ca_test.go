package ca

import (
	"os"
	"testing"

	"github.com/salacoste/siberia/siberia/config"
)

func TestEnsureCA(t *testing.T) {
	// Setup Temp Dir
	tmpDir, err := os.MkdirTemp("", "siberia_ca_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &config.AppConfig{
		AppDataDir: tmpDir,
	}

	svc := NewService(cfg)

	// 1. EnsureCA should create files
	err = svc.EnsureCA()
	if err != nil {
		t.Fatalf("First EnsureCA failed: %v", err)
	}

	certPath := svc.GetCAPath()
	keyPath := svc.GetKeyPath()

	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		t.Error("Cert file was not created")
	}
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		t.Error("Key file was not created")
	}

	// 2. EnsureCA should be idempotent (and load existing)
	err = svc.EnsureCA()
	if err != nil {
		t.Fatalf("Second EnsureCA failed: %v", err)
	}

	// 3. GetCAPair should work
	pair, err := svc.GetCAPair()
	if err != nil {
		t.Fatalf("GetCAPair failed: %v", err)
	}
	if len(pair.Certificate) == 0 {
		t.Error("Certificate chain is empty")
	}
}
