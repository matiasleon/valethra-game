// Package core provides game orchestration and use cases.
package core

import (
	"valethra-game/core/entities"
	"valethra-game/core/ports"
)

// GameMaster orchestrates quest progression and combat execution.
type GameMaster struct {
	Hero         *entities.Character
	HeroMaxAttrs entities.Attributes
	Quests       []entities.Quest
	CurrentQuest int
}

func NewGameMaster(hero *entities.Character, quests []entities.Quest) *GameMaster {
	return &GameMaster{
		Hero:         hero,
		HeroMaxAttrs: hero.Attributes,
		Quests:       quests,
		CurrentQuest: 0,
	}
}

// ExecuteRound runs one round of combat: hero attacks, then enemy attacks if alive.
func (gm *GameMaster) ExecuteRound(hero, enemy *entities.Character) ports.RoundResult {
	result := ports.RoundResult{}

	result.HeroAction = hero.Attack(enemy)

	if enemy.Attributes.Health <= 0 {
		result.EnemyDead = true
		return result
	}

	result.EnemyAction = enemy.Attack(hero)

	if hero.Attributes.Health <= 0 {
		result.HeroDead = true
	}

	return result
}

func (gm *GameMaster) CurrentQuestData() entities.Quest {
	return gm.Quests[gm.CurrentQuest]
}

func (gm *GameMaster) HasNextQuest() bool {
	return gm.CurrentQuest+1 < len(gm.Quests)
}

func (gm *GameMaster) AdvanceQuest() {
	gm.CurrentQuest++
}

func (gm *GameMaster) IsGameOver() bool {
	return gm.CurrentQuest >= len(gm.Quests)
}

// HealHero restores 50% of missing health and 25% of original armor between quests.
func (gm *GameMaster) HealHero() {
	missingHealth := gm.HeroMaxAttrs.Health - gm.Hero.Attributes.Health
	gm.Hero.Attributes.Health += missingHealth / 2

	armorRestore := gm.HeroMaxAttrs.Armor / 4
	gm.Hero.Attributes.Armor += armorRestore
	if gm.Hero.Attributes.Armor > gm.HeroMaxAttrs.Armor {
		gm.Hero.Attributes.Armor = gm.HeroMaxAttrs.Armor
	}
}
