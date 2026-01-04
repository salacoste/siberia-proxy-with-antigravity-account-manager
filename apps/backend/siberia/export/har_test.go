package export

import (
	"encoding/json"
	"testing"

	"github.com/salacoste/siberia/siberia/types"
)

func TestToHAR(t *testing.T) {
	// 1. Setup Mock Event
	event := types.ProxyRequestEvent{
		Method:   "POST",
		URL:      "https://api.example.com/v1/data",
		Status:   201,
		Duration: 150,
		Time:     "12:00:00",
		Size:     42,
		ReqHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		RespHeaders: map[string]string{
			"Server": "Go",
		},
		ReqBody:  `{"foo":"bar"}`,
		RespBody: `{"status":"ok"}`,
	}

	// 2. Convert
	harJSON, err := ToHAR(event)
	if err != nil {
		t.Fatalf("ToHAR failed: %v", err)
	}

	// 3. Verify Structure
	var har HAR
	if err := json.Unmarshal([]byte(harJSON), &har); err != nil {
		t.Fatalf("Failed to parse HAR JSON: %v", err)
	}

	// 4. Assertions
	if har.Log.Version != "1.2" {
		t.Errorf("Expected version 1.2, got %s", har.Log.Version)
	}
	if len(har.Log.Entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(har.Log.Entries))
	}

	entry := har.Log.Entries[0]
	if entry.Request.Method != "POST" {
		t.Errorf("Expected method POST, got %s", entry.Request.Method)
	}
	if entry.Request.PostData.Text != `{"foo":"bar"}` {
		t.Errorf("PostData mismatch")
	}
	if entry.Response.Status != 201 {
		t.Errorf("Status mismatch")
	}
}
