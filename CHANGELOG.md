# Changelog

## [Unreleased]

## [0.1.0] - 2026-04-07

### Features
- Real-time system monitoring TUI dashboard with Catppuccin-inspired dark theme
- CPU usage percentage and logical core count monitoring (2s polling)
- Memory (RAM) monitoring with total, used, available, and pagefile stats (2s polling)
- Disk capacity and I/O monitoring with read/write bytes and operation counts (1s polling)
- Color-coded progress bars with dynamic thresholds (accent < 75%, yellow 75-89%, red >= 90%)
- 3-hour sparkline trend charts per metric with peak display
- Three-column responsive card layout (CPU, Memory, Disk) with horizontal and vertical centering
- Scrollable info modal (F1) with detailed metric glossary and custom scrollbar
- Concurrent metric collection with goroutine-per-collector architecture
- Graceful shutdown via SIGINT/SIGTERM with context-based cancellation
- Cross-platform builds for Linux, macOS, and Windows (amd64/arm64)
- Static binary with CGO_ENABLED=0
