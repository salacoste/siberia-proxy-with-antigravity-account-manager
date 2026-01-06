package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/salacoste/siberia/siberia/accounts"
	"github.com/salacoste/siberia/siberia/analytics"
	"github.com/salacoste/siberia/siberia/ca"
	"github.com/salacoste/siberia/siberia/cloud"
	"github.com/salacoste/siberia/siberia/config"
	"github.com/salacoste/siberia/siberia/db"
	"github.com/salacoste/siberia/siberia/ide"
	"github.com/salacoste/siberia/siberia/logger"
	"github.com/salacoste/siberia/siberia/mcp"
	"github.com/salacoste/siberia/siberia/migration"
	"github.com/salacoste/siberia/siberia/modules/injection"
	"github.com/salacoste/siberia/siberia/modules/process"
	"github.com/salacoste/siberia/siberia/proxy"
	"github.com/salacoste/siberia/siberia/proxy/middleware"
	"github.com/salacoste/siberia/siberia/quota"
	"github.com/salacoste/siberia/siberia/share"
	"github.com/salacoste/siberia/siberia/tray"
	"github.com/salacoste/siberia/siberia/types"
	"github.com/salacoste/siberia/siberia/updater"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// RegisterTypes is a dummy method to force Wails to generate bindings for these types
// They are used in Events but not in other methods, so Wails v2 skips them otherwise.
func (a *App) RegisterTypes() (proxy.PendingRequest, proxy.WebSocketFrame, types.ProxyRequestEvent, middleware.MapLocalRule) {
	return proxy.PendingRequest{}, proxy.WebSocketFrame{}, types.ProxyRequestEvent{}, middleware.MapLocalRule{}
}

// App struct
type App struct {
	ctx              context.Context
	config           *config.Manager
	proxyService     *proxy.Service
	database         *db.Database
	processService   *process.Service
	injectionService *injection.Service
	updaterService   *updater.UpdateService
	accountService   *accounts.Service
	caService        *ca.Service
	shareService     *share.Service
	cloudService     *cloud.Service
	AnalyticsService *analytics.AnalyticsService
	updateService    *updater.UpdateService
	trayManager      *tray.Manager
	migrationService *migration.Service
	mcpServer        *mcp.McpServer
}

// Version is the current application version
const Version = "v1.2.0"

// NewApp creates a new App application struct
func NewApp(cfg *config.Manager) *App {
	// Initialize DB
	database, err := db.Init(cfg.ConfigDir(), cfg.Config.MasterKey)
	if err != nil {
		log.Printf("Fatal: Failed to init DB: %v\n", err)
	}

	// Initialize Logger
	logger.InitAccessLogger(cfg.ConfigDir())

	// Initialize CA Service
	caSvc := ca.NewService(&cfg.Config)
	if err := caSvc.EnsureCA(); err != nil {
		log.Printf("Fatal: Failed to ensure CA: %v\n", err)
	}

	// Initialize Analytics
	analyticsEngine := analytics.NewAnalyticsEngine()
	analyticsSvc := analytics.NewAnalyticsService(analyticsEngine)

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
	log.Printf("Initializing Share Service with MinIO at %s...\n", minioEndpoint)
	s3Provider, err := share.NewS3Provider(minioEndpoint, minioAccess, minioSecret, minioBucket, false)
	if err == nil {
		shareProvider = s3Provider
	} else {
		log.Printf("Warning: Failed to init MinIO: %v. Falling back to Mock Share Provider.\n", err)
		shareProvider = &share.MockProvider{}
	}

	shareSvc := share.NewService(shareProvider)

	// Initialize Cloud Service
	// We create a new logger instance for Cloud
	cloudLog := logger.New("CLOUD")
	cloudSvc := cloud.NewService(cfg, cloudLog)

	// Initialize Account Service
	accountSvc := accounts.NewService(database, cfg)

	// Initialize Migration Service
	migrationSvc := migration.NewService(accountSvc)

	// Initialize Update Service
	// TODO: Get repo from config or hardcode for now
	updateSvc := updater.NewUpdateService(Version, "salacoste/siberia-proxy-with-antigravity-account-manager")

	// Initialize Tray Manager
	trayMgr := tray.NewManager()

	// Initialize MCP Server
	mcpSvc := mcp.NewServer(cfg)

	return &App{
		config:           cfg,
		proxyService:     proxy.NewService(&cfg.Config, caSvc, analyticsEngine, accountSvc),
		database:         database,
		accountService:   accountSvc,
		processService:   process.NewService(),
		injectionService: injection.NewService(),
		updaterService:   updateSvc,
		caService:        caSvc,
		shareService:     shareSvc,
		cloudService:     cloudSvc,
		AnalyticsService: analyticsSvc,
		trayManager:      trayMgr,
		migrationService: migrationSvc,
		mcpServer:        mcpSvc,
	}
}

// ... existing code ...

// === Migration API ===

// CheckLegacyData checks if there are legacy accounts to import
func (a *App) CheckLegacyData() (migration.LegacyStatus, error) {
	if a.migrationService == nil {
		return migration.LegacyStatus{}, fmt.Errorf("migration service not initialized")
	}
	return a.migrationService.CheckLegacyData()
}

// PerformLegacyImport executes the import
func (a *App) PerformLegacyImport() (int, error) {
	if a.migrationService == nil {
		return 0, fmt.Errorf("migration service not initialized")
	}
	return a.migrationService.PerformImport()
}

// OpenProjectInIDE opens the current working directory in the configured IDE
func (a *App) OpenProjectInIDE() error {
	// 1. Get Configured IDE
	target := a.config.Get().TargetIDE
	if target == "" {
		target = "cursor" // Default to Cursor if not set? Or vscode.
	}

	// 2. Get Current Working Directory (Project Root)
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// 3. Open
	opener := ide.NewOpener()
	return opener.Open(target, cwd, 0)
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

func (i *InMemoryProvider) SignUp(email, password string) error {
	return nil
}

func (i *InMemoryProvider) SignIn(email, password string) error {
	return nil
}

func (i *InMemoryProvider) GetUser() string {
	return "mock-user-id"
}

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

	fmt.Println("Siberia Backend Started (v" + Version + ")")

	// Ensure CA
	if err := a.caService.EnsureCA(); err != nil { // Corrected to use a.caService
		log.Printf("Warning: Failed to ensure CA: %v", err)
	}

	// Initialize DB (already done in NewApp, but if there's a re-init logic, it would go here)
	// For now, assuming DB is ready from NewApp.

	// Initialize MinIO (already done in NewApp, but if there's a re-init logic, it would go here)
	// For now, assuming S3 is ready from NewApp.

	// Setup Tray
	a.trayManager.Setup(ctx)

	// Auto-Import Legacy Data
	if count, err := a.migrationService.PerformImport(); err == nil && count > 0 {
		runtime.LogInfo(a.ctx, fmt.Sprintf("Imported %d legacy accounts", count))
	} else if err != nil {
		runtime.LogErrorf(a.ctx, "Legacy import failed: %v", err)
	}

	// Listen for Proxy Events to update Tray
	// Inject Context into Quota Service for event emission
	if a.accountService != nil && a.accountService.GetQuotaService() != nil {
		a.accountService.GetQuotaService().SetContext(ctx)
	}

	runtime.EventsOn(ctx, types.EventQuotaUpdate, func(optionalData ...interface{}) {
		if len(optionalData) > 0 {
			// Wails events come as []interface{}
			// We need to marshal/unmarshal or assume structure?
			// Actually, runtime.EventsEmit sends data.
			// Let's assert map[string]interface{} or just handle the struct.
			// Wait, EventQuotaUpdate emits *quota.Stats.
			// Wails JS->Go might be different, but Go->Go is direct pass?
			// No, Wails Event Bus serializes to JSON usually if crossing boundaries.
			// But creating a Go-only listener might be cleaner with a direct hook.
			// For now, let's assume we can cast if it's in-process.
			// Actually, wails `EventsOn` is for JS-compat.
			// Let's cast safely.
			if stats, ok := optionalData[0].(*quota.Stats); ok {
				a.trayManager.UpdateQuota(stats.Models)
			} else {
				// Fallback if serialization happened (unlikely for Go-Go in same process but possible)
				// Log warning
				// runtime.LogWarning(ctx, "Received invalid quota update event data")
			}
		}
	})

	// Start Proxy Service
	if err := a.proxyService.Start(ctx); err != nil {
		runtime.LogErrorf(a.ctx, "Failed to start proxy: %v", err)
	}

	runtime.LogInfo(a.ctx, fmt.Sprintf("Proxy started on port %d", a.config.Get().ProxyPort))

	// Start MCP Server if enabled
	mcpCfg := a.config.Get().ZaiMcpConfig
	if mcpCfg.Enabled {
		// If reusing main port, we might mount handler differently,
		// but since Wails controls main http server, we likely need a separate goroutine if dedicated port,
		// OR we rely on Wails binding (but that's not std http).
		// For MVP: Start on dedicated port as planned in config.
		port := mcpCfg.Port
		if port == 0 {
			port = 6200
		}
		go func() {
			runtime.LogInfo(a.ctx, fmt.Sprintf("Starting MCP Server on :%d", port))
			if err := http.ListenAndServe(fmt.Sprintf(":%d", port), a.mcpServer); err != nil {
				runtime.LogErrorf(a.ctx, "MCP Server failed: %v", err)
			}
		}()
	}
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Save Window Size
	if width, height := runtime.WindowGetSize(ctx); width > 0 && height > 0 {
		config := a.config.Get()
		config.WindowWidth = width
		config.WindowHeight = height
		// We use Update which saves to disk
		_ = a.config.Update(config)
	}

	if a.proxyService != nil {
		a.proxyService.Stop(ctx)
	}
}

// SaveWindowSize persists the window dimensions (called from frontend on resize)
// We ignore the passed arguments because frontend sometimes reports 0 on some platforms/webviews.
// Instead we ask the Wails runtime for the authoritative window size.
func (a *App) SaveWindowSize(frontendWidth, frontendHeight int) {
	width, height := runtime.WindowGetSize(a.ctx)

	if width > 0 && height > 0 {
		config := a.config.Get()
		config.WindowWidth = width
		config.WindowHeight = height

		// Use Update to save to disk
		_ = a.config.Update(config)
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
// This implementation delegates to the service
func (a *App) CheckForUpdates() (*updater.UpdateInfo, error) {
	if a.updaterService == nil {
		return nil, fmt.Errorf("update service not initialized")
	}
	return a.updaterService.CheckForUpdates()
}

// GetVersion returns the current app version
func (a *App) GetVersion() string {
	return Version
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

// GetBreakpointRules returns all active breakpoint rules
func (a *App) GetBreakpointRules() []proxy.BreakpointRule {
	return a.proxyService.BreakpointManager.GetRules()
}

// ResumeRequest resumes a paused request
func (a *App) ResumeRequest(id string, mod proxy.ModifiedRequest) bool {
	return a.proxyService.BreakpointManager.ResumeRequest(id, mod)
}

// === Map Local API ===

// AddMapLocalRule adds a rule to map a URL to a local file
func (a *App) AddMapLocalRule(rule middleware.MapLocalRule) error {
	if a.proxyService.MapLocal == nil {
		return fmt.Errorf("map local middleware not initialized")
	}
	return a.proxyService.MapLocal.AddRule(rule)
}

// DeleteMapLocalRule removes a rule
func (a *App) DeleteMapLocalRule(id string) {
	if a.proxyService.MapLocal != nil {
		a.proxyService.MapLocal.RemoveRule(id)
	}
}

// GetMapLocalRules returns all active rules
func (a *App) GetMapLocalRules() []middleware.MapLocalRule {
	if a.proxyService.MapLocal == nil {
		return []middleware.MapLocalRule{}
	}
	// Extract rules from sync.Map
	var rules []middleware.MapLocalRule
	a.proxyService.MapLocal.RangeRules(func(r middleware.MapLocalRule) {
		rules = append(rules, r)
	})
	return rules
}

// UploadSession exports a request as HAR and uploads it
func (a *App) UploadSession(event types.ProxyRequestEvent) (string, error) {
	return a.shareService.UploadSession(event)
}

// SyncPush triggers a push of the current profile
func (a *App) CloudSync() error {
	return a.cloudService.Sync(a.ctx)
}

// SyncSignUp registers a user
func (a *App) CloudSignUp(email, password string) error {
	_, err := a.cloudService.Client().SignUp(email, password)
	// Login immediately after signup?
	if err == nil {
		return a.cloudService.Login(a.ctx, email, password)
	}
	return err
}

// SyncSignIn logs in a user
func (a *App) CloudLogin(email, password string) error {
	return a.cloudService.Login(a.ctx, email, password)
}

// CloudLogout
func (a *App) CloudLogout() error {
	return a.cloudService.Logout(a.ctx)
}

// CloudGetStatus
func (a *App) CloudGetStatus() map[string]interface{} {
	cfg := a.config.Get()
	return map[string]interface{}{
		"enabled":   cfg.CloudEnabled,
		"email":     cfg.CloudEmail,
		"last_sync": cfg.CloudLastSync,
	}
}

// === Scripting API ===

// UpdateScript updates the traffic interception script
func (a *App) UpdateScript(script string, active bool) error {
	if a.proxyService.ScriptEngine == nil {
		return fmt.Errorf("script engine not initialized")
	}
	a.proxyService.ScriptEngine.UpdateScript(script, active)
	return nil
}

type ScriptState struct {
	Code   string `json:"code"`
	Active bool   `json:"active"`
}

// GetScript returns the current script and active state
func (a *App) GetScript() (ScriptState, error) {
	if a.proxyService.ScriptEngine == nil {
		return ScriptState{}, fmt.Errorf("script engine not initialized")
	}
	return ScriptState{
		Code:   a.proxyService.ScriptEngine.Script,
		Active: a.proxyService.ScriptEngine.Active,
	}, nil
}
