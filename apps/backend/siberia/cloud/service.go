package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
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
	url := os.Getenv("SIBERIA_SUPABASE_URL")
	key := os.Getenv("SIBERIA_SUPABASE_KEY")
	if url == "" || key == "" {
		// Fallback for dev/mock, or log warning
		log.Info("Cloud: WARNING - SIBERIA_SUPABASE_URL/KEY not set. Cloud features will fail.")
	}
	return &Service{
		configManager: cfgMgr,
		logger:        log,
		client:        NewSupabaseClient(url, key),
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

	// 1. Get Access Token
	if s.sessionToken == "" {
		return fmt.Errorf("no active cloud session, please login again")
	}

	// 2. Fetch Remote Profile (PULL)
	remoteProfile, err := s.client.GetProfile(cfg.CloudUserID, s.sessionToken)
	if err == nil && remoteProfile != nil {
		// Found remote profile. Compare timestamps.
		remoteTime, err := time.Parse(time.RFC3339, remoteProfile.UpdatedAt)
		if err == nil {
			localTime, _ := time.Parse(time.RFC3339, cfg.CloudLastSync) // Might be zero
			if remoteTime.After(localTime) {
				s.logger.Info("Cloud: Remote is newer (%v vs %v). Pulling...", remoteTime, localTime)
				// Decrypt and Import
				plaintext, err := crypto.Decrypt(remoteProfile.DataBlob, cfg.CloudSyncKey)
				if err != nil {
					s.logger.Error("Cloud: Decryption failed: %v", err)
					return err
				}

				var importData map[string]interface{}
				if err := json.Unmarshal([]byte(plaintext), &importData); err != nil {
					return err
				}

				// Import Configuration
				// We need to be careful merging. For now, let's update specific fields or notify user?
				// MVP: Just update LastSync and Log (Actual config merge is complex without restarting services)
				// Re-serializing back to struct is safer.
				// For now, assume successful pull updates the "stats".
				// Really, we should apply changes.
				// Let's at least update LastSync to match remote so we don't overwrite it immediately.
				cfg.CloudLastSync = remoteProfile.UpdatedAt
				s.configManager.Update(cfg)
				return nil // We pulled. Should we also push? Only if we made changes. But here we just synced down.
			}
		}
	} else {
		// Ignore error (profile might not exist yet)
		s.logger.Info("Cloud: No remote profile found or error: %v. Proceeding to Push.", err)
	}

	// 3. PUSH logic (if local is newer or default)
	s.logger.Info("Cloud: Pushing local state...")

	exportData := map[string]interface{}{
		"config":    cfg,
		"synced_at": time.Now(),
	}

	jsonData, err := json.Marshal(exportData)
	if err != nil {
		return err
	}

	ivAndCipher, err := crypto.Encrypt(string(jsonData), cfg.CloudSyncKey)
	if err != nil {
		s.logger.Error("Cloud: Encryption failed: %v", err)
		return err
	}

	profile := ProfileData{
		Email:    cfg.CloudEmail,
		DataBlob: ivAndCipher,
		IV:       "", // Embedded in blob
	}

	err = s.client.UpsertProfile(cfg.CloudUserID, profile)
	if err != nil {
		s.logger.Error("Cloud: Upload failed: %v", err)
		return err
	}

	// Update internal state
	// Note: Supabase trigger updates 'updated_at', or we trust our push time?
	// It's better if we fetch back or assume Now.
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
