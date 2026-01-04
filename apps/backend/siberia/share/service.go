package share

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/salacoste/siberia/siberia/export"
	"github.com/salacoste/siberia/siberia/types"
)

type Service struct {
	provider StorageProvider
}

// NewService creates a new Share Service with a given provider
func NewService(provider StorageProvider) *Service {
	return &Service{
		provider: provider,
	}
}

// UploadSession takes the raw event from frontend, converts to HAR, and "uploads" it.
func (s *Service) UploadSession(event types.ProxyRequestEvent) (string, error) {
	// 1. Convert to HAR
	harJSON, err := export.ToHAR(event)
	if err != nil {
		return "", fmt.Errorf("failed to generate HAR: %w", err)
	}

	// 2. Generate Key
	hash := md5.Sum([]byte(harJSON + time.Now().String()))
	id := hex.EncodeToString(hash[:])[:8]
	key := fmt.Sprintf("%s.har", id)

	// 3. Upload via Provider
	return s.provider.Upload(key, []byte(harJSON), "application/json")
}

// MockProvider for fallback when no S3/MinIO is configured
type MockProvider struct{}

func (m *MockProvider) Upload(key string, data []byte, contentType string) (string, error) {
	// Mock Link
	return fmt.Sprintf("https://share.siberia.dev/log/MOCK-%s", key), nil
}
