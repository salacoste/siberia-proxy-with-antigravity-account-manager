package accounts

import (
	"time"

	"github.com/salacoste/siberia/siberia/db"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(database *db.Database) *Service {
	return &Service{
		db: database.Conn,
	}
}

// AccountDTO is safe for frontend consumption (no passwords/tokens)
type AccountDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	// Password and SessionToken are explicitly excluded
	ProxyGroup string `json:"proxy_group"`
	IsActive   bool   `json:"is_active"`
	Stats      string `json:"stats"` // JSON string
}

func (s *Service) ListAccounts() ([]AccountDTO, error) {
	var accounts []db.Account
	if err := s.db.Find(&accounts).Error; err != nil {
		return nil, err
	}

	dtos := make([]AccountDTO, len(accounts))
	for i, acc := range accounts {
		dtos[i] = AccountDTO{
			ID:         acc.ID,
			CreatedAt:  acc.CreatedAt,
			UpdatedAt:  acc.UpdatedAt,
			Email:      acc.Email,
			ProxyGroup: acc.ProxyGroup,
			IsActive:   acc.IsActive,
			Stats:      acc.Stats,
		}
	}
	return dtos, nil
}

func (s *Service) DeleteAccount(id uint) error {
	return s.db.Delete(&db.Account{}, id).Error
}
