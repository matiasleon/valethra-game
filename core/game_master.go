package core

import (
	"saturday-chill/core/entities"
	"saturday-chill/core/models"
)

type GameMaster struct {
	Quests []entities.Quest
}

func (gm *GameMaster) StartQuest(quest models.Quest) {
	quest.Start()
}
