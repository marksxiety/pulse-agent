package ui

import (
	"time"

	"pulse-agent/models"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	CPU         models.CPUPayload
	Mem         models.MemoryPayload
	Disk        models.DiskPayload
	cpuReady    bool
	memReady    bool
	diskReady   bool
	startedAt   time.Time
	termW       int
	termH       int
	showInfo    bool
	modalScroll int

	cpuHistory  *models.MetricHistory
	memHistory  *models.MetricHistory
	diskHistory *models.MetricHistory
	lastSample  time.Time
}

func InitialModel() Model {
	return Model{
		termW:       120,
		termH:       40,
		startedAt:   time.Now(),
		cpuHistory:  models.NewMetricHistory(),
		memHistory:  models.NewMetricHistory(),
		diskHistory: models.NewMetricHistory(),
	}
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
		// Always let ctrl+c / q quit regardless of modal state.
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		// Use msg.Type for special keys — terminals send escape sequences,
		// not literal strings like "f1". msg.Type is what bubbletea parses
		// those sequences into.
		switch msg.Type {
		case tea.KeyF1:
			m.showInfo = !m.showInfo
			if m.showInfo {
				m.modalScroll = 0
			}
		case tea.KeyEsc, tea.KeyEnter:
			m.showInfo = false
		case tea.KeyUp:
			if m.showInfo && m.modalScroll > 0 {
				m.modalScroll--
			}
		case tea.KeyDown:
			if m.showInfo {
				m.modalScroll++
			}
		}
	}
	return m, nil
}
