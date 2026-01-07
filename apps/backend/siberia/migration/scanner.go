package migration

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

type IDEProfile struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	ConfigPath  string `json:"config_path"`
	DBPath      string `json:"db_path"`
	ProcessName string `json:"process_name"`
	IsRunning   bool   `json:"is_running"`
}

// ScanIDEProfiles searches for known JetBrains/Android Studio directories.
func ScanIDEProfiles() ([]IDEProfile, error) {
	var profiles []IDEProfile

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home dir: %w", err)
	}

	// Paths vary by OS
	// macOS: ~/Library/Application Support/Google/AndroidStudio*
	//        ~/Library/Application Support/JetBrains/IntelliJIdea*
	// Linux: ~/.config/Google/AndroidStudio*
	//        ~/.config/JetBrains/IntelliJIdea*

	var searchPaths []string

	switch runtime.GOOS {
	case "darwin":
		appSupport := filepath.Join(homeDir, "Library", "Application Support")
		searchPaths = append(searchPaths,
			filepath.Join(appSupport, "Google"),
			filepath.Join(appSupport, "JetBrains"),
		)
	case "linux":
		configDir := filepath.Join(homeDir, ".config")
		searchPaths = append(searchPaths,
			filepath.Join(configDir, "Google"),
			filepath.Join(configDir, "JetBrains"),
		)
	default:
		// Windows is distinct, add later if needed.
		return nil, nil // Not fully supported yet for MVP scanning
	}

	for _, baseDir := range searchPaths {
		entries, err := os.ReadDir(baseDir)
		if err != nil {
			continue // Skip if dir doesn't exist
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			name := entry.Name()
			fullPath := filepath.Join(baseDir, name)

			// Heuristic: Check for AndroidStudio or IntelliJ or PyCharm
			if isSupportedIDE(name) {
				// Look for storage.db?
				// Typically in `options/other.xml` OR a central `storage.db`?
				// Ref App says "storage.db" implies SQLite.
				// JetBrains moved config to SQLite in recent versions?
				// Or is it `onediff-local-storage.db`?
				// Assuming standard JetBrains storage layout:
				// Sometimes `workspace` directory.
				// Actually, for "Deep State" tokens:
				// It's often in a plugin specific DB or the main IDE config.
				// Let's assume requirements imply finding `storage.db`.
				// If not found, we skip.

				// Common location for modern InteliJ/AS is the config root directly or specific subdir?
				// Actually, many plugins dump to `options/other.xml`.
				// But migrating to SQLite is common.
				// Let's look for ANY `.db` file or assume it's in root?
				// Let's look for `storage.db` specifically.

				// Note: `storage.db` is often at the root of the config dir for that version.
				dbPath := filepath.Join(fullPath, "storage.db")
				if _, err := os.Stat(dbPath); err == nil {
					// Found!
					profiles = append(profiles, IDEProfile{
						Name:        name,
						ConfigPath:  fullPath,
						DBPath:      dbPath,
						ProcessName: mapProcessName(name),
					})
				} else {
					// Fallback: Searching recursively?
					// No, keep it fast.
				}
			}
		}
	}

	return profiles, nil
}

func isSupportedIDE(name string) bool {
	lower := strings.ToLower(name)
	return strings.Contains(lower, "androidstudio") ||
		strings.Contains(lower, "intellij") ||
		strings.Contains(lower, "pycharm") ||
		strings.Contains(lower, "goland") ||
		strings.Contains(lower, "webstorm")
}

func mapProcessName(dirName string) string {
	lower := strings.ToLower(dirName)
	if strings.Contains(lower, "androidstudio") {
		return "studio"
	}
	if strings.Contains(lower, "intellij") {
		return "idea"
	}
	if strings.Contains(lower, "pycharm") {
		return "pycharm"
	}
	if strings.Contains(lower, "goland") {
		return "goland"
	}
	if strings.Contains(lower, "webstorm") {
		return "webstorm"
	}
	return "java" // Generic fallback
}

// CheckRunningProcesses updates the IsRunning status
func (s *IDEProfile) CheckIsRunning() bool {
	procs, _ := process.Processes()
	for _, p := range procs {
		name, _ := p.Name()
		if strings.Contains(strings.ToLower(name), s.ProcessName) {
			s.IsRunning = true
			return true
		}
	}
	return false
}
