package process

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

// Manager handles external process control
type Manager interface {
	Kill(name string) error
	Start(path string) error
}

type Service struct {
	DryRun bool
}

func NewService() *Service {
	return &Service{DryRun: false} // Default to real implementation
}

func (m *Service) Kill(name string) error {
	if m.DryRun {
		fmt.Printf("[Process] [DRY RUN] Kill process matching: %s\n", name)
		return nil
	}

	procs, err := process.Processes()
	if err != nil {
		return fmt.Errorf("failed to list processes: %w", err)
	}

	killedCount := 0
	for _, p := range procs {
		pName, err := p.Name()
		if err != nil {
			continue // Skip processes we can't read
		}

		// Case-insensitive containment check
		if strings.Contains(strings.ToLower(pName), strings.ToLower(name)) {
			fmt.Printf("[Process] Found PID %d: %s. Initiating graceful termination...\n", p.Pid, pName)

			// 1. Terminate (SIGTERM)
			if err := p.Terminate(); err != nil {
				// If Terminate fails, it might be already dead or permission error.
				// We proceed to check if it's running.
			}

			// 2. Wait up to 5s
			timeout := time.After(5 * time.Second)
			ticker := time.NewTicker(500 * time.Millisecond)
			terminated := false

			for !terminated {
				select {
				case <-timeout:
					fmt.Printf("[Process] PID %d timed out. Force Killing (SIGKILL).\n", p.Pid)
					_ = p.Kill()
					terminated = true
				case <-ticker.C:
					running, err := p.IsRunning()
					if err != nil || !running {
						terminated = true
					}
				}
			}
			ticker.Stop()
			killedCount++
		}
	}

	if killedCount == 0 {
		fmt.Printf("[Process] No process found matching: %s\n", name)
	}

	return nil
}

func (m *Service) Start(path string) error {
	if m.DryRun {
		fmt.Printf("[Process] [DRY RUN] Start process: %s\n", path)
		return nil
	}

	var cmd *exec.Cmd
	// Simple heuristic for macOS App Bundles or Application Names
	if strings.HasSuffix(path, ".app") || !strings.Contains(path, "/") {
		// e.g. "Cursor.app" or "Cursor"
		cmd = exec.Command("open", "-a", path)
	} else {
		// Direct binary path
		cmd = exec.Command(path)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start process %s: %w", path, err)
	}

	fmt.Printf("[Process] Started %s (PID: %d)\n", path, cmd.Process.Pid)
	return nil
}
