package core

import (
	"saturday-chill/core/entities"
)

type GameMaster struct {
	Quests []entities.Quest
}

func (gm *GameMaster) StartQuest(quest entities.Quest) {
}
