package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/salacoste/siberia/siberia/crypto"
)

type AppConfig struct {
	AppDataDir       string `json:"-"` // Added for internal use
	Theme            string `json:"theme"`
	Language         string `json:"language"`
	AutoRefresh      bool   `json:"auto_refresh"`
	RefreshInterval  int    `json:"refresh_interval"`
	DbSync           bool   `json:"db_sync"`
	ProxyPort        int    `json:"proxy_port"`
	UpstreamProxy    string `json:"upstream_proxy"`
	TargetIDE        string `json:"target_ide"`      // "vscode", "cursor", "windsurf"
	LegacyIDEPath    string `json:"legacy_ide_path"` // Path to state.vscdb
	WindowWidth      int    `json:"window_width"`
	WindowHeight     int    `json:"window_height"`
	MaxLogHistory    int    `json:"max_log_history"` // Default 5000
	AccessLogEnabled bool   `json:"access_log_enabled"`

	// z.ai Provider
	ZaiEnabled bool   `json:"zai_enabled"`
	ZaiBaseURL string `json:"zai_base_url"`
	ZaiApiKey  string `json:"zai_api_key"`
	// Dispatch Mode: "off", "exclusive", "pooled", "fallback"
	ZaiDispatchMode string            `json:"zai_dispatch_mode"`
	ZaiModelMapping map[string]string `json:"zai_model_mapping"`

	// Security
	MitmEnabled bool   `json:"mitm_enabled"`
	AuthEnabled bool   `json:"auth_enabled"`
	AuthToken   string `json:"auth_token"`
	MasterKey   string `json:"master_key"` // 32-byte hex encoded key

	// Cloud Sync
	CloudEnabled  bool   `json:"cloud_enabled"`
	CloudUserID   string `json:"cloud_user_id"`
	CloudEmail    string `json:"cloud_email"`
	CloudLastSync string `json:"cloud_last_sync"` // RFC3339 Timestamp
	CloudSyncKey  string `json:"cloud_sync_key"`  // 32-byte hex key for encrypting cloud blob

	// MCP Server
	ZaiMcpConfig ZaiMcpConfig `json:"zai_mcp_config"`

	// Scheduling
	SchedulingMode string `json:"scheduling_mode"` // "PerformanceFirst" or "CacheFirst"
}

type ZaiMcpConfig struct {
	Enabled bool `json:"enabled"`
	Port    int  `json:"port"` // Optional, if 0 uses embedded or default
}

const (
	SchedulePerformanceFirst = "PerformanceFirst"
	ScheduleCacheFirst       = "CacheFirst"
)

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
			Theme:            "system",
			Language:         "en",
			AutoRefresh:      true,
			RefreshInterval:  15,
			DbSync:           true,
			ProxyPort:        7100,
			UpstreamProxy:    "",
			TargetIDE:        "vscode", // Default
			WindowWidth:      1024,
			WindowHeight:     768,
			MaxLogHistory:    5000,
			AccessLogEnabled: true,
			ZaiEnabled:       false,
			ZaiBaseURL:       "https://api.z.ai/v1",
			ZaiApiKey:        "",
			ZaiDispatchMode:  "off",
			ZaiModelMapping:  map[string]string{},
			AuthEnabled:      false,
			AuthToken:        "",
			MasterKey:        "",
			CloudEnabled:     false,
			CloudUserID:      "",
			CloudEmail:       "",
			CloudLastSync:    "",
			CloudSyncKey:     "",
			AppDataDir:       appDir,
			ZaiMcpConfig: ZaiMcpConfig{
				Enabled: false,
				Port:    6200,
			},
			SchedulingMode: SchedulePerformanceFirst,
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

// NewTestManager creates a manager with a specific directory for testing
func NewTestManager(dir string) (*Manager, error) {
	appDir := filepath.Join(dir, "siberia")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return nil, err
	}
	mgr := &Manager{
		configPath: filepath.Join(appDir, "config.json"),
		Config: AppConfig{
			AppDataDir: appDir,
		},
	}
	// Init default keys
	key, _ := crypto.GenerateKey()
	mgr.Config.MasterKey = key
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
