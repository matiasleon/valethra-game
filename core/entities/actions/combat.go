// Package actions defines concrete quest action implementations.
package actions

import "github.com/matiasleon/valethra-game/core/entities"

// CombatAction represents a combat encounter with one or more enemies.
type CombatAction struct {
	Enemies []entities.Character
	Desc    string // optional custom description
}

// Execute executes the combat action.
// This is a placeholder - the actual combat execution happens in CombatScene.
func (c *CombatAction) Execute(ctx *entities.ActionContext) entities.ActionResult {
	// Combat execution is handled by CombatScene/CombatEngine
	// This method is for future use when we decouple combat from scenes
	return entities.ActionResult{
		Success: true,
		Message: "Combat initiated",
	}
}

// Type returns ActionCombat.
func (c *CombatAction) Type() entities.ActionType {
	return entities.ActionCombat
}

// Description returns a description of the combat.
func (c *CombatAction) Description() string {
	if c.Desc != "" {
		return c.Desc
	}
	if len(c.Enemies) == 1 {
		return "Enfrentar a " + c.Enemies[0].Name
	}
	return "Enfrentar a múltiples enemigos"
}
