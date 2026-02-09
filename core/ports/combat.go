// Package ports defines interfaces for core use cases.
package ports

import "saturday-chill/core/entities"

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

// CombatState represents the current state of the combat flow.
type CombatState int

const (
	StateFighting CombatState = iota
	StateEventVictory
	StateQuestVictory
	StateDefeat
)

// CombatAction represents user intent from the UI layer.
type CombatAction int

const (
	ActionConfirm CombatAction = iota
	ActionQuit
)

// CombatStateDTO is a snapshot of combat state for the UI to render.
type CombatStateDTO struct {
	Quest      entities.Quest
	Hero       *entities.Character
	Enemy      *entities.Character
	HeroMax    entities.Attributes
	EnemyMax   entities.Attributes
	Round      int
	State      CombatState
	CombatLog  []string
	EventIndex int
}

// CombatSession orchestrates combat flow for a quest.
type CombatSession interface {
	State() CombatStateDTO
	Handle(action CombatAction) CombatStateDTO
}
