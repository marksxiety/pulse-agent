//go:build windows

package main

import (
	"context"
	"testing"
	"time"

	"pulse-agent/models"
	"pulse-agent/theme"
	"pulse-agent/types"
	"pulse-agent/utils"
)

// testInterval is deliberately short for the tests that exercise the sampling
// cadence through apply; the cadence itself is pinned deterministically in
// tests/models.
const testInterval = 5 * time.Second

// newTestApp builds an App with just enough wiring for the pure state machine.
// The Wails runtime, collectors, and webview are all absent on purpose: nothing
// under test here touches them.
func newTestApp(interval time.Duration) *App {
	return &App{
		sampler:   models.NewSampler(interval),
		startedAt: time.Now(),
	}
}

// withTempHome points the shared config at a throwaway directory and reloads it,
// so widget tests never touch the real ~/.pulse-agent/config.yaml.
//
// os.UserHomeDir reads USERPROFILE on Windows and HOME elsewhere; both are set
// so this keeps working if the suite is ever run on another platform.
func withTempHome(t *testing.T) {
	t.Helper()

	dir := t.TempDir()
	t.Setenv("USERPROFILE", dir)
	t.Setenv("HOME", dir)

	if err := utils.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
}

func cpuMetric(percent float64, cores int) models.Metric {
	return models.Metric{Source: types.CPU, Data: models.CPUPayload{Percentage: percent, CoreCount: cores}}
}

func memMetric(used, total uint64) models.Metric {
	return models.Metric{Source: types.MEM, Data: models.MemoryPayload{Used: used, Total: total}}
}

func diskMetric(used, total uint64) models.Metric {
	return models.Metric{Source: types.DISK, Data: models.DiskPayload{Used: used, Total: total}}
}

func TestApplyDispatchesMetricsToTheirFields(t *testing.T) {
	app := newTestApp(testInterval)

	app.apply(cpuMetric(30, 8))
	if app.cpu.Percentage != 30 || app.cpu.CoreCount != 8 {
		t.Errorf("cpu = %+v, want 30%% across 8 cores", app.cpu)
	}
	if !app.cpuReady {
		t.Error("cpuReady = false after a CPU reading")
	}

	app.apply(memMetric(50, 100))
	if app.memory.Used != 50 || app.memory.Total != 100 {
		t.Errorf("memory = %+v, want 50/100", app.memory)
	}
	if !app.memReady {
		t.Error("memReady = false after a memory reading")
	}

	app.apply(diskMetric(25, 100))
	if app.disk.Used != 25 || app.disk.Total != 100 {
		t.Errorf("disk = %+v, want 25/100", app.disk)
	}
	if !app.diskReady {
		t.Error("diskReady = false after a disk reading")
	}
}

// The first sample lands as soon as any collector reports, so the collectors
// that have not spoken yet must not contribute a zero to the trend.
func TestApplyDoesNotSampleCollectorsThatHaveNotReported(t *testing.T) {
	app := newTestApp(testInterval)

	app.apply(cpuMetric(30, 8))

	if app.sampler.CPU.Len() != 1 {
		t.Errorf("CPU history has %d samples, want 1", app.sampler.CPU.Len())
	}
	if app.sampler.Memory.Len() != 0 {
		t.Errorf("Memory history has %d samples, want 0 before its collector reports",
			app.sampler.Memory.Len())
	}
	if app.sampler.Disk.Len() != 0 {
		t.Errorf("Disk history has %d samples, want 0 before its collector reports",
			app.sampler.Disk.Len())
	}
}

// apply records a sample from whatever state exists at the moment the interval
// elapses, so a metric updated immediately afterwards waits for the next
// interval. That is the same cadence the terminal UI uses.
func TestApplySamplesCurrentStateWhenTheIntervalElapses(t *testing.T) {
	app := newTestApp(50 * time.Millisecond)

	// These three land inside one interval, so together they produce a single
	// sample: the first records CPU alone, the other two arrive too soon.
	app.apply(cpuMetric(30, 8))
	app.apply(memMetric(50, 100))
	app.apply(diskMetric(25, 100))

	time.Sleep(60 * time.Millisecond)

	// This reading lands past the interval, so it records all three metrics
	// from the state above.
	app.apply(cpuMetric(60, 8))

	if got := app.sampler.CPU.Len(); got != 2 {
		t.Errorf("CPU history has %d samples, want 2", got)
	}
	if got := app.sampler.Memory.Len(); got != 1 {
		t.Errorf("Memory history has %d samples, want 1", got)
	}
	if got := app.sampler.Disk.Len(); got != 1 {
		t.Errorf("Disk history has %d samples, want 1", got)
	}

	if got := app.sampler.CPU.Max(); got != 60 {
		t.Errorf("CPU peak = %v, want the triggering reading 60", got)
	}
	if got := app.sampler.Memory.Max(); got != 50 {
		t.Errorf("Memory peak = %v, want the value held when the sample was taken, 50", got)
	}
	if got := app.sampler.Disk.Max(); got != 25 {
		t.Errorf("Disk peak = %v, want the value held when the sample was taken, 25", got)
	}
}

func TestConsumeDrainsTheChannelIntoState(t *testing.T) {
	app := newTestApp(testInterval)
	app.dataPipe = make(chan models.Metric, 2)
	app.dataPipe <- cpuMetric(30, 8)
	app.dataPipe <- memMetric(50, 100)
	close(app.dataPipe)

	app.consume()

	if !app.cpuReady || app.cpu.Percentage != 30 {
		t.Errorf("cpu = %+v (ready=%v), want 30", app.cpu, app.cpuReady)
	}
	if !app.memReady || app.memory.Total != 100 {
		t.Errorf("memory = %+v (ready=%v), want total 100", app.memory, app.memReady)
	}
}

func TestSnapshotAssemblesStateAndPercentages(t *testing.T) {
	app := newTestApp(testInterval)
	app.cpu = models.CPUPayload{Percentage: 42, CoreCount: 16}
	app.memory = models.MemoryPayload{Used: 8_000_000_000, Total: 16_000_000_000}
	app.disk = models.DiskPayload{Used: 250_000_000_000, Total: 1_000_000_000_000}
	app.startedAt = time.Now().Add(-90 * time.Second)

	got := app.snapshot()

	if got.CPU.Percent != 42 {
		t.Errorf("cpu percent = %v, want 42", got.CPU.Percent)
	}
	if got.CPUCores != 16 {
		t.Errorf("cpu cores = %d, want 16", got.CPUCores)
	}
	if got.Memory.Percent != 50 {
		t.Errorf("memory percent = %v, want 50 for 8/16 GB", got.Memory.Percent)
	}
	if got.Disk.Percent != 25 {
		t.Errorf("disk percent = %v, want 25 for 250/1000 GB", got.Disk.Percent)
	}
	if got.MemoryUsed != 8_000_000_000 || got.MemoryTotal != 16_000_000_000 {
		t.Errorf("memory bytes = %d/%d, want the raw values through",
			got.MemoryUsed, got.MemoryTotal)
	}
	if got.DiskUsed != 250_000_000_000 || got.DiskTotal != 1_000_000_000_000 {
		t.Errorf("disk bytes = %d/%d, want the raw values through",
			got.DiskUsed, got.DiskTotal)
	}
	if got.UptimeSecs < 89 || got.UptimeSecs > 95 {
		t.Errorf("uptime = %ds, want about 90s", got.UptimeSecs)
	}
}

func TestSnapshotWithNoCapacityYet(t *testing.T) {
	app := newTestApp(testInterval)

	got := app.snapshot()

	if got.Memory.Percent != 0 || got.Disk.Percent != 0 {
		t.Errorf("percentages = %v/%v, want 0/0 before any capacity is known",
			got.Memory.Percent, got.Disk.Percent)
	}
	if got.Memory.Trend != nil || got.Disk.Trend != nil {
		t.Error("trends should be nil for empty histories")
	}
}

func TestGetBootstrapReportsDefaults(t *testing.T) {
	withTempHome(t)

	got := NewApp().GetBootstrap()

	if got.Theme != utils.DefaultTheme {
		t.Errorf("theme = %q, want %q", got.Theme, utils.DefaultTheme)
	}
	if len(got.Themes) != len(theme.Presets) {
		t.Errorf("themes has %d entries, want the %d presets", len(got.Themes), len(theme.Presets))
	}
	if got.Version != utils.Version {
		t.Errorf("version = %q, want %q", got.Version, utils.Version)
	}
	if !got.AlwaysOnTop {
		t.Error("alwaysOnTop = false, want the widget pinned by default")
	}
}

func TestSetThemePersistsTheChoice(t *testing.T) {
	withTempHome(t)
	app := NewApp()

	applied, err := app.SetTheme("nord")
	if err != nil {
		t.Fatalf("SetTheme() error = %v", err)
	}
	if applied.Name != "nord" {
		t.Errorf("applied theme = %q, want %q", applied.Name, "nord")
	}

	if err := utils.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if got := utils.GetConfig().Theme; got != "nord" {
		t.Errorf("persisted theme = %q, want %q", got, "nord")
	}
}

func TestSetThemeUnknownNameFallsBackToOriginal(t *testing.T) {
	withTempHome(t)
	app := NewApp()

	applied, err := app.SetTheme("not-a-theme")
	if err != nil {
		t.Fatalf("SetTheme() error = %v", err)
	}
	if applied.Name != theme.Original.Name {
		t.Errorf("applied theme = %q, want the original fallback", applied.Name)
	}

	if err := utils.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	// The fallback must be what lands on disk, not the rejected name.
	if got := utils.GetConfig().Theme; got != theme.Original.Name {
		t.Errorf("persisted theme = %q, want %q", got, theme.Original.Name)
	}
}

// A theme change from the widget must not disturb the geometry the same widget
// saved earlier; both live in one shared file.
func TestSetThemeKeepsPersistedWindowState(t *testing.T) {
	withTempHome(t)

	if err := utils.SaveConfig(utils.Config{
		Theme:        "original",
		WindowX:      1556,
		WindowY:      24,
		WindowWidth:  340,
		WindowHeight: 236,
		AlwaysOnTop:  true,
	}); err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	app := NewApp()
	if _, err := app.SetTheme("dracula"); err != nil {
		t.Fatalf("SetTheme() error = %v", err)
	}

	if err := utils.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	got := utils.GetConfig()
	if got.Theme != "dracula" {
		t.Errorf("theme = %q, want %q", got.Theme, "dracula")
	}
	if got.WindowX != 1556 || got.WindowY != 24 {
		t.Errorf("window position = (%d,%d), want (1556,24)", got.WindowX, got.WindowY)
	}
	if got.WindowWidth != 340 || got.WindowHeight != 236 {
		t.Errorf("window size = %dx%d, want 340x236", got.WindowWidth, got.WindowHeight)
	}
}

// Shutdown runs from the Wails close path, which can fire without a preceding
// startup; it must not panic on the nil pipe.
func TestShutdownWithoutStartupIsSafe(t *testing.T) {
	app := NewApp()

	app.shutdown(context.Background())
}

// primaryScreenSize is a thin syscall wrapper, but a wrong SM_ constant would
// silently park the widget off-screen, so assert it reports a real display.
func TestPrimaryScreenSize(t *testing.T) {
	width, height := primaryScreenSize()

	if width <= 0 || height <= 0 {
		t.Errorf("primaryScreenSize() = (%d,%d), want a positive display size", width, height)
	}
}
