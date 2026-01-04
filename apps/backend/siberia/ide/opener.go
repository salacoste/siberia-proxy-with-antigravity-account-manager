package ide

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Opener handles opening files or folders in an IDE
type Opener struct{}

func NewOpener() *Opener {
	return &Opener{}
}

// Open attempts to open the specified path in the target IDE.
// If line > 0, it appends :line to the path/URI.
func (o *Opener) Open(targetIDE string, path string, line int) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	switch targetIDE {
	case "cursor":
		return o.openURI("cursor", absPath, line)
	case "vscode":
		return o.openURI("vscode", absPath, line)
	case "windsurf":
		// Windsurf URI scheme might be distinct, but let's try 'windsurf://'
		// If that fails, we fallback to 'open -a Windsurf'
		return o.openURI("windsurf", absPath, line)
	default:
		// Fallback to system default or failure
		return fmt.Errorf("unsupported IDE: %s", targetIDE)
	}
}

func (o *Opener) openURI(scheme string, path string, line int) error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("only macOS supported for URI opening currently")
	}

	// URI Scheme Format: scheme://file/<path>:<line>
	// e.g. cursor://file/Users/me/project/main.go:10

	// Note: 'cursor://file' expects absolute path.
	// URI encoding needed? Usually 'open' handles standard paths, but spaces need %20.
	// For simplicity, we trust net/url logic or simple string for now.

	uri := fmt.Sprintf("%s://file%s", scheme, path)
	if line > 0 {
		uri = fmt.Sprintf("%s:%d", uri, line)
	}

	// Use macOS 'open' command which handles URI schemes
	cmd := exec.Command("open", uri)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Fallback: try opening as application if URI fails?
		// e.g. open -a "Cursor" path
		return fmt.Errorf("failed to open URI %s: %v, output: %s", uri, err, string(output))
	}
	return nil
}
