//go:build windows

// Package main is the Pulse Agent desktop widget. It is Windows-only: the
// build constraint keeps it out of the module graph on other platforms, which
// is what lets the CI lint and test jobs run on Linux without GTK headers.
package main

import (
	"embed"
	"log"

	"pulse-agent/utils"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

const (
	defaultWidth  = 340
	defaultHeight = 236

	// widgetMargin keeps the first-run placement clear of the screen edge.
	widgetMargin = 24
)

func main() {
	if err := utils.LoadConfig(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := wails.Run(widgetOptions(utils.GetConfig(), NewApp())); err != nil {
		log.Fatalf("wails.Run: %v", err)
	}
}

// widgetOptions builds the window contract for the widget.
//
// It is separate from main purely so the contract can be asserted in tests —
// these settings are what make the app a widget rather than a window, and they
// are easy to break by accident.
func widgetOptions(cfg utils.Config, app *App) *options.App {
	return &options.App{
		Title: "Pulse Agent",

		// The frame is fixed: this is a widget to park and leave alone, and a
		// fixed size also stops the three cards from reflowing mid-glance. A
		// saved size is still honoured, it just cannot be dragged afterwards.
		Width:         sizeOrDefault(cfg.WindowWidth, defaultWidth),
		Height:        sizeOrDefault(cfg.WindowHeight, defaultHeight),
		DisableResize: true,
		Frameless:     true,
		AlwaysOnTop:   cfg.AlwaysOnTop,

		// A fully transparent window is what lets the CSS card provide its
		// own rounded corners over whatever is behind the widget.
		BackgroundColour: options.NewRGBA(0, 0, 0, 0),

		AssetServer: &assetserver.Options{Assets: assets},

		OnStartup:     app.startup,
		OnDomReady:    app.domReady,
		OnBeforeClose: app.beforeClose,
		OnShutdown:    app.shutdown,

		Bind: []interface{}{app},

		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "pulse-agent-widget",
			OnSecondInstanceLaunch: func(options.SecondInstanceData) {
				app.focus()
			},
		},

		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			DisableWindowIcon:    true,

			// Let the CSS card own the corners and shadow. Leaving the
			// frameless decorations on would fight the translucent background
			// with a DWM-drawn rounded rect.
			DisableFramelessWindowDecorations: true,

			// A widget has no use for pinch-zoom, and it is easy to trigger
			// accidentally while dragging the window.
			DisablePinchZoom: true,
		},
	}
}

func sizeOrDefault(value, fallback int) int {
	if value <= 0 {
		return fallback
	}
	return value
}
