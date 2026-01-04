package accounts

import (
	"fmt"
	"time"

	"github.com/salacoste/siberia/siberia/config"
	"github.com/salacoste/siberia/siberia/db"
	"github.com/salacoste/siberia/siberia/ide"
	"github.com/salacoste/siberia/siberia/modules/injection"
	"github.com/salacoste/siberia/siberia/modules/process"
	"gorm.io/gorm"
)

type Service struct {
	db        *gorm.DB
	process   *process.Service
	injector  *injection.Service
	config    *config.Manager
	configDir string
}

func NewService(database *db.Database, cfg *config.Manager) *Service {
	return &Service{
		db:        database.Conn,
		process:   process.NewService(),
		injector:  injection.NewService(),
		config:    cfg,
		configDir: cfg.ConfigDir(),
	}
}

// ... existing code ...

func (s *Service) ActivateAccount(id uint) error {
	var acc db.Account
	if err := s.db.First(&acc, id).Error; err != nil {
		return err
	}

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

	fmt.Printf("[Account] Activating for IDE: %s (Process: %s)\n", profile.Name, profile.ProcessName)

	// 2. Kill Target Process
	if err := s.process.Kill(profile.ProcessName); err != nil {
		return fmt.Errorf("failed to kill process %s: %w", profile.ProcessName, err)
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

	// 5. Start Process
	// We can't strictly "Start" the app easily from binary path across all OSes consistently without knowing where the App bundle is.
	// IdeProfile doesn't store App Path yet.
	// For now, we only Kill. Starting usually requires user to click the icon, specifically for "Code Helper" which is a background process of VS Code.
	// VS Code relaunch is tricky programmatically (`code .` works if in PATH).
	// Let's try to start if we know the path, otherwise just log instructions.
	// Since we don't store App Path in Profile yet, we'll keep the process.Start generic logic or skip it.
	// "Code Helper" is internal. We usually want to start the Main App?
	// The original code tried: s.process.Start(targetApp).
	// If targetApp was "Code Helper", that likely failed or started a helper.
	// Let's Skip Auto-Start for now (Wait for user action) OR simple `open -a "Visual Studio Code"` on Mac.
	// Let's leave Start empty or commented out for safety until we add AppPath to Profile.
	fmt.Println("[Account] Please restart your IDE manually if it doesn't open.")

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
}

func (s *Service) ListAccounts() ([]AccountDTO, error) {
	var accounts []db.Account
	if err := s.db.Find(&accounts).Error; err != nil {
		return nil, err
	}

	dtos := make([]AccountDTO, len(accounts))
	for i, acc := range accounts {
		dtos[i] = AccountDTO{
			ID:            acc.ID,
			CreatedAt:     acc.CreatedAt,
			UpdatedAt:     acc.UpdatedAt,
			Email:         acc.Email,
			RecoveryEmail: acc.RecoveryEmail,
			ProxyGroup:    acc.ProxyGroup,
			IsActive:      acc.IsActive,
			Stats:         acc.Stats,
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

	return s.db.Create(acc).Error
}
