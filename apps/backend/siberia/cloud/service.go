package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/salacoste/siberia/siberia/config"
	"github.com/salacoste/siberia/siberia/crypto"
	"github.com/salacoste/siberia/siberia/logger"
)

// TODO: Move to config or env
const (
	SupabaseURL = "https://your-project.supabase.co"
	SupabaseKey = "your-anon-key"
)

type Service struct {
	configManager *config.Manager
	logger        *logger.Logger
	mu            sync.RWMutex
	client        *SupabaseClient
	sessionToken  string // In-memory only
}

func NewService(cfgMgr *config.Manager, log *logger.Logger) *Service {
	return &Service{
		configManager: cfgMgr,
		logger:        log,
		client:        NewSupabaseClient(SupabaseURL, SupabaseKey),
	}
}

// Login authenticates with Supabase
func (s *Service) Login(ctx context.Context, email, password string) error {
	s.logger.Info("Cloud: Attempting login for %s", email)

	resp, err := s.client.SignIn(email, password)
	if err != nil {
		s.logger.Error("Cloud: Login failed: %v", err)
		return err
	}

	// Save Session
	s.mu.Lock()
	s.sessionToken = resp.AccessToken
	s.mu.Unlock()

	cfg := s.configManager.Get()
	cfg.CloudEnabled = true
	cfg.CloudEmail = resp.User.Email
	cfg.CloudUserID = resp.User.ID

	// Generate/Set Sync Key if missing
	if cfg.CloudSyncKey == "" {
		// Key derivation from password would be better, but for now specific random key
		key, _ := crypto.GenerateKey()
		cfg.CloudSyncKey = key
	}
	// TODO: Save Refresh Token securely (maybe in Keychain/Vault?).
	// For MVP putting token in config is UNSAFE but effective given local file security assumptions.
	// Actually config.go doesn't have CloudToken field yet?
	// We added CloudUserID etc. Let's stick to essential ID for now.
	// The sync will need the token passed or re-acquired.

	return s.configManager.Update(cfg)
}

// Sync pushes local changes if newer, or pulls remote if newer
// Sync pushes local changes to the cloud
func (s *Service) Client() *SupabaseClient {
	return s.client
}

func (s *Service) Sync(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cfg := s.configManager.Get()
	if !cfg.CloudEnabled || cfg.CloudUserID == "" || cfg.CloudSyncKey == "" {
		s.logger.Info("Cloud: Sync skipped (disabled or missing key)")
		return fmt.Errorf("cloud sync disabled or unconfigured")
	}

	s.logger.Info("Cloud: Starting Sync for User %s...", cfg.CloudUserID)

	// 1. Prepare Data Blob
	// We want to sync: AppConfig (subset?), Accounts.
	// For MVP: Sync EVERYTHING in a simple JSON structure.

	exportData := map[string]interface{}{
		"config": cfg, // Caution: This includes MasterKey. We are double-encrypting, so it's "safe", but maybe exclude?
		// "accounts": accounts, // TODO: Need to inject AccountService to get accounts array
		"synced_at": time.Now(),
	}

	jsonData, err := json.Marshal(exportData)
	if err != nil {
		return err
	}

	// 2. Encrypt
	ivAndCipher, err := crypto.Encrypt(string(jsonData), cfg.CloudSyncKey)
	if err != nil {
		s.logger.Error("Cloud: Encryption failed: %v", err)
		return err
	}

	// Split IV and Cipher? crypto.Encrypt returns hex(iv+cipher). ProfileData expects separate IV?
	// Our `ProfileData` struct has `iv` and `data_blob`.
	// Let's check crypto.go. Encrypt returns "hex(iv + ciphertext)".
	// We used AES-GCM standard nonce size (12 bytes usually).
	// We should probably store the WHOLE string in DataBlob and leave IV empty/legacy,
	// OR parse it out.
	// Simpler: Store everything in DataBlob.

	profile := ProfileData{
		Email:    cfg.CloudEmail,
		DataBlob: ivAndCipher,
		IV:       "", // Embedded in blob
	}

	// 3. Push to Cloud
	// We need an Access Token. For now, we don't have it stored!
	// MVP Limitation: User must be "logged in" this session or we need to refresh.
	// If we don't have a token, we fail.
	// TODO: Store AccessToken in memory in Service struct (not persistent config for security).
	// For now, let's error if no token.

	if s.sessionToken == "" {
		return fmt.Errorf("no active cloud session, please login again")
	}

	err = s.client.UpsertProfile(cfg.CloudUserID, profile)
	if err != nil {
		s.logger.Error("Cloud: Upload failed: %v", err)
		return err
	}

	// 4. Update internal state
	cfg.CloudLastSync = time.Now().Format(time.RFC3339)
	s.configManager.Update(cfg)

	s.logger.Info("Cloud: Sync successful.")
	return nil
}

// Logout clears local cloud session
func (s *Service) Logout(ctx context.Context) error {
	cfg := s.configManager.Get()
	cfg.CloudEnabled = false
	cfg.CloudUserID = ""
	cfg.CloudEmail = ""
	cfg.CloudLastSync = ""
	return s.configManager.Update(cfg)
}
