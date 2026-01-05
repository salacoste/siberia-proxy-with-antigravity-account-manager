package injection

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

func TestInjector_Inject_Integration(t *testing.T) {
	// 1. Setup Temp Dir and DB
	tempDir, err := os.MkdirTemp("", "siberia_injector_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "state.vscdb")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("Failed to open temp db: %v", err)
	}
	defer db.Close()

	// 2. Init DB Schema and Data
	initSQL := `
	CREATE TABLE ItemTable (key TEXT PRIMARY KEY, value BLOB);
	INSERT INTO ItemTable (key, value) VALUES (?, ?);
	`
	// Create a fake blob: Field 1 (08 96 01)
	fakeBlob := []byte{0x08, 0x96, 0x01}
	targetKey := "github.authentication/github.auth"

	_, err = db.Exec(initSQL, targetKey, fakeBlob)
	if err != nil {
		t.Fatalf("Failed to init db: %v", err)
	}
	db.Close() // Close so Injector can open it

	// 3. Run Injector
	svc := NewService()
	newToken := "TEST_ACCESS_TOKEN"
	expiry := time.Now().Add(1 * time.Hour)

	err = svc.Inject(dbPath, newToken, "refresh_token", expiry)
	if err != nil {
		t.Fatalf("Inject failed: %v", err)
	}

	// 4. Verify Backup Exists
	if _, err := os.Stat(dbPath + ".bak"); os.IsNotExist(err) {
		t.Errorf("Backup file was not created")
	}

	// 5. Verify Data Modified
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("Failed to reopen db: %v", err)
	}
	defer db.Close()

	var newBlob []byte
	err = db.QueryRow("SELECT value FROM ItemTable WHERE key = ?", targetKey).Scan(&newBlob)
	if err != nil {
		t.Fatalf("Failed to read back value: %v", err)
	}

	// Check if larger than original (simple heuristic)
	// Original (3 bytes) + Field 6 tag (1) + len (1) + val (17 bytes) ~= 22 bytes
	if len(newBlob) <= len(fakeBlob) {
		t.Errorf("Blob did not grow. Len: %d", len(newBlob))
	}

	// Warning: We are not decoding the proto here to verify exact contents,
	// relying on unit tests for that. Here we verify DB interaction.
	t.Logf("Injection successful. Blob changed from %d to %d bytes.", len(fakeBlob), len(newBlob))
}
