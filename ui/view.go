package ui

import (
	"fmt"
	"strings"
	"time"

	"pulse-agent/components"
	"pulse-agent/utils"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	cardsW := (components.CardWidth+2+2)*3 + components.CardGap*2

	minRequired := cardsW + components.OuterPadH*2
	if m.termW > 0 && m.termW < minRequired {
		return lipgloss.NewStyle().
			Foreground(components.ColorWarn).
			Padding(components.OuterPadV, components.OuterPadH).
			Render(fmt.Sprintf(
				"◆ PULSE AGENT\n\n"+
					"Terminal too narrow (%d cols).\n"+
					"Resize to at least %d cols.",
				m.termW, minRequired,
			))
	}

	effectiveW := m.termW
	if effectiveW > components.MaxWidth {
		effectiveW = components.MaxWidth
	}
	dashW := cardsW
	if dashW > effectiveW-components.OuterPadH*2 {
		dashW = effectiveW - components.OuterPadH*2
	}

	now := time.Now().Format("15:04:05")
	titleStr := lipgloss.NewStyle().Bold(true).Foreground(components.ColorCPUAccent).Render(components.CardIcons().Header+"PULSE") +
		lipgloss.NewStyle().Foreground(components.ColorSubtle).Render(" AGENT (F1) ")
	clockStr := lipgloss.NewStyle().Foreground(components.ColorDim).Render(now) +
		components.DimStyle.Render(fmt.Sprintf("  up %s", utils.FormatUptime(time.Since(m.startedAt))))

	hl := components.HeaderStyle.Render(titleStr)
	hr := components.FooterStyle.Render(clockStr)
	hGap := dashW - lipgloss.Width(hl) - lipgloss.Width(hr)
	if hGap < 0 {
		hGap = 0
	}
	header := components.CentreBlock(hl+strings.Repeat(" ", hGap)+hr, m.termW)

	gapStr := strings.Repeat(" ", components.CardGap)
	cards := lipgloss.JoinHorizontal(lipgloss.Top,
		components.CPUCard(m.cpuReady, m.CPU, m.cpuHistory),
		gapStr,
		components.MemCard(m.memReady, m.Mem, m.memHistory),
		gapStr,
		components.DiskCard(m.diskReady, m.Disk, m.diskHistory),
	)
	centeredCards := components.CentreBlock(cards, m.termW)

	quitStr := components.FooterStyle.Render("  q  quit")
	themeStr := components.FooterStyle.Render(fmt.Sprintf("  t  theme [%s]", currentThemeName()))
	versionStr := components.DimStyle.Render(utils.Version)
	fGap := dashW - lipgloss.Width(quitStr) - lipgloss.Width(themeStr) - lipgloss.Width(versionStr)
	if fGap < 0 {
		fGap = 0
	}
	footer := components.CentreBlock(quitStr+strings.Repeat(" ", fGap/2)+themeStr+strings.Repeat(" ", fGap-fGap/2)+versionStr, m.termW)

	headerLines := 1
	cardBlockLines := lipgloss.Height(centeredCards)
	footerLines := 1

	usedLines := headerLines + 1 + cardBlockLines + 1 + footerLines

	remaining := m.termH - usedLines
	if remaining < 0 {
		remaining = 0
	}
	topPad := remaining / 2
	bottomPad := remaining - topPad

	top := strings.Repeat("\n", topPad)
	mid := strings.Repeat("\n", bottomPad)

	content := top + header + "\n\n" + centeredCards + "\n" + mid + footer

	// Wrap in terminal dimensions so bubbletea clears the full frame on resize.
	frame := lipgloss.NewStyle().Width(m.termW).Height(m.termH).Render(content)

	// Modal overlay — rendered AFTER the background frame using lipgloss.Place.
	// lipgloss.Place understands ANSI codes and positions correctly, unlike
	// raw string slicing which corrupts escape sequences.
	if m.showInfo {
		modal := components.InfoModal(m.modalScroll)
		frame = lipgloss.Place(
			m.termW, m.termH,
			lipgloss.Center, lipgloss.Center,
			modal,
			lipgloss.WithWhitespaceBackground(components.ColorOverlay),
		)
	}

	if m.showQuitDialog {
		dialog := components.QuitConfirmModal()
		frame = lipgloss.Place(
			m.termW, m.termH,
			lipgloss.Center, lipgloss.Center,
			dialog,
			lipgloss.WithWhitespaceBackground(components.ColorOverlay),
		)
	}

	if m.showThemePicker {
		modal := components.ThemePickerModal(m.themeCursor, currentThemeName())
		frame = lipgloss.Place(
			m.termW, m.termH,
			lipgloss.Center, lipgloss.Center,
			modal,
			lipgloss.WithWhitespaceBackground(components.ColorOverlay),
		)
	}

	return frame
}
