package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:pipeline-editor/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "Stack Agent",
		Width:  1440,
		Height: 900,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind:      []interface{}{app},
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHiddenInset(),
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
