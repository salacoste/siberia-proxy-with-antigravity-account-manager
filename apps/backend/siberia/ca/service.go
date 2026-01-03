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

	if _, err := os.Stat(certPath); err == nil {
		// Verify we can load it
		_, err := s.GetCAPair()
		if err == nil {
			return nil // Already exists and valid
		}
		// If load failed, regenerate
		fmt.Println("CA load failed, regenerating...")
	}

	return s.generateCA(certPath, keyPath)
}

func (s *Service) generateCA(certPath, keyPath string) error {
	// Generate Key
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create Template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1), // In production, use random serial
		Subject: pkix.Name{
			Organization: []string{"Siberia"},
			CommonName:   "Siberia Proxy CA",
		},
		NotBefore: time.Now().Add(-1 * time.Minute),
		NotAfter:  time.Now().Add(10 * 365 * 24 * time.Hour), // 10 years

		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// Sign Cert
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	// Save Cert
	certOut, err := os.Create(certPath)
	if err != nil {
		return fmt.Errorf("failed to open cert.pem for writing: %w", err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	// Save Key
	keyOut, err := os.Create(keyPath)
	if err != nil {
		return fmt.Errorf("failed to open key.pem for writing: %w", err)
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return fmt.Errorf("unable to marshal private key: %w", err)
	}
	pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	keyOut.Close()

	fmt.Printf("Generated new Root CA at: %s\n", certPath)
	return nil
}
