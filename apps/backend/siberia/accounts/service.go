package accounts

import (
	"fmt"
	"time"

	"github.com/salacoste/siberia/siberia/db"
	"github.com/salacoste/siberia/siberia/modules/injection"
	"github.com/salacoste/siberia/siberia/modules/process"
	"gorm.io/gorm"
)

type Service struct {
	db       *gorm.DB
	process  process.Manager
	injector injection.Injector
}

func NewService(database *db.Database) *Service {
	return &Service{
		db:       database.Conn,
		process:  process.NewMockManager(true), // Dry-run by default
		injector: injection.NewMockInjector(),
	}
}

// ... existing code ...

func (s *Service) ActivateAccount(id uint) error {
	var acc db.Account
	if err := s.db.First(&acc, id).Error; err != nil {
		return err
	}

	// 1. Decrypt credentials (happens automatically via Value(), but here we access the raw string from Scanner?
	// Actually EncryptedString has Value() for DB write, and Scan() for DB read.
	// Verify: When we read `acc`, `acc.Password` is already the `EncryptedString` type,
	// but the underlying string (if we cast it) is the decrypted value?
	// Wait, EncryptedString is `string` type alias.
	// db.go implementation: Scan() decrypts bytes -> string. So `acc.Password` holds DECRYPTED string.
	// Correct.

	// 2. Kill Target Process
	// Hardcoded for now, move to config later
	targetApp := "Code Helper"
	if err := s.process.Kill(targetApp); err != nil {
		return fmt.Errorf("failed to kill process: %w", err)
	}

	// 3. Inject Token
	// Hardcoded path for now
	targetDB := "/Users/r2d2/Library/Application Support/Code/User/globalStorage/state.vscdb"

	// Use Password as Access Token for MVP/Demo purposes since we don't have separate token fields yet
	// In real world, we'd use SessionToken or specific fields
	err := s.injector.Inject(targetDB, string(acc.Password), string(acc.SessionToken), time.Now().Add(1*time.Hour))
	if err != nil {
		return fmt.Errorf("failed to inject: %w", err)
	}

	// 4. Start Process
	if err := s.process.Start(targetApp); err != nil {
		return fmt.Errorf("failed to start process: %w", err)
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
