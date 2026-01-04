package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/salacoste/siberia/siberia/accounts"
	"github.com/salacoste/siberia/siberia/ca"
	"github.com/salacoste/siberia/siberia/config"
	"github.com/salacoste/siberia/siberia/db"
	"github.com/salacoste/siberia/siberia/logger"
	"github.com/salacoste/siberia/siberia/modules/injection"
	"github.com/salacoste/siberia/siberia/modules/process"
	"github.com/salacoste/siberia/siberia/modules/sync"
	"github.com/salacoste/siberia/siberia/modules/vault"
	"github.com/salacoste/siberia/siberia/proxy"
	"github.com/salacoste/siberia/siberia/share"
	"github.com/salacoste/siberia/siberia/updater"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// RegisterTypes is a dummy method to force Wails to generate bindings for these types
// They are used in Events but not in other methods, so Wails v2 skips them otherwise.
func (a *App) RegisterTypes() (proxy.PendingRequest, proxy.WebSocketFrame) {
	return proxy.PendingRequest{}, proxy.WebSocketFrame{}
}

// App struct
type App struct {
	ctx              context.Context
	config           *config.Manager
	proxyService     *proxy.Service
	database         *db.Database
	processService   *process.Service
	injectionService *injection.Service
	updaterService   *updater.Service
	accountService   *accounts.Service
	caService        *ca.Service
	shareService     *share.Service
	syncManager      *sync.Manager
}

// NewApp creates a new App application struct
func NewApp(cfg *config.Manager) *App {
	// Initialize DB
	database, err := db.Init(cfg.ConfigDir(), cfg.Config.MasterKey)
	if err != nil {
		fmt.Printf("Fatal: Failed to init DB: %v\n", err)
	}

	// Initialize Logger
	logger.InitAccessLogger(cfg.ConfigDir())

	// Initialize CA Service
	caSvc := ca.NewService(&cfg.Config)
	if err := caSvc.EnsureCA(); err != nil {
		fmt.Printf("Fatal: Failed to ensure CA: %v\n", err)
	}

	// Initialize Share Service (MinIO / S3)
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	if minioEndpoint == "" {
		minioEndpoint = "localhost:9000" // Default for local docker
	}
	minioAccess := os.Getenv("MINIO_ROOT_USER") // or MINIO_ACCESS_KEY
	if minioAccess == "" {
		minioAccess = "minioadmin"
	}
	minioSecret := os.Getenv("MINIO_ROOT_PASSWORD") // or MINIO_SECRET_KEY
	if minioSecret == "" {
		minioSecret = "minioadmin"
	}
	minioBucket := os.Getenv("MINIO_BUCKET")
	if minioBucket == "" {
		minioBucket = "siberia-shares"
	}

	var shareProvider share.StorageProvider
	// Try connecting to MinIO
	fmt.Printf("Initializing Share Service with MinIO at %s...\n", minioEndpoint)
	s3Provider, err := share.NewS3Provider(minioEndpoint, minioAccess, minioSecret, minioBucket, false)
	if err == nil {
		shareProvider = s3Provider
	} else {
		fmt.Printf("Warning: Failed to init MinIO: %v. Falling back to Mock Share Provider.\n", err)
		shareProvider = &share.MockProvider{}
	}

	shareSvc := share.NewService(shareProvider)

	// Initialize Sync Manager (Supabase)
	// Try loading from .env (local dev) or environment
	_ = godotenv.Load() // Ignore error for prod/built binaries

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	var syncMgr *sync.Manager
	if supabaseURL != "" && supabaseKey != "" {
		fmt.Println("Initializing Real Supabase Backend...")
		// For MVP, we use a hardcoded user ID since we don't have a full auth flow yet
		client := sync.NewSupabaseClient(supabaseURL, supabaseKey, "demo-user-001")
		syncMgr = sync.NewManager(client)
	} else {
		fmt.Println("Warning: SUPABASE credentials not found. Falling back to in-memory Mock.")
		// Fallback to Mock if no credentials (so app doesn't crash in dev without .env)
		// We adapt MockServer (which is an HTTP server) to the CloudProvider?
		// Actually, MockServer was a full HTTP server.
		// Since we refactored Manager to use an interface, we can't easily plug the old MockServer (URL-based)
		// without an adapter.
		// For now, let's just warn and leave syncMgr nil, or use a dummy provider.
		// Re-using the Mock logic efficiently requires rewriting MockServer to implement CloudProvider
		// instead of being an HTTP Handler.
		// Let's implement a simple InMemoryProvider here for fallback.
		syncMgr = sync.NewManager(NewInMemoryProvider())
	}

	return &App{
		config:           cfg,
		proxyService:     proxy.NewService(&cfg.Config, caSvc),
		database:         database,
		accountService:   accounts.NewService(database),
		processService:   process.NewService(),
		injectionService: injection.NewService(),
		updaterService:   updater.NewService("v1.0.1"),
		caService:        caSvc,
		shareService:     shareSvc,
		syncManager:      syncMgr,
	}
}

// InMemoryProvider is a fallback for when keys are missing
type InMemoryProvider struct {
	data string
}

func NewInMemoryProvider() *InMemoryProvider {
	return &InMemoryProvider{}
}

func (i *InMemoryProvider) Push(data string) error {
	i.data = data
	return nil
}

func (i *InMemoryProvider) Pull() (string, error) {
	return i.data, nil
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

// InstallCert installs the root CA into the OS trust store
func (a *App) InstallCert() error {
	return a.caService.InstallCert()
}

// CheckCertTrust checks if the root CA is trusted
func (a *App) CheckCertTrust() bool {
	return a.caService.CheckTrust()
}

// CheckForUpdates checks for updates
func (a *App) CheckForUpdates() updater.UpdateInfo {
	return a.updaterService.CheckForUpdates()
}

// GetAppVersion returns the current app version
func (a *App) GetVersion() string {
	return "v1.0.1"
}

// === Breakpoint API ===

// AddBreakpointRule adds a rule to pause requests
func (a *App) AddBreakpointRule(pattern, method string) proxy.BreakpointRule {
	return a.proxyService.BreakpointManager.AddRule(pattern, method)
}

// DeleteBreakpointRule removes a rule
func (a *App) DeleteBreakpointRule(id string) {
	a.proxyService.BreakpointManager.DeleteRule(id)
}

// ResumeRequest resumes a paused request
func (a *App) ResumeRequest(id string, mod proxy.ModifiedRequest) bool {
	return a.proxyService.BreakpointManager.ResumeRequest(id, mod)
}

// UploadSession exports a request as HAR and uploads it
func (a *App) UploadSession(event proxy.ProxyRequestEvent) (string, error) {
	return a.shareService.UploadSession(event)
}

// === Sync & Vault API ===

// SyncPush triggers a push of the current profile
// For MVP: We hardcode a dummy profileID/password since we don't have the full UI for it yet
// In real impl, these come from the Vault State.
func (a *App) SyncPush(password string) error {
	// 1. Encrypt Data (Mocking data for now)
	vault := vault.NewVault()
	data := []byte("{\"mock\": \"profile data\"}")
	blob, err := vault.Encrypt(data, password)
	if err != nil {
		return err
	}

	// 2. Push
	return a.syncManager.Push("default-profile", blob)
}

// SyncPull triggers a pull
func (a *App) SyncPull(password string) (string, error) {
	payload, err := a.syncManager.Pull("default-profile")
	if err != nil {
		return "", err
	}
	if payload == nil {
		return "", fmt.Errorf("no remote data")
	}

	// Decrypt
	vault := vault.NewVault()
	data, err := vault.Decrypt(payload.Blob, password)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
