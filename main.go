package main

import (
	"embed"
	"helptools/internal/service"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Help Tools",
		Width:  800,
		Height: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},

		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 0, A: 0},
		OnStartup:        app.startup,
		Bind:             append([]interface{}{app}, service.Bind()...),
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
