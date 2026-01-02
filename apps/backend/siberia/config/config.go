package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type AppConfig struct {
	Theme           string `json:"theme"`
	Language        string `json:"language"`
	AutoRefresh     bool   `json:"auto_refresh"`
	RefreshInterval int    `json:"refresh_interval"`
	DbSync          bool   `json:"db_sync"`
}

type Manager struct {
	configPath string
	mu         sync.RWMutex
	Config     AppConfig
}

func NewManager() (*Manager, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	appDir := filepath.Join(configDir, "siberia")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return nil, err
	}

	return &Manager{
		configPath: filepath.Join(appDir, "config.json"),
		Config: AppConfig{
			Theme:           "system",
			Language:        "en",
			AutoRefresh:     true,
			RefreshInterval: 15,
			DbSync:          true,
		},
	}, nil
}

func (m *Manager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(m.configPath)
	if os.IsNotExist(err) {
		return m.saveInternal()
	}
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &m.Config)
}

func (m *Manager) Save() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.saveInternal()
}

// Public thread-safe Save
func (m *Manager) SaveConfig() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.saveInternal()
}

func (m *Manager) saveInternal() error {
	data, err := json.MarshalIndent(m.Config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.configPath, data, 0644)
}

func (m *Manager) Update(newConfig AppConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Config = newConfig
	return m.saveInternal()
}

func (m *Manager) Get() AppConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Config
}
