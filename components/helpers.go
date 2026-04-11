package components

import (
	"fmt"
	"strings"

	"pulse-agent/models"

	"github.com/charmbracelet/lipgloss"
)

func HistorySection(h *models.MetricHistory, sparkWidth int, accent lipgloss.Color) []string {
	if h.Len() < 2 {
		return []string{
			"",
			Separator(ColorMuted),
			"",
			DimStyle.Render("  Collecting data…"),
		}
	}

	sampled := h.Downsample(sparkWidth)
	spark := Sparkline(sampled, accent)

	return []string{
		"",
		Separator(ColorMuted),
		"",
		lipgloss.NewStyle().Foreground(ColorSubtle).Bold(true).Background(ColorSurface).Render("  Trend (3h)"),
		spark,
		"",
		lipgloss.NewStyle().Foreground(ColorSubtle).Background(ColorSurface).Render(
			fmt.Sprintf("  Peak %5.1f%%", h.Max())),
	}
}

func Row(label, val string) string {
	inner := CardWidth - 4
	l := LabelStyle.Render(label)
	v := ValueStyle.Render(val)
	gap := inner - lipgloss.Width(l) - lipgloss.Width(v)
	if gap < 1 {
		gap = 1
	}
	gapStr := lipgloss.NewStyle().Background(ColorSurface).Render(strings.Repeat(" ", gap))
	return l + gapStr + v
}

func Separator(c lipgloss.Color) string {
	ch := "-"
	if termCap == TermNerdFont || termCap == TermUnicode {
		ch = "─"
	}
	return lipgloss.NewStyle().Foreground(c).Render(strings.Repeat(ch, CardWidth-4))
}

func OverflowGuard(lines []string, maxLines int) []string {
	if len(lines) <= maxLines {
		return lines
	}
	hidden := len(lines) - (maxLines - 1)
	arrow := "v"
	if termCap == TermNerdFont || termCap == TermUnicode {
		arrow = "↕"
	}
	clipped := lines[:maxLines-1]
	return append(clipped, DimStyle.Render(fmt.Sprintf("  %s  %d lines hidden", arrow, hidden)))
}
