package actions

import "github.com/matiasleon/valethra/core/entities"

// DialogueAction represents a dialogue with an NPC (future implementation).
type DialogueAction struct {
	NPC     string
	Text    string
	Options []DialogueOption
}

// DialogueOption represents a choice in a dialogue.
type DialogueOption struct {
	Text   string
	Effect func() entities.ActionResult
}

// Execute executes the dialogue action (placeholder).
func (d *DialogueAction) Execute(ctx *entities.ActionContext) entities.ActionResult {
	return entities.ActionResult{
		Success: true,
		Message: "Dialogue completed",
	}
}

// Type returns ActionDialogue.
func (d *DialogueAction) Type() entities.ActionType {
	return entities.ActionDialogue
}

// Description returns the dialogue description.
func (d *DialogueAction) Description() string {
	return "Hablar con " + d.NPC
}
