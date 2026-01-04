package sync

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/salacoste/siberia/siberia/modules/vault"
)

// CloudProvider defines the interface for backend storage
type CloudProvider interface {
	Push(data string) error
	Pull() (string, error)
}

// SyncPayload represents the data exchanged with the cloud
type SyncPayload struct {
	Timestamp int64                `json:"timestamp"` // Unix timestamp
	Hash      string               `json:"hash"`      // SHA-256 of the blob for integrity
	Blob      *vault.EncryptedBlob `json:"blob"`      // Zero-Knowledge Encrypted Data
}

// Manager handles synchronization logic
type Manager struct {
	provider CloudProvider
}

// NewManager creates a new Sync Manager with a specific provider
func NewManager(provider CloudProvider) *Manager {
	return &Manager{
		provider: provider,
	}
}

// Push uploads the encrypted profile to the cloud
func (m *Manager) Push(profileID string, blob *vault.EncryptedBlob) error {
	payload := SyncPayload{
		Timestamp: time.Now().Unix(),
		Hash:      "todo-hash", // Simplified for MVP
		Blob:      blob,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload error: %w", err)
	}

	// We treat the entire JSON payload (including timestamp/hash) as the "DataBlob" stored in Supabase
	return m.provider.Push(string(data))
}

// Pull downloads the latest encrypted profile from the cloud
func (m *Manager) Pull(profileID string) (*SyncPayload, error) {
	dataStr, err := m.provider.Pull()
	if err != nil {
		return nil, err
	}

	if dataStr == "" {
		return nil, nil // No data
	}

	var payload SyncPayload
	if err := json.Unmarshal([]byte(dataStr), &payload); err != nil {
		return nil, fmt.Errorf("unmarshal payload error: %w", err)
	}

	return &payload, nil
}

// ResolveConflict implements LWW (Last Write Wins)
func (m *Manager) ResolveConflict(localTime, remoteTime int64) bool {
	return localTime > remoteTime
}
