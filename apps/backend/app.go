package main

import (
	"context"
	"fmt"

	"github.com/salacoste/siberia/siberia/accounts"
	"github.com/salacoste/siberia/siberia/config"
	"github.com/salacoste/siberia/siberia/db"
	"github.com/salacoste/siberia/siberia/proxy"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx            context.Context
	config         *config.Manager
	proxyService   *proxy.Service
	database       *db.Database
	accountService *accounts.Service
}

// NewApp creates a new App application struct
func NewApp(cfg *config.Manager) *App {
	// Initialize DB
	database, err := db.Init(cfg.ConfigDir(), cfg.Config.MasterKey)
	if err != nil {
		fmt.Printf("Fatal: Failed to init DB: %v\n", err)
	}

	return &App{
		config:         cfg,
		proxyService:   proxy.NewService(&cfg.Config),
		database:       database,
		accountService: accounts.NewService(database),
	}
}

// ... existing methods ...

// ListAccounts returns all accounts (safe DTOs)
func (a *App) ListAccounts() ([]accounts.AccountDTO, error) {
	return a.accountService.ListAccounts()
}

// DeleteAccount deletes an account by ID
func (a *App) DeleteAccount(id uint) error {
	return a.accountService.DeleteAccount(id)
}

// CreateAccount creates a new account
func (a *App) CreateAccount(email, password, recovery, proxyGroup string) error {
	return a.accountService.CreateAccount(email, password, recovery, proxyGroup)
}

// ActivateAccount Orchestrates the switch to this account
func (a *App) ActivateAccount(id uint) error {
	return a.accountService.ActivateAccount(id)
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Start Proxy Service
	if err := a.proxyService.Start(ctx); err != nil {
		runtime.LogErrorf(a.ctx, "Failed to start proxy: %v", err)
	}

	runtime.LogInfo(a.ctx, fmt.Sprintf("Proxy started on port %d", a.config.Config.ProxyPort))
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	if a.proxyService != nil {
		a.proxyService.Stop(ctx)
	}
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
