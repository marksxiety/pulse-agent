//go:build windows

package main

import (
	"context"
	"log"
	"sync"
	"time"

	"pulse-agent/collector"
	"pulse-agent/models"
	"pulse-agent/theme"
	"pulse-agent/utils"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	// EventSnapshot carries one tick of metrics to the webview.
	EventSnapshot = "pulse:snapshot"

	// historyPoints is how many samples the chart receives per tick. The
	// backend holds up to models.HistoryCapacity samples, so downsampling here
	// keeps the event payload small instead of shipping three hours of raw
	// data every second.
	historyPoints = 120

	// historyInterval matches the terminal UI: history is sampled every 5s, so
	// HistoryCapacity samples covers three hours of trend.
	historyInterval = 5 * time.Second

	// emitInterval is how often the merged snapshot is pushed to the webview.
	emitInterval = time.Second
)

// Snapshot is everything the widget renders for a single tick.
type Snapshot struct {
	CPU    models.Series `json:"cpu"`
	Memory models.Series `json:"memory"`
	Disk   models.Series `json:"disk"`

	CPUCores    int    `json:"cpuCores"`
	MemoryUsed  uint64 `json:"memoryUsed"`
	MemoryTotal uint64 `json:"memoryTotal"`
	DiskUsed    uint64 `json:"diskUsed"`
	DiskTotal   uint64 `json:"diskTotal"`
	UptimeSecs  int64  `json:"uptimeSecs"`
}

// Bootstrap is the one-off state the widget needs before the first snapshot.
type Bootstrap struct {
	Theme       string        `json:"theme"`
	Themes      []theme.Theme `json:"themes"`
	Version     string        `json:"version"`
	AlwaysOnTop bool          `json:"alwaysOnTop"`
}

// App is the Wails-bound backend. Its exported methods are callable from the
// webview; everything else stays behind the mutex.
type App struct {
	ctx    context.Context
	cancel context.CancelFunc

	// producers tracks the collector goroutines and the emitter.
	producers sync.WaitGroup

	live     collector.LiveSystemInfo
	dataPipe chan models.Metric

	mu        sync.RWMutex
	cpu       models.CPUPayload
	memory    models.MemoryPayload
	disk      models.DiskPayload
	cpuReady  bool
	memReady  bool
	diskReady bool
	startedAt time.Time
	sampler   *models.Sampler
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.startedAt = time.Now()
	a.sampler = models.NewSampler(historyInterval)

	runCtx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel
	a.dataPipe = make(chan models.Metric)

	a.producers.Add(4)
	go func() {
		defer a.producers.Done()
		collector.NewCPUCollector(a.live).Collect(runCtx, a.dataPipe)
	}()
	go func() {
		defer a.producers.Done()
		collector.NewMemoryCollector(a.live).Collect(runCtx, a.dataPipe)
	}()
	go func() {
		defer a.producers.Done()
		collector.NewDiskCollector(a.live).Collect(runCtx, a.dataPipe)
	}()
	go a.emit(runCtx)

	// The reader is deliberately not part of producers: shutdown closes the
	// pipe only after every producer has stopped, which is what lets the
	// collectors finish a blocked send and return.
	go a.consume()
}

// domReady restores the saved window position once the webview is alive.
func (a *App) domReady(ctx context.Context) {
	cfg := utils.GetConfig()

	if cfg.WindowX == 0 && cfg.WindowY == 0 {
		// No position has ever been saved, so park the widget inside the
		// top-right of the primary display rather than letting it land under
		// the taskbar at 0,0.
		screenW, _ := primaryScreenSize()
		runtime.WindowSetPosition(ctx, screenW-sizeOrDefault(cfg.WindowWidth, defaultWidth)-widgetMargin, widgetMargin)
		return
	}

	runtime.WindowSetPosition(ctx, cfg.WindowX, cfg.WindowY)
}

// beforeClose persists geometry so the widget reappears where it was left.
func (a *App) beforeClose(context.Context) bool {
	a.persistWindowState()
	return false
}

func (a *App) shutdown(context.Context) {
	if a.cancel == nil {
		return
	}
	a.cancel()
	a.producers.Wait()
	close(a.dataPipe)
	a.cancel = nil
}

// GetBootstrap returns the one-off state the widget needs at startup.
func (a *App) GetBootstrap() Bootstrap {
	cfg := utils.GetConfig()
	return Bootstrap{
		Theme:       cfg.Theme,
		Themes:      theme.Presets,
		Version:     utils.Version,
		AlwaysOnTop: cfg.AlwaysOnTop,
	}
}

// SetTheme persists the chosen palette and returns it so the webview can apply
// it without a second round trip. The terminal UI picks the change up on its
// next launch because both frontends share one config file.
func (a *App) SetTheme(name string) (theme.Theme, error) {
	chosen := theme.Get(name)
	if err := utils.SaveTheme(chosen.Name); err != nil {
		return chosen, err
	}
	return chosen, nil
}

// SetAlwaysOnTop pins or unpins the widget above other windows.
func (a *App) SetAlwaysOnTop(pinned bool) error {
	runtime.WindowSetAlwaysOnTop(a.ctx, pinned)
	return utils.UpdateConfig(func(c *utils.Config) { c.AlwaysOnTop = pinned })
}

// Quit closes the widget, persisting window state first.
func (a *App) Quit() {
	a.persistWindowState()
	runtime.Quit(a.ctx)
}

func (a *App) focus() {
	if a.ctx == nil {
		return
	}
	runtime.WindowShow(a.ctx)
	runtime.WindowSetAlwaysOnTop(a.ctx, utils.GetConfig().AlwaysOnTop)
}

// consume drains the collector channel into the shared state.
func (a *App) consume() {
	for m := range a.dataPipe {
		a.apply(m)
	}
}

func (a *App) apply(m models.Metric) {
	a.mu.Lock()
	defer a.mu.Unlock()

	switch p := m.Data.(type) {
	case models.CPUPayload:
		a.cpu, a.cpuReady = p, true
	case models.MemoryPayload:
		a.memory, a.memReady = p, true
	case models.DiskPayload:
		a.disk, a.diskReady = p, true
	}

	a.sampler.Record(time.Now(), a.readyCPU(), a.readyMemory(), a.readyDisk())
}

// The ready helpers pass nil for a collector that has not reported yet, which
// is how the sampler tells "no data" apart from a genuine zero reading.

func (a *App) readyCPU() *models.CPUPayload {
	if !a.cpuReady {
		return nil
	}
	return &a.cpu
}

func (a *App) readyMemory() *models.MemoryPayload {
	if !a.memReady {
		return nil
	}
	return &a.memory
}

func (a *App) readyDisk() *models.DiskPayload {
	if !a.diskReady {
		return nil
	}
	return &a.disk
}

func (a *App) emit(ctx context.Context) {
	defer a.producers.Done()

	ticker := time.NewTicker(emitInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if a.ctx != nil {
				runtime.EventsEmit(a.ctx, EventSnapshot, a.snapshot())
			}
		}
	}
}

func (a *App) snapshot() Snapshot {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return Snapshot{
		CPU:         models.NewSeries(a.sampler.CPU, a.cpu.Percentage, historyPoints),
		Memory:      models.NewSeries(a.sampler.Memory, models.Percent(a.memory.Used, a.memory.Total), historyPoints),
		Disk:        models.NewSeries(a.sampler.Disk, models.Percent(a.disk.Used, a.disk.Total), historyPoints),
		CPUCores:    a.cpu.CoreCount,
		MemoryUsed:  a.memory.Used,
		MemoryTotal: a.memory.Total,
		DiskUsed:    a.disk.Used,
		DiskTotal:   a.disk.Total,
		UptimeSecs:  int64(time.Since(a.startedAt).Seconds()),
	}
}

func (a *App) persistWindowState() {
	if a.ctx == nil {
		return
	}

	x, y := runtime.WindowGetPosition(a.ctx)
	w, h := runtime.WindowGetSize(a.ctx)

	err := utils.UpdateConfig(func(c *utils.Config) {
		c.WindowX, c.WindowY = x, y
		c.WindowWidth, c.WindowHeight = w, h
	})
	if err != nil {
		log.Printf("[widget] failed to persist window state: %v", err)
	}
}
