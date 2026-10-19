package main

import (
	"app-dev/internal/driver"
	"app-dev/internal/service"
	"app-dev/internal/store"
	"embed"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/out
var assets embed.FS

func main() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	storePath := filepath.Join(configDir, "dbstudio-lite", "connections.json")
	connStore := store.NewFileConnectionStore(storePath)

	manager := service.NewConnectionManager(driver.NewMultiDriverConnector())
	app := NewApp(manager, connStore)

	err = wails.Run(&options.App{
		Title:  "DBStudio Lite",
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