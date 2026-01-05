package ide

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type IdeProfile struct {
	ID   string
	Name string
}

// GetProcessName returns the OS-specific process name to kill
func (p IdeProfile) GetProcessName() string {
	switch p.ID {
	case "vscode":
		switch runtime.GOOS {
		case "windows":
			return "Code.exe"
		case "linux":
			return "code"
		default:
			return "Code Helper" // macOS
		}
	case "cursor":
		switch runtime.GOOS {
		case "windows":
			return "Cursor.exe"
		case "linux":
			return "cursor"
		default:
			return "Cursor"
		}
	case "windsurf":
		switch runtime.GOOS {
		case "windows":
			return "Windsurf.exe"
		case "linux":
			return "windsurf"
		default:
			return "Windsurf"
		}
	default:
		return ""
	}
}

var Registry = map[string]IdeProfile{
	"vscode": {
		ID:   "vscode",
		Name: "VS Code",
	},
	"cursor": {
		ID:   "cursor",
		Name: "Cursor",
	},
	"windsurf": {
		ID:   "windsurf",
		Name: "Windsurf",
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

	var suffix string
	switch p.ID {
	case "vscode":
		switch runtime.GOOS {
		case "windows":
			suffix = "AppData/Roaming/Code/User/globalStorage/state.vscdb"
		case "linux":
			suffix = ".config/Code/User/globalStorage/state.vscdb"
		default: // darwin
			suffix = "Library/Application Support/Code/User/globalStorage/state.vscdb"
		}
	case "cursor":
		switch runtime.GOOS {
		case "windows":
			suffix = "AppData/Roaming/Cursor/User/globalStorage/state.vscdb"
		case "linux":
			suffix = ".config/Cursor/User/globalStorage/state.vscdb"
		default: // darwin
			suffix = "Library/Application Support/Cursor/User/globalStorage/state.vscdb"
		}
	case "windsurf":
		switch runtime.GOOS {
		case "windows":
			suffix = "AppData/Roaming/Windsurf/User/globalStorage/state.vscdb"
		case "linux":
			suffix = ".config/Windsurf/User/globalStorage/state.vscdb"
		default: // darwin
			suffix = "Library/Application Support/Windsurf/User/globalStorage/state.vscdb"
		}
	}

	if suffix == "" {
		return "", fmt.Errorf("unsupported OS or IDE for path resolution")
	}

	return filepath.Join(home, suffix), nil
}
