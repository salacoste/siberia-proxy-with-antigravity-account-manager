package ide

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type IdeProfile struct {
	ID           string
	Name         string
	ProcessName  string
	DBPathSuffix string // Relative to User Home
}

var Registry = map[string]IdeProfile{
	"vscode": {
		ID:           "vscode",
		Name:         "VS Code",
		ProcessName:  "Code Helper", // macOS specific, might need OS check later
		DBPathSuffix: "Library/Application Support/Code/User/globalStorage/state.vscdb",
	},
	"cursor": {
		ID:           "cursor",
		Name:         "Cursor",
		ProcessName:  "Cursor",
		DBPathSuffix: "Library/Application Support/Cursor/User/globalStorage/state.vscdb",
	},
	"windsurf": {
		ID:           "windsurf",
		Name:         "Windsurf",
		ProcessName:  "Windsurf",
		DBPathSuffix: "Library/Application Support/Windsurf/User/globalStorage/state.vscdb",
	},
}

func GetProfile(id string) (IdeProfile, error) {
	profile, ok := Registry[id]
	if !ok {
		return IdeProfile{}, fmt.Errorf("unknown IDE profile: %s", id)
	}
	return profile, nil
}

func (p IdeProfile) GetDBPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// TODO: Handle Windows/Linux paths if needed in future. Assuming macOS for now based on project context.
	if runtime.GOOS != "darwin" {
		// Fallback or specific logic for other OS
	}
	return filepath.Join(home, p.DBPathSuffix), nil
}
