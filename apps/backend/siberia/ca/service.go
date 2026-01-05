package ca

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/salacoste/siberia/siberia/config"
)

type Service struct {
	config *config.AppConfig
	caCert *tls.Certificate
}

func NewService(cfg *config.AppConfig) *Service {
	return &Service{
		config: cfg,
	}
}

// GetCAPath returns the absolute path to the CA certificate file
func (s *Service) GetCAPath() string {
	return filepath.Join(s.config.AppDataDir, "ca.pem")
}

// GetKeyPath returns the absolute path to the CA private key file
func (s *Service) GetKeyPath() string {
	return filepath.Join(s.config.AppDataDir, "ca.key")
}

// GetCAPair returns the loaded TLS certificate pair for the CA
func (s *Service) GetCAPair() (*tls.Certificate, error) {
	if s.caCert != nil {
		return s.caCert, nil
	}

	certPath := s.GetCAPath()
	keyPath := s.GetKeyPath()

	// Check if files exist
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("CA certificate not found at %s", certPath)
	}

	pair, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load CA key pair: %w", err)
	}

	s.caCert = &pair
	return s.caCert, nil
}

// EnsureCA checks if CA exists, otherwise generates a new one
func (s *Service) EnsureCA() error {
	certPath := s.GetCAPath()
	keyPath := s.GetKeyPath()

	_, certErr := os.Stat(certPath)
	_, keyErr := os.Stat(keyPath)

	if certErr == nil && keyErr == nil {
		// Verify we can load it
		_, err := s.GetCAPair()
		if err == nil {
			return nil // Already exists and valid
		}
		// If load failed, regenerate
		s.caCert = nil // Clear cache
		fmt.Printf("CA load failed (%v), regenerating...\n", err)
	}

	return s.generateCA(certPath, keyPath)
}

func (s *Service) generateCA(certPath, keyPath string) error {
	// 1. Generate Key (RSA 2048)
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	// 2. Create Template
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Siberia Proxy"},
			CommonName:   "Siberia Proxy CA",
		},
		NotBefore: time.Now().Add(-1 * time.Minute),
		NotAfter:  time.Now().Add(10 * 365 * 24 * time.Hour), // 10 years

		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// 3. Sign Cert
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	// 4. Save Cert
	certOut, err := os.Create(certPath)
	if err != nil {
		return fmt.Errorf("failed to open cert.pem for writing: %w", err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		certOut.Close()
		return fmt.Errorf("failed to encode cert: %w", err)
	}
	if err := certOut.Close(); err != nil {
		return fmt.Errorf("failed to close cert file: %w", err)
	}

	// 5. Save Key (Securely 0600)
	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open key.pem for writing: %w", err)
	}
	privBytes := x509.MarshalPKCS1PrivateKey(priv)
	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes}); err != nil {
		keyOut.Close()
		return fmt.Errorf("failed to encode key: %w", err)
	}
	if err := keyOut.Close(); err != nil {
		return fmt.Errorf("failed to close key file: %w", err)
	}

	fmt.Printf("Generated new Root CA at: %s\n", certPath)
	return nil
}
