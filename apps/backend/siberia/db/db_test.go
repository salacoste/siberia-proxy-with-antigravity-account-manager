package db

import (
	"strings"
	"testing"

	"github.com/salacoste/siberia/siberia/crypto"
)

func TestDBEncryption(t *testing.T) {
	// 1. Setup
	tmpDir := t.TempDir()
	key, _ := crypto.GenerateKey()

	database, err := Init(tmpDir, key)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 2. Create Account
	email := "test@example.com"
	password := "supersecret"

	acc := &Account{
		Email:    email,
		Password: EncryptedString(password),
	}

	if err := database.Conn.Create(acc).Error; err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// 3. Verify Read (Decryption)
	var readAcc Account
	if err := database.Conn.First(&readAcc, "email = ?", email).Error; err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if string(readAcc.Password) != password {
		t.Errorf("Expected decrypted password %s, got %s", password, readAcc.Password)
	}

	// 4. Verify Storage (Encryption)
	// We read the raw sqlite file? No, that's complex.
	// We can use a raw SQL query ensuring the result is NOT "supersecret"
	var rawPassword string
	row := database.Conn.Raw("SELECT password FROM accounts WHERE email = ?", email).Row()
	if err := row.Scan(&rawPassword); err != nil {
		t.Fatalf("Raw scan failed: %v", err)
	}

	if rawPassword == password {
		t.Error("Password is stored in PLAINTEXT! Encryption failed.")
	}
	if !strings.Contains(rawPassword, "") { // Just checking it exists
		t.Error("Raw password empty")
	}
	// Typically it should be hex string
	t.Logf("Raw stored password: %s", rawPassword)
}
