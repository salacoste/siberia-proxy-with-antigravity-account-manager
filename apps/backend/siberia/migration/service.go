package migration

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/salacoste/siberia/siberia/accounts"
	"github.com/salacoste/siberia/siberia/logger"
)

type Service struct {
	accounts *accounts.Service
}

func NewService(accSvc *accounts.Service) *Service {
	return &Service{
		accounts: accSvc,
	}
}

// CheckLegacyData scans for known legacy config paths
func (s *Service) CheckLegacyData() (LegacyStatus, error) {
	path, err := s.getLegacyPath()
	if err != nil {
		return LegacyStatus{Found: false, Count: 0}, nil
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return LegacyStatus{Found: false, Count: 0}, nil
	}

	// Try to parse to get count
	config, err := s.parseFile(path)
	if err != nil {
		logger.New("Migration").Error(fmt.Sprintf("Failed to parse legacy config at %s: %v", path, err))
		return LegacyStatus{Found: true, Count: 0}, nil // Found file but failed to parse
	}

	return LegacyStatus{Found: true, Count: len(config.Accounts)}, nil
}

// PerformImport reads the legacy file and imports accounts into Siberia
func (s *Service) PerformImport() (int, error) {
	path, err := s.getLegacyPath()
	if err != nil {
		logger.New("Migration").Error(fmt.Sprintf("Path resolution failed: %v", err))
		return 0, fmt.Errorf("could not resolve legacy path: %v", err)
	}
	logger.New("Migration").Info(fmt.Sprintf("Reading legacy file from: %s", path))

	config, err := s.parseFile(path)
	if err != nil {
		logger.New("Migration").Error(fmt.Sprintf("Parse failed: %v", err))
		return 0, fmt.Errorf("failed to parse legacy config: %v", err)
	}
	logger.New("Migration").Info(fmt.Sprintf("Found %d accounts in legacy config", len(config.Accounts)))

	importedCount := 0
	for _, acc := range config.Accounts {
		logger.New("Migration").Info(fmt.Sprintf("Processing account: %s", acc.Email))
		if acc.Email == "" || acc.RefreshToken == "" {
			logger.New("Migration").Info("Skipping invalid entry (missing email or token)")
			continue // Skip invalid entries
		}

		// We use RefreshToken as the "password" since the old agent structure often swapped them,
		// or we treat it as an actionable token.
		// In Siberia's CreateAccount, password is encrypted.
		// NOTE: Detailed logic might depend on how CreateAccount handles existence checks.
		// CreateAccount currently returns error if it exists? Or just overwrites?
		// Based on previous analysis, we should try-catch.

		// Determine a ProxyGroup - legacy didn't have this, default to "legacy-import" or "default"
		err := s.accounts.CreateAccount(acc.Email, acc.RefreshToken, "", "legacy-import")
		if err == nil {
			importedCount++
			logger.New("Migration").Info(fmt.Sprintf("Imported legacy account: %s", acc.Email))
		} else {
			logger.New("Migration").Error(fmt.Sprintf("Failed to import %s: %v", acc.Email, err))
		}
	}

	return importedCount, nil
}

func (s *Service) getLegacyPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Priority 1: ~/.antigravity-agent/antigravity_accounts.json
	p1 := filepath.Join(home, ".antigravity-agent", "antigravity_accounts.json")
	if _, err := os.Stat(p1); err == nil {
		return p1, nil
	}

	// Priority 2: ~/antigravity_accounts.json
	p2 := filepath.Join(home, "antigravity_accounts.json")
	if _, err := os.Stat(p2); err == nil {
		return p2, nil
	}

	// Not found, return default p1 for CheckLegacyData to fail gracefully on
	return p1, fmt.Errorf("legacy config not found")
}

func (s *Service) parseFile(path string) (*LegacyConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config LegacyConfig
	// Attempt root object parsing first
	if err := json.Unmarshal(data, &config); err != nil {
		// Fallback: Maybe it's just a raw list of accounts?
		// Or maybe the keys are different.
		// For MVP, we assume the specific struct defined.
		return nil, err
	}
	return &config, nil
}

// === Story-62: IDE Scanning Logic ===

// ScanForIDEAccounts scans the system for supported IDEs and tries to extract tokens
func (s *Service) ScanForIDEAccounts() ([]DiscoveredAccount, error) {
	profiles, err := ScanIDEProfiles()
	if err != nil {
		return nil, fmt.Errorf("failed to scan profiles: %w", err)
	}

	var results []DiscoveredAccount

	for _, p := range profiles {
		// Try to read token
		blob, err := ReadStateToken(p.DBPath)
		if err != nil {
			// Just log and skip, maybe they haven't logged in
			logger.New("Migration").Info(fmt.Sprintf("WARN: Could not read DB for %s: %v", p.Name, err))
			continue
		}

		token, err := ExtractRefreshToken(blob)
		if err != nil {
			logger.New("Migration").Info(fmt.Sprintf("WARN: Could not extract token for %s: %v", p.Name, err))
			continue
		}

		if token == "" {
			continue
		}

		// Mask Token
		masked := ""
		if len(token) > 8 {
			masked = token[:4] + "..." + token[len(token)-4:]
		} else {
			masked = "***"
		}

		results = append(results, DiscoveredAccount{
			Name:        p.Name,
			Path:        p.ConfigPath,
			RawToken:    token,
			MaskedToken: masked,
			Source:      "IDE", // could refine based on p.Name
		})
	}

	return results, nil
}

// ImportAccount imports a single account by raw token
func (s *Service) ImportAccount(rawToken string) error {
	// Generate a name
	// In functional requirements, we don't have email.
	// We'll use a placeholder. user can rename later.
	// Or maybe we can fetch it? for now, minimal.
	// Use current time to make unique if needed, or just "Imported Account"
	// CreateAccount does NOT enforce unique email on DB level? standard gorm struct usually does if tagged.
	// Check account entity...
	// But assuming we can handle duplicates or let's generate a unique one.

	// Better: "Imported Account"
	// Actually, `CreateAccount` takes email.
	// If we use same email "Imported Account", duplicate key error?
	// Let's append random suffix or check if we can get email.
	// Since we can't reliably get email without querying an API (and we might be offline or blocked),
	// we use a placeholder. TODO: User can rename.

	importName := fmt.Sprintf("Imported Account %s", GenerateRandomSuffix(4))

	// Use token as password (since we use password field for tokens in this MVP/proxy logic)
	// CreateAccount(email, password, recovery, proxyGroup)
	return s.accounts.CreateAccount(importName, rawToken, "", "imported")
}

func GenerateRandomSuffix(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "X"
	}
	return hex.EncodeToString(bytes)
}
