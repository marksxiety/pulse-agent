package components

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Name       string
	Bg         string
	Surface    string
	Border     string
	Muted      string
	Subtle     string
	Text       string
	Dim        string
	CPUAccent  string
	MemAccent  string
	DiskAccent string
	Warn       string
	Danger     string
}

func (t Theme) Colors() (bg, surface, border, muted, subtle, text, dim, cpu, mem, disk, warn, dang lipgloss.Color) {
	return lipgloss.Color(t.Bg), lipgloss.Color(t.Surface), lipgloss.Color(t.Border),
		lipgloss.Color(t.Muted), lipgloss.Color(t.Subtle), lipgloss.Color(t.Text), lipgloss.Color(t.Dim),
		lipgloss.Color(t.CPUAccent), lipgloss.Color(t.MemAccent), lipgloss.Color(t.DiskAccent),
		lipgloss.Color(t.Warn), lipgloss.Color(t.Danger)
}

var ThemePresets = []Theme{
	CatppuccinMochaTheme,
	DraculaTheme,
	NordTheme,
	GruvboxTheme,
	TokyoNightTheme,
}

var ThemeNames []string

func init() {
	ThemeNames = make([]string, len(ThemePresets))
	for i, t := range ThemePresets {
		ThemeNames[i] = t.Name
	}
}

func GetTheme(name string) Theme {
	for _, t := range ThemePresets {
		if t.Name == name {
			return t
		}
	}
	return CatppuccinMochaTheme
}

var CatppuccinMochaTheme = Theme{
	Name:       "catppuccin-mocha",
	Bg:         "#0d0f14",
	Surface:    "#141720",
	Border:     "#1e2330",
	Muted:      "#3b4260",
	Subtle:     "#6272a4",
	Text:       "#cdd6f4",
	Dim:        "#585b70",
	CPUAccent:  "#89b4fa",
	MemAccent:  "#a6e3a1",
	DiskAccent: "#cba6f7",
	Warn:       "#f9e2af",
	Danger:     "#f38ba8",
}

var DraculaTheme = Theme{
	Name:       "dracula",
	Bg:         "#1e1f29",
	Surface:    "#282a36",
	Border:     "#44475a",
	Muted:      "#44475a",
	Subtle:     "#6272a4",
	Text:       "#f8f8f2",
	Dim:        "#6272a4",
	CPUAccent:  "#8be9fd",
	MemAccent:  "#50fa7b",
	DiskAccent: "#bd93f9",
	Warn:       "#f1fa8c",
	Danger:     "#ff5555",
}

var NordTheme = Theme{
	Name:       "nord",
	Bg:         "#2e3440",
	Surface:    "#3b4252",
	Border:     "#434c5e",
	Muted:      "#4c566a",
	Subtle:     "#7b88a1",
	Text:       "#eceff4",
	Dim:        "#616e88",
	CPUAccent:  "#88c0d0",
	MemAccent:  "#a3be8c",
	DiskAccent: "#b48ead",
	Warn:       "#ebcb8b",
	Danger:     "#bf616a",
}

var GruvboxTheme = Theme{
	Name:       "gruvbox",
	Bg:         "#1d2021",
	Surface:    "#282828",
	Border:     "#3c3836",
	Muted:      "#504945",
	Subtle:     "#928374",
	Text:       "#ebdbb2",
	Dim:        "#665c54",
	CPUAccent:  "#83a598",
	MemAccent:  "#b8bb26",
	DiskAccent: "#d3869b",
	Warn:       "#fabd2f",
	Danger:     "#fb4934",
}

var TokyoNightTheme = Theme{
	Name:       "tokyo-night",
	Bg:         "#1a1b26",
	Surface:    "#24283b",
	Border:     "#3b4261",
	Muted:      "#3b4261",
	Subtle:     "#565f89",
	Text:       "#c0caf5",
	Dim:        "#565f89",
	CPUAccent:  "#7dcfff",
	MemAccent:  "#9ece6a",
	DiskAccent: "#bb9af7",
	Warn:       "#e0af68",
	Danger:     "#f7768e",
}
