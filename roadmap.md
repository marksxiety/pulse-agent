# Pulse-Agent Learning Roadmap

This roadmap tracks the Go fundamentals required to transform this CLI into a production-grade industrial tool.

## Phase 0: Learning Go Fundamentals

- [x] **Go Basics** — Learn the core syntax: variables, types, functions, and control flow  
  - Understand `var`, `const`, short variable declaration `:=`, and basic data types  
    - Notes:
      - When a variable is declared without an initial value, Go assigns a "zero value"
      - For `int`, it is `0`; for `string`, it is `""`; and for `bool`, it is `false`
      - `:=` automatically infers the variable type
- [x] **Packages & Modules** — Learn how `go mod`, imports, and the Go module system work  
  - Understand `go.mod`, `go.sum`, and how to organize code across packages  
  - Notes:
    - **Folder Name = Package Name**: If your folder is named `calculator`, the files inside should start with `package calculator`
    - **The `main` Package**: This is special. Every executable program must have one folder (usually the root) labeled `package main`. This is where `func main()` resides
    - **Exporting (Public vs. Private)**:
      - If a function starts with a capital letter (e.g., `Add`), it is exported and can be used by other packages
      - If it starts with a lowercase letter (e.g., `add`), it is private to its own package
    - **One Package per Folder**: Every `.go` file inside the same folder must have the exact same package name at the top
    - **`main` Package as the Entry Point**: To create a runnable program (an executable or binary), your starting file must use `package main` and contain `func main()`
    - **No Repeats**: You do not need to import a package to use functions from another file in the same folder
      - Example: If `file_a.go` and `file_b.go` are both in the `calculator` folder, `file_a` can use functions from `file_b` automatically
- [x] **Slices & Maps** — Go's primary data structures for collections  
  - Learn `append`, `make`, `len`, `cap`, and iterating with `range`.  
  - Notes:
    - `make()` is used to create and initialize certain built-in reference types.
    - Slices are for arrays, maps are for objects, and channels are used for pipelines or queues.

- [x] **Functions** — Multiple return values, named returns, and `defer`.  
  - Understand why Go uses the `func (args) (returnType, error)` pattern.
    - Notes:
      - Multiple returns (standard) return both the value and the error itself.
      - Named returns improve self-documentation by explicitly declaring return variables and their types.
      - `defer` ensures a function is executed regardless of what happens in the surrounding function. It is commonly used for closing database connections or cleaning up resources.

- [x] **Go Tooling** — Essential commands for development  
  - `go run`, `go build`, `go test`, `go fmt`, `go vet`
    - Notes:
      - `go run` — Compiles and runs the code (no file saved).
      - `go build` — Compiles the code into a standalone `.exe` or binary.
      - `go fmt` — Rewrites code with proper formatting and spacing.
      - `go vet` — Analyzes code and reports potential logic issues.
      - `go test` — Executes tests and outputs a "PASS" or "FAIL" report.

## Phase 1: Data Modeling (The "What")

- [x] **Structs & Tags** -- Learn how to define the `Metric` struct.
  - Understand `json:"value"` tags for future API integration.
    - Notes:
      - `Structs` are **groupings of related data**. They are usually used for forms or structured data. Think of them as the **"container" for your data logic**.
      - `Tags` are used to label struct fields and define where or how they will be used, such as `json:` or `db:`. Think of them as **"mapping instructions"** that tell external libraries how to handle that data.
- [x] **Custom Types** -- Creating a `type Source string` to replace raw strings for better safety.
  - Notes:
    - Custom types act as a constraint that restricts data to a defined type. It's the equivalent of `Enum` in TypeScript, ensuring only specific values can be assigned to a variable.
- [x] **Pointers (`*`)** -- Understanding when to pass a reference to a metric instead of a copy.
  - Notes:
    - Using pointers depends on how large the struct is and how frequently it is passed around.
    - Pointers add overhead (heap allocation, GC tracking, nil checks). For small structs (like `Metric` at ~24 bytes), passing by value is often cheaper and simpler.
    - Use pointers when: (1) the struct is large and copying is expensive, (2) you need to mutate the original, or (3) a method requires a pointer receiver to satisfy an interface.

## Phase 2: Logic & Abstraction (The "How")

- [x] **Interfaces** -- **CRITICAL:** Creating a `type Collector interface { Collect() Metric }`.
  - This allows you to add any new hardware (GPU, Network) without changing `main.go`.
  - Notes:
    - It describes the **behavior**.
    - This serves as the ability of the type
    - For the analogy, `Struct` is the noun and `interface` is the verb.
- [x] **Methods** -- Attaching logic directly to your collectors (e.g., `func (c CPUCollector) Collect()`).
  - Notes:
    - Methods are the **bridge between structs and interfaces**. A struct defines the data, an interface defines the expected behavior, and a method is the actual implementation that satisfies it.
    - When a struct implements all the methods defined by an interface, it is said to **satisfy** that interface — no explicit declaration needed.
    - **Value receivers** (`func (c CPUCollector)`) operate on a copy of the struct. **Pointer receivers** (`func (c *CPUCollector)`) operate on the original and can mutate it.
- [x] **Error Handling** -- Mastering the `if err != nil` pattern to handle OS permission denials.
  - Notes:
    - **The Pattern**: Call a function that returns a value and an error, then immediately check `if err != nil` to handle failures before proceeding (e.g., `result, err := function(); if err != nil { ... }`).
    - **Permission Errors**: Use `os.IsPermission(err)` to detect when the OS denies access, so you can inform the user that elevated privileges are required.
    - **Interfaces**: Always include `error` in your interface method signatures (e.g., `Collect() (Metric, error)`) so the caller knows whether the returned data is valid.

## Phase 3: High-Performance Concurrency (The "Engine")

- [x] **Goroutines** -- Mastering the `go` keyword and managing the lifecycle of background workers.
- Notes:
  - Managed by `Go runtime` (not OS directly)
  - `Go` is not single thread by default. It uses 'GMP' scheduler
- [x] **Channels** -- Buffered vs. Unbuffered channels.
  - Using `select` to handle multiple data streams and "Quit" signals.
  - Notes:
    - This serve as the `medium` of goroutines.
    - 'Sending' `Mychannel <- element`
    - 'Recieving' `element := <-Mychannel`
- [x] **Context (`context` package)** -- Implementing graceful shutdowns so the agent finishes its last task before exiting.
- Notes:
  - provides a mechanism to control the lifecycle, cancellation, and propagation of requests across multiple goroutines.
  - It can add also in http requests, database operations

## Phase 4: Production Features (The "Where")

- [x] **Environment Variables** -- Using `os.Getenv` to configure polling intervals without re-compiling.
  - Notes:
    - `Setenv` used to override the current env usually for testing or debugging
    - `Getenv` is the retrieving or most common approach for environment variables
    - `LookupEnv` better approach to ensure that the enviroment variables are set before using in the logic/processes.
    - Go **does not** read `.env` files natively. Use `github.com/joho/godotenv` to auto-load `.env` into the process environment via `godotenv.Load()` (typically in `init()`).
- [ ] **TUI (Terminal UI)** -- Using `Bubble Tea` to create a visual dashboard with real-time graphs.

## Phase 5: Metric Collection Goals

- [ ] **CPU Metrics** -- Fetch CPU usage percentage, core count, and load averages.
  - Use `github.com/shirou/gopsutil/v3/cpu` or parse `/proc/stat` on Linux.
- [ ] **Memory Metrics** -- Fetch total, used, available memory and swap usage.
  - Use `github.com/shirou/gopsutil/v3/mem` or call OS-specific APIs.
- [ ] **Disk Metrics** -- Fetch disk usage (total, used, free) and I/O stats.
  - Use `github.com/shirou/gopsutil/v3/disk` or parse `df`/`wmic` output.
