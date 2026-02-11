package entities

// Quest is a quest that the player can complete.
type Quest struct {
	Title        string
	Introduction string
	Events       []Event
}

// Event is a single event in the quest.
// An event can contain one or more actions (combat, dialogue, puzzle, etc.).
type Event struct {
	Description string
	Actions     []ActionEvent // polymorphic actions (combat, dialogue, puzzle)
	Enemies     []Character   // deprecated: use Actions with actions.CombatAction instead (kept for backward compatibility)
}

// GetFirstActionByType returns the first action with the given type.
func (e *Event) GetFirstActionByType(actionType ActionType) ActionEvent {
	for _, action := range e.Actions {
		if action.Type() == actionType {
			return action
		}
	}
	return nil
}

// HasCombat returns true if this event contains a combat action.
func (e *Event) HasCombat() bool {
	// Legacy compatibility: keep treating Enemies as combat.
	return len(e.Enemies) > 0 || e.GetFirstActionByType(ActionCombat) != nil
}

// GetActionsByType returns all actions of the given type.
func (e *Event) GetActionsByType(actionType ActionType) []ActionEvent {
	var result []ActionEvent
	for _, action := range e.Actions {
		if action.Type() == actionType {
			result = append(result, action)
		}
	}
	return result
}
