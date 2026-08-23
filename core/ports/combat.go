// Package ports defines interfaces for core use cases.
package ports

import "github.com/matiasleon/valethra-game/core/entities"

// RoundResult contains the outcome of a single combat round for the UI to render.
type RoundResult struct {
	HeroAction  string
	EnemyAction string
	EnemyDead   bool
	HeroDead    bool
}

// CombatEngine executes a round of combat between two characters.
type CombatEngine interface {
	ExecuteRound(hero, enemy *entities.Character) RoundResult
}
