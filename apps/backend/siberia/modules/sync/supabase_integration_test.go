//go:build integration

package sync

import (
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestSupabaseIntegration(t *testing.T) {
	// 1. Load Env (Assuming running from apps/backend usually, but here run from package dir)
	// We might need to look up for .env
	_ = godotenv.Load("../../../.env") // Walk up from siberia/modules/sync to apps/backend
	// Or try absolute path if known, or just expect env vars

	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")

	if url == "" || key == "" {
		t.Skip("Skipping integration test: SUPABASE credentials not found")
	}

	client := NewSupabaseClient(url, key, "integration-test-user")

	// 2. Test Push
	testBlob := "encrypted-data-blob-" + time.Now().String()
	err := client.Push(testBlob)
	if err != nil {
		t.Fatalf("Push failed: %v", err)
	}

	// 3. Test Pull
	fetchedBlob, err := client.Pull()
	if err != nil {
		t.Fatalf("Pull failed: %v", err)
	}

	if fetchedBlob != testBlob {
		t.Errorf("Mismatch. Expected %s, got %s", testBlob, fetchedBlob)
	}
}
