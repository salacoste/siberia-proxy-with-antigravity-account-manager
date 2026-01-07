package ca

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

// GenerateCA creates a new self-signed Root CA certificate and private key.
// It uses RSA 2048 for compatibility and security.
func GenerateCA() (certPEM []byte, keyPEM []byte, err error) {
	// 1. Generate Key (RSA 2048)
	// ECDSA P256 is faster, but RSA is widely supported for root stores.
	// Story requirements mention both, choosing RSA for broad compatibility by default unless performance is critical.
	// Given this is run once every 10 years, generation speed is irrelevant.
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	// 2. Create Certificate Template
	validity := 10 * 365 * 24 * time.Hour // 10 Years
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   "Siberia Proxy CA",
			Organization: []string{"Siberia Proxy"},
			Country:      []string{"US"},
		},
		NotBefore: time.Now().Add(-10 * time.Minute), // Buffer for clock skew
		NotAfter:  time.Now().Add(validity),

		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// 3. Self-Sign
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return nil, nil, err
	}

	// 4. Encode to PEM
	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes}) //nolint:staticcheck // Why is staticcheck complaining about Type? Use standard string.
	keyBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}
	keyPEM = pem.EncodeToMemory(keyBlock)

	return certPEM, keyPEM, nil
}
