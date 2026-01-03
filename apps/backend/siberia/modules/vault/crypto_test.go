package vault

import (
	"bytes"
	"testing"
)

func TestVault_EncryptDecrypt(t *testing.T) {
	v := NewVault()
	password := "correct-horse-battery-staple"
	data := []byte("Sensitive User Profile Data")

	// 1. Encrypt
	blob, err := v.Encrypt(data, password)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	if len(blob.Ciphertext) == 0 {
		t.Error("Ciphertext is empty")
	}
	if len(blob.Salt) != 16 {
		t.Error("Invalid salt length")
	}

	// 2. Decrypt Success
	plaintext, err := v.Decrypt(blob, password)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if !bytes.Equal(plaintext, data) {
		t.Errorf("Decrypted data mismatch. Got %s, want %s", plaintext, data)
	}

	// 3. Decrypt Failure (Wrong Password)
	_, err = v.Decrypt(blob, "wrong-password")
	if err == nil {
		t.Error("Expected error for wrong password, got nil")
	}
}

func TestVault_Serialization(t *testing.T) {
	v := NewVault()
	password := "test"
	data := []byte("data")

	blob, _ := v.Encrypt(data, password)

	// To JSON
	jsonBytes, err := blob.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	// From JSON
	restoredBlob, err := FromJSON(jsonBytes)
	if err != nil {
		t.Fatalf("FromJSON failed: %v", err)
	}

	// Verify Equivalency
	if !bytes.Equal(blob.Ciphertext, restoredBlob.Ciphertext) {
		t.Error("Restored ciphertext mismatch")
	}
	if !bytes.Equal(blob.Nonce, restoredBlob.Nonce) {
		t.Error("Restored nonce mismatch")
	}
}
