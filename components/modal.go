package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	ModalWidth      = 76
	ModalHeight     = 28
	scrollBarW      = 1
	modalInnerW     = ModalWidth - 4 - scrollBarW
	modalPinnedRows = 2                                 // title row + blank divider, always visible
	modalBodyRows   = ModalHeight - 4 - modalPinnedRows // scrollable area height

	themeModalW        = 42
	themeModalH        = 16
	ThemeModalBodyRows = 8
)

func QuitConfirmModal(cursor int) string {
	dialogW := 42
	dialogH := 10

	title := TitleStyle.Render("  Quit?")
	closeHint := TitleStyle.Render("esc close")
	hGap := dialogW - 4 - lipgloss.Width(title) - lipgloss.Width(closeHint)
	if hGap < 1 {
		hGap = 1
	}
	header := title + strings.Repeat(" ", hGap) + closeHint

	sep := lipgloss.NewStyle().
		Foreground(ColorSubtle).
		Background(ColorSurface).
		Render(strings.Repeat("-", dialogW-4))

	body := lipgloss.NewStyle().Foreground(ColorText).Background(ColorSurface).Render("  Are you sure you want to quit?")

	options := []struct {
		label  string
		accent lipgloss.Color
	}{
		{"yes", ColorWarn},
		{"no", ColorSubtle},
	}

	var lines []string
	for i, opt := range options {
		if i == cursor {
			indicator := lipgloss.NewStyle().Foreground(opt.accent).Background(ColorSurface).Render("▸ ")
			selected := lipgloss.NewStyle().Foreground(opt.accent).Background(ColorSurface).Render(opt.label)
			lines = append(lines, indicator+selected)
		} else {
			entry := lipgloss.NewStyle().Foreground(ColorDim).Background(ColorSurface).Render("  " + opt.label)
			lines = append(lines, entry)
		}
	}

	choices := strings.Join(lines, "\n")

	inner := header + "\n" + sep + "\n" + body + "\n" + choices + "\n" + sep + "\n" +
		lipgloss.NewStyle().Foreground(ColorDim).Background(ColorSurface).Render("  ↑↓ navigate · enter select")

	borderStyle := lipgloss.NewStyle().
		Width(dialogW).
		Height(dialogH).
		Background(ColorSurface).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ColorSubtle).
		Padding(1, 2)

	return borderStyle.Render(inner)
}

func InfoModal(scroll int) string {
	trackStyle := lipgloss.NewStyle().Foreground(ColorMuted).Background(ColorSurface)

	sepLine := trackStyle.Render(strings.Repeat("-", modalInnerW))

	// Pinned header (never scrolls)
	title := TitleStyle.Render("  Metric Glossary")
	closeHint := TitleStyle.Render("F1/Esc close")
	hGap := modalInnerW - lipgloss.Width(title) - lipgloss.Width(closeHint)
	if hGap < 1 {
		hGap = 1
	}
	headerRow := title + strings.Repeat(" ", hGap) + closeHint
	// Pinned separator between header and scrollable body
	pinnedSep := lipgloss.NewStyle().
		Foreground(ColorSubtle).
		Background(ColorSurface).
		Render(strings.Repeat("-", modalInnerW))

	// Scrollable content lines (header is NOT in here)
	var lines []string

	lines = append(lines, section(modalInnerW, ColorCPUAccent, "CPU  PROCESSOR",
		[]metricEntry{
			{
				label: "Cores",
				desc:  "Total number of logical processors reported by the OS.",
				subs: []string{
					"Equals physical cores x threads per core (SMT / Hyper-Threading).",
					"Each logical core can execute one hardware thread at a time.",
				},
			},
			{
				label: "Usage (%)",
				desc:  "Percentage of total CPU time spent in non-idle state.",
				subs: []string{
					"Aggregated across ALL logical cores over the sample interval (2 s).",
					"100% means every core is fully saturated, not a single core.",
					"Calculated as: (1 - idle_time / total_time) x 100.",
					"Sustained > 80% indicates CPU-bound workload; > 95% = bottleneck.",
				},
			},
		},
	)...)

	lines = append(lines, "", sepLine, "")

	lines = append(lines, section(modalInnerW, ColorMemAccent, "MEM  RAM",
		[]metricEntry{
			{
				label: "Used",
				desc:  "Physical memory currently allocated to processes + kernel.",
				subs: []string{
					"Includes active pages (recently accessed) and standby / cached pages.",
					"Windows may report high usage even when standby is reclaimable.",
					"Rising trend without new processes suggests a memory leak.",
				},
			},
			{
				label: "Total",
				desc:  "Total physical RAM installed on the motherboard.",
				subs: []string{
					"Fixed value, does not change at runtime.",
					"Determines the upper ceiling for in-memory workloads.",
				},
			},
			{
				label: "Available",
				desc:  "Memory the OS can allocate without swapping to pagefile.",
				subs: []string{
					"= Total - Used (includes zero, free, and standby lists).",
					"Low available (< 500 MB) causes aggressive paging and thrashing.",
				},
			},
			{
				label: "Pagefile",
				desc:  "Virtual memory currently backed by the swap file on disk.",
				subs: []string{
					"When physical RAM is exhausted, Windows pages least-used memory here.",
					"Pagefile I/O is orders of magnitude slower than RAM (~100 ns vs ~10 ms).",
					"High pagefile usage correlates with poor application responsiveness.",
				},
			},
		},
	)...)

	lines = append(lines, "", sepLine, "")

	lines = append(lines, section(modalInnerW, ColorDiskAccent, "DISK  STORAGE",
		[]metricEntry{
			{
				label: "Used / Total / Available",
				desc:  "Disk space consumed, total capacity, and remaining free space.",
				subs: []string{
					"Measured across the primary system volume (usually C:\\).",
					"Does NOT include external drives or network mounts.",
					"Running below 10% free can fragment the MFT and degrade I/O.",
				},
			},
			{
				label: "I/O Reads / Writes",
				desc:  "Cumulative bytes transferred to/from disk since agent start.",
				subs: []string{
					"Reads = data fetched from storage into RAM.",
					"Writes = data flushed from RAM to persistent storage.",
					"Both include OS-level caching unless direct / unbuffered I/O is used.",
				},
			},
			{
				label: "Read Ops / Write Ops",
				desc:  "Total count of individual read and write system calls.",
				subs: []string{
					"High ops with low byte count = many small random I/O (worst for HDDs).",
					"Low ops with high byte count = fewer large sequential transfers (optimal).",
					"SSDs handle random I/O far better than spinning disks (no seek time).",
				},
			},
		},
	)...)

	lines = append(lines, "", sepLine, "")

	lines = append(lines, section(modalInnerW, ColorSubtle, "TREND  (3h)",
		[]metricEntry{
			{
				label: "Sparkline",
				desc:  "Unicode bar chart showing sampled values over the last 3 hours.",
				subs: []string{
					"Uses a circular buffer of 2160 samples (one every 5 s).",
					"Downsampled to fit the card width for rendering.",
					"Runes map to min-max range within the window.",
				},
			},
			{
				label: "Peak",
				desc:  "Highest recorded percentage in the 3-hour retention window.",
				subs: []string{
					"Useful for spotting transient spikes invisible to the current reading.",
					"Compared against usage % to gauge workload variability.",
				},
			},
			{
				label: "Progress Bar Thresholds",
				desc:  "Color-coded fill bar with dynamic thresholds.",
				subs: []string{
					"< 75%  -> card accent color (blue / green / purple).",
					"75-89% -> yellow warning, approaching capacity limit.",
					">= 90% -> red danger, critically high, immediate attention.",
				},
			},
		},
	)...)

	lines = append(lines, "", sepLine, "")

	lines = append(lines, section(modalInnerW, ColorSubtle, "KEYBOARD  SHORTCUTS",
		[]metricEntry{
			{
				label: "t",
				desc:  "Open the theme picker to switch color themes.",
				subs: []string{
					"Choose from 11 built-in themes (original, catppuccin-mocha, dracula, etc.).",
					"Your selection is saved and persists across restarts.",
				},
			},
			{
				label: "F1",
				desc:  "Toggle this info / glossary modal.",
				subs:  []string{"Press F1 or Esc to close."},
			},
			{
				label: "q",
				desc:  "Open the quit confirmation dialog.",
				subs:  []string{"Press q again or Enter on yes to quit."},
			},
		},
	)...)

	totalLines := len(lines)

	// Clamp scroll
	maxScroll := totalLines - modalBodyRows
	if maxScroll < 0 {
		maxScroll = 0
	}
	if scroll > maxScroll {
		scroll = maxScroll
	}
	if scroll < 0 {
		scroll = 0
	}

	end := scroll + modalBodyRows
	if end > totalLines {
		end = totalLines
	}
	visible := lines[scroll:end]

	// Scrollbar sized to the scrollable area only
	vScroll := verticalScrollbar(scroll, totalLines, modalBodyRows)

	// Render pinned header rows (no scrollbar char, just padded to full width)
	renderPinned := func(content string) string {
		w := lipgloss.Width(content)
		pad := modalInnerW - w
		if pad < 0 {
			pad = 0
		}
		// The scrollbar column on pinned rows shows nothing (space)
		return content + trackStyle.Render(strings.Repeat(" ", pad)) + trackStyle.Render(" ")
	}

	// Build final body: pinned rows first, then scrollable area
	totalRendered := modalPinnedRows + modalBodyRows
	allLines := make([]string, totalRendered)

	// Row 0: title / close hint (pinned)
	allLines[0] = renderPinned(headerRow)
	// Row 1: separator (pinned)
	allLines[1] = renderPinned(pinnedSep)

	// Rows 2+: scrollable content zipped with scrollbar
	for i := 0; i < modalBodyRows; i++ {
		var content string
		if i < len(visible) {
			content = visible[i]
		}

		contentW := lipgloss.Width(content)
		pad := modalInnerW - contentW
		if pad < 0 {
			pad = 0
		}
		paddedContent := content + trackStyle.Render(strings.Repeat(" ", pad))

		var bar string
		if i < len(vScroll) {
			bar = vScroll[i]
		} else {
			bar = trackStyle.Render(" ")
		}

		allLines[modalPinnedRows+i] = paddedContent + bar
	}

	body := strings.Join(allLines, "\n")

	borderStyle := lipgloss.NewStyle().
		Width(ModalWidth).
		Height(ModalHeight).
		Background(ColorSurface).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ColorSubtle).
		Padding(1, 2)

	return borderStyle.Render(body)
}

func verticalScrollbar(scroll, totalLines, visibleRows int) []string {
	trackStyle := lipgloss.NewStyle().Foreground(ColorMuted).Background(ColorSurface)
	thumbStyle := lipgloss.NewStyle().Foreground(ColorSubtle).Background(ColorSurface)

	gutter := make([]string, visibleRows)

	if totalLines <= visibleRows {
		for i := range gutter {
			gutter[i] = trackStyle.Render(" ")
		}
		return gutter
	}

	thumbH := visibleRows * visibleRows / totalLines
	if thumbH < 1 {
		thumbH = 1
	}

	maxScroll := totalLines - visibleRows
	thumbTop := 0
	if maxScroll > 0 {
		thumbTop = scroll * (visibleRows - thumbH) / maxScroll
	}

	thumbChar, trackChar := "#", "-"
	if termCap == TermNerdFont || termCap == TermUnicode {
		thumbChar, trackChar = "█", "░"
	}

	for i := range gutter {
		if i >= thumbTop && i < thumbTop+thumbH {
			gutter[i] = thumbStyle.Render(thumbChar)
		} else {
			gutter[i] = trackStyle.Render(trackChar)
		}
	}
	return gutter
}

type metricEntry struct {
	label string
	desc  string
	subs  []string
}

func section(w int, accent lipgloss.Color, heading string, entries []metricEntry) []string {
	icon := lipgloss.NewStyle().Foreground(accent).Bold(true).Background(ColorSurface).Render("●")
	head := lipgloss.NewStyle().Foreground(accent).Bold(true).Background(ColorSurface).Render(heading)
	var lines []string
	lines = append(lines, fmt.Sprintf(" %s %s", icon, head))

	descMax := w - 4
	subMax := w - 9

	for _, e := range entries {
		lbl := lipgloss.NewStyle().Bold(true).Foreground(ColorText).Background(ColorSurface).Render("  " + e.label)
		lines = append(lines, lbl)
		for _, dl := range wrapText(e.desc, descMax) {
			lines = append(lines,
				lipgloss.NewStyle().Foreground(ColorSubtle).Background(ColorSurface).Render("    "+dl))
		}
		for _, s := range e.subs {
			for _, sl := range wrapText(s, subMax) {
				lines = append(lines,
					lipgloss.NewStyle().Foreground(ColorDim).Background(ColorSurface).Render("      .  "+sl))
			}
		}
	}
	return lines
}

func ThemePickerModal(cursor, scroll int, currentTheme string) string {
	title := TitleStyle.Render("  Theme")
	closeHint := TitleStyle.Render("esc close")
	hGap := themeModalW - 4 - lipgloss.Width(title) - lipgloss.Width(closeHint)
	if hGap < 1 {
		hGap = 1
	}
	header := title + strings.Repeat(" ", hGap) + closeHint

	sep := lipgloss.NewStyle().
		Foreground(ColorSubtle).
		Background(ColorSurface).
		Render(strings.Repeat("-", themeModalW-4))

	total := len(ThemeNames)
	maxScroll := total - ThemeModalBodyRows
	if maxScroll < 0 {
		maxScroll = 0
	}
	if scroll > maxScroll {
		scroll = maxScroll
	}
	if scroll < 0 {
		scroll = 0
	}
	if cursor < scroll {
		scroll = cursor
	}
	if cursor >= scroll+ThemeModalBodyRows {
		scroll = cursor - ThemeModalBodyRows + 1
	}

	end := scroll + ThemeModalBodyRows
	if end > total {
		end = total
	}

	var lines []string
	for i := scroll; i < end; i++ {
		name := ThemeNames[i]
		suffix := ""
		if name == currentTheme {
			suffix = " ●"
		}
		if i == cursor {
			accent := lipgloss.Color(GetTheme(name).CPUAccent)
			indicator := lipgloss.NewStyle().Foreground(accent).Background(ColorSurface).Render("▸ ")
			selected := lipgloss.NewStyle().Foreground(accent).Background(ColorSurface).Render(name + suffix)
			lines = append(lines, indicator+selected)
		} else {
			entry := lipgloss.NewStyle().Foreground(ColorDim).Background(ColorSurface).Render("  " + name + suffix)
			lines = append(lines, entry)
		}
	}

	for len(lines) < ThemeModalBodyRows {
		lines = append(lines, strings.Repeat(" ", themeModalW-4))
	}

	scrollHint := ""
	if total > ThemeModalBodyRows {
		scrollHint = fmt.Sprintf(" [%d/%d]", scroll+1, maxScroll+1)
	}

	body := strings.Join(lines, "\n")

	inner := header + "\n" + sep + "\n" + body + "\n" + sep + "\n" +
		lipgloss.NewStyle().Foreground(ColorDim).Background(ColorSurface).Render("  ↑↓ navigate · enter select"+scrollHint)

	borderStyle := lipgloss.NewStyle().
		Width(themeModalW).
		Height(themeModalH).
		Background(ColorSurface).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ColorSubtle).
		Padding(1, 2)

	return borderStyle.Render(inner)
}

func wrapText(text string, maxLen int) []string {
	if maxLen < 1 {
		maxLen = 1
	}
	if lipgloss.Width(text) <= maxLen {
		return []string{text}
	}
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}
	var result []string
	var current strings.Builder
	currentLen := 0
	for _, word := range words {
		wordLen := lipgloss.Width(word)
		if currentLen == 0 {
			current.WriteString(word)
			currentLen = wordLen
		} else if currentLen+1+wordLen <= maxLen {
			current.WriteByte(' ')
			current.WriteString(word)
			currentLen += 1 + wordLen
		} else {
			result = append(result, current.String())
			current.Reset()
			current.WriteString(word)
			currentLen = wordLen
		}
	}
	if currentLen > 0 {
		result = append(result, current.String())
	}
	return result
}
