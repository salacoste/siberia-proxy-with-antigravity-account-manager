//go:build integration

package share_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/salacoste/siberia/siberia/share"
)

func TestMinIOIntegration(t *testing.T) {
	// 1. Setup - Params match docker-compose
	endpoint := "localhost:9000"
	accessKey := "minioadmin"
	secretKey := "minioadmin"
	bucket := "siberia-shares"

	provider, err := share.NewS3Provider(endpoint, accessKey, secretKey, bucket, false)
	if err != nil {
		t.Fatalf("Failed to create S3Provider: %v", err)
	}

	// 2. Upload
	testContent := "This is a test content for MinIO Integration " + time.Now().String()
	key := fmt.Sprintf("test-%d.txt", time.Now().Unix())

	link, err := provider.Upload(key, []byte(testContent), "text/plain")
	if err != nil {
		t.Fatalf("Failed to upload: %v", err)
	}

	t.Logf("Upload successful. Link: %s", link)

	// 3. Verify Download (Public Link)
	// MinIO takes a moment to propagate policy sometimes, but usually instant
	resp, err := http.Get(link)
	if err != nil {
		t.Fatalf("Failed to download from link: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}

	if string(body) != testContent {
		t.Errorf("Content mismatch. Got %s, want %s", string(body), testContent)
	}
}
