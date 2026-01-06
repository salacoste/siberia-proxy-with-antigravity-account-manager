package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/salacoste/siberia/siberia/config"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Initialize Config
	cfgManager, err := config.NewManager()
	if err != nil {
		println("Error initializing config:", err.Error())
		return
	}
	if err := cfgManager.Load(); err != nil {
		println("Error loading config:", err.Error())
	}

	// CLI Commands
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "set-env":
			port := cfgManager.Config.ProxyPort
			if port == 0 {
				port = 8888
			}
			fmt.Printf("export HTTP_PROXY=http://127.0.0.1:%d HTTPS_PROXY=http://127.0.0.1:%d ALL_PROXY=http://127.0.0.1:%d\n", port, port, port)
			os.Exit(0)
		case "help":
			fmt.Println("Siberia Proxy Helper")
			fmt.Println("Usage:")
			fmt.Println("  siberia              Start the GUI application")
			fmt.Println("  siberia set-env      Output shell export commands")
			os.Exit(0)
		}
	}

	// Create an instance of the app structure
	app := NewApp(cfgManager)

	// Start System Tray (Non-blocking external loop)
	// We access it via the app instance where it was initialized
	app.trayManager.Run()

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "siberia",
		Width:  cfgManager.Config.WindowWidth,
		Height: cfgManager.Config.WindowHeight,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
			app.AnalyticsService,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
