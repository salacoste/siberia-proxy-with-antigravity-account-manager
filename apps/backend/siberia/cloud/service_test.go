package cloud

import (
	"context"
	"os"
	"testing"

	"github.com/salacoste/siberia/siberia/config"
	"github.com/salacoste/siberia/siberia/logger"
)

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
