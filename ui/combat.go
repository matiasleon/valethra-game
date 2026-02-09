// Package ui provides terminal UI models and views.
package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"saturday-chill/core/entities"
	"saturday-chill/core/ports"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			PaddingTop(1).
			MarginBottom(1)

	heroStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("82"))

	enemyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("196"))

	healthBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("82"))

	armorBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39"))

	damageBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	logStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			PaddingLeft(2).
			MarginTop(0)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)

	victoryStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("82")).
			MarginTop(1)

	defeatStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("196")).
			MarginTop(1)

	eventVictoryStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("226")).
				MarginTop(1)
)

type CombatModel struct {
	quest      entities.Quest
	hero       *entities.Character
	enemy      *entities.Character
	heroMax    entities.Attributes
	enemyMax   entities.Attributes
	round      int
	combatLog  []string
	state      ports.CombatState
	eventIndex int
	session    ports.CombatSession
	width      int
	height     int
}

func NewCombatModel(session ports.CombatSession) CombatModel {
	model := CombatModel{session: session}
	return model.applyState(session.State())
}

func (m CombatModel) HeroWon() bool {
	return m.state == ports.StateQuestVictory
}

func (m CombatModel) Init() tea.Cmd {
	return nil
}

func (m CombatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case " ", "enter", "a":
			switch m.state {
			case ports.StateQuestVictory, ports.StateDefeat:
				return m, tea.Quit
			default:
				return m.applyState(m.session.Handle(ports.ActionConfirm)), nil
			}
		}
	}
	return m, nil
}

func (m CombatModel) applyState(state ports.CombatStateDTO) CombatModel {
	m.quest = state.Quest
	m.hero = state.Hero
	m.enemy = state.Enemy
	m.heroMax = state.HeroMax
	m.enemyMax = state.EnemyMax
	m.round = state.Round
	m.state = state.State
	m.combatLog = state.CombatLog
	m.eventIndex = state.EventIndex
	return m
}

func (m CombatModel) View() string {
	var b strings.Builder

	title := fmt.Sprintf("=== %s ===", m.quest.Title)
	if len(m.quest.Events) > 1 {
		title += fmt.Sprintf("  [Encuentro %d/%d]", m.eventIndex+1, len(m.quest.Events))
	}
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n\n")

	if m.eventIndex < len(m.quest.Events) {
		desc := m.quest.Events[m.eventIndex].Description
		maxWidth := m.width - 4
		if maxWidth < 40 {
			maxWidth = 80
		}
		b.WriteString(lipgloss.NewStyle().Width(maxWidth).Render(desc))
		b.WriteString("\n\n")
	}

	b.WriteString(m.renderCharacter(m.hero, &m.heroMax, heroStyle, "HÉROE"))
	b.WriteString("\n\n")
	b.WriteString(m.renderCharacter(m.enemy, &m.enemyMax, enemyStyle, "ENEMIGO"))
	b.WriteString("\n")

	logStart := 0
	if len(m.combatLog) > 6 {
		logStart = len(m.combatLog) - 6
	}
	for _, log := range m.combatLog[logStart:] {
		b.WriteString(logStyle.Render(log))
		b.WriteString("\n")
	}

	switch m.state {
	case ports.StateQuestVictory:
		b.WriteString("\n")
		b.WriteString(victoryStyle.Render("¡VICTORIA! " + m.hero.Name + " ha completado la misión."))
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("Presiona cualquier tecla para continuar."))
	case ports.StateEventVictory:
		b.WriteString("\n")
		b.WriteString(eventVictoryStyle.Render(m.enemy.Name + " derrotado! Pero el camino continúa..."))
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("Presiona cualquier tecla para el siguiente encuentro."))
	case ports.StateDefeat:
		b.WriteString("\n")
		b.WriteString(defeatStyle.Render("DERROTA. " + m.hero.Name + " ha caído."))
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("Presiona cualquier tecla para salir."))
	case ports.StateFighting:
		b.WriteString(helpStyle.Render("[A/Enter/Space] Atacar • [Q] Salir"))
	}

	return b.String()
}

func (m CombatModel) renderCharacter(c *entities.Character, max *entities.Attributes, style lipgloss.Style, label string) string {
	var b strings.Builder

	b.WriteString(style.Render(fmt.Sprintf("%s: %s", label, c.Name)))
	b.WriteString("\n")

	healthPct := float64(c.Attributes.Health) / float64(max.Health)
	armorPct := float64(c.Attributes.Armor) / float64(max.Armor)
	if max.Armor == 0 {
		armorPct = 0
	}

	b.WriteString(fmt.Sprintf("  HP:    %s %d/%d\n",
		renderBar(healthPct, 20, healthBarStyle),
		c.Attributes.Health, max.Health))

	b.WriteString(fmt.Sprintf("  Armor: %s %d/%d\n",
		renderBar(armorPct, 20, armorBarStyle),
		c.Attributes.Armor, max.Armor))

	b.WriteString(fmt.Sprintf("  ATK: %d", c.Attributes.AttackPower))

	return b.String()
}

func renderBar(pct float64, width int, style lipgloss.Style) string {
	if pct < 0 {
		pct = 0
	}
	if pct > 1 {
		pct = 1
	}

	filled := int(pct * float64(width))
	empty := width - filled

	bar := style.Render(strings.Repeat("█", filled)) +
		damageBarStyle.Render(strings.Repeat("░", empty))

	return "[" + bar + "]"
}
