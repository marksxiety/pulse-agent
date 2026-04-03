# Pulse-Check Learning Roadmap

This roadmap tracks the Go fundamentals required to transform this CLI into a production-grade industrial tool.

## Phase 0: Learning Go Fundamentals

- [x] **Go Basics** — Learn the core syntax: variables, types, functions, and control flow  
  - Understand `var`, `const`, short variable declaration `:=`, and basic data types  
    - Notes:
      - When a variable is declared without an initial value, Go assigns a "zero value"
      - For `int`, it is `0`; for `string`, it is `""`; and for `bool`, it is `false`
      - `:=` automatically infers the variable type
- [ ] **Packages & Modules** -- How `go mod`, imports, and the Go module system work.
  - Understand `go.mod`, `go.sum`, and organizing code across packages.
- [ ] **Slices & Maps** -- Go's primary data structures for collections.
  - Learn `append`, `make`, `len`, `cap`, and iterating with `range`.
- [ ] **Functions** -- Multiple return values, named returns, and defer.
  - Understand why Go uses `func (args) (returnType, error)` pattern.
- [ ] **Go Tooling** -- Essential commands for development.
  - `go run`, `go build`, `go test`, `go fmt`, `go vet`.

## Phase 1: Data Modeling (The "What")

- [ ] **Structs & Tags** -- Learn how to define the `Metric` struct.
  - Understand `json:"value"` tags for future API integration.
- [ ] **Custom Types** -- Creating a `type Source string` to replace raw strings for better safety.
- [ ] **Pointers (`*`)** -- Understanding when to pass a reference to a metric instead of a copy to save RAM.

## Phase 2: Logic & Abstraction (The "How")

- [ ] **Interfaces** -- **CRITICAL:** Creating a `type Collector interface { Collect() Metric }`.
  - This allows you to add any new hardware (GPU, Network) without changing `main.go`.
- [ ] **Methods** -- Attaching logic directly to your collectors (e.g., `func (c CPUCollector) Collect()`).
- [ ] **Error Handling** -- Mastering the `if err != nil` pattern to handle OS permission denials.

## Phase 3: High-Performance Concurrency (The "Engine")

- [ ] **Goroutines** -- Mastering the `go` keyword and managing the lifecycle of background workers.
- [ ] **Channels** -- Buffered vs. Unbuffered channels.
  - Using `select` to handle multiple data streams and "Quit" signals.
- [ ] **Context (`context` package)** -- Implementing graceful shutdowns so the agent finishes its last task before exiting.

## Phase 4: Production Features (The "Where")

- [ ] **Standard Library (`net/http`)** -- Sending your metrics to a central server (like KRONOS) via POST requests.
- [ ] **Environment Variables** -- Using `os.Getenv` to configure polling intervals without re-compiling.
- [ ] **TUI (Terminal UI)** -- Using `Bubble Tea` to create a visual dashboard with real-time graphs.

## Phase 5: Metric Collection Goals

- [ ] **CPU Metrics** -- Fetch CPU usage percentage, core count, and load averages.
  - Use `github.com/shirou/gopsutil/v3/cpu` or parse `/proc/stat` on Linux.
- [ ] **Memory Metrics** -- Fetch total, used, available memory and swap usage.
  - Use `github.com/shirou/gopsutil/v3/mem` or call OS-specific APIs.
- [ ] **Disk Metrics** -- Fetch disk usage (total, used, free) and I/O stats.
  - Use `github.com/shirou/gopsutil/v3/disk` or parse `df`/`wmic` output.
