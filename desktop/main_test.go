//go:build windows

package main

import (
	"testing"

	"pulse-agent/utils"
)

func TestSizeOrDefault(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		fallback int
		want     int
	}{
		{"zero uses the fallback", 0, 340, 340},
		{"negative uses the fallback", -10, 340, 340},
		{"positive wins", 500, 340, 500},
		{"one is a real size", 1, 340, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sizeOrDefault(tt.value, tt.fallback); got != tt.want {
				t.Errorf("sizeOrDefault(%d, %d) = %d, want %d", tt.value, tt.fallback, got, tt.want)
			}
		})
	}
}

// The widget's whole identity lives in these options, so pin them. Resizing in
// particular is easy to switch back on without realising it changes what the
// app is.
func TestWidgetOptions_WindowContract(t *testing.T) {
	opts := widgetOptions(utils.Config{Theme: utils.DefaultTheme, AlwaysOnTop: true}, NewApp())

	if !opts.DisableResize {
		t.Error("DisableResize = false, want the widget frame fixed")
	}
	if !opts.Frameless {
		t.Error("Frameless = false, want no OS title bar")
	}
	if !opts.AlwaysOnTop {
		t.Error("AlwaysOnTop = false, want the widget pinned by default")
	}
	if opts.BackgroundColour == nil || opts.BackgroundColour.A != 0 {
		t.Errorf("BackgroundColour = %+v, want fully transparent", opts.BackgroundColour)
	}
	if opts.Windows == nil {
		t.Fatal("Windows options are nil, so transparency would not be enabled")
	}
	if !opts.Windows.WebviewIsTransparent {
		t.Error("WebviewIsTransparent = false, the CSS card would sit on an opaque surface")
	}
	if !opts.Windows.WindowIsTranslucent {
		t.Error("WindowIsTranslucent = false, the corners would not be transparent")
	}
	if !opts.Windows.DisableFramelessWindowDecorations {
		t.Error("DisableFramelessWindowDecorations = false, DWM would draw its own corners")
	}
}

func TestWidgetOptions_AlwaysOnTopFollowsConfig(t *testing.T) {
	opts := widgetOptions(utils.Config{Theme: utils.DefaultTheme, AlwaysOnTop: false}, NewApp())

	if opts.AlwaysOnTop {
		t.Error("AlwaysOnTop = true, want the unpinned config value honoured")
	}
}

func TestWidgetOptions_UsesSavedSizeAndFallsBackToDefaults(t *testing.T) {
	tests := []struct {
		name       string
		cfg        utils.Config
		wantWidth  int
		wantHeight int
	}{
		{
			name:       "no saved size",
			cfg:        utils.Config{Theme: utils.DefaultTheme},
			wantWidth:  defaultWidth,
			wantHeight: defaultHeight,
		},
		{
			name:       "saved size wins",
			cfg:        utils.Config{Theme: utils.DefaultTheme, WindowWidth: 500, WindowHeight: 400},
			wantWidth:  500,
			wantHeight: 400,
		},
		{
			name:       "negative saved size falls back",
			cfg:        utils.Config{Theme: utils.DefaultTheme, WindowWidth: -1, WindowHeight: -1},
			wantWidth:  defaultWidth,
			wantHeight: defaultHeight,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := widgetOptions(tt.cfg, NewApp())

			if opts.Width != tt.wantWidth {
				t.Errorf("Width = %d, want %d", opts.Width, tt.wantWidth)
			}
			if opts.Height != tt.wantHeight {
				t.Errorf("Height = %d, want %d", opts.Height, tt.wantHeight)
			}
		})
	}
}

// Without the lifecycle hooks and the binding the widget renders nothing and
// cannot be driven from the webview, so assert they are all wired.
func TestWidgetOptions_WiresLifecycleAndBindings(t *testing.T) {
	app := NewApp()
	opts := widgetOptions(utils.Config{Theme: utils.DefaultTheme}, app)

	if opts.OnStartup == nil || opts.OnDomReady == nil ||
		opts.OnBeforeClose == nil || opts.OnShutdown == nil {
		t.Error("one or more lifecycle hooks are nil")
	}
	if len(opts.Bind) != 1 {
		t.Fatalf("Bind has %d entries, want the App bound", len(opts.Bind))
	}
	if opts.Bind[0] != app {
		t.Error("Bind does not carry the App instance, so the webview could not call it")
	}
	if opts.AssetServer == nil {
		t.Error("AssetServer is nil, the frontend would not be served")
	}
}
