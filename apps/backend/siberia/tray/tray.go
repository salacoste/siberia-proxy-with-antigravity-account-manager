package tray

import (
	"context"
	"fmt"
	"time"

	"github.com/energye/systray"
)

type Manager struct {
	ctx           context.Context
	OnToggleShow  func()
	OnQuit        func()
	OnToggleProxy func()

	mShow   *systray.MenuItem
	mToggle *systray.MenuItem
	mQuit   *systray.MenuItem
}

func NewManager() *Manager {
	return &Manager{}
}

// RunWithExternalLoop starts the tray logic without blocking the main thread
// It must be called before the main application loop
func (m *Manager) Run() {
	// Start systray in a non-blocking way compatible with Wails
	systray.RunWithExternalLoop(m.onReady, m.onExit)
}

func (m *Manager) Stop() {
	systray.Quit()
}

func (m *Manager) onReady() {
	systray.SetIcon(IconInactive)
	systray.SetTitle("Siberia")
	systray.SetTooltip("Siberia Proxy: Inactive")

	m.mShow = systray.AddMenuItem("Show Siberia", "Open the main window")
	m.mToggle = systray.AddMenuItem("Start Proxy", "Enable/Disable the proxy server")
	systray.AddSeparator()
	m.mQuit = systray.AddMenuItem("Quit", "Quit the application")

	// Handle menu clicks via callbacks (energye/systray specific)
	m.mShow.Click(func() {
		if m.OnToggleShow != nil {
			m.OnToggleShow()
		}
	})

	m.mToggle.Click(func() {
		if m.OnToggleProxy != nil {
			m.OnToggleProxy()
		}
	})

	m.mQuit.Click(func() {
		if m.OnQuit != nil {
			m.OnQuit()
		}
	})
}

func (m *Manager) onExit() {
	// Cleanup if needed
	fmt.Println("Tray Exiting...")
}

func (m *Manager) UpdateMenu(proxyActive bool) {
	// Needs to be safe to call from any goroutine, systray handles this usually
	// But mostly we should check if mToggle is initialized
	if m.mToggle == nil {
		return
	}

	if proxyActive {
		systray.SetIcon(IconActive)
		systray.SetTooltip("Siberia Proxy: Active")
		m.mToggle.SetTitle("Stop Proxy")
		m.mToggle.Check()
	} else {
		systray.SetIcon(IconInactive)
		systray.SetTooltip("Siberia Proxy: Inactive")
		m.mToggle.SetTitle("Start Proxy")
		m.mToggle.Uncheck()
	}
}

// Setup is called from App.startup to pass the context if needed
func (m *Manager) Setup(ctx context.Context) {
	m.ctx = ctx
	// Update initial state after a small delay to ensure Tray is ready
	go func() {
		time.Sleep(500 * time.Millisecond)
		m.UpdateMenu(false)
	}()
}
