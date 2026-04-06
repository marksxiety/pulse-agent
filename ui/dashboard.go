package ui

import (
	"fmt"
	"pulse-agent/models"

	tea "github.com/charmbracelet/bubbletea"
)

func bytesToGB(b uint64) string {
	return fmt.Sprintf("%.2f GB", float64(b)/1e9)
}

type Model struct {
	CPU       models.CPUPayload
	Mem       models.MemoryPayload
	Disk      models.DiskPayload
	cpuReady  bool
	memReady  bool
	diskReady bool
}

func InitialModel() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	s := "PULSE AGENT\n"
	s += "--------------------------\n"
	// CPU
	if m.cpuReady {
		s += fmt.Sprintf("CPU:  [%.2f%%] (%d cores)\n", m.CPU.Percentage, m.CPU.CoreCount)
	} else {
		s += "CPU:  Initializing...\n"
	}
	// Memory
	if m.memReady {
		pct := float64(m.Mem.Used) / float64(m.Mem.Total) * 100
		s += fmt.Sprintf("MEM:  [%s/%s] [%.2f%%]\n",
			bytesToGB(m.Mem.Used), bytesToGB(m.Mem.Total), pct)
		s += fmt.Sprintf("      Pagefile/Swap: %s\n", bytesToGB(m.Mem.PagefileUsage))
	} else {
		s += "MEM:  Initializing...\n"
	}
	// Disk
	if m.diskReady {
		pct := float64(m.Disk.Used) / float64(m.Disk.Total) * 100
		s += fmt.Sprintf("DISK: [%s/%s] [%.2f%%]\n",
			bytesToGB(m.Disk.Used), bytesToGB(m.Disk.Total), pct)
		s += fmt.Sprintf("      I/O  R: %s | W: %s\n",
			bytesToGB(m.Disk.IOStats.ReadBytes), bytesToGB(m.Disk.IOStats.WriteBytes))
	} else {
		s += "DISK: Initializing...\n"
	}
	return s + "\n(press q to quit)"
}
