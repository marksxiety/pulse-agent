package components

import "github.com/charmbracelet/lipgloss"

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

	SparklineRunes = []rune("▁▂▃▄▅▆▇█")
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
	ApplyTheme(CatppuccinMochaTheme)
}

func ApplyTheme(t Theme) {
	ColorBg, ColorSurface, ColorBorder, ColorMuted, ColorSubtle, ColorText, ColorDim,
		ColorCPUAccent, ColorMemAccent, ColorDiskAccent, ColorWarn, ColorDang = t.Colors()

	BaseCard = lipgloss.NewStyle().
		Width(CardWidth).
		Padding(1, 2).
		Background(ColorSurface).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder)

	TitleStyle = lipgloss.NewStyle().Bold(true).Foreground(ColorText)
	LabelStyle = lipgloss.NewStyle().Bold(true).Foreground(ColorText).Background(ColorSurface)
	ValueStyle = lipgloss.NewStyle().Bold(true).Foreground(ColorText).Background(ColorSurface)
	DimStyle = lipgloss.NewStyle().Foreground(ColorDim)
	HeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(ColorText).Padding(0, 1)
	FooterStyle = lipgloss.NewStyle().Foreground(ColorDim).Padding(0, 1)
}
