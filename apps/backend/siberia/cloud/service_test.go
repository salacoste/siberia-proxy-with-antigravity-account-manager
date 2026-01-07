package cloud

import (
	"context"
	"os"
	"testing"

	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/salacoste/siberia/siberia/config"
	"github.com/salacoste/siberia/siberia/crypto"
	"github.com/salacoste/siberia/siberia/logger"
)

// MockTransport allows mocking HTTP responses
type MockTransport struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

func TestService_LoginLogout(t *testing.T) {
	// Setup Temp Config
	tmpDir, err := os.MkdirTemp("", "siberia_cloud_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Use NewTestManager to ensure valid config path
	configMgr, err := config.NewTestManager(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Create dummy Logger
	log := logger.New("TEST")

	svc := NewService(configMgr, log)

	// Test Logout (Should clear fields)
	cfg := configMgr.Get()
	cfg.CloudEnabled = true
	cfg.CloudUserID = "user123"
	configMgr.Update(cfg)

	if err := svc.Logout(context.Background()); err != nil {
		t.Fatalf("Logout failed: %v", err)
	}

	cfg = configMgr.Get()
	if cfg.CloudEnabled {
		t.Error("CloudEnabled should be false after Logout")
	}
	if cfg.CloudUserID != "" {
		t.Error("CloudUserID should be empty after Logout")
	}

	// Login requires Mock Client or Real network.
	// We skip Login network test here, relying on manual verify or mock if implemented.
}

func TestService_Sync_Pull(t *testing.T) {
	// Setup
	tmpDir, _ := os.MkdirTemp("", "siberia_sync_test")
	defer os.RemoveAll(tmpDir)

	configMgr, _ := config.NewTestManager(tmpDir)
	log := logger.New("TEST")
	svc := NewService(configMgr, log)

	// Inject Session & Config
	svc.sessionToken = "fake-token"
	cfg := configMgr.Get()
	cfg.CloudEnabled = true
	cfg.CloudUserID = "user_test"
	cfg.CloudSyncKey = "12345678901234567890123456789012"                   // 32 bytes
	cfg.CloudLastSync = time.Now().Add(-2 * time.Hour).Format(time.RFC3339) // Old time
	configMgr.Update(cfg)

	// Prepare Mock Response
	// We want Remote to be NEWER (-1 hour)
	newerTime := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
	exportData := map[string]interface{}{
		"config":    map[string]interface{}{"CloudUserID": "user_test"},
		"synced_at": newerTime,
	}
	jsonBytes, _ := json.Marshal(exportData)
	enc, _ := crypto.Encrypt(string(jsonBytes), cfg.CloudSyncKey)

	mockResponse := []ProfileData{
		{
			UserID:    "user_test",
			Email:     "test@example.com",
			DataBlob:  enc,
			UpdatedAt: newerTime,
		},
	}
	respBody, _ := json.Marshal(mockResponse)

	// Mock Transport
	svc.client.Client.Transport = &MockTransport{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			// Verify Headers
			if req.Header.Get("Authorization") != "Bearer fake-token" {
				t.Error("Missing Auth Header")
			}

			// Expect GET profile (Pull)
			if req.Method == "GET" {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBuffer(respBody)),
					Header:     make(http.Header),
				}, nil
			}
			return &http.Response{StatusCode: 404, Body: io.NopCloser(bytes.NewBufferString(""))}, nil
		},
	}

	// EXECUTE
	err := svc.Sync(context.Background())
	if err != nil {
		t.Fatalf("Sync Pull failed: %v", err)
	}

	// VERIFY
	// Local LastSync should be updated to matches Remote (newerTime)
	newCfg := configMgr.Get()
	if newCfg.CloudLastSync != newerTime {
		t.Errorf("Expected LastSync %s, got %s", newerTime, newCfg.CloudLastSync)
	}
}

func TestService_Sync_Push(t *testing.T) {
	// Setup
	tmpDir, _ := os.MkdirTemp("", "siberia_sync_push")
	defer os.RemoveAll(tmpDir)

	configMgr, _ := config.NewTestManager(tmpDir)
	log := logger.New("TEST")
	svc := NewService(configMgr, log)

	svc.sessionToken = "fake-token"
	cfg := configMgr.Get()
	cfg.CloudEnabled = true
	cfg.CloudUserID = "user_test"
	cfg.CloudSyncKey = "12345678901234567890123456789012"
	cfg.CloudLastSync = time.Now().Add(-1 * time.Minute).Format(time.RFC3339) // Past time
	configMgr.Update(cfg)

	// Mock Response: Remote is Older or Empty
	svc.client.Client.Transport = &MockTransport{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method == "GET" {
				// Return Empty list (Empty remote)
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString("[]")),
					Header:     make(http.Header),
				}, nil
			}
			if req.Method == "POST" {
				// PUSH
				return &http.Response{
					StatusCode: 201,
					Body:       io.NopCloser(bytes.NewBufferString("")),
					Header:     make(http.Header),
				}, nil
			}
			return &http.Response{StatusCode: 400}, nil
		},
	}

	err := svc.Sync(context.Background())
	if err != nil {
		t.Fatalf("Sync Push failed: %v", err)
	}

	// Verify Local LastSync Updated (to Now, roughly)
	newCfg := configMgr.Get()
	// Should be strictly greater than start time, but hard to test exact millisecond.
	// We trust it updated.
	if newCfg.CloudLastSync == cfg.CloudLastSync {
		// It updates to Now() after push.
		t.Error("LastSync should have updated after Push")
	}
}
