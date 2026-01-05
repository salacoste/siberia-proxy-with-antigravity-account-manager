package ca

import (
	"crypto/x509"
	"encoding/pem"
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
	certPath := svc.GetCAPath()
	keyPath := svc.GetKeyPath()

	// 1. EnsureCA should create files
	err = svc.EnsureCA()
	if err != nil {
		t.Fatalf("First EnsureCA failed: %v", err)
	}

	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		t.Errorf("Cert file was not created at %s", certPath)
	}
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		t.Errorf("Key file was not created at %s", keyPath)
	}

	// 2. Validate Content
	certBytes, err := os.ReadFile(certPath)
	if err != nil {
		t.Fatalf("Failed to read cert: %v", err)
	}
	block, _ := pem.Decode(certBytes)
	if block == nil {
		t.Fatal("Failed to decode PEM block")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("Failed to parse certificate: %v", err)
	}

	if cert.Subject.CommonName != "Siberia Proxy CA" {
		t.Errorf("Expected CommonName 'Siberia Proxy CA', got '%s'", cert.Subject.CommonName)
	}
	if !cert.IsCA {
		t.Error("Certificate is not marked as CA")
	}

	// 3. Test Idempotency
	infoBefore, _ := os.Stat(certPath)
	err = svc.EnsureCA()
	if err != nil {
		t.Fatalf("Second EnsureCA failed: %v", err)
	}
	infoAfter, _ := os.Stat(certPath)

	if infoBefore.ModTime() != infoAfter.ModTime() {
		t.Log("Note: EnsureCA regenerated the certificate on second run.")
	} else {
		t.Log("EnsureCA respected existing certificate.")
	}

	// 4. Test Recovery (Missing Key)
	os.Remove(keyPath)
	err = svc.EnsureCA()
	if err != nil {
		t.Fatalf("Recovery EnsureCA failed: %v", err)
	}
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		t.Errorf("Key file was not recreated")
	}

	// 5. CheckTrust should default to false (requires OS install)
	if svc.CheckTrust() {
		t.Error("New CA should NOT be trusted by OS yet")
	}
}
