package logger

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLogAccess(t *testing.T) {
	// 1. Setup Temp Dir
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "logs", "access.log")

	// 2. Init Logger
	InitAccessLogger(tempDir)

	// 3. Log an entry
	entry := AccessEntry{
		Time:       time.Now().Format("15:04:05"),
		Timestamp:  time.Now().Unix(),
		Method:     "GET",
		URL:        "http://example.com",
		Status:     200,
		DurationMs: 50,
		ClientIP:   "127.0.0.1",
		UserAgent:  "Go-Test-Agent",
	}

	LogAccess(entry)

	// 4. Wait for async write (lumberjack might buffer slightly or goroutine needs yield)
	time.Sleep(100 * time.Millisecond)

	// 5. Verify Content
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	var readEntry AccessEntry
	if err := json.Unmarshal(content, &readEntry); err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v. Content: %s", err, string(content))
	}

	if readEntry.Method != "GET" {
		t.Errorf("Expected Method GET, got %s", readEntry.Method)
	}
}
