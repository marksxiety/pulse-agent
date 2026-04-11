package components

import (
	"fmt"
	"strings"

	"pulse-agent/models"
	"pulse-agent/utils"

	"github.com/charmbracelet/lipgloss"
)

func cardHeader(icon, title, badge string, accent lipgloss.Color) string {
	iconStr := lipgloss.NewStyle().Foreground(accent).Render(icon)
	titleStr := TitleStyle.Render(title)
	badgeStr := lipgloss.NewStyle().
		Foreground(ColorBg).Background(accent).
		Padding(0, 1).Bold(true).Render(badge)

	left := iconStr + titleStr
	gap := CardWidth - 4 - lipgloss.Width(left) - lipgloss.Width(badgeStr)
	if gap < 0 {
		gap = 0
	}
	return left + strings.Repeat(" ", gap) + badgeStr
}

func CPUCard(ready bool, cpu models.CPUPayload, history *models.MetricHistory) string {
	accent := ColorCPUAccent
	header := cardHeader(cardIcons.CPU, "CPU", "PROCESSOR", accent)
	sparkW := CardWidth - 4

	var lines []string
	if !ready {
		lines = []string{header, "", DimStyle.Render("  Waiting for data…")}
	} else {
		bar := ProgressBar(cpu.Percentage, CardWidth-4-8, accent)
		lines = []string{
			header,
			"",
			Separator(accent),
			"",
			bar,
			"",
			Row("Cores", fmt.Sprintf("%d", cpu.CoreCount)),
			Row("Usage", fmt.Sprintf("%.2f%%", cpu.Percentage)),
		}
		hist := HistorySection(history, sparkW, accent)
		for i := len(lines) + len(hist); i < CardBodyLines; i++ {
			lines = append(lines, "")
		}
		lines = append(lines, hist...)
	}

	lines = utils.PadToHeight(OverflowGuard(lines, CardBodyLines), CardBodyLines)
	bc := ColorBorder
	if ready {
		bc = accent
	}
	return BaseCard.BorderForeground(bc).Render(strings.Join(lines, "\n"))
}

func MemCard(ready bool, mem models.MemoryPayload, history *models.MetricHistory) string {
	accent := ColorMemAccent
	header := cardHeader(cardIcons.Mem, "Memory", "RAM", accent)
	sparkW := CardWidth - 4

	var lines []string
	if !ready {
		lines = []string{header, "", DimStyle.Render("  Waiting for data…")}
	} else {
		pct := float64(mem.Used) / float64(mem.Total) * 100
		bar := ProgressBar(pct, CardWidth-4-8, accent)
		lines = []string{
			header,
			"",
			Separator(accent),
			"",
			bar,
			"",
			Row("Used", utils.BytesToGB(mem.Used)),
			Row("Total", utils.BytesToGB(mem.Total)),
			Row("Available", utils.BytesToGB(mem.Available)),
			Row("Pagefile", utils.BytesToMB(mem.PagefileUsage)),
		}
		hist := HistorySection(history, sparkW, accent)
		for i := len(lines) + len(hist); i < CardBodyLines; i++ {
			lines = append(lines, "")
		}
		lines = append(lines, hist...)
	}

	lines = utils.PadToHeight(OverflowGuard(lines, CardBodyLines), CardBodyLines)
	bc := ColorBorder
	if ready {
		bc = accent
	}
	return BaseCard.BorderForeground(bc).Render(strings.Join(lines, "\n"))
}

func DiskCard(ready bool, disk models.DiskPayload, history *models.MetricHistory) string {
	accent := ColorDiskAccent
	header := cardHeader(cardIcons.Disk, "Disk", "STORAGE", accent)
	sparkW := CardWidth - 4

	var lines []string
	if !ready {
		lines = []string{header, "", DimStyle.Render("  Waiting for data…")}
	} else {
		pct := float64(disk.Used) / float64(disk.Total) * 100
		bar := ProgressBar(pct, CardWidth-4-8, accent)
		io := disk.IOStats
		lines = []string{
			header,
			"",
			Separator(accent),
			"",
			bar,
			"",
			Row("Used", utils.BytesToGB(disk.Used)),
			Row("Total", utils.BytesToGB(disk.Total)),
			Row("Available", utils.BytesToGB(disk.Available)),
			"",
			Separator(ColorMuted),
			"",
			Row("I/O Reads", utils.BytesToMB(io.ReadBytes)),
			Row("I/O Writes", utils.BytesToMB(io.WriteBytes)),
			Row("Read Ops", fmt.Sprintf("%d", io.ReadCount)),
			Row("Write Ops", fmt.Sprintf("%d", io.WriteCount)),
		}
		hist := HistorySection(history, sparkW, accent)
		for i := len(lines) + len(hist); i < CardBodyLines; i++ {
			lines = append(lines, "")
		}
		lines = append(lines, hist...)
	}

	lines = utils.PadToHeight(OverflowGuard(lines, CardBodyLines), CardBodyLines)
	bc := ColorBorder
	if ready {
		bc = accent
	}
	return BaseCard.BorderForeground(bc).Render(strings.Join(lines, "\n"))
}
