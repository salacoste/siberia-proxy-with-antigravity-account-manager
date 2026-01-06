//go:build !darwin

package tray

import (
	"fmt"
	"time"

	"github.com/energye/systray"
)

type trayState struct {
	mShow   *systray.MenuItem
	mToggle *systray.MenuItem
	mQuit   *systray.MenuItem
}

func newTrayState() trayState {
	return trayState{}
}

func (m *Manager) runImpl() {
	// Start systray in a non-blocking way compatible with Wails (on Linux/Windows)
	systray.RunWithExternalLoop(m.onReady, m.onExit)
}

func (m *Manager) stopImpl() {
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
	// Note: We need to launch goroutines or handle safely if these block
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
	fmt.Println("Tray Exiting...")
}

func (m *Manager) setupImpl() {
	// Update initial state after a small delay to ensure Tray is ready
	go func() {
		time.Sleep(500 * time.Millisecond)
		m.UpdateMenu(false)
	}()
}

func (m *Manager) updateMenuImpl(proxyActive bool) {
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

func (m *Manager) updateQuotaImpl(stats map[string]int32) {
	// Initialize map if nil
	if m.QuotaMenus == nil {
		m.QuotaMenus = make(map[string]*TrayItem)
	}

	for model, pct := range stats {
		label := fmt.Sprintf("%s: %d%%", model, pct)

		if item, exists := m.QuotaMenus[model]; exists {
			// Update existing
			if menuItem, ok := item.Item.(*systray.MenuItem); ok {
				menuItem.SetTitle(label)
			}
		} else {
			// Create new
			newItem := systray.AddMenuItem(label, "Quota Usage")
			// Store opaque wrapper
			m.QuotaMenus[model] = &TrayItem{Item: newItem}
		}
	}
}
