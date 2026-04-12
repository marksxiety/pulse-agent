# Testing

Pulse-Agent uses Go's standard `testing` package. All tests live under
`tests/` and use **external test packages** (`package xxx_test`) so they
only exercise the public API.

## Running Tests

```bash
# Run all tests
go test ./tests/...

# Verbose output
go test ./tests/... -v

# Coverage summary
go test ./tests/... -cover

# Coverage report (HTML)
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Structure

```text
tests/
  models/
    history_test.go      # Circular buffer: push, wrap, min, max, avg
    metric_test.go       # JSON marshaling per payload type
    payload_test.go      # SourceType() for CPU, MEM, DISK payloads
  collector/
    collector_test.go    # CPU, Memory, Disk with mocked SystemInfo
  utils/
    format_test.go       # BytesToGB, BytesToMB, FormatUptime, PadToHeight
  components/
    render_test.go       # ProgressBar, Sparkline, CentreBlock
    helpers_test.go      # OverflowGuard, Separator, Row
    modal_test.go        # InfoModal scroll behavior and edge cases
```

## Test Approach

### Models

Pure functions — no mocking needed. Table-driven where applicable.
Tests cover empty state, single values, multiple values, and circular
buffer wrap-around.

### Collectors

A `mockSystemInfo` struct implements the `SystemInfo` interface,
replacing real `gopsutil` calls. Tests verify data flow, error handling,
and graceful shutdown via context cancellation.

### Utils

Table-driven tests for format conversions. Edge cases: zero input,
rounding, boundary values.

### Components

Rendering functions are tested for non-empty output, edge-case inputs
(zero width, empty data, negative scroll), and structural correctness
(overflow guard clipping).

## Conventions

- External test packages only — unexported internals are not directly tested.
- No third-party test libraries or assertion frameworks.
- No test helpers beyond `t.Fatal` / `t.Errorf`.
- Collector tests use a timeout to avoid hanging on blocking channel reads.
