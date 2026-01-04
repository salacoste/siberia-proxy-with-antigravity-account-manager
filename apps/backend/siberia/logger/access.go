package logger

import (
	"encoding/json"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

type AccessEntry struct {
	Time       string `json:"time"`
	Timestamp  int64  `json:"timestamp"`
	Method     string `json:"method"`
	URL        string `json:"url"`
	Status     int    `json:"status"`
	DurationMs int64  `json:"duration_ms"`
	ClientIP   string `json:"client_ip"`
	UserAgent  string `json:"user_agent"`
	Size       int64  `json:"size"` // Response size if captured, or 0
}

var accessLogger *lumberjack.Logger

// InitAccessLogger initializes the rotating file logger
func InitAccessLogger(configDir string) {
	logPath := filepath.Join(configDir, "logs", "access.log")
	accessLogger = &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10, // megabytes
		MaxBackups: 5,
		MaxAge:     28,   // days
		Compress:   true, // disabled by default
	}
}

// LogAccess writes an entry to the access log.
// It is safely called from the TelemetryWorker so it doesn't need its own goroutine per request.
func LogAccess(entry AccessEntry) {
	if accessLogger == nil {
		return
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return
	}
	// Append newline
	data = append(data, '\n')
	accessLogger.Write(data)
}
