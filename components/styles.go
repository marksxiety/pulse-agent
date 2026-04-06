package components

import "github.com/charmbracelet/lipgloss"

var (
	ColorBg      = lipgloss.Color("#0d0f14")
	ColorSurface = lipgloss.Color("#141720")
	ColorBorder  = lipgloss.Color("#1e2330")
	ColorMuted   = lipgloss.Color("#3b4260")
	ColorSubtle  = lipgloss.Color("#6272a4")
	ColorText    = lipgloss.Color("#cdd6f4")
	ColorDim     = lipgloss.Color("#585b70")

	ColorCPUAccent  = lipgloss.Color("#89b4fa")
	ColorMemAccent  = lipgloss.Color("#a6e3a1")
	ColorDiskAccent = lipgloss.Color("#cba6f7")

	ColorWarn = lipgloss.Color("#f9e2af")
	ColorDang = lipgloss.Color("#f38ba8")

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
	BaseCard = lipgloss.NewStyle().
			Width(CardWidth).
			Padding(1, 2).
			Background(ColorSurface).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder)

	TitleStyle  = lipgloss.NewStyle().Bold(true).Foreground(ColorText)
	LabelStyle  = lipgloss.NewStyle().Bold(true).Foreground(ColorText).Background(ColorSurface)
	ValueStyle  = lipgloss.NewStyle().Bold(true).Foreground(ColorText).Background(ColorSurface)
	DimStyle    = lipgloss.NewStyle().Foreground(ColorDim)
	HeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(ColorText).Padding(0, 1)
	FooterStyle = lipgloss.NewStyle().Foreground(ColorDim).Padding(0, 1)
)
