package injection

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

func TestService_Inject(t *testing.T) {
	// 1. Create a Temp DB
	tempDir, err := os.MkdirTemp("", "siberia_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // Cleanup

	dbPath := filepath.Join(tempDir, "state.vscdb")

	// 2. Initialize DB with Schema
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("Failed to create db: %v", err)
	}

	_, err = db.Exec("CREATE TABLE ItemTable (key TEXT PRIMARY KEY, value TEXT)")
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	db.Close() // Close so Injector can open it

	// 3. Run Injector
	svc := NewService()
	err = svc.Inject(dbPath, "test-access", "test-refresh", time.Now().Add(1*time.Hour))
	if err != nil {
		t.Fatalf("Inject failed: %v", err)
	}

	// 4. Verify Backup Exists
	if _, err := os.Stat(dbPath + ".bak"); os.IsNotExist(err) {
		t.Errorf("Backup file was not created")
	}

	// 5. Verify Data Inserted
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("Failed to re-open db: %v", err)
	}
	defer db.Close()

	var value string
	err = db.QueryRow("SELECT value FROM ItemTable WHERE key = ?", "siberia.auth_token").Scan(&value)
	if err != nil {
		t.Fatalf("Failed to query inserted token: %v", err)
	}

	t.Logf("Injected Value: %s", value)

	// Basic string check
	if value == "" {
		t.Error("Injected value is empty")
	}
}
