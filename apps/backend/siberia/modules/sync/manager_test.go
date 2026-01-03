package sync

import (
	"testing"
	"time"

	"github.com/salacoste/siberia/siberia/modules/vault"
)

func TestSyncManager_PushPull(t *testing.T) {
	// 1. Start Mock Server
	mock := NewMockServer()
	ts := mock.Start()
	defer ts.Close()

	// 2. Init Manager
	manager := NewManager(ts.URL)
	profileID := "user-123"

	// 3. Prepare Data
	blob := &vault.EncryptedBlob{
		Ciphertext: []byte("encrypted-data"),
		Nonce:      []byte("nonce"),
		Salt:       []byte("salt"),
	}

	// 4. Test Push
	err := manager.Push(profileID, blob)
	if err != nil {
		t.Fatalf("Push failed: %v", err)
	}

	// 5. Test Pull
	pulled, err := manager.Pull(profileID)
	if err != nil {
		t.Fatalf("Pull failed: %v", err)
	}

	if string(pulled.Blob.Ciphertext) != "encrypted-data" {
		t.Errorf("Pulled data mismatch")
	}
}

func TestSyncManager_Conflict(t *testing.T) {
	// 1. Start Mock Server
	mock := NewMockServer()
	ts := mock.Start()
	defer ts.Close()

	manager := NewManager(ts.URL)
	profileID := "conflict-user"

	// 2. Simulate Server State (Newer)
	mock.Store[profileID] = SyncPayload{
		Timestamp: time.Now().Add(1 * time.Hour).Unix(),
		Blob: &vault.EncryptedBlob{
			Ciphertext: []byte("remote-newer"),
		},
	}

	// 3. Try Push with Older Data
	blob := &vault.EncryptedBlob{
		Ciphertext: []byte("local-older"),
	}

	err := manager.Push(profileID, blob)

	// 4. Expect Conflict
	if err == nil || err.Error() != "conflict: remote version is newer" {
		t.Fatalf("Expected conflict error, got: %v", err)
	}

	// 5. Client Logic: Should Pull
	pulled, _ := manager.Pull(profileID)
	if string(pulled.Blob.Ciphertext) != "remote-newer" {
		t.Error("Should be able to pull remote data after conflict")
	}
}
