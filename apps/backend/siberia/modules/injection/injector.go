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
	const targetKey = "github.authentication/github.auth" // Real VS Code Key
	const selectSQL = `SELECT value FROM ItemTable WHERE key = ?`
	const updateSQL = `INSERT OR REPLACE INTO ItemTable (key, value) VALUES (?, ?)`

	// A. Read existing blob
	var existingBlob []byte
	err = db.QueryRow(selectSQL, targetKey).Scan(&existingBlob)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to read existing blob: %w", err)
	}

	// B. Prepare new token payload (Access Token)
	// In the Reference App, we just inject the Access Token as the value of Field 6.
	// The rest of the structure (expiry, scopes) is handled by VS Code or we reuse existing if we were smarter,
	// but for now, we follow the "Insert New Field 6" strategy.
	newTokenPayload := []byte(accessToken)

	// C. Modify Protobuf
	var newBlob []byte
	if len(existingBlob) == 0 {
		// If no record exists, we probably need a valid skeleton.
		// For now, let's create a minimal valid blob with just Field 6,
		// though typically VS Code creates the structure first.
		// Strategy: Create empty blob + Field 6.
		newBlob, err = ReplaceField6([]byte{}, newTokenPayload)
	} else {
		newBlob, err = ReplaceField6(existingBlob, newTokenPayload)
	}

	if err != nil {
		return fmt.Errorf("failed to modify protobuf: %w", err)
	}

	// D. Write back
	_, err = db.Exec(updateSQL, targetKey, newBlob)
	if err != nil {
		return fmt.Errorf("failed to inject token: %w", err)
	}

	fmt.Printf("[Injector] Successfully injected token for key: %s (Size: %d -> %d)\n", targetKey, len(existingBlob), len(newBlob))
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
