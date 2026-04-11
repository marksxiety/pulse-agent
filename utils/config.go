package utils

import (
	"os"
	"path/filepath"
	"sync"
)

const Version = "v1.0.0"

type Config struct {
	Theme string `yaml:"theme"`
}

var (
	cfgOnce sync.Once
	cfg     Config
)

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
	cfgOnce.Do(func() {
		cfg = Config{Theme: "catppuccin-mocha"}
	})
	return cfg
}

func SaveTheme(name string) error {
	cfg = GetConfig()
	cfg.Theme = name

	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	content := "theme: " + name + "\n"
	return os.WriteFile(ConfigPath(), []byte(content), 0o644)
}

func LoadConfig() error {
	data, err := os.ReadFile(ConfigPath())
	if err != nil {
		if os.IsNotExist(err) {
			cfg = Config{Theme: "catppuccin-mocha"}
			return nil
		}
		return err
	}

	theme := parseThemeField(data)
	if theme != "" {
		cfg = Config{Theme: theme}
	} else {
		cfg = Config{Theme: "catppuccin-mocha"}
	}
	return nil
}

func parseThemeField(data []byte) string {
	for i := 0; i < len(data); i++ {
		lineStart := i
		for i < len(data) && data[i] != '\n' {
			i++
		}
		line := string(data[lineStart:i])
		key, val, ok := parseYAMLLine(line)
		if ok && key == "theme" {
			return val
		}
	}
	return ""
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
