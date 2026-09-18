package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const Version = "v1.0.0"

// DefaultTheme is used when no theme has been persisted yet.
const DefaultTheme = "original"

// Config is the persisted user configuration. It is shared by every frontend
// and stored as a flat "key: value" file at ConfigPath.
//
// The window fields are only meaningful to the desktop widget; the terminal UI
// ignores them but must preserve them when it writes.
type Config struct {
	Theme string `yaml:"theme"`

	WindowX      int  `yaml:"window_x"`
	WindowY      int  `yaml:"window_y"`
	WindowWidth  int  `yaml:"window_width"`
	WindowHeight int  `yaml:"window_height"`
	AlwaysOnTop  bool `yaml:"always_on_top"`
}

// Both frontends share one file, and the widget reaches the config from several
// goroutines at once — webview-bound methods, the close handler, and
// second-instance launches. Every access therefore goes through cfgMu.
var (
	cfgMu sync.RWMutex
	cfg   = defaults()
)

// defaults returns the configuration used when no file has been written yet.
func defaults() Config {
	return Config{
		Theme: DefaultTheme,
		// The desktop widget is pinned above other windows unless the user
		// explicitly unpins it.
		AlwaysOnTop: true,
	}
}

func ConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, ".pulse-agent")
}

func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config.yaml")
}

func GetConfig() Config {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	return cfg
}

// SaveConfig replaces the entire configuration, both in memory and on disk.
func SaveConfig(c Config) error {
	if c.Theme == "" {
		c.Theme = DefaultTheme
	}

	cfgMu.Lock()
	cfg = c
	cfgMu.Unlock()

	return writeConfigFile(c)
}

// UpdateConfig re-reads the configuration from disk, applies mutate, and writes
// the result back.
//
// Re-reading first is what stops one frontend from clobbering settings written
// by the other while it was running: the terminal UI changing the theme must
// not revert the window position the widget saved at some earlier point, and
// vice versa. Merging against the on-disk state rather than a possibly stale
// in-memory copy is the whole point of this function.
func UpdateConfig(mutate func(*Config)) error {
	cfgMu.Lock()
	defer cfgMu.Unlock()

	current, err := readConfigFile()
	if err != nil {
		return err
	}

	mutate(&current)
	if current.Theme == "" {
		current.Theme = DefaultTheme
	}

	cfg = current
	return writeConfigFile(current)
}

// SaveTheme persists only the theme so that unrelated settings — such as the
// widget's window geometry — survive a theme change.
func SaveTheme(name string) error {
	return UpdateConfig(func(c *Config) { c.Theme = name })
}

func LoadConfig() error {
	c, err := readConfigFile()
	if err != nil {
		return err
	}

	cfgMu.Lock()
	cfg = c
	cfgMu.Unlock()

	return nil
}

// readConfigFile returns the on-disk configuration, falling back to the
// defaults when no file exists yet.
func readConfigFile() (Config, error) {
	data, err := os.ReadFile(ConfigPath())
	if err != nil {
		if os.IsNotExist(err) {
			return defaults(), nil
		}
		return Config{}, err
	}

	return unmarshalConfig(data), nil
}

// writeConfigFile persists c. Callers are responsible for holding cfgMu.
func writeConfigFile(c Config) error {
	if err := os.MkdirAll(ConfigDir(), 0o755); err != nil {
		return err
	}
	return os.WriteFile(ConfigPath(), []byte(marshalConfig(c)), 0o644)
}

func marshalConfig(c Config) string {
	var b strings.Builder
	fmt.Fprintf(&b, "theme: %s\n", c.Theme)
	fmt.Fprintf(&b, "window_x: %d\n", c.WindowX)
	fmt.Fprintf(&b, "window_y: %d\n", c.WindowY)
	fmt.Fprintf(&b, "window_width: %d\n", c.WindowWidth)
	fmt.Fprintf(&b, "window_height: %d\n", c.WindowHeight)
	fmt.Fprintf(&b, "always_on_top: %t\n", c.AlwaysOnTop)
	return b.String()
}

func unmarshalConfig(data []byte) Config {
	c := defaults()

	for i := 0; i < len(data); i++ {
		lineStart := i
		for i < len(data) && data[i] != '\n' {
			i++
		}
		applyConfigLine(&c, string(data[lineStart:i]))
	}

	if c.Theme == "" {
		c.Theme = DefaultTheme
	}
	return c
}

// applyConfigLine folds a single "key: value" line into c. Unknown keys are
// ignored, so a config written by a newer build degrades gracefully.
func applyConfigLine(c *Config, line string) {
	key, value, ok := parseYAMLLine(line)
	if !ok {
		return
	}

	switch key {
	case "theme":
		c.Theme = value
	case "window_x":
		c.WindowX = parseIntOr(value, c.WindowX)
	case "window_y":
		c.WindowY = parseIntOr(value, c.WindowY)
	case "window_width":
		c.WindowWidth = parseIntOr(value, c.WindowWidth)
	case "window_height":
		c.WindowHeight = parseIntOr(value, c.WindowHeight)
	case "always_on_top":
		c.AlwaysOnTop = parseBoolOr(value, c.AlwaysOnTop)
	}
}

func parseIntOr(value string, fallback int) int {
	n, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return n
}

func parseBoolOr(value string, fallback bool) bool {
	b, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return b
}

func parseYAMLLine(line string) (key, value string, ok bool) {
	for i := 0; i < len(line); i++ {
		if line[i] == ':' {
			key = line[:i]
			rest := line[i+1:]
			for len(rest) > 0 && (rest[0] == ' ' || rest[0] == '\t') {
				rest = rest[1:]
			}
			if len(rest) > 0 && rest[len(rest)-1] == '\r' {
				rest = rest[:len(rest)-1]
			}
			return key, rest, true
		}
	}
	return "", "", false
}
