package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"

	"golang.org/x/crypto/argon2"
)

// EncryptedBlob represents the stored format of the encrypted data
type EncryptedBlob struct {
	Salt       []byte `json:"salt"`       // Salt used for Argon2id
	Nonce      []byte `json:"nonce"`      // Nonce for AES-GCM
	Ciphertext []byte `json:"ciphertext"` // Encrypted Payload
}

// Vault manages encryption and decryption logic
type Vault struct{}

// NewVault creates a new Vault instance
func NewVault() *Vault {
	return &Vault{}
}

// DeriveKey derives a 32-byte key from a password and salt using Argon2id
func (v *Vault) DeriveKey(password string, salt []byte) []byte {
	// Argon2id parameters (OWASP recommendations)
	// Time: 1, Memory: 64MB (64*1024), Threads: 4, KeyLen: 32
	return argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
}

// GenerateSalt creates a random 16-byte salt
func (v *Vault) GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	return salt, nil
}

// Encrypt encrypts data using AES-256-GCM with a key derived from the password
func (v *Vault) Encrypt(data []byte, password string) (*EncryptedBlob, error) {
	// 1. Generate Salt
	salt, err := v.GenerateSalt()
	if err != nil {
		return nil, err
	}

	// 2. Derive Key
	key := v.DeriveKey(password, salt)

	// 3. Create Cipher Block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 4. Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 5. Generate Nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// 6. Encrypt
	ciphertext := gcm.Seal(nil, nonce, data, nil)

	return &EncryptedBlob{
		Salt:       salt,
		Nonce:      nonce,
		Ciphertext: ciphertext,
	}, nil
}

// Decrypt decrypts an EncryptedBlob using the password
func (v *Vault) Decrypt(blob *EncryptedBlob, password string) ([]byte, error) {
	// 1. Derive Key using stored salt
	key := v.DeriveKey(password, blob.Salt)

	// 2. Create Cipher Block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 3. Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 4. Decrypt
	plaintext, err := gcm.Open(nil, blob.Nonce, blob.Ciphertext, nil)
	if err != nil {
		return nil, errors.New("decryption failed: invalid password or corrupted data")
	}

	return plaintext, nil
}

// Serialize converts EncryptedBlob to JSON
func (b *EncryptedBlob) ToJSON() ([]byte, error) {
	return json.Marshal(b)
}

// Deserialize parses JSON into EncryptedBlob
func FromJSON(data []byte) (*EncryptedBlob, error) {
	var blob EncryptedBlob
	err := json.Unmarshal(data, &blob)
	return &blob, err
}
