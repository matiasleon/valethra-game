package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	transitionTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("205")).
				MarginBottom(1)

	transitionTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))
)

type TransitionModel struct {
	title   string
	message string
	width   int
	height  int
}

func NewTransitionModel(title, message string) TransitionModel {
	return TransitionModel{
		title:   title,
		message: message,
	}
}

func (m TransitionModel) Init() tea.Cmd {
	return nil
}

func (m TransitionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", " ", "enter":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m TransitionModel) View() string {
	var b strings.Builder

	b.WriteString(transitionTitleStyle.Render(fmt.Sprintf("=== %s ===", m.title)))
	b.WriteString("\n\n")
	maxWidth := m.width - 4
	if maxWidth < 40 {
		maxWidth = 80
	}
	b.WriteString(transitionTextStyle.Width(maxWidth).Render(m.message))
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("Presiona cualquier tecla para continuar."))

	return b.String()
}
