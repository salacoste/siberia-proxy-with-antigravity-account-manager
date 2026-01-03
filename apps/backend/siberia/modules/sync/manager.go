package sync

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/salacoste/siberia/siberia/modules/vault"
)

// SyncPayload represents the data exchanged with the cloud
type SyncPayload struct {
	Timestamp int64                `json:"timestamp"` // Unix timestamp
	Hash      string               `json:"hash"`      // SHA-256 of the blob for integrity
	Blob      *vault.EncryptedBlob `json:"blob"`      // Zero-Knowledge Encrypted Data
}

// Manager handles synchronization logic
type Manager struct {
	client  *http.Client
	baseURL string
}

// NewManager creates a new Sync Manager
func NewManager(baseURL string) *Manager {
	return &Manager{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: baseURL,
	}
}

// Push uploads the encrypted profile to the cloud
// Returns error if network fails or 409 Conflict occurs (requiring a Pull first)
func (m *Manager) Push(profileID string, blob *vault.EncryptedBlob) error {
	payload := SyncPayload{
		Timestamp: time.Now().Unix(),
		Hash:      "todo-hash", // Simplified for MVP
		Blob:      blob,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/sync/%s", m.baseURL, profileID)
	resp, err := m.client.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		return errors.New("conflict: remote version is newer")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("sync failed with status: %d", resp.StatusCode)
	}

	return nil
}

// Pull downloads the latest encrypted profile from the cloud
func (m *Manager) Pull(profileID string) (*SyncPayload, error) {
	url := fmt.Sprintf("%s/sync/%s", m.baseURL, profileID)
	resp, err := m.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // No remote data
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pull failed with status: %d", resp.StatusCode)
	}

	var payload SyncPayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

// ResolveConflict implements LWW (Last Write Wins)
// For MVP, if timestamps differ, we bias towards the one with the later timestamp.
// Returns true if local wins (needs Push), false if remote wins (needs Apply).
func (m *Manager) ResolveConflict(localTime, remoteTime int64) bool {
	return localTime > remoteTime
}
