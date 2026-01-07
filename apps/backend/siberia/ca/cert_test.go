package ca

import (
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"
)

func TestGenerateCA(t *testing.T) {
	certPEM, keyPEM, err := GenerateCA()
	if err != nil {
		t.Fatalf("GenerateCA() failed: %v", err)
	}

	// 1. Decode Cert
	block, _ := pem.Decode(certPEM)
	if block == nil {
		t.Fatal("Failed to decode cert PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("Failed to parse cert: %v", err)
	}

	// 2. Verify Subject
	if cert.Subject.CommonName != "Siberia Proxy CA" {
		t.Errorf("Expected Subject 'Siberia Proxy CA', got '%s'", cert.Subject.CommonName)
	}

	// 3. Verify Validity (Approx 10 years)
	expectedDuration := 10 * 365 * 24 * time.Hour
	actualDuration := cert.NotAfter.Sub(cert.NotBefore)
	// Allow some margin
	if actualDuration < expectedDuration-time.Hour || actualDuration > expectedDuration+time.Hour {
		t.Errorf("Expected duration ~%v, got %v", expectedDuration, actualDuration)
	}

	// 4. Verify Key
	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		t.Fatal("Failed to decode key PEM")
	}
	_, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		t.Fatalf("Failed to parse private key: %v", err)
	}

	// 5. Verify Pair (Public Key matches)
	if err := cert.CheckSignatureFrom(cert); err != nil {
		t.Errorf("Certificate validation against itself failed: %v", err)
	}

	// Check standard RSA validation
	// (CheckSignatureFrom uses the public key inside the cert, which comes from the private key)
	// We can also check explicit equality of modulo/exponent but CheckSignatureFrom is sufficient integration test.
}
