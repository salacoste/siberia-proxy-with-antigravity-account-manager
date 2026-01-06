//go:build darwin

package tray

// No-op implementation for macOS to avoid Wails conflict
type trayState struct{}

func newTrayState() trayState {
	return trayState{}
}

func (m *Manager) runImpl() {
	// No-op on macOS
}

func (m *Manager) stopImpl() {
	// No-op
}

func (m *Manager) setupImpl() {
	// No-op
}

func (m *Manager) updateMenuImpl(proxyActive bool) {
	// No-op
}

func (m *Manager) updateQuotaImpl(stats map[string]int32) {
	// No-op
}
