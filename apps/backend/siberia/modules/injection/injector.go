package injection

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	// Import the pure Go SQLite driver
	_ "github.com/glebarez/go-sqlite"
)

// Injector handles token injection into external app DBs
type Injector interface {
	Inject(dbPath string, accessToken, refreshToken string, expiry time.Time) error
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (i *Service) Inject(dbPath string, accessToken, refreshToken string, expiry time.Time) error {
	fmt.Printf("[Injector] Injecting into DB: %s\n", dbPath)

	// 1. Check if file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return fmt.Errorf("target database not found: %s", dbPath)
	}

	// 2. Create Backup
	backupPath := dbPath + ".bak"
	if err := copyFile(dbPath, backupPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}
	fmt.Printf("[Injector] Backup created at: %s\n", backupPath)

	// 3. Open Database
	// Using generic "sqlite" driver name which glebarez/go-sqlite registers
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// 4. Update Token
	// Using a generic key for demonstration. In reality, this matches the target extension's key.
	// VS Code state.vscdb table is 'ItemTable' (key TEXT, value TEXT)
	const targetKey = "siberia.auth_token"
	const updateSQL = `INSERT OR REPLACE INTO ItemTable (key, value) VALUES (?, ?)`

	// Value is typically a JSON string in VS Code
	tokenValue := fmt.Sprintf(`{"accessToken":"%s","refreshToken":"%s"}`, accessToken, refreshToken)

	_, err = db.Exec(updateSQL, targetKey, tokenValue)
	if err != nil {
		return fmt.Errorf("failed to inject token: %w", err)
	}

	fmt.Printf("[Injector] Successfully injected token for key: %s\n", targetKey)
	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create directory if needed
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
