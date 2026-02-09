// Package game provides ebiten runtime and scenes.
package game

import (
	"saturday-chill/core/entities"
	"saturday-chill/core/ports"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

type Game struct {
	combat *CombatScene
}

func NewGame(quest entities.Quest, hero, enemy *entities.Character, engine ports.CombatEngine) *Game {
	return &Game{
		combat: NewCombatScene(quest, hero, enemy, engine),
	}
}

func (g *Game) Update() error {
	return g.combat.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.combat.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
