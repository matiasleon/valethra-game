// Package usecases provides application-specific business logic.
package usecases

import (
	"fmt"
	"strings"

	"saturday-chill/core/entities"
	"saturday-chill/core/ports"
)

type combatSession struct {
	quest      entities.Quest
	hero       *entities.Character
	enemy      *entities.Character
	heroMax    entities.Attributes
	enemyMax   entities.Attributes
	engine     ports.CombatEngine
	round      int
	state      ports.CombatState
	combatLog  []string
	eventIndex int
}

// NewCombatSession creates a new combat session for a quest.
func NewCombatSession(quest entities.Quest, hero *entities.Character, engine ports.CombatEngine) ports.CombatSession {
	session := &combatSession{
		quest:      quest,
		hero:       hero,
		engine:     engine,
		round:      1,
		state:      ports.StateFighting,
		combatLog:  []string{"El combate comienza..."},
		eventIndex: 0,
	}

	if len(quest.Events) > 0 && len(quest.Events[0].Enemies) > 0 {
		enemyCopy := quest.Events[0].Enemies[0]
		session.enemy = &enemyCopy
		session.heroMax = hero.Attributes
		session.enemyMax = enemyCopy.Attributes
	}

	return session
}

func (s *combatSession) State() ports.CombatStateDTO {
	return ports.CombatStateDTO{
		Quest:      s.quest,
		Hero:       s.hero,
		Enemy:      s.enemy,
		HeroMax:    s.heroMax,
		EnemyMax:   s.enemyMax,
		Round:      s.round,
		State:      s.state,
		CombatLog:  append([]string(nil), s.combatLog...),
		EventIndex: s.eventIndex,
	}
}

func (s *combatSession) Handle(action ports.CombatAction) ports.CombatStateDTO {
	switch action {
	case ports.ActionConfirm:
		switch s.state {
		case ports.StateFighting:
			s.executeRound()
		case ports.StateEventVictory:
			s.advanceToNextEvent()
		case ports.StateQuestVictory, ports.StateDefeat:
			// No-op; UI can decide to exit.
		}
	case ports.ActionQuit:
		// No-op; UI handles quit.
	}

	return s.State()
}

func (s *combatSession) executeRound() {
	s.combatLog = append(s.combatLog, fmt.Sprintf("-- Ronda %d --", s.round))

	result := s.engine.ExecuteRound(s.hero, s.enemy)
	s.combatLog = append(s.combatLog, cleanLog(result.HeroAction))

	if result.EnemyDead {
		s.combatLog = append(s.combatLog, fmt.Sprintf("%s ha sido derrotado!", s.enemy.Name))
		if s.eventIndex+1 < len(s.quest.Events) {
			s.state = ports.StateEventVictory
		} else {
			s.state = ports.StateQuestVictory
		}
		return
	}

	s.combatLog = append(s.combatLog, cleanLog(result.EnemyAction))

	if result.HeroDead {
		s.state = ports.StateDefeat
		s.combatLog = append(s.combatLog, fmt.Sprintf("%s ha caído en combate!", s.hero.Name))
		return
	}

	s.round++
}

func (s *combatSession) advanceToNextEvent() {
	s.eventIndex++
	if s.eventIndex >= len(s.quest.Events) {
		s.state = ports.StateQuestVictory
		return
	}

	event := s.quest.Events[s.eventIndex]
	if len(event.Enemies) > 0 {
		enemyCopy := event.Enemies[0]
		s.enemy = &enemyCopy
		s.enemyMax = enemyCopy.Attributes
	}

	s.round = 1
	s.state = ports.StateFighting
	s.combatLog = append(s.combatLog, "")
	s.combatLog = append(s.combatLog, fmt.Sprintf("--- %s ---", event.Description))
	s.combatLog = append(s.combatLog, "El combate comienza...")
}

func cleanLog(s string) string {
	return strings.TrimSpace(s)
}
