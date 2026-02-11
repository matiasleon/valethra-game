package entities

// ActionType represents the type of action in an event.
type ActionType string

const (
	ActionCombat   ActionType = "combat"
	ActionDialogue ActionType = "dialogue"
	ActionPuzzle   ActionType = "puzzle"
	ActionCutscene ActionType = "cutscene"
)

// ActionContext holds all context needed for an action to execute.
type ActionContext struct {
	Hero   *Character
	Quest  *Quest
	Engine interface{} // CombatEngine, DialogueEngine, etc.
}

// ActionResult is the result of executing an action.
type ActionResult struct {
	Success bool
	Message string
	Data    interface{} // extra data (loot, dialogue choices, etc.)
}

// ActionEvent is an interface for any action that can happen in an event.
// Examples: combat, dialogue, puzzle, cutscene, trade, etc.
type ActionEvent interface {
	// Execute runs the action with the given context.
	Execute(ctx *ActionContext) ActionResult
	// Type returns the action type (combat, dialogue, etc.).
	Type() ActionType
	// Description returns a human-readable description.
	Description() string
}
