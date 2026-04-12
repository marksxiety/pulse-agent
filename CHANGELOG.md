# Changelog

## [1.1.0] - 2026-04-12

### Features (1.1.0)

- Theme system with 11 color presets: Original, Catppuccin Mocha, Dracula,
  Nord, Gruvbox, Tokyo Night, One Dark Pro, GitHub Dark, Ayu Mirage,
  Monokai Pro, and Synthwave
- Theme picker overlay (`t` key) with j/k and arrow key navigation, live
  theme switching, and full-width row highlight
- Persistent user config (`~/.pulse-agent/config.yaml`) with YAML-based
  theme saving and loading
- Terminal capability detection (Basic, Unicode, NerdFont) with automatic
  fallback for card icons, borders, sparklines, progress bars, scrollbars,
  and separators
- Platform-aware icon sets for dashboard cards and view title header
- Modal backdrop dimming via `ColorOverlay` for all modal overlays

### Enhancements

- Decoupled color palette from hardcoded values; theme applied globally
  via `ApplyTheme`
- Quit confirm modal restyled with yes/no selector navigation
  (j/k + arrow keys)
- Theme picker supports scrollable list with page window and position hint
  for large theme sets
- Footer hint shows active theme name at a glance
- Theme picker toggles on repeated `t` presses instead of only opening
- Added screenshot demo to README

### Fixes

- Theme picker scroll state synced with cursor position for correct
  rendering
- Theme picker modal correctly receives scroll offset for scroll-aware
  rendering
- `SaveTheme` error handled gracefully — closes theme picker on write
  failure
- `LoadConfig` error handled with fatal log exit on failure
- Empty or missing theme defaults to Original instead of Catppuccin Mocha
- Quit confirm modal receives cursor state for live selector rendering

## [1.0.0] - 2026-04-09

### Features (1.0.0)

- Quit confirmation dialog on `q` with Enter to confirm and Esc to cancel
  (Ctrl+C still quits immediately)
- Cross-platform disk root path resolution in DiskCollector (Windows `/`,
  Linux `/`, macOS `/`)
- Centralized `utils.Version` constant replacing hardcoded version strings
  across the codebase

### Testing

- External test suite for `models` (history, metric, payload) covering
  normalization, sparkline aggregation, threshold validation, and JSON
  serialization
- External test suite for `collector` (disk, CPU, memory) covering path
  resolution, percentage calculations, and error handling
- External test suite for `components` (modal, render, helpers) covering
  modal behavior, color formatting, and style application
- External test suite for `utils` (format) covering byte/percentage
  formatting, progress bar generation, and threshold coloring

### Documentation

- Redesigned README with ASCII banner, badges, expanded component
  descriptions, and contributing section
- Added setup guide (`docs/setup.md`) covering build, run, and keyboard
  controls
- Added testing guide (`docs/tests.md`) covering test structure, strategy,
  and conventions
- Removed standalone roadmap in favor of consolidated docs

### CI/CD

- Added GitHub Actions build workflow with cross-platform matrix (Linux,
  macOS, Windows × amd64, arm64)
- Added GitHub Actions lint workflow with `go vet`, `gofmt`, and
  `golangci-lint` v1.64
- Added GitHub Actions test workflow with coverage reporting
- Configured GoReleaser to skip auto-changelog generation

### Other

- Added MIT license

## [0.1.0] - 2026-04-07

### Features (0.1.0)

- Real-time system monitoring TUI dashboard with Catppuccin-inspired dark
  theme
- CPU usage percentage and logical core count monitoring (2s polling)
- Memory (RAM) monitoring with total, used, available, and pagefile stats
  (2s polling)
- Disk capacity and I/O monitoring with read/write bytes and operation
  counts (1s polling)
- Color-coded progress bars with dynamic thresholds (accent < 75%,
  yellow 75-89%, red >= 90%)
- 3-hour sparkline trend charts per metric with peak display
- Three-column responsive card layout (CPU, Memory, Disk) with horizontal
  and vertical centering
- Scrollable info modal (F1) with detailed metric glossary and custom
  scrollbar
- Concurrent metric collection with goroutine-per-collector architecture
- Graceful shutdown via SIGINT/SIGTERM with context-based cancellation
- Cross-platform builds for Linux, macOS, and Windows (amd64/arm64)
- Static binary with CGO_ENABLED=0
