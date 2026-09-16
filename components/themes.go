package components

import "pulse-agent/theme"

// The palettes themselves live in the UI-free theme package so the desktop
// widget can consume them without pulling in lipgloss. These aliases keep the
// terminal UI's call sites unchanged.
var (
	OriginalTheme = theme.Original
	ThemeNames    = theme.Names
)

// GetTheme returns the named palette, falling back to the original.
func GetTheme(name string) theme.Theme {
	return theme.Get(name)
}
