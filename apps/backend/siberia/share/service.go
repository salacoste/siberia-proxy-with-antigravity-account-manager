package share

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/salacoste/siberia/siberia/export"
	"github.com/salacoste/siberia/siberia/proxy"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

// UploadSession takes the raw event from frontend, converts to HAR, and "uploads" it.
func (s *Service) UploadSession(event proxy.ProxyRequestEvent) (string, error) {
	// 1. Convert to HAR
	harJSON, err := export.ToHAR(event)
	if err != nil {
		return "", fmt.Errorf("failed to generate HAR: %w", err)
	}

	// 2. Mock Upload (In real implementation, PUT to S3/R2)
	// For MVP: We assume it succeeded.
	// We generate a deterministic ID based on content to simulate a unique link.
	hash := md5.Sum([]byte(harJSON + time.Now().String()))
	id := hex.EncodeToString(hash[:])[:8]

	// 3. Return Mock Link
	mockLink := fmt.Sprintf("https://share.siberia.dev/log/%s", id)
	return mockLink, nil
}
