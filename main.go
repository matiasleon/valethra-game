package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"saturday-chill/core"
	"saturday-chill/core/entities"
	"saturday-chill/core/usecases"
	"saturday-chill/ui"
)

func main() {
	hero := createHero()
	quests := []entities.Quest{
		createTrollQuest(),
		createForestQuest(),
	}

	gm := core.NewGameMaster(hero, quests)

	// Pantalla de intro
	intro := ui.NewTransitionModel(
		"Valethra",
		"En las tierras de Valethra, donde los caminos son tan peligrosos como las bestias que los acechan, un mercenario sin tierra acepta los contratos que otros rechazan.\n\nSu nombre es Aldric. No busca gloria ni redención — solo monedas suficientes para sobrevivir otra semana.\n\nPero esta semana será diferente.",
	)
	if _, err := tea.NewProgram(intro).Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Intro de la primera quest
	firstQuest := gm.CurrentQuestData()
	firstQuestIntro := ui.NewTransitionModel(firstQuest.Title, firstQuest.Introduction)
	if _, err := tea.NewProgram(firstQuestIntro).Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	for !gm.IsGameOver() {
		quest := gm.CurrentQuestData()

		session := usecases.NewCombatSession(quest, hero, gm)
		combatModel := ui.NewCombatModel(session)
		finalModel, err := tea.NewProgram(combatModel).Run()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		result, ok := finalModel.(ui.CombatModel)
		if !ok {
			fmt.Println("Error: tipo de modelo inesperado")
			os.Exit(1)
		}

		if !result.HeroWon() {
			defeat := ui.NewTransitionModel(
				"Fin del Camino",
				"Aldric cae de rodillas. La oscuridad lo envuelve mientras el frío del acero enemigo se desvanece.\n\nSu historia termina aquí, en un camino olvidado de Valethra. Pero las tierras no olvidan a quienes luchan — aunque caigan.",
			)
			if _, err := tea.NewProgram(defeat).Run(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			break
		}

		if gm.HasNextQuest() {
			gm.HealHero()
			gm.AdvanceQuest()

			transition := ui.NewTransitionModel(
				gm.CurrentQuestData().Title,
				gm.CurrentQuestData().Introduction,
			)
			if _, err := tea.NewProgram(transition).Run(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		} else {
			gm.AdvanceQuest()
			victory := ui.NewTransitionModel(
				"Victoria",
				"Aldric limpia su espada y observa el horizonte. El bosque de Varnwood queda atrás, y con él las criaturas que lo acechaban.\n\nLos mercaderes podrán cruzar el puente. El bosque sanará con el tiempo. Y Aldric tendrá monedas para otra semana.\n\nPero en Valethra, la calma nunca dura. Su leyenda apenas comienza.",
			)
			if _, err := tea.NewProgram(victory).Run(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
	}
}

func createHero() *entities.Character {
	return &entities.Character{
		ID:          "hero-1",
		Name:        "Aldric",
		Description: "Un mercenario sin tierra que sobrevive aceptando contratos que otros rechazan.",
		Attributes: entities.Attributes{
			Health:       100,
			AttackPower:  25,
			Armor:        40,
			Agility:      15,
			Intelligence: 10,
			Willpower:    12,
			Speed:        14,
			Level:        3,
		},
	}
}

func createTrollQuest() entities.Quest {
	return entities.Quest{
		Title:        "El Peaje de Grukh",
		Introduction: "Los mercaderes de Varnock llevan semanas sin poder cruzar el puente sur. Un troll se ha asentado debajo y exige un peaje absurdo a todo el que pase. El gremio de comerciantes ofrece una recompensa por eliminarlo. Nadie pregunta si el troll tiene razones para estar ahí.",
		Events: []entities.Event{
			{
				Description: "Aldric llega al puente de Varnock al anochecer. El olor a carne quemada y huesos apilados confirma los rumores. Grukh emerge de las sombras bajo el puente, con un garrote improvisado.",
				Enemies: []entities.Character{
					{
						ID:          "enemy-1",
						Name:        "Grukh",
						Description: "Un troll exiliado de su clan que se instaló bajo el puente de Varnock, cobrando peaje a los viajeros. Algunos dicen que solo quiere sobrevivir. Otros dicen que devora a quienes no pagan.",
						Attributes: entities.Attributes{
							Health:       150,
							AttackPower:  35,
							Armor:        60,
							Agility:      5,
							Intelligence: 4,
							Willpower:    8,
							Speed:        6,
							Level:        4,
						},
					},
				},
			},
		},
	}
}

func createForestQuest() entities.Quest {
	return entities.Quest{
		Title:        "La Senda de Varnwood",
		Introduction: "Tras cruzar el puente, Aldric se adentra en el bosque de Varnwood. Los mercaderes advirtieron sobre criaturas que acechan entre los árboles desde que la patrulla del barón dejó de recorrer el camino. Nadie sabe si fue negligencia o un acuerdo con algo peor.",
		Events: []entities.Event{
			{
				Description: "Una sombra se desliza entre los troncos. Un lobo con los ojos inyectados en una luz violeta antinatural bloquea el sendero. No es un animal común — algo lo ha transformado.",
				Enemies: []entities.Character{
					{
						ID:          "enemy-2",
						Name:        "Lobo Corrompido",
						Description: "Un lobo del bosque de Varnwood corrompido por magia oscura. Sus ojos brillan con una luz violeta antinatural.",
						Attributes: entities.Attributes{
							Health:       80,
							AttackPower:  20,
							Armor:        15,
							Agility:      25,
							Intelligence: 3,
							Willpower:    5,
							Speed:        22,
							Level:        3,
						},
					},
				},
			},
			{
				Description: "El bosque se abre en un claro donde los árboles están muertos. En el centro, una figura encapuchada murmura sobre un círculo de runas. Al notar a Aldric, se yergue. Un desterrado del Cónclave de Reth, practicando magia prohibida con los restos de los animales del bosque.",
				Enemies: []entities.Character{
					{
						ID:          "enemy-3",
						Name:        "Hereje de Reth",
						Description: "Un mago desterrado del Cónclave de Reth por practicar artes prohibidas. Se refugió en Varnwood para continuar sus experimentos con la corrupción de la vida.",
						Attributes: entities.Attributes{
							Health:       70,
							AttackPower:  30,
							Armor:        25,
							Agility:      8,
							Intelligence: 22,
							Willpower:    18,
							Speed:        10,
							Level:        4,
						},
					},
				},
			},
		},
	}
}
