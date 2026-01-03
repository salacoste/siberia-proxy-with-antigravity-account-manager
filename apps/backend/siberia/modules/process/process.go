package process

import (
	"fmt"
	"time"
)

// Manager handles external process control
type Manager interface {
	Kill(name string) error
	Start(path string) error
}

type MockManager struct {
	DryRun bool
}

func NewMockManager(dryRun bool) *MockManager {
	return &MockManager{DryRun: dryRun}
}

func (m *MockManager) Kill(name string) error {
	fmt.Printf("[ProcessManager] Killing process: %s (DryRun=%v)\n", name, m.DryRun)
	// Simulate work
	time.Sleep(500 * time.Millisecond)
	if m.DryRun {
		return nil
	}
	// TODO: Real implementation would use os/exec or syscall
	return nil
}

func (m *MockManager) Start(path string) error {
	fmt.Printf("[ProcessManager] Starting process: %s (DryRun=%v)\n", path, m.DryRun)
	time.Sleep(500 * time.Millisecond)
	if m.DryRun {
		return nil
	}
	// TODO: Real implementation
	return nil
}
