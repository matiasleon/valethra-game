// Package game provides ebiten runtime and scenes.
package game

import (
	"saturday-chill/core"
	"saturday-chill/core/entities"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

// GameConfig holds all configuration data needed to run the game.
type GameConfig struct {
	Hero       *entities.Character
	Quests     []entities.Quest
	IntroTitle string
	IntroText  string
	VictoryTitle string
	VictoryText  string
	DefeatTitle  string
	DefeatText   string
}

type gamePhase int

const (
	phaseIntro gamePhase = iota
	phaseQuestIntro
	phaseCombat
	phaseResult // victory or defeat transition
)

type Game struct {
	currentScene Scene
	phase        gamePhase
	gm           *core.GameMaster
	config       GameConfig

	// Pre-loaded sprites
	heroSprites  *CharacterSprites
	enemySprites *CharacterSprites
}

// NewGame creates the game, loads sprites, and starts with the intro scene.
func NewGame(cfg GameConfig) (*Game, error) {
	gm := core.NewGameMaster(cfg.Hero, cfg.Quests)

	heroSprites, err := LoadCharacterSprites("assets/soldier", false)
	if err != nil {
		return nil, err
	}

	enemySprites, err := LoadCharacterSprites("assets/orc", true)
	if err != nil {
		return nil, err
	}

	g := &Game{
		gm:           gm,
		config:       cfg,
		heroSprites:  heroSprites,
		enemySprites: enemySprites,
		phase:        phaseIntro,
	}

	g.currentScene = NewIntroScene(cfg.IntroTitle, cfg.IntroText)
	return g, nil
}

func (g *Game) Update() error {
	result := g.currentScene.Update()
	if !result.Done {
		return nil
	}

	// Exit requested from any scene
	if result.NextTag == "exit" {
		return ebiten.Termination
	}

	switch g.phase {
	case phaseIntro:
		// After intro -> first quest intro
		g.phase = phaseQuestIntro
		quest := g.gm.CurrentQuestData()
		g.currentScene = NewQuestIntroScene(quest.Title, quest.Introduction)

	case phaseQuestIntro:
		// After quest intro -> combat
		g.phase = phaseCombat
		quest := g.gm.CurrentQuestData()
		g.currentScene = NewCombatScene(quest, g.config.Hero, g.gm, g.heroSprites, g.enemySprites)

	case phaseCombat:
		switch result.NextTag {
		case "victory":
			if g.gm.HasNextQuest() {
				// Heal and advance to next quest
				g.gm.HealHero()
				g.gm.AdvanceQuest()
				g.phase = phaseQuestIntro
				quest := g.gm.CurrentQuestData()
				g.currentScene = NewQuestIntroScene(quest.Title, quest.Introduction)
			} else {
				// All quests completed
				g.gm.AdvanceQuest()
				g.phase = phaseResult
				g.currentScene = NewVictoryScene(g.config.VictoryTitle, g.config.VictoryText)
			}
		case "defeat":
			g.phase = phaseResult
			g.currentScene = NewDefeatScene(g.config.DefeatTitle, g.config.DefeatText)
		default:
			return ebiten.Termination
		}

	case phaseResult:
		// After victory/defeat screen -> exit
		return ebiten.Termination
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.currentScene.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
