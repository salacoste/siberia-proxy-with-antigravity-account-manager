package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/salacoste/siberia/siberia/crypto"
)

type AppConfig struct {
	AppDataDir      string `json:"-"` // Added for internal use
	Theme           string `json:"theme"`
	Language        string `json:"language"`
	AutoRefresh     bool   `json:"auto_refresh"`
	RefreshInterval int    `json:"refresh_interval"`
	DbSync          bool   `json:"db_sync"`
	ProxyPort       int    `json:"proxy_port"`
	UpstreamProxy   string `json:"upstream_proxy"`

	// z.ai Provider
	ZaiEnabled bool   `json:"zai_enabled"`
	ZaiBaseURL string `json:"zai_base_url"`
	ZaiApiKey  string `json:"zai_api_key"`

	// Security
	MitmEnabled bool   `json:"mitm_enabled"`
	AuthEnabled bool   `json:"auth_enabled"`
	AuthToken   string `json:"auth_token"`
	MasterKey   string `json:"master_key"` // 32-byte hex encoded key
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

	mgr := &Manager{
		configPath: filepath.Join(appDir, "config.json"),
		Config: AppConfig{
			Theme:           "system",
			Language:        "en",
			AutoRefresh:     true,
			RefreshInterval: 15,
			DbSync:          true,
			ProxyPort:       3000,
			UpstreamProxy:   "",
			ZaiEnabled:      false,
			ZaiBaseURL:      "https://api.z.ai/v1",
			ZaiApiKey:       "",
			AuthEnabled:     false,
			AuthToken:       "",
			MasterKey:       "",
			AppDataDir:      appDir,
		},
	}

	// Ensure MasterKey exists immediately
	if mgr.Config.MasterKey == "" {
		key, err := crypto.GenerateKey()
		if err != nil {
			return nil, err
		}
		mgr.Config.MasterKey = key
	}

	return mgr, nil
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

	if err := json.Unmarshal(data, &m.Config); err != nil {
		return err
	}

	// Ensure MasterKey exists if loaded config didn't have it
	if m.Config.MasterKey == "" {
		key, err := crypto.GenerateKey()
		if err != nil {
			return err
		}
		m.Config.MasterKey = key
		// Save immediately to persist the new key
		return m.saveInternal()
	}

	return nil
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

func (m *Manager) ConfigDir() string {
	return filepath.Dir(m.configPath)
}
