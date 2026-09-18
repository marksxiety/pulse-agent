package utils_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"pulse-agent/utils"
)

// withTempHome points the config directory at a throwaway directory so tests
// never touch the real ~/.pulse-agent/config.yaml.
//
// os.UserHomeDir reads USERPROFILE on Windows and HOME everywhere else, so both
// are set. The assertion is the point: without it, dropping either variable
// would silently start writing to the real home directory on Linux CI.
func withTempHome(t *testing.T) {
	t.Helper()

	dir := t.TempDir()
	t.Setenv("USERPROFILE", dir)
	t.Setenv("HOME", dir)

	if got, want := utils.ConfigDir(), filepath.Join(dir, ".pulse-agent"); got != want {
		t.Fatalf("test isolation failed: ConfigDir() = %q, want %q", got, want)
	}
}

func TestLoadConfig_MissingFileDefaultsTheme(t *testing.T) {
	withTempHome(t)

	if err := utils.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	got := utils.GetConfig()
	if got.Theme != utils.DefaultTheme {
		t.Errorf("Theme = %q, want %q", got.Theme, utils.DefaultTheme)
	}
	if !got.AlwaysOnTop {
		t.Error("AlwaysOnTop should default to true so the widget is pinned")
	}
	if got.WindowWidth != 0 || got.WindowHeight != 0 {
		t.Errorf("expected zero geometry, got %dx%d", got.WindowWidth, got.WindowHeight)
	}
}

func TestSaveConfig_RoundTripsAllFields(t *testing.T) {
	withTempHome(t)

	want := utils.Config{
		Theme:        "dracula",
		WindowX:      120,
		WindowY:      -40,
		WindowWidth:  360,
		WindowHeight: 220,
		AlwaysOnTop:  true,
	}
	if err := utils.SaveConfig(want); err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	if err := utils.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if got := utils.GetConfig(); got != want {
		t.Errorf("GetConfig() = %+v, want %+v", got, want)
	}
}

// TestSaveTheme_PreservesWindowFields guards the regression where the terminal
// UI's theme write truncated the file to a single line and destroyed whatever
// the desktop widget had persisted.
func TestSaveTheme_PreservesWindowFields(t *testing.T) {
	withTempHome(t)

	before := utils.Config{
		Theme:        "original",
		WindowX:      800,
		WindowY:      24,
		WindowWidth:  340,
		WindowHeight: 200,
		AlwaysOnTop:  true,
	}
	if err := utils.SaveConfig(before); err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	if err := utils.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if err := utils.SaveTheme("nord"); err != nil {
		t.Fatalf("SaveTheme() error = %v", err)
	}

	if err := utils.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	got := utils.GetConfig()
	if got.Theme != "nord" {
		t.Errorf("Theme = %q, want %q", got.Theme, "nord")
	}
	if got.WindowX != before.WindowX || got.WindowY != before.WindowY {
		t.Errorf("window position lost: got (%d,%d), want (%d,%d)",
			got.WindowX, got.WindowY, before.WindowX, before.WindowY)
	}
	if got.WindowWidth != before.WindowWidth || got.WindowHeight != before.WindowHeight {
		t.Errorf("window size lost: got %dx%d, want %dx%d",
			got.WindowWidth, got.WindowHeight, before.WindowWidth, before.WindowHeight)
	}
	if got.AlwaysOnTop != before.AlwaysOnTop {
		t.Errorf("AlwaysOnTop = %v, want %v", got.AlwaysOnTop, before.AlwaysOnTop)
	}
}

// TestLoadConfig_LegacyThemeOnlyFile covers a config written by an older build.
func TestLoadConfig_LegacyThemeOnlyFile(t *testing.T) {
	withTempHome(t)

	if err := os.MkdirAll(utils.ConfigDir(), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(utils.ConfigPath(), []byte("theme: gruvbox\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if err := utils.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	got := utils.GetConfig()
	if got.Theme != "gruvbox" {
		t.Errorf("Theme = %q, want %q", got.Theme, "gruvbox")
	}
	if got.WindowWidth != 0 || got.WindowHeight != 0 {
		t.Errorf("expected zero geometry, got %+v", got)
	}
	if !got.AlwaysOnTop {
		t.Error("a config omitting always_on_top should keep the pinned default")
	}
}

func TestLoadConfig_MalformedValuesFallBackToDefaults(t *testing.T) {
	withTempHome(t)

	body := strings.Join([]string{
		"theme: nord",
		"window_x: not-a-number",
		"window_width: 12pt",
		"always_on_top: yes-please",
		"",
	}, "\n")
	if err := os.MkdirAll(utils.ConfigDir(), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(utils.ConfigPath(), []byte(body), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if err := utils.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	got := utils.GetConfig()
	if got.Theme != "nord" {
		t.Errorf("Theme = %q, want %q", got.Theme, "nord")
	}
	if got.WindowX != 0 || got.WindowWidth != 0 {
		t.Errorf("malformed ints should fall back to 0, got x=%d w=%d", got.WindowX, got.WindowWidth)
	}
	if !got.AlwaysOnTop {
		t.Error("malformed bool should fall back to the pinned default")
	}
}

func TestLoadConfig_UnknownKeysIgnored(t *testing.T) {
	withTempHome(t)

	body := "theme: nord\nfuture_setting: whatever\n"
	if err := os.MkdirAll(utils.ConfigDir(), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(utils.ConfigPath(), []byte(body), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if err := utils.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if got := utils.GetConfig(); got.Theme != "nord" {
		t.Errorf("Theme = %q, want %q", got.Theme, "nord")
	}
}

func TestSaveConfig_EmptyThemeFallsBackToDefault(t *testing.T) {
	withTempHome(t)

	if err := utils.SaveConfig(utils.Config{}); err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}
	if err := utils.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if got := utils.GetConfig(); got.Theme != utils.DefaultTheme {
		t.Errorf("Theme = %q, want %q", got.Theme, utils.DefaultTheme)
	}
}

func TestSaveConfig_HandlesNegativeCoordinates(t *testing.T) {
	withTempHome(t)

	// A widget parked on a secondary monitor to the left of the primary one
	// legitimately has negative X.
	want := utils.Config{Theme: "nord", WindowX: -1920, WindowY: -100, WindowWidth: 320, WindowHeight: 180}
	if err := utils.SaveConfig(want); err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}
	if err := utils.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if got := utils.GetConfig(); got != want {
		t.Errorf("GetConfig() = %+v, want %+v", got, want)
	}
}

// TestLoadConfig_ReadsWidgetWrittenFile covers the shared-config contract in
// the direction that actually ships: the desktop widget writes geometry
// alongside the theme, and the terminal UI must still read the theme from it.
func TestLoadConfig_ReadsWidgetWrittenFile(t *testing.T) {
	withTempHome(t)

	body := strings.Join([]string{
		"theme: synthwave",
		"window_x: 1556",
		"window_y: 24",
		"window_width: 340",
		"window_height: 236",
		"always_on_top: true",
		"",
	}, "\n")

	if err := os.MkdirAll(utils.ConfigDir(), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(utils.ConfigPath(), []byte(body), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if err := utils.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	got := utils.GetConfig()
	if got.Theme != "synthwave" {
		t.Errorf("Theme = %q, want %q", got.Theme, "synthwave")
	}
	if got.WindowX != 1556 || got.WindowY != 24 {
		t.Errorf("window position not parsed: (%d,%d)", got.WindowX, got.WindowY)
	}
	if got.WindowWidth != 340 || got.WindowHeight != 236 {
		t.Errorf("window size not parsed: %dx%d", got.WindowWidth, got.WindowHeight)
	}
	if !got.AlwaysOnTop {
		t.Error("always_on_top not parsed")
	}
}

// TestSaveTheme_DoesNotClobberSettingsWrittenByAnotherProcess is the
// cross-process half of the clobbering bug. The in-memory copy is deliberately
// left stale, standing in for the other frontend writing the file behind our
// back after we last read it.
func TestSaveTheme_DoesNotClobberSettingsWrittenByAnotherProcess(t *testing.T) {
	withTempHome(t)

	// This process loads a config that has no geometry yet.
	if err := utils.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Meanwhile the widget saves its window position straight to disk.
	if err := os.MkdirAll(utils.ConfigDir(), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	body := strings.Join([]string{
		"theme: original",
		"window_x: 1556",
		"window_y: 24",
		"window_width: 340",
		"window_height: 236",
		"always_on_top: false",
		"",
	}, "\n")
	if err := os.WriteFile(utils.ConfigPath(), []byte(body), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// Now this process changes the theme, from its stale view.
	if err := utils.SaveTheme("dracula"); err != nil {
		t.Fatalf("SaveTheme() error = %v", err)
	}

	if err := utils.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	got := utils.GetConfig()
	if got.Theme != "dracula" {
		t.Errorf("Theme = %q, want %q", got.Theme, "dracula")
	}
	if got.WindowX != 1556 || got.WindowY != 24 {
		t.Errorf("window position clobbered: (%d,%d), want (1556,24)", got.WindowX, got.WindowY)
	}
	if got.WindowWidth != 340 || got.WindowHeight != 236 {
		t.Errorf("window size clobbered: %dx%d, want 340x236", got.WindowWidth, got.WindowHeight)
	}
	if got.AlwaysOnTop {
		t.Error("AlwaysOnTop clobbered: the persisted false should survive")
	}
}

func TestUpdateConfig_PreservesUnrelatedKeys(t *testing.T) {
	withTempHome(t)

	if err := utils.SaveConfig(utils.Config{
		Theme:        "nord",
		WindowX:      10,
		WindowY:      20,
		WindowWidth:  300,
		WindowHeight: 200,
		AlwaysOnTop:  true,
	}); err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	if err := utils.UpdateConfig(func(c *utils.Config) { c.AlwaysOnTop = false }); err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	if err := utils.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	got := utils.GetConfig()
	if got.AlwaysOnTop {
		t.Error("AlwaysOnTop = true, want false")
	}
	if got.Theme != "nord" {
		t.Errorf("Theme = %q, want %q", got.Theme, "nord")
	}
	if got.WindowX != 10 || got.WindowWidth != 300 {
		t.Errorf("geometry lost: %+v", got)
	}
}
