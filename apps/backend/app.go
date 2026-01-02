package main

import (
	"context"
	"fmt"

	"github.com/salacoste/siberia/siberia/config"
)

// App struct
type App struct {
	ctx    context.Context
	config *config.Manager
}

// NewApp creates a new App application struct
func NewApp(cfg *config.Manager) *App {
	return &App{
		config: cfg,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetAppConfig returns the current configuration
func (a *App) GetAppConfig() config.AppConfig {
	return a.config.Get()
}

// UpdateAppConfig updates the configuration and persists it
func (a *App) UpdateAppConfig(newConfig config.AppConfig) error {
	return a.config.Update(newConfig)
}

// GetAppVersion returns the current app version
func (a *App) GetAppVersion() string {
	return "v0.1.0" // TODO: Load from compilation time
}
