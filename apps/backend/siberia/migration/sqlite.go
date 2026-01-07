package migration

import (
	"database/sql"
	"fmt"

	_ "github.com/glebarez/go-sqlite"
)

// ReadStateToken opens the target SQLite file in read-only mode and extracts the state blob.
func ReadStateToken(dbPath string) (string, error) {
	// Construct DSN for Read-Only mode
	// glebarez/go-sqlite uses normal file paths, but query params for mode?
	// It doesn't strictly support ?mode=ro in the same way CGO sqlite3 does,
	// but let's try standard checking or just open.
	// Actually, preventing locks is key.
	// "file:path?mode=ro" is standard URI format.

	dsn := fmt.Sprintf("%s?mode=ro&_pragma=query_only(true)", dbPath)

	// Since dbPath might contain spaces, we might need to url encode, but local file paths are tricky.
	// Let's rely on the driver handling paths. typical logic:
	// If the driver supports URI filenames, we should use them.
	// glebarez/go-sqlite supports standard behavior.

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return "", fmt.Errorf("failed to open db: %w", err)
	}
	defer db.Close()

	// Verify connection
	if err := db.Ping(); err != nil {
		return "", fmt.Errorf("failed to ping db: %w", err)
	}

	// Query ItemTable
	// The key is 'jetskiStateSync.agentManagerInitState'
	query := `SELECT value FROM ItemTable WHERE key = ? LIMIT 1`

	// We might need to check if table exists first?
	// Or just try query and handle error.

	var value string
	err = db.QueryRow(query, "jetskiStateSync.agentManagerInitState").Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("token key not found")
		}
		return "", fmt.Errorf("query error: %w", err)
	}

	return value, nil
}
