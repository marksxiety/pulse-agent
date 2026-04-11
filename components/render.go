package components

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func ProgressBar(pct float64, width int, accent lipgloss.Color) string {
	pct = math.Max(0, math.Min(100, pct))
	filled := int(math.Round(pct / 100 * float64(width)))
	empty := width - filled

	barColor := accent
	if pct >= 90 {
		barColor = ColorDang
	} else if pct >= 75 {
		barColor = ColorWarn
	}

	barFilled, barEmpty := "#", "-"
	if termCap == TermNerdFont || termCap == TermUnicode {
		barFilled, barEmpty = "█", "░"
	}

	bar := lipgloss.NewStyle().Foreground(barColor).Render(strings.Repeat(barFilled, filled)) +
		lipgloss.NewStyle().Foreground(ColorMuted).Render(strings.Repeat(barEmpty, empty))
	pctStr := lipgloss.NewStyle().Foreground(barColor).Bold(true).Render(fmt.Sprintf("%5.1f%%", pct))
	return bar + " " + pctStr
}

func Sparkline(values []float64, accent lipgloss.Color) string {
	if len(values) == 0 {
		return ""
	}

	min, max := values[0], values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	span := max - min
	if span == 0 {
		return lipgloss.NewStyle().Foreground(accent).Render(
			strings.Repeat(string(SparklineRunes[4]), len(values)))
	}

	levels := len(SparklineRunes) - 1
	var sb strings.Builder
	for _, v := range values {
		idx := int((v-min)/span*float64(levels) + 0.5)
		if idx > levels {
			idx = levels
		}
		sb.WriteRune(SparklineRunes[idx])
	}

	return lipgloss.NewStyle().Foreground(accent).Render(sb.String())
}

func CentreBlock(s string, termW int) string {
	blockW := lipgloss.Width(strings.SplitN(s, "\n", 2)[0])
	margin := (termW - blockW) / 2
	if margin < 0 {
		margin = 0
	}
	return lipgloss.NewStyle().MarginLeft(margin).Render(s)
}
