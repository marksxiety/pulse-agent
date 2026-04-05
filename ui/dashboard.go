package ui

import (
	"fmt"
	"pulse-agent/models"
	"pulse-agent/types"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Stats map[types.Source]float64
}

func InitialModel() Model {
	return Model{
		Stats: make(map[types.Source]float64),
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case models.NewDataMsg:
		// Update the map with the incoming metric
		m.Stats[msg.Source] = msg.Value
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	s := "PULSE AGENT 🖥️\n"
	s += "--------------------------\n"

	// Define order so the UI doesn't jump around
	order := []types.Source{types.CPU, types.MEM, types.DISK}

	for _, source := range order {
		val, ok := m.Stats[source]
		if !ok {
			s += fmt.Sprintf("%-5s: Initializing...\n", source)
			continue
		}
		s += fmt.Sprintf("%-5s: [%.2f%%]\n", source, val)
	}

	return s + "\n(press q to quit)"
}
