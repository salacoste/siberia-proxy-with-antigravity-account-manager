package accounts

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/salacoste/siberia/siberia/config"
	"github.com/salacoste/siberia/siberia/db"
	"github.com/salacoste/siberia/siberia/ide"
	"github.com/salacoste/siberia/siberia/modules/injection"
	"github.com/salacoste/siberia/siberia/modules/process"
	"github.com/salacoste/siberia/siberia/quota"
	"gorm.io/gorm"
)

type Service struct {
	db        *gorm.DB
	process   *process.Service
	injector  *injection.Service
	quota     *quota.Service
	config    *config.Manager
	configDir string
}

func NewService(database *db.Database, cfg *config.Manager) *Service {
	return &Service{
		db:        database.Conn,
		process:   process.NewService(),
		injector:  injection.NewService(),
		quota:     quota.NewService(),
		config:    cfg,
		configDir: cfg.ConfigDir(),
	}
}

func (s *Service) GetQuotaService() *quota.Service {
	return s.quota
}

// ... existing code ...

func (s *Service) ActivateAccount(id uint) error {
	var acc db.Account
	if err := s.db.First(&acc, id).Error; err != nil {
		return err
	}

	return s.performActivation(&acc)
}

func (s *Service) performActivation(acc *db.Account) error {
	// 1. Resolve IDE Profile
	cfg := s.config.Get()
	targetIDE := cfg.TargetIDE
	if targetIDE == "" {
		targetIDE = "vscode" // Default fallback
	}

	profile, err := ide.GetProfile(targetIDE)
	if err != nil {
		return fmt.Errorf("failed to get IDE profile for '%s': %w", targetIDE, err)
	}

	fmt.Printf("[Account] Activating for IDE: %s (Process: %s)\n", profile.Name, profile.GetProcessName())

	// 2. Kill Target Process
	if err := s.process.Kill(profile.GetProcessName()); err != nil {
		return fmt.Errorf("failed to kill process %s: %w", profile.GetProcessName(), err)
	}

	// 3. Resolve DB Path
	dbPath, err := profile.GetDBPath()
	if err != nil {
		return fmt.Errorf("failed to resolve DB path: %w", err)
	}

	// 4. Inject Token
	// Using Password as Access Token (MVP)
	err = s.injector.Inject(dbPath, string(acc.Password), string(acc.SessionToken), time.Now().Add(1*time.Hour))
	if err != nil {
		return fmt.Errorf("failed to inject into %s: %w", dbPath, err)
	}

	return nil
}

func (s *Service) DeleteAccount(id uint) error {
	return s.db.Delete(&db.Account{}, id).Error
}

// AccountDTO is safe for frontend consumption (no passwords/tokens)
type AccountDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	// Password and SessionToken are explicitly excluded
	RecoveryEmail string `json:"recovery_email"`
	ProxyGroup    string `json:"proxy_group"`
	IsActive      bool   `json:"is_active"`
	Stats         string `json:"stats"` // JSON string
	Tier          string `json:"tier"`  // Extracted from stats
}

func (s *Service) ListAccounts() ([]AccountDTO, error) {
	var accounts []db.Account
	if err := s.db.Find(&accounts).Error; err != nil {
		return nil, err
	}

	dtos := make([]AccountDTO, len(accounts))
	for i, acc := range accounts {
		// Parse Tier from Stats
		tier := "Free"
		if acc.Stats != "" {
			var statsMap map[string]interface{}
			if err := json.Unmarshal([]byte(acc.Stats), &statsMap); err == nil {
				if t, ok := statsMap["tier"].(string); ok {
					tier = t
				}
			}
		}

		dtos[i] = AccountDTO{
			ID:            acc.ID,
			CreatedAt:     acc.CreatedAt,
			UpdatedAt:     acc.UpdatedAt,
			Email:         acc.Email,
			RecoveryEmail: acc.RecoveryEmail,
			ProxyGroup:    acc.ProxyGroup,
			IsActive:      acc.IsActive,
			Stats:         acc.Stats,
			Tier:          tier,
		}
	}
	return dtos, nil
}

func (s *Service) CreateAccount(email, password, recovery, proxyGroup string) error {
	if proxyGroup == "" {
		proxyGroup = "default"
	}

	acc := &db.Account{
		Email:         email,
		Password:      db.EncryptedString(password),
		RecoveryEmail: recovery,
		ProxyGroup:    proxyGroup,
		IsActive:      true,
	}

	// Start Transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 1. Create Record
	if err := tx.Create(acc).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. Validate/Activate (Inject into IDE)
	if err := s.performActivation(acc); err != nil {
		tx.Rollback()
		// Wrap error to be friendly to UI
		return fmt.Errorf("validation failed: could not activate account in IDE: %v", err)
	}

	// 3. Commit
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// GetRotatingToken returns a valid session token from the pool of active accounts.
func (s *Service) GetRotatingToken() (string, error) {
	var accounts []db.Account
	if err := s.db.Where("is_active = ?", true).Find(&accounts).Error; err != nil {
		return "", err
	}

	if len(accounts) == 0 {
		return "", fmt.Errorf("no active accounts available")
	}

	idx := rand.Intn(len(accounts))

	// Check SessionToken first, fallback to Password if SessionToken is empty
	token := string(accounts[idx].SessionToken)
	if token == "" {
		token = string(accounts[idx].Password)
	}

	if token == "" {
		return "", fmt.Errorf("selected account %d has no tokens", accounts[idx].ID)
	}

	return token, nil
}

// UpdateQuota forces a refresh of the account's quota and tier status
func (s *Service) UpdateQuota(id uint) error {
	var acc db.Account
	if err := s.db.First(&acc, id).Error; err != nil {
		return err
	}

	// Used decrypted session token
	token := string(acc.SessionToken)
	if token == "" {
		token = string(acc.Password) // Fallback for MVP
	}

	stats, err := s.quota.FetchAccountStats(token, acc.Email)
	if err != nil {
		if err.Error() == "403 Forbidden" {
			// Mark inactive if forbidden
			acc.IsActive = false
			s.db.Save(&acc)
			return fmt.Errorf("account marked inactive: 403 forbidden")
		}
		return err
	}

	// Serialize stats to JSON
	statsJSON, err := json.Marshal(stats)
	if err == nil {
		acc.Stats = string(statsJSON)
		s.db.Save(&acc)
	}

	return nil
}
