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

const spriteScale = 3
const (
	nameYOffset       = -20
	healthBarOffset   = 10
	healthTextOffset  = 25
	armorBarOffset    = 45
	armorTextOffset   = 60
	attackTextOffset  = 80
	hudPadding        = 20
	logHeaderHeight   = 20
	logLineHeight     = 16
	statusBarHeight   = 30
	minLogLines       = 2
)

type CombatScene struct {
	quest        entities.Quest
	hero         *entities.Character
	enemy        *entities.Character
	heroMax      entities.Attributes
	enemyMax     entities.Attributes
	engine       ports.CombatEngine
	round        int
	state        CombatState
	combatLog    []string
	eventIndex   int // current event within the quest
	waitingNext  bool // waiting for key press before next event or exit

	// Sprite animators
	heroSprites  *CharacterSprites
	enemySprites *CharacterSprites
	// Loader for enemy sprites per event
	enemySpriteDir string
	enemyFlipX     bool
}

func NewCombatScene(quest entities.Quest, hero *entities.Character, engine ports.CombatEngine, heroSprites, enemySprites *CharacterSprites) *CombatScene {
	// Start with first event's first enemy
	var enemy *entities.Character
	if len(quest.Events) > 0 && len(quest.Events[0].Enemies) > 0 {
		enemy = &quest.Events[0].Enemies[0]
	}

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
	if len(event.Enemies) > 0 {
		c.enemy = &c.quest.Events[c.eventIndex].Enemies[0]
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

	// Title
	title := fmt.Sprintf("=== %s ===", c.quest.Title)
	ebitenutil.DebugPrintAt(screen, title, 20, 20)

	// Event description (current event)
	if c.eventIndex < len(c.quest.Events) {
		desc := wrapText(c.quest.Events[c.eventIndex].Description, 90)
		ebitenutil.DebugPrintAt(screen, desc, 20, 50)
	}

	// Hero (left side)
	heroX := 100
	enemyX := 500
	spriteY := 180
	c.drawCharacterWithSprite(screen, c.hero, &c.heroMax, heroX, spriteY, c.heroSprites)

	// Enemy (right side)
	c.drawCharacterWithSprite(screen, c.enemy, &c.enemyMax, enemyX, spriteY, c.enemySprites)

	// Combat log (dinámico según altura disponible)
	_, heroH := spriteDimensions(c.heroSprites)
	_, enemyH := spriteDimensions(c.enemySprites)
	maxSpriteHeight := heroH
	if enemyH > maxSpriteHeight {
		maxSpriteHeight = enemyH
	}
	hudBottom := spriteY + maxSpriteHeight + attackTextOffset
	logTop := hudBottom + hudPadding
	statusTop := ScreenHeight - statusBarHeight
	availableLogHeight := statusTop - logTop - hudPadding
	maxLines := (availableLogHeight - logHeaderHeight) / logLineHeight
	if maxLines < minLogLines {
		maxLines = minLogLines
		logTop = statusTop - (logHeaderHeight + (maxLines * logLineHeight)) - hudPadding
	}
	c.drawCombatLog(screen, 20, logTop, maxLines)

	// Controls / End state
	c.drawStatus(screen)
}

func (c *CombatScene) drawCharacterWithSprite(screen *ebiten.Image, char *entities.Character, max *entities.Attributes, x, y int, sprites *CharacterSprites) {
	frameWidth := 100
	frameHeight := 100

	// Dibujar sprite animado
	if sprites != nil {
		frame := sprites.Animator.CurrentFrame()
		if frame != nil {
			op := &ebiten.DrawImageOptions{}
			frameWidth = frame.Bounds().Dx()
			frameHeight = frame.Bounds().Dy()

			// Si FlipX está activado, espejar horizontalmente
			if sprites.FlipX {
				op.GeoM.Scale(-spriteScale, spriteScale)              // Escala negativa en X = espejo
				op.GeoM.Translate(float64(frameWidth)*spriteScale, 0) // Compensar porque el espejo mueve la imagen
			} else {
				op.GeoM.Scale(spriteScale, spriteScale)
			}

			op.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(frame, op)
		}
	} else {
		// Fallback: rectángulo si no hay sprites
		drawRect(screen, x, y, int(float64(frameWidth)*spriteScale), int(float64(frameHeight)*spriteScale), color.RGBA{100, 100, 100, 255})
	}

	spriteWidth := int(float64(frameWidth) * spriteScale)
	spriteHeight := int(float64(frameHeight) * spriteScale)
	barWidth := spriteWidth
	nameX := x + (spriteWidth / 2) - (len(char.Name) * 3)

	// Name (centrado sobre el sprite)
	ebitenutil.DebugPrintAt(screen, char.Name, nameX, y+nameYOffset)

	// Health bar
	healthPct := float64(char.Attributes.Health) / float64(max.Health)
	c.drawBar(screen, x, y+spriteHeight+healthBarOffset, barWidth, 12, healthPct, color.RGBA{80, 180, 80, 255}, color.RGBA{60, 60, 60, 255})
	healthText := fmt.Sprintf("HP: %d/%d", char.Attributes.Health, max.Health)
	ebitenutil.DebugPrintAt(screen, healthText, x, y+spriteHeight+healthTextOffset)

	// Armor bar
	armorPct := 0.0
	if max.Armor > 0 {
		armorPct = float64(char.Attributes.Armor) / float64(max.Armor)
	}
	c.drawBar(screen, x, y+spriteHeight+armorBarOffset, barWidth, 12, armorPct, color.RGBA{80, 140, 200, 255}, color.RGBA{60, 60, 60, 255})
	armorText := fmt.Sprintf("Armor: %d/%d", char.Attributes.Armor, max.Armor)
	ebitenutil.DebugPrintAt(screen, armorText, x, y+spriteHeight+armorTextOffset)

	// Attack power
	atkText := fmt.Sprintf("ATK: %d", char.Attributes.AttackPower)
	ebitenutil.DebugPrintAt(screen, atkText, x, y+spriteHeight+attackTextOffset)
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
	ebitenutil.DebugPrintAt(screen, "--- Combat Log ---", x, y)

	start := 0
	if len(c.combatLog) > maxLines {
		start = len(c.combatLog) - maxLines
	}

	for i, log := range c.combatLog[start:] {
		// Truncate long lines
		displayLog := log
		if len(displayLog) > 95 {
			displayLog = displayLog[:92] + "..."
		}
		ebitenutil.DebugPrintAt(screen, displayLog, x, y+logHeaderHeight+(i*logLineHeight))
	}
}

func spriteDimensions(sprites *CharacterSprites) (int, int) {
	if sprites == nil {
		size := int(100 * spriteScale)
		return size, size
	}
	frame := sprites.Animator.CurrentFrame()
	if frame == nil {
		size := int(100 * spriteScale)
		return size, size
	}
	w := int(float64(frame.Bounds().Dx()) * spriteScale)
	h := int(float64(frame.Bounds().Dy()) * spriteScale)
	return w, h
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
