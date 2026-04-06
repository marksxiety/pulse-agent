package ui

import (
	"fmt"
	"math"
	"strings"
	"time"

	"pulse-agent/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ── Palette ───────────────────────────────────────────────────────────────────

var (
	colorBg      = lipgloss.Color("#0d0f14")
	colorSurface = lipgloss.Color("#141720")
	colorBorder  = lipgloss.Color("#1e2330")
	colorMuted   = lipgloss.Color("#3b4260")
	colorSubtle  = lipgloss.Color("#6272a4")
	colorText    = lipgloss.Color("#cdd6f4")
	colorDim     = lipgloss.Color("#585b70")

	colorCPUAccent  = lipgloss.Color("#89b4fa")
	colorMemAccent  = lipgloss.Color("#a6e3a1")
	colorDiskAccent = lipgloss.Color("#cba6f7")

	colorWarn = lipgloss.Color("#f9e2af")
	colorDang = lipgloss.Color("#f38ba8")
)

// ── Layout constants ──────────────────────────────────────────────────────────

const (
	cardWidth = 36  // inner content width per card
	cardGap   = 2   // spaces between cards
	maxWidth  = 120 // max dashboard width; mirrors CSS max-width
	outerPadH = 2   // left/right outer padding (cols)
	outerPadV = 1   // top/bottom outer padding (lines)

	// All cards are padded to this exact line count so they share one height.
	// Tallest card is Disk with 16 body lines — CPU/Mem get blank lines appended.
	cardBodyLines = 16
)

// ── Styles ────────────────────────────────────────────────────────────────────

var (
	baseCard = lipgloss.NewStyle().
			Width(cardWidth).
			Padding(1, 2).
			Background(colorSurface).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder)

	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(colorText)
	labelStyle  = lipgloss.NewStyle().Foreground(colorSubtle)
	valueStyle  = lipgloss.NewStyle().Bold(true).Foreground(colorText)
	dimStyle    = lipgloss.NewStyle().Foreground(colorDim)
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(colorText).Padding(0, 1)
	footerStyle = lipgloss.NewStyle().Foreground(colorDim).Padding(0, 1)
)

// ── Helpers ───────────────────────────────────────────────────────────────────

func bytesToGB(b uint64) string { return fmt.Sprintf("%.2f GB", float64(b)/1e9) }
func bytesToMB(b uint64) string { return fmt.Sprintf("%.1f MB", float64(b)/1e6) }

func formatUptime(d time.Duration) string {
	d = d.Round(time.Second)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %02dm %02ds", h, m, s)
	}
	return fmt.Sprintf("%dm %02ds", m, s)
}

func progressBar(pct float64, width int, accent lipgloss.Color) string {
	pct = math.Max(0, math.Min(100, pct))
	filled := int(math.Round(pct / 100 * float64(width)))
	empty := width - filled

	barColor := accent
	if pct >= 90 {
		barColor = colorDang
	} else if pct >= 75 {
		barColor = colorWarn
	}

	bar := lipgloss.NewStyle().Foreground(barColor).Render(strings.Repeat("█", filled)) +
		lipgloss.NewStyle().Foreground(colorMuted).Render(strings.Repeat("░", empty))
	pctStr := lipgloss.NewStyle().Foreground(barColor).Bold(true).Render(fmt.Sprintf("%5.1f%%", pct))
	return bar + " " + pctStr
}

// row renders a label + right-aligned value within the card's inner width.
func row(label, val string) string {
	inner := cardWidth - 4
	l := valueStyle.Render(label)
	v := valueStyle.Render(val)
	gap := inner - lipgloss.Width(l) - lipgloss.Width(v)
	if gap < 1 {
		gap = 1
	}
	return l + strings.Repeat(" ", gap) + v
}

func separator(c lipgloss.Color) string {
	return lipgloss.NewStyle().Foreground(c).Render(strings.Repeat("─", cardWidth-4))
}

func cardHeader(icon, title, badge string, accent lipgloss.Color) string {
	iconStr := lipgloss.NewStyle().Foreground(accent).Render(icon)
	titleStr := titleStyle.Render(title)
	badgeStr := lipgloss.NewStyle().
		Foreground(colorBg).Background(accent).
		Padding(0, 1).Bold(true).Render(badge)

	left := iconStr + titleStr
	gap := cardWidth - 4 - lipgloss.Width(left) - lipgloss.Width(badgeStr)
	if gap < 0 {
		gap = 0
	}
	return left + strings.Repeat(" ", gap) + badgeStr
}

// padToHeight appends blank lines until the slice reaches exactly n lines.
// This is the TUI equivalent of CSS `align-items: stretch` / equal card height.
func padToHeight(lines []string, n int) []string {
	for len(lines) < n {
		lines = append(lines, "")
	}
	return lines
}

// overflowGuard clips lines that exceed maxLines and shows a count indicator.
func overflowGuard(lines []string, maxLines int) []string {
	if len(lines) <= maxLines {
		return lines
	}
	hidden := len(lines) - (maxLines - 1)
	clipped := lines[:maxLines-1]
	return append(clipped, dimStyle.Render(fmt.Sprintf("  ↕  %d lines hidden", hidden)))
}

// centreBlock centres a block of text horizontally within termW by applying
// a left margin via lipgloss (safe for bubbletea's frame diffing).
func centreBlock(s string, termW int) string {
	blockW := lipgloss.Width(strings.SplitN(s, "\n", 2)[0])
	margin := (termW - blockW) / 2
	if margin < 0 {
		margin = 0
	}
	return lipgloss.NewStyle().MarginLeft(margin).Render(s)
}

// ── Card builders ─────────────────────────────────────────────────────────────

func cpuCard(m Model) string {
	accent := colorCPUAccent
	header := cardHeader("󰻠 ", "CPU", "PROCESSOR", accent)

	var lines []string
	if !m.cpuReady {
		lines = []string{header, "", dimStyle.Render("  Waiting for data…")}
	} else {
		bar := progressBar(m.CPU.Percentage, cardWidth-4-8, accent)
		lines = []string{
			header,
			"",
			separator(accent),
			"",
			bar,
			"",
			row("Cores", fmt.Sprintf("%d", m.CPU.CoreCount)),
			row("Usage", fmt.Sprintf("%.2f%%", m.CPU.Percentage)),
		}
	}

	lines = padToHeight(overflowGuard(lines, cardBodyLines), cardBodyLines)
	bc := colorBorder
	if m.cpuReady {
		bc = accent
	}
	return baseCard.BorderForeground(bc).Render(strings.Join(lines, "\n"))
}

func memCard(m Model) string {
	accent := colorMemAccent
	header := cardHeader("󰍛 ", "Memory", "RAM", accent)

	var lines []string
	if !m.memReady {
		lines = []string{header, "", dimStyle.Render("  Waiting for data…")}
	} else {
		pct := float64(m.Mem.Used) / float64(m.Mem.Total) * 100
		bar := progressBar(pct, cardWidth-4-8, accent)
		lines = []string{
			header,
			"",
			separator(accent),
			"",
			bar,
			"",
			row("Used", bytesToGB(m.Mem.Used)),
			row("Total", bytesToGB(m.Mem.Total)),
			row("Available", bytesToGB(m.Mem.Available)),
			row("Pagefile", bytesToMB(m.Mem.PagefileUsage)),
		}
	}

	lines = padToHeight(overflowGuard(lines, cardBodyLines), cardBodyLines)
	bc := colorBorder
	if m.memReady {
		bc = accent
	}
	return baseCard.BorderForeground(bc).Render(strings.Join(lines, "\n"))
}

func diskCard(m Model) string {
	accent := colorDiskAccent
	header := cardHeader("󰋊 ", "Disk", "STORAGE", accent)

	var lines []string
	if !m.diskReady {
		lines = []string{header, "", dimStyle.Render("  Waiting for data…")}
	} else {
		pct := float64(m.Disk.Used) / float64(m.Disk.Total) * 100
		bar := progressBar(pct, cardWidth-4-8, accent)
		io := m.Disk.IOStats
		lines = []string{
			header,
			"",
			separator(accent),
			"",
			bar,
			"",
			row("Used", bytesToGB(m.Disk.Used)),
			row("Total", bytesToGB(m.Disk.Total)),
			row("Available", bytesToGB(m.Disk.Available)),
			"",
			separator(colorMuted),
			"",
			row("I/O Reads", bytesToMB(io.ReadBytes)),
			row("I/O Writes", bytesToMB(io.WriteBytes)),
			row("Read Ops", fmt.Sprintf("%d", io.ReadCount)),
			row("Write Ops", fmt.Sprintf("%d", io.WriteCount)),
		}
	}

	lines = padToHeight(overflowGuard(lines, cardBodyLines), cardBodyLines)
	bc := colorBorder
	if m.diskReady {
		bc = accent
	}
	return baseCard.BorderForeground(bc).Render(strings.Join(lines, "\n"))
}

// ── Model ─────────────────────────────────────────────────────────────────────

type Model struct {
	CPU       models.CPUPayload
	Mem       models.MemoryPayload
	Disk      models.DiskPayload
	cpuReady  bool
	memReady  bool
	diskReady bool
	startedAt time.Time
	termW     int
	termH     int
}

func InitialModel() Model {
	return Model{termW: 120, termH: 40, startedAt: time.Now()}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termW = msg.Width
		m.termH = msg.Height
	case models.NewDataMsg:
		switch p := msg.Data.(type) {
		case models.CPUPayload:
			m.CPU = p
			m.cpuReady = true
		case models.MemoryPayload:
			m.Mem = p
			m.memReady = true
		case models.DiskPayload:
			m.Disk = p
			m.diskReady = true
		}
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

// ── View ──────────────────────────────────────────────────────────────────────

func (m Model) View() string {
	// Total rendered width of the 3-card row:
	// each card = cardWidth + border(2) + padding(4) = cardWidth+6 cols
	// but lipgloss Width(card) counts the inner+border, so it's cardWidth+2.
	// Plus 2 gap strings between cards.
	cardsW := (cardWidth+2+2)*3 + cardGap*2

	// ── Overflow guard: terminal narrower than minimum? show fallback ──
	minRequired := cardsW + outerPadH*2
	if m.termW > 0 && m.termW < minRequired {
		return lipgloss.NewStyle().
			Foreground(colorWarn).
			Padding(outerPadV, outerPadH).
			Render(fmt.Sprintf(
				"◆ PULSE AGENT\n\n"+
					"Terminal too narrow (%d cols).\n"+
					"Resize to at least %d cols.",
				m.termW, minRequired,
			))
	}

	// dashW = the content span for header/footer (capped at maxWidth).
	effectiveW := m.termW
	if effectiveW > maxWidth {
		effectiveW = maxWidth
	}
	dashW := cardsW
	if dashW > effectiveW-outerPadH*2 {
		dashW = effectiveW - outerPadH*2
	}

	// ── Header ──
	now := time.Now().Format("15:04:05")
	titleStr := lipgloss.NewStyle().Bold(true).Foreground(colorCPUAccent).Render("◆ PULSE") +
		lipgloss.NewStyle().Foreground(colorSubtle).Render(" AGENT")
	clockStr := lipgloss.NewStyle().Foreground(colorDim).Render(now) +
		dimStyle.Render(fmt.Sprintf("  up %s", formatUptime(time.Since(m.startedAt))))

	hl := headerStyle.Render(titleStr)
	hr := footerStyle.Render(clockStr)
	hGap := dashW - lipgloss.Width(hl) - lipgloss.Width(hr)
	if hGap < 0 {
		hGap = 0
	}
	header := centreBlock(hl+strings.Repeat(" ", hGap)+hr, m.termW)

	// ── Cards ──
	gapStr := strings.Repeat(" ", cardGap)
	cards := lipgloss.JoinHorizontal(lipgloss.Top,
		cpuCard(m),
		gapStr,
		memCard(m),
		gapStr,
		diskCard(m),
	)
	centeredCards := centreBlock(cards, m.termW)

	// ── Footer (fixed to bottom) ──
	quitStr := footerStyle.Render("  q  quit")
	versionStr := dimStyle.Render("pulse-agent v0.1")
	fGap := dashW - lipgloss.Width(quitStr) - lipgloss.Width(versionStr)
	if fGap < 0 {
		fGap = 0
	}
	footer := centreBlock(quitStr+strings.Repeat(" ", fGap)+versionStr, m.termW)

	// ── Vertical centering ──
	// Measure how many lines the top section (header + blank + cards) occupies.
	headerLines := 1
	cardBlockLines := lipgloss.Height(centeredCards)
	footerLines := 1

	// Total lines consumed by header + gap + cards + gap + footer
	usedLines := headerLines + 1 + cardBlockLines + 1 + footerLines

	// Remaining space is split: half above header, half below cards (above footer).
	remaining := m.termH - usedLines
	if remaining < 0 {
		remaining = 0
	}
	topPad := remaining / 2
	// bottomPad fills the gap between cards and the pinned footer.
	bottomPad := remaining - topPad

	top := strings.Repeat("\n", topPad)
	mid := strings.Repeat("\n", bottomPad)

	// Assemble: top padding → header → cards → elastic middle → pinned footer.
	content := top + header + "\n\n" + centeredCards + "\n" + mid + footer

	// Wrap in terminal dimensions so bubbletea diffs cleanly on every resize.
	return lipgloss.NewStyle().
		Width(m.termW).
		Height(m.termH).
		Render(content)
}
