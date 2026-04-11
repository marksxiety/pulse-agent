package ui

import (
	"time"

	"pulse-agent/components"
	"pulse-agent/models"
	"pulse-agent/utils"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	CPU            models.CPUPayload
	Mem            models.MemoryPayload
	Disk           models.DiskPayload
	cpuReady       bool
	memReady       bool
	diskReady      bool
	startedAt      time.Time
	termW          int
	termH          int
	showInfo       bool
	modalScroll    int
	showQuitDialog bool
	quitCursor     int

	cpuHistory  *models.MetricHistory
	memHistory  *models.MetricHistory
	diskHistory *models.MetricHistory
	lastSample  time.Time

	showThemePicker bool
	themeCursor     int
}

func InitialModel() Model {
	cfg := utils.GetConfig()

	themeCursor := 0
	for i, name := range components.ThemeNames {
		if name == cfg.Theme {
			themeCursor = i
			break
		}
	}

	return Model{
		termW:       120,
		termH:       40,
		startedAt:   time.Now(),
		cpuHistory:  models.NewMetricHistory(),
		memHistory:  models.NewMetricHistory(),
		diskHistory: models.NewMetricHistory(),
		themeCursor: themeCursor,
	}
}

func (m Model) Init() tea.Cmd {
	cfg := utils.GetConfig()
	theme := components.GetTheme(cfg.Theme)
	components.ApplyTheme(theme)
	return nil
}

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
		if time.Since(m.lastSample) >= 5*time.Second {
			m.lastSample = time.Now()
			if m.cpuReady {
				m.cpuHistory.Push(m.CPU.Percentage)
			}
			if m.memReady && m.Mem.Total > 0 {
				m.memHistory.Push(float64(m.Mem.Used) / float64(m.Mem.Total) * 100)
			}
			if m.diskReady && m.Disk.Total > 0 {
				m.diskHistory.Push(float64(m.Disk.Used) / float64(m.Disk.Total) * 100)
			}
		}
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if m.showThemePicker {
			switch msg.Type {
			case tea.KeyEsc:
				m.showThemePicker = false
			case tea.KeyUp:
				if m.themeCursor > 0 {
					m.themeCursor--
				}
			case tea.KeyDown:
				if m.themeCursor < len(components.ThemeNames)-1 {
					m.themeCursor++
				}
			case tea.KeyRunes:
				switch msg.String() {
				case "k":
					if m.themeCursor > 0 {
						m.themeCursor--
					}
				case "j":
					if m.themeCursor < len(components.ThemeNames)-1 {
						m.themeCursor++
					}
				}
			case tea.KeyEnter:
				selected := components.ThemeNames[m.themeCursor]
				theme := components.GetTheme(selected)
				components.ApplyTheme(theme)
				utils.SaveTheme(selected)
				m.showThemePicker = false
			}
			return m, nil
		}

		switch msg.Type {
		case tea.KeyEsc:
			if m.showQuitDialog {
				m.showQuitDialog = false
				m.quitCursor = 0
				break
			}
			m.showInfo = false
		case tea.KeyUp:
			if m.showQuitDialog {
				if m.quitCursor > 0 {
					m.quitCursor--
				}
				break
			}
			if m.showInfo && m.modalScroll > 0 {
				m.modalScroll--
			}
		case tea.KeyDown:
			if m.showQuitDialog {
				if m.quitCursor < 1 {
					m.quitCursor++
				}
				break
			}
			if m.showInfo {
				m.modalScroll++
			}
		case tea.KeyF1:
			if m.showQuitDialog {
				break
			}
			m.showInfo = !m.showInfo
			if m.showInfo {
				m.modalScroll = 0
			}
		case tea.KeyRunes:
			if m.showQuitDialog {
				switch msg.String() {
				case "k":
					if m.quitCursor > 0 {
						m.quitCursor--
					}
				case "j":
					if m.quitCursor < 1 {
						m.quitCursor++
					}
				}
				break
			}
			if m.showInfo {
				break
			}
			for _, r := range msg.String() {
				if r == 'q' {
					m.showQuitDialog = true
					m.quitCursor = 0
					break
				}
				if r == 't' {
					m.showThemePicker = !m.showThemePicker
					break
				}
			}
		case tea.KeyEnter:
			if m.showQuitDialog {
				if m.quitCursor == 0 {
					return m, tea.Quit
				}
				m.showQuitDialog = false
				m.quitCursor = 0
				break
			}
			m.showInfo = false
		}
	}
	return m, nil
}

func currentThemeName() string {
	cfg := utils.GetConfig()
	return cfg.Theme
}
