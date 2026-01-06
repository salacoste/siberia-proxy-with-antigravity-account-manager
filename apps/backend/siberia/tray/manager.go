package tray

import (
	"context"
)

// Manager handles the system tray
// The struct fields are defined in platform_*.go to avoid importing systray on macOS
type TrayItem struct {
	// Opaque wrapper
	Item interface{}
}

type Manager struct {
	ctx           context.Context
	OnToggleShow  func()
	OnQuit        func()
	OnToggleProxy func()

	// Platform specific fields are embedded or handled via opaque state in impl files
	// For simplicity in this refactor, we'll use an 'internal' state interface or just methods
	// identifying the specific items.
	// Actually, Go doesn't allow splitting struct fields across files.
	// We will use an interface{} or specific pointer that is platform dependent,
	// OR better: we define 'trayState' in platform files and embed it.
	// Quota Menus
	QuotaMenus map[string]*TrayItem // Abstracted item

	trayState
}

func NewManager() *Manager {
	return &Manager{
		trayState: newTrayState(),
	}
}

// Setup is called from App.startup to pass the context if needed
func (m *Manager) Setup(ctx context.Context) {
	m.ctx = ctx
	m.setupImpl()
}

func (m *Manager) Run() {
	m.runImpl()
}

func (m *Manager) Stop() {
	m.stopImpl()
}

func (m *Manager) UpdateMenu(proxyActive bool) {
	m.updateMenuImpl(proxyActive)
}

func (m *Manager) UpdateQuota(stats map[string]int32) {
	m.updateQuotaImpl(stats)
}
