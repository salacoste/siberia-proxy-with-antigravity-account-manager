package db

import (
	"database/sql/driver"
	"errors"
	"time"

	"github.com/salacoste/siberia/siberia/crypto"
	"gorm.io/gorm"
)

var masterKey string

// SetMasterKey sets the encryption key for the package
func SetMasterKey(key string) {
	masterKey = key
}

// EncryptedString handles transparent encryption/decryption
type EncryptedString string

// Value interface for database driver
func (s EncryptedString) Value() (driver.Value, error) {
	if masterKey == "" {
		return nil, errors.New("master key not set")
	}
	if s == "" {
		return "", nil
	}
	return crypto.Encrypt(string(s), masterKey)
}

// Scan interface for database driver
func (s *EncryptedString) Scan(value interface{}) error {
	if masterKey == "" {
		return errors.New("master key not set")
	}

	bytes, ok := value.([]byte)
	if !ok {
		str, ok := value.(string)
		if !ok {
			return errors.New("failed to unmarshal JSONB value")
		}
		bytes = []byte(str)
	}

	if len(bytes) == 0 {
		*s = ""
		return nil
	}

	decrypted, err := crypto.Decrypt(string(bytes), masterKey)
	if err != nil {
		return err
	}
	*s = EncryptedString(decrypted)
	return nil
}

type Account struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Email         string          `gorm:"uniqueIndex" json:"email"`
	Password      EncryptedString `json:"-"` // Never expose via JSON
	RecoveryEmail string          `json:"recovery_email"`
	SessionToken  EncryptedString `json:"-"` // Never expose via JSON
	ProxyGroup    string          `gorm:"default:default" json:"proxy_group"`
	IsActive      bool            `gorm:"default:true" json:"is_active"`
	Stats         string          `json:"stats"` // JSON string
}
