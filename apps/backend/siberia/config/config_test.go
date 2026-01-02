package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigPersistence(t *testing.T) {
	// Setup temporary directory
	tempDir, err := os.MkdirTemp("", "siberia_config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Override standard config path logic for testing by manually constructing manager
	// Since NewManager uses os.UserConfigDir, we can't easily mock it without refactoring or env vars.
	// However, we can test the Save/Load logic if we allow injecting path or just test NewManager if we assume environment is writable.
	// Better: Refactor NewManager to take an optional path, or just test logic on a Manager struct created manually.

	configFile := filepath.Join(tempDir, "config.json")
	manager := &Manager{
		configPath: configFile,
		Config: AppConfig{
			Theme: "system",
		},
	}

	// Test Save
	if err := manager.SaveConfig(); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Fatalf("Config file was not created at %s", configFile)
	}

	// Test Load
	newManager := &Manager{
		configPath: configFile,
	}
	if err := newManager.Load(); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if newManager.Config.Theme != "system" {
		t.Errorf("Expected Theme 'system', got '%s'", newManager.Config.Theme)
	}

	// Test Update
	newConfig := newManager.Config
	newConfig.Theme = "dark"
	if err := newManager.Update(newConfig); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify persistence logic again
	finalManager := &Manager{configPath: configFile}
	finalManager.Load()
	if finalManager.Config.Theme != "dark" {
		t.Errorf("Expected Theme 'dark' after update, got '%s'", finalManager.Config.Theme)
	}
}
