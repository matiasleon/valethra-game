package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"saturday-chill/core/entities"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
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
			MarginTop(1)

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
)

type CombatModel struct {
	quest      entities.Quest
	hero       *entities.Character
	enemy      *entities.Character
	heroMax    entities.Attributes
	enemyMax   entities.Attributes
	round      int
	combatLog  []string
	gameOver   bool
	heroWon    bool
}

func NewCombatModel(quest entities.Quest, hero *entities.Character, enemy *entities.Character) CombatModel {
	return CombatModel{
		quest:     quest,
		hero:      hero,
		enemy:     enemy,
		heroMax:   hero.Attributes,
		enemyMax:  enemy.Attributes,
		round:     1,
		combatLog: []string{"El combate comienza..."},
		gameOver:  false,
		heroWon:   false,
	}
}

func (m CombatModel) Init() tea.Cmd {
	return nil
}

func (m CombatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case " ", "enter", "a":
			if m.gameOver {
				return m, tea.Quit
			}
			return m.executeRound(), nil
		}
	}
	return m, nil
}

func (m CombatModel) executeRound() CombatModel {
	m.combatLog = append(m.combatLog, fmt.Sprintf("-- Ronda %d --", m.round))

	result := m.hero.Attack(m.enemy)
	m.combatLog = append(m.combatLog, cleanLog(result))

	if m.enemy.Attributes.Health <= 0 {
		m.gameOver = true
		m.heroWon = true
		m.combatLog = append(m.combatLog, fmt.Sprintf("%s ha sido derrotado!", m.enemy.Name))
		return m
	}

	result = m.enemy.Attack(m.hero)
	m.combatLog = append(m.combatLog, cleanLog(result))

	if m.hero.Attributes.Health <= 0 {
		m.gameOver = true
		m.heroWon = false
		m.combatLog = append(m.combatLog, fmt.Sprintf("%s ha caído en combate!", m.hero.Name))
		return m
	}

	m.round++
	return m
}

func cleanLog(s string) string {
	return strings.TrimSpace(s)
}

func (m CombatModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(fmt.Sprintf("=== %s ===", m.quest.Title)))
	b.WriteString("\n\n")

	if len(m.quest.Events) > 0 {
		b.WriteString(m.quest.Events[0].Description)
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

	if m.gameOver {
		b.WriteString("\n")
		if m.heroWon {
			b.WriteString(victoryStyle.Render("¡VICTORIA! " + m.hero.Name + " ha triunfado."))
		} else {
			b.WriteString(defeatStyle.Render("DERROTA. " + m.hero.Name + " ha caído."))
		}
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("Presiona cualquier tecla para salir."))
	} else {
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
