package sync

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/salacoste/siberia/siberia/modules/vault"
)

// MockProvider for testing logic without HTTP
type MockProvider struct {
	LastPushData string
	PullData     string
	PushError    error
	PullError    error
}

func (m *MockProvider) Push(data string) error {
	if m.PushError != nil {
		return m.PushError
	}
	m.LastPushData = data
	return nil
}

func (m *MockProvider) Pull() (string, error) {
	if m.PullError != nil {
		return "", m.PullError
	}
	return m.PullData, nil
}

func TestManager_Push(t *testing.T) {
	mock := &MockProvider{}
	m := NewManager(mock)

	blob := &vault.EncryptedBlob{
		Ciphertext: []byte("test"),
		Nonce:      []byte("nonce"),
		Salt:       []byte("salt"),
	}
	err := m.Push("ignored-id", blob)
	if err != nil {
		t.Fatalf("Push failed: %v", err)
	}

	if mock.LastPushData == "" {
		t.Fatal("Mock provider did not receive data")
	}

	// Verify payload structure
	var payload SyncPayload
	if err := json.Unmarshal([]byte(mock.LastPushData), &payload); err != nil {
		t.Fatalf("Failed to unmarshal push payload: %v", err)
	}

	if string(payload.Blob.Ciphertext) != "test" {
		t.Errorf("Expected ciphertext 'test', got %s", payload.Blob.Ciphertext)
	}
}

func TestManager_Pull(t *testing.T) {
	blob := &vault.EncryptedBlob{
		Ciphertext: []byte("remote"),
		Nonce:      []byte("n"),
		Salt:       []byte("s"),
	}
	payload := SyncPayload{Timestamp: 123, Blob: blob}
	data, _ := json.Marshal(payload)

	mock := &MockProvider{PullData: string(data)}
	m := NewManager(mock)

	pPayload, err := m.Pull("ignored-id")
	if err != nil {
		t.Fatalf("Pull failed: %v", err)
	}

	if !bytes.Equal(pPayload.Blob.Ciphertext, []byte("remote")) {
		t.Errorf("Expected 'remote', got %s", pPayload.Blob.Ciphertext)
	}
}
