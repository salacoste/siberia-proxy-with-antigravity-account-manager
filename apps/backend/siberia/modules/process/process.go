package process

import (
	"fmt"
	"os/exec"
	"strings"

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

		// Simple case-insensitive containment check
		// e.g., "Code Helper" matches "Code Helper (GPU)"
		if strings.Contains(strings.ToLower(pName), strings.ToLower(name)) {
			fmt.Printf("[Process] Killing PID %d: %s\n", p.Pid, pName)
			if err := p.Terminate(); err != nil {
				// Try Kill (Force) if Terminate fails
				_ = p.Kill()
			}
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

	// Start detached process
	cmd := exec.Command(path)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start process %s: %w", path, err)
	}

	fmt.Printf("[Process] Started %s (PID: %d)\n", path, cmd.Process.Pid)
	return nil
}
