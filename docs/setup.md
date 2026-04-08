# Setup

## Prerequisites

- [Go 1.21+](https://go.dev/dl/)

## Clone

```bash
git clone https://github.com/<org>/pulse-agent.git
cd pulse-agent
```

## Install Dependencies

```bash
go mod download
```

## Build

```bash
go build -o pulse-agent ./agent
```

This produces a `pulse-agent` binary (or `pulse-agent.exe` on Windows).

## Run

```bash
./pulse-agent
```

### Controls

| Key | Action |
|-----|--------|
| `q` / `Ctrl+C` | Quit |
| `F1` | Toggle info modal |
| `Esc` | Close info modal |
| `Enter` | Close info modal |
| `Up` / `Down` | Scroll info modal |
