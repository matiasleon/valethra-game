package game

import (
	"fmt"
	"image/color"
	"strings"

	"saturday-chill/core/entities"
	"saturday-chill/core/ports"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type CombatState int

const (
	StatePlaying CombatState = iota
	StateVictory
	StateDefeat
)

type CombatScene struct {
	quest     entities.Quest
	hero      *entities.Character
	enemy     *entities.Character
	heroMax   entities.Attributes
	enemyMax  entities.Attributes
	engine    ports.CombatEngine
	round     int
	state     CombatState
	combatLog []string

	// Animation state
	heroFlash  int
	enemyFlash int
}

func NewCombatScene(quest entities.Quest, hero, enemy *entities.Character, engine ports.CombatEngine) *CombatScene {
	return &CombatScene{
		quest:     quest,
		hero:      hero,
		enemy:     enemy,
		heroMax:   hero.Attributes,
		enemyMax:  enemy.Attributes,
		engine:    engine,
		round:     1,
		state:     StatePlaying,
		combatLog: []string{"El combate comienza... [SPACE/ENTER para atacar]"},
	}
}

func (c *CombatScene) Update() error {
	// Decrease flash timers
	if c.heroFlash > 0 {
		c.heroFlash--
	}
	if c.enemyFlash > 0 {
		c.enemyFlash--
	}

	// Handle input
	if c.state != StatePlaying {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
			inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
			inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			return ebiten.Termination
		}
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
		inpututil.IsKeyJustPressed(ebiten.KeyA) {
		c.executeRound()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) ||
		inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

	return nil
}

func (c *CombatScene) executeRound() {
	c.combatLog = append(c.combatLog, fmt.Sprintf("-- Ronda %d --", c.round))

	// Hero attacks
	result := c.engine.ExecuteRound(c.hero, c.enemy)
	c.combatLog = append(c.combatLog, cleanLog(result.HeroAction))
	c.enemyFlash = 15 // Flash for 15 frames

	if result.EnemyDead {
		c.state = StateVictory
		c.combatLog = append(c.combatLog, fmt.Sprintf("%s ha sido derrotado!", c.enemy.Name))
		return
	}

	// Enemy attacks
	c.combatLog = append(c.combatLog, cleanLog(result.EnemyAction))
	c.heroFlash = 15

	if result.HeroDead {
		c.state = StateDefeat
		c.combatLog = append(c.combatLog, fmt.Sprintf("%s ha caido en combate!", c.hero.Name))
		return
	}

	c.round++
}

func cleanLog(s string) string {
	return strings.TrimSpace(s)
}

func (c *CombatScene) Draw(screen *ebiten.Image) {
	// Background
	screen.Fill(color.RGBA{30, 30, 40, 255})

	// Title
	title := fmt.Sprintf("=== %s ===", c.quest.Title)
	ebitenutil.DebugPrintAt(screen, title, 20, 20)

	// Event description
	if len(c.quest.Events) > 0 {
		desc := wrapText(c.quest.Events[0].Description, 90)
		ebitenutil.DebugPrintAt(screen, desc, 20, 50)
	}

	// Hero (left side)
	c.drawCharacter(screen, c.hero, &c.heroMax, 50, 180, c.heroFlash > 0, color.RGBA{60, 120, 60, 255})

	// Enemy (right side)
	c.drawCharacter(screen, c.enemy, &c.enemyMax, 450, 180, c.enemyFlash > 0, color.RGBA{120, 60, 60, 255})

	// Combat log
	c.drawCombatLog(screen, 20, 420)

	// Controls / End state
	c.drawStatus(screen)
}

func (c *CombatScene) drawCharacter(screen *ebiten.Image, char *entities.Character, max *entities.Attributes, x, y int, flashing bool, baseColor color.RGBA) {
	// Character sprite placeholder (colored rectangle)
	spriteColor := baseColor
	if flashing {
		spriteColor = color.RGBA{255, 255, 255, 255}
	}
	if char.Attributes.Health <= 0 {
		spriteColor = color.RGBA{80, 80, 80, 255}
	}

	drawRect(screen, x, y, 100, 120, spriteColor)

	// Name
	ebitenutil.DebugPrintAt(screen, char.Name, x, y-20)

	// Health bar
	healthPct := float64(char.Attributes.Health) / float64(max.Health)
	c.drawBar(screen, x, y+130, 100, 12, healthPct, color.RGBA{80, 180, 80, 255}, color.RGBA{60, 60, 60, 255})
	healthText := fmt.Sprintf("HP: %d/%d", char.Attributes.Health, max.Health)
	ebitenutil.DebugPrintAt(screen, healthText, x, y+145)

	// Armor bar
	armorPct := 0.0
	if max.Armor > 0 {
		armorPct = float64(char.Attributes.Armor) / float64(max.Armor)
	}
	c.drawBar(screen, x, y+165, 100, 12, armorPct, color.RGBA{80, 140, 200, 255}, color.RGBA{60, 60, 60, 255})
	armorText := fmt.Sprintf("Armor: %d/%d", char.Attributes.Armor, max.Armor)
	ebitenutil.DebugPrintAt(screen, armorText, x, y+180)

	// Attack power
	atkText := fmt.Sprintf("ATK: %d", char.Attributes.AttackPower)
	ebitenutil.DebugPrintAt(screen, atkText, x, y+200)
}

func (c *CombatScene) drawBar(screen *ebiten.Image, x, y, width, height int, pct float64, fillColor, bgColor color.Color) {
	if pct < 0 {
		pct = 0
	}
	if pct > 1 {
		pct = 1
	}

	// Background
	drawRect(screen, x, y, width, height, bgColor)

	// Fill
	fillWidth := int(float64(width) * pct)
	if fillWidth > 0 {
		drawRect(screen, x, y, fillWidth, height, fillColor)
	}
}

func (c *CombatScene) drawCombatLog(screen *ebiten.Image, x, y int) {
	ebitenutil.DebugPrintAt(screen, "--- Combat Log ---", x, y)

	start := 0
	if len(c.combatLog) > 8 {
		start = len(c.combatLog) - 8
	}

	for i, log := range c.combatLog[start:] {
		// Truncate long lines
		displayLog := log
		if len(displayLog) > 95 {
			displayLog = displayLog[:92] + "..."
		}
		ebitenutil.DebugPrintAt(screen, displayLog, x, y+20+(i*16))
	}
}

func (c *CombatScene) drawStatus(screen *ebiten.Image) {
	var status string
	statusColor := color.RGBA{200, 200, 200, 255}

	switch c.state {
	case StatePlaying:
		status = "[A/SPACE/ENTER] Atacar  |  [Q/ESC] Salir"
	case StateVictory:
		status = "VICTORIA! " + c.hero.Name + " ha triunfado. [Presiona cualquier tecla]"
		statusColor = color.RGBA{100, 255, 100, 255}
	case StateDefeat:
		status = "DERROTA. " + c.hero.Name + " ha caido. [Presiona cualquier tecla]"
		statusColor = color.RGBA{255, 100, 100, 255}
	}

	// Draw status bar background
	drawRect(screen, 0, ScreenHeight-30, ScreenWidth, 30, color.RGBA{40, 40, 50, 255})

	_ = statusColor // TODO: use colored text when available
	ebitenutil.DebugPrintAt(screen, status, 20, ScreenHeight-22)
}

func drawRect(screen *ebiten.Image, x, y, width, height int, clr color.Color) {
	rect := ebiten.NewImage(width, height)
	rect.Fill(clr)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(rect, op)
}

func wrapText(text string, maxWidth int) string {
	if len(text) <= maxWidth {
		return text
	}

	var result strings.Builder
	words := strings.Fields(text)
	lineLen := 0

	for i, word := range words {
		if i > 0 {
			if lineLen+1+len(word) > maxWidth {
				result.WriteString("\n")
				lineLen = 0
			} else {
				result.WriteString(" ")
				lineLen++
			}
		}
		result.WriteString(word)
		lineLen += len(word)
	}

	return result.String()
}
