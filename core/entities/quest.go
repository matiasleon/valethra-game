package entities

// Event is a single event in the quest.
type Event struct {
	Description string
	Enemies     []Character
}

// Quest is a quest that the player can complete.
type Quest struct {
	Title        string
	Introduction string
	Events       []Event
}
