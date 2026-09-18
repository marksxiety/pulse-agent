package components

import (
	"pulse-agent/theme"

	"github.com/charmbracelet/lipgloss"
)

var (
	ColorBg      lipgloss.Color
	ColorSurface lipgloss.Color
	ColorBorder  lipgloss.Color
	ColorMuted   lipgloss.Color
	ColorSubtle  lipgloss.Color
	ColorText    lipgloss.Color
	ColorDim     lipgloss.Color

	ColorCPUAccent  lipgloss.Color
	ColorMemAccent  lipgloss.Color
	ColorDiskAccent lipgloss.Color

	ColorWarn lipgloss.Color
	ColorDang lipgloss.Color

	SparklineRunes []rune
	ColorOverlay   lipgloss.Color
)

const (
	CardWidth     = 36
	CardGap       = 2
	MaxWidth      = 120
	OuterPadH     = 2
	OuterPadV     = 1
	CardBodyLines = 23
)

var (
	BaseCard    lipgloss.Style
	TitleStyle  lipgloss.Style
	LabelStyle  lipgloss.Style
	ValueStyle  lipgloss.Style
	DimStyle    lipgloss.Style
	HeaderStyle lipgloss.Style
	FooterStyle lipgloss.Style
)

func init() {
	ApplyTheme(OriginalTheme)
}

func ApplyTheme(t theme.Theme) {
	ColorBg, ColorSurface, ColorBorder, ColorMuted, ColorSubtle, ColorText, ColorDim,
		ColorCPUAccent, ColorMemAccent, ColorDiskAccent, ColorWarn, ColorDang = palette(t)

	ColorOverlay = lipgloss.Color("#06060a")

	BaseCard = lipgloss.NewStyle().
		Width(CardWidth).
		Padding(1, 2).
		Background(ColorSurface).
		BorderStyle(ActiveBorder()).
		BorderForeground(ColorBorder)

	TitleStyle = lipgloss.NewStyle().Bold(true).Foreground(ColorText)
	LabelStyle = lipgloss.NewStyle().Bold(true).Foreground(ColorText).Background(ColorSurface)
	ValueStyle = lipgloss.NewStyle().Bold(true).Foreground(ColorText).Background(ColorSurface)
	DimStyle = lipgloss.NewStyle().Foreground(ColorDim)
	HeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(ColorText).Padding(0, 1)
	FooterStyle = lipgloss.NewStyle().Foreground(ColorDim).Padding(0, 1)
}

// palette adapts a UI-free theme.Theme into lipgloss colours.
func palette(t theme.Theme) (bg, surface, border, muted, subtle, text, dim, cpu, mem, disk, warn, dang lipgloss.Color) {
	return lipgloss.Color(t.Bg), lipgloss.Color(t.Surface), lipgloss.Color(t.Border),
		lipgloss.Color(t.Muted), lipgloss.Color(t.Subtle), lipgloss.Color(t.Text), lipgloss.Color(t.Dim),
		lipgloss.Color(t.CPUAccent), lipgloss.Color(t.MemAccent), lipgloss.Color(t.DiskAccent),
		lipgloss.Color(t.Warn), lipgloss.Color(t.Danger)
}
