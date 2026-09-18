// Package theme holds the colour palettes shared by every frontend.
//
// It deliberately has no UI dependencies: the terminal UI adapts these values
// into lipgloss styles, while the desktop widget serialises them straight to
// the webview as JSON.
package theme

// Theme is a named colour palette. Every field is a hex colour string.
type Theme struct {
	Name       string `json:"name"`
	Bg         string `json:"bg"`
	Surface    string `json:"surface"`
	Border     string `json:"border"`
	Muted      string `json:"muted"`
	Subtle     string `json:"subtle"`
	Text       string `json:"text"`
	Dim        string `json:"dim"`
	CPUAccent  string `json:"cpu_accent"`
	MemAccent  string `json:"mem_accent"`
	DiskAccent string `json:"disk_accent"`
	Warn       string `json:"warn"`
	Danger     string `json:"danger"`
}

var Original = Theme{
	Name:       "original",
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

var CatppuccinMocha = Theme{
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

var Dracula = Theme{
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

var Nord = Theme{
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

var Gruvbox = Theme{
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

var TokyoNight = Theme{
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

var OneDarkPro = Theme{
	Name:       "one-dark-pro",
	Bg:         "#21252b",
	Surface:    "#282c34",
	Border:     "#3e4451",
	Muted:      "#3e4451",
	Subtle:     "#5c6370",
	Text:       "#abb2bf",
	Dim:        "#4b5263",
	CPUAccent:  "#61afef",
	MemAccent:  "#98c379",
	DiskAccent: "#c678dd",
	Warn:       "#e5c07b",
	Danger:     "#e06c75",
}

var GitHubDark = Theme{
	Name:       "github-dark",
	Bg:         "#0d1117",
	Surface:    "#161b22",
	Border:     "#30363d",
	Muted:      "#21262d",
	Subtle:     "#8b949e",
	Text:       "#e6edf3",
	Dim:        "#484f58",
	CPUAccent:  "#58a6ff",
	MemAccent:  "#3fb950",
	DiskAccent: "#bc8cff",
	Warn:       "#d29922",
	Danger:     "#f85149",
}

var AyuMirage = Theme{
	Name:       "ayu-mirage",
	Bg:         "#1a1f29",
	Surface:    "#242936",
	Border:     "#343d4d",
	Muted:      "#3d4752",
	Subtle:     "#5c6773",
	Text:       "#cccac2",
	Dim:        "#4a5568",
	CPUAccent:  "#5ccfe6",
	MemAccent:  "#bae67e",
	DiskAccent: "#d4bfff",
	Warn:       "#ffd580",
	Danger:     "#ff6666",
}

var MonokaiPro = Theme{
	Name:       "monokai-pro",
	Bg:         "#19181a",
	Surface:    "#221f22",
	Border:     "#3a3a3c",
	Muted:      "#403e41",
	Subtle:     "#727072",
	Text:       "#fcfcfa",
	Dim:        "#5b595c",
	CPUAccent:  "#78dce8",
	MemAccent:  "#a9dc76",
	DiskAccent: "#ab9df2",
	Warn:       "#ffd866",
	Danger:     "#ff6188",
}

var Synthwave = Theme{
	Name:       "synthwave",
	Bg:         "#1a1333",
	Surface:    "#241b4d",
	Border:     "#3a2d6e",
	Muted:      "#3d2d7a",
	Subtle:     "#8b6fc8",
	Text:       "#f0eff1",
	Dim:        "#5a4a8a",
	CPUAccent:  "#36f9f6",
	MemAccent:  "#72f1b8",
	DiskAccent: "#fe45e6",
	Warn:       "#fede5d",
	Danger:     "#fe4450",
}

// Presets is the ordered list shown in the UI. Order is user-visible: the
// theme picker cursors through it by index.
var Presets = []Theme{
	Original,
	CatppuccinMocha,
	Dracula,
	Nord,
	Gruvbox,
	TokyoNight,
	OneDarkPro,
	GitHubDark,
	AyuMirage,
	MonokaiPro,
	Synthwave,
}

// Names holds the preset names, in Presets order.
var Names []string

func init() {
	Names = make([]string, len(Presets))
	for i, t := range Presets {
		Names[i] = t.Name
	}
}

// Get returns the preset with the given name, falling back to Original.
func Get(name string) Theme {
	for _, t := range Presets {
		if t.Name == name {
			return t
		}
	}
	return Original
}
