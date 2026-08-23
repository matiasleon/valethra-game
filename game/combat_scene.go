package game

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/matiasleon/valethra/core/entities"
	entityactions "github.com/matiasleon/valethra/core/entities/actions"
	"github.com/matiasleon/valethra/core/ports"

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
	quest       entities.Quest
	hero        *entities.Character
	enemy       *entities.Character
	heroMax     entities.Attributes
	enemyMax    entities.Attributes
	engine      ports.CombatEngine
	round       int
	state       CombatState
	combatLog   []string
	eventIndex  int  // current event within the quest
	waitingNext bool // waiting for key press before next event or exit

	// Sprite animators
	heroSprites  *CharacterSprites
	enemySprites *CharacterSprites
	// Loader for enemy sprites per event
	enemySpriteDir string
	enemyFlipX     bool

	// Layout
	layoutCfg LayoutConfig
	layoutMgr *LayoutManager
}

func NewCombatScene(quest entities.Quest, hero *entities.Character, engine ports.CombatEngine, heroSprites, enemySprites *CharacterSprites) *CombatScene {
	// Get combat action from first event (supports both new Actions and legacy Enemies)
	var enemy *entities.Character
	if len(quest.Events) > 0 {
		enemies := combatEnemiesFromEvent(quest.Events[0])
		if len(enemies) > 0 {
			enemy = &enemies[0]
		}
	}

	layoutCfg := DefaultLayoutConfig(ScreenWidth, ScreenHeight)
	cs := &CombatScene{
		quest:        quest,
		hero:         hero,
		enemy:        enemy,
		heroMax:      hero.Attributes,
		engine:       engine,
		round:        1,
		state:        StatePlaying,
		combatLog:    []string{"El combate comienza... [SPACE/ENTER para atacar]"},
		heroSprites:  heroSprites,
		enemySprites: enemySprites,
		eventIndex:   0,
		layoutCfg:    layoutCfg,
		layoutMgr:    NewLayoutManager(layoutCfg),
	}
	if enemy != nil {
		cs.enemyMax = enemy.Attributes
	}
	return cs
}

func (c *CombatScene) Update() SceneResult {
	// Actualizar animaciones de sprites
	if c.heroSprites != nil {
		c.heroSprites.Animator.Update()
		if c.heroSprites.Animator.IsFinished() && c.heroSprites.Animator.State() != StateDead {
			c.heroSprites.Animator.SetState(StateIdle)
		}
	}
	if c.enemySprites != nil {
		c.enemySprites.Animator.Update()
		if c.enemySprites.Animator.IsFinished() && c.enemySprites.Animator.State() != StateDead {
			c.enemySprites.Animator.SetState(StateIdle)
		}
	}

	keyPressed := inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEnter)

	// If waiting for key after victory/defeat, handle transition
	if c.waitingNext && keyPressed {
		switch c.state {
		case StateVictory:
			// Check if there are more events in this quest
			if c.eventIndex+1 < len(c.quest.Events) {
				c.advanceToNextEvent()
				return SceneResult{}
			}
			return SceneResult{Done: true, NextTag: "victory"}
		case StateDefeat:
			return SceneResult{Done: true, NextTag: "defeat"}
		}
	}

	if c.state != StatePlaying {
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) ||
			inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			return SceneResult{Done: true, NextTag: "exit"}
		}
		return SceneResult{}
	}

	if keyPressed || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		c.executeRound()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) ||
		inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return SceneResult{Done: true, NextTag: "exit"}
	}

	return SceneResult{}
}

func (c *CombatScene) executeRound() {
	c.combatLog = append(c.combatLog, fmt.Sprintf("-- Ronda %d --", c.round))

	// Hero ataca -> animación de ataque
	if c.heroSprites != nil {
		c.heroSprites.Animator.SetState(StateAttack)
	}

	result := c.engine.ExecuteRound(c.hero, c.enemy)
	c.combatLog = append(c.combatLog, cleanLog(result.HeroAction))

	// Enemigo recibe daño -> animación de hurt o death
	if c.enemySprites != nil {
		if result.EnemyDead {
			c.enemySprites.Animator.SetState(StateDead)
		} else {
			c.enemySprites.Animator.SetState(StateHurt)
		}
	}

	if result.EnemyDead {
		c.combatLog = append(c.combatLog, fmt.Sprintf("%s ha sido derrotado!", c.enemy.Name))
		c.state = StateVictory
		c.waitingNext = true
		if c.eventIndex+1 < len(c.quest.Events) {
			c.combatLog = append(c.combatLog, "Un nuevo enemigo aparece... [SPACE/ENTER]")
		}
		return
	}

	// Enemy ataca -> animación de ataque
	if c.enemySprites != nil {
		c.enemySprites.Animator.SetState(StateAttack)
	}

	c.combatLog = append(c.combatLog, cleanLog(result.EnemyAction))

	// Héroe recibe daño -> animación de hurt o death
	if c.heroSprites != nil {
		if result.HeroDead {
			c.heroSprites.Animator.SetState(StateDead)
		} else {
			c.heroSprites.Animator.SetState(StateHurt)
		}
	}

	if result.HeroDead {
		c.combatLog = append(c.combatLog, fmt.Sprintf("%s ha caido en combate!", c.hero.Name))
		c.state = StateDefeat
		c.waitingNext = true
		return
	}

	c.round++
}

func (c *CombatScene) advanceToNextEvent() {
	c.eventIndex++
	event := c.quest.Events[c.eventIndex]

	// Get combat action from the new event.
	enemies := combatEnemiesFromEvent(event)
	if len(enemies) > 0 {
		c.enemy = &enemies[0]
		c.enemyMax = c.enemy.Attributes
	}

	c.round = 1
	c.state = StatePlaying
	c.waitingNext = false
	c.combatLog = append(c.combatLog, "")
	c.combatLog = append(c.combatLog, fmt.Sprintf("=== %s aparece! ===", c.enemy.Name))
	c.combatLog = append(c.combatLog, "[SPACE/ENTER para atacar]")

	// Reset enemy sprite to idle
	if c.enemySprites != nil {
		c.enemySprites.Animator.SetState(StateIdle)
	}
	// Reset hero sprite to idle
	if c.heroSprites != nil {
		c.heroSprites.Animator.SetState(StateIdle)
	}
}

func combatEnemiesFromEvent(event entities.Event) []entities.Character {
	// New model: actions list.
	action := event.GetFirstActionByType(entities.ActionCombat)
	if combat, ok := action.(*entityactions.CombatAction); ok {
		return combat.Enemies
	}
	// Legacy fallback.
	return event.Enemies
}

// HeroWon returns true if the combat ended in victory.
func (c *CombatScene) HeroWon() bool {
	return c.state == StateVictory
}

func cleanLog(s string) string {
	return strings.TrimSpace(s)
}

func (c *CombatScene) Draw(screen *ebiten.Image) {
	// Background
	screen.Fill(color.RGBA{30, 30, 40, 255})

	// Compute layout for this frame
	title := c.quest.Title
	layout := c.layoutMgr.Compute(len(title), c.heroSprites, c.enemySprites)

	// Title
	ebitenutil.DebugPrintAt(screen, title, layout.TitleX, layout.TitleY)

	// Event description
	if c.eventIndex < len(c.quest.Events) {
		desc := wrapText(c.quest.Events[c.eventIndex].Description, c.layoutCfg.WrapWidth)
		ebitenutil.DebugPrintAt(screen, desc, layout.EventDescX, layout.EventDescY)
	}

	// Characters
	c.drawCharacter(screen, c.hero, &c.heroMax, layout.Hero, c.heroSprites)
	c.drawCharacter(screen, c.enemy, &c.enemyMax, layout.Enemy, c.enemySprites)

	// Combat log
	c.drawCombatLog(screen, layout.LogX, layout.LogY, layout.LogMaxLines)

	// Status bar
	c.drawStatus(screen, layout.StatusBarY, layout.StatusTextY)
}

func (c *CombatScene) drawCharacter(screen *ebiten.Image, char *entities.Character, max *entities.Attributes, layout CharacterLayout, sprites *CharacterSprites) {
	// Draw sprite
	if sprites != nil {
		frame := sprites.Animator.CurrentFrame()
		if frame != nil {
			op := &ebiten.DrawImageOptions{}
			if sprites.FlipX {
				op.GeoM.Scale(-c.layoutCfg.SpriteScale, c.layoutCfg.SpriteScale)
				op.GeoM.Translate(float64(layout.SpriteW), 0)
			} else {
				op.GeoM.Scale(c.layoutCfg.SpriteScale, c.layoutCfg.SpriteScale)
			}
			op.GeoM.Translate(float64(layout.X), float64(layout.Y))
			screen.DrawImage(frame, op)
		}
	} else {
		// Fallback rectangle
		drawRect(screen, layout.X, layout.Y, layout.SpriteW, layout.SpriteH, color.RGBA{100, 100, 100, 255})
	}

	// Name (centered above sprite)
	nameX := layout.X + (layout.SpriteW / 2) - (len(char.Name) * c.layoutCfg.CharWidthPX / 2)
	nameY := layout.Y + c.layoutCfg.NameYOffset
	ebitenutil.DebugPrintAt(screen, char.Name, nameX, nameY)

	// Health bar
	healthPct := float64(char.Attributes.Health) / float64(max.Health)
	c.drawBar(screen, layout.X, layout.HealthBarY, layout.BarWidth, c.layoutCfg.BarHeight, healthPct, color.RGBA{80, 180, 80, 255}, color.RGBA{60, 60, 60, 255})
	healthText := fmt.Sprintf("HP: %d/%d", char.Attributes.Health, max.Health)
	ebitenutil.DebugPrintAt(screen, healthText, layout.X, layout.HealthTextY)

	// Armor bar
	armorPct := 0.0
	if max.Armor > 0 {
		armorPct = float64(char.Attributes.Armor) / float64(max.Armor)
	}
	c.drawBar(screen, layout.X, layout.ArmorBarY, layout.BarWidth, c.layoutCfg.BarHeight, armorPct, color.RGBA{80, 140, 200, 255}, color.RGBA{60, 60, 60, 255})
	armorText := fmt.Sprintf("Armor: %d/%d", char.Attributes.Armor, max.Armor)
	ebitenutil.DebugPrintAt(screen, armorText, layout.X, layout.ArmorTextY)

	// Attack power
	atkText := fmt.Sprintf("ATK: %d", char.Attributes.AttackPower)
	ebitenutil.DebugPrintAt(screen, atkText, layout.X, layout.AttackTextY)
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

func (c *CombatScene) drawCombatLog(screen *ebiten.Image, x, y, maxLines int) {
	header := "--- Combat Log ---"
	headerX := (c.layoutCfg.ScreenWidth / 2) - (len(header) * c.layoutCfg.CharWidthPX / 2)
	ebitenutil.DebugPrintAt(screen, header, headerX, y)

	start := 0
	if len(c.combatLog) > maxLines {
		start = len(c.combatLog) - maxLines
	}

	for i, log := range c.combatLog[start:] {
		// Truncate long lines
		displayLog := log
		if len(displayLog) > c.layoutCfg.LogMaxLineLen {
			displayLog = displayLog[:c.layoutCfg.LogMaxLineLen-3] + "..."
		}
		ebitenutil.DebugPrintAt(screen, displayLog, x, y+c.layoutCfg.LogHeaderHeight+(i*c.layoutCfg.LogLineHeight))
	}
}

func (c *CombatScene) drawStatus(screen *ebiten.Image, statusBarY, statusTextY int) {
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
	drawRect(screen, 0, statusBarY, c.layoutCfg.ScreenWidth, c.layoutCfg.StatusBarHeight, color.RGBA{40, 40, 50, 255})

	_ = statusColor // TODO: use colored text when available
	statusX := (c.layoutCfg.ScreenWidth / 2) - (len(status) * c.layoutCfg.CharWidthPX / 2)
	ebitenutil.DebugPrintAt(screen, status, statusX, statusTextY)
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
