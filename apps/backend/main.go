package main

import (
	"embed"

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

	// Create an instance of the app structure
	app := NewApp(cfgManager)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "siberia",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
