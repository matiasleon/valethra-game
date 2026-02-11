package main

import (
	"fmt"
	"log"

	"saturday-chill/core/entities"
	entityactions "saturday-chill/core/entities/actions"
	"saturday-chill/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	cfg := game.GameConfig{
		Hero:   createHero(),
		Quests: []entities.Quest{createTrollQuest(), createForestQuest()},
		IntroTitle: "Valethra",
		IntroText: "En las tierras de Valethra, donde los caminos son tan peligrosos como las bestias que los acechan, un mercenario sin tierra acepta los contratos que otros rechazan.\n\nSu nombre es Aldric. No busca gloria ni redención — solo monedas suficientes para sobrevivir otra semana.\n\nPero esta semana será diferente.",
		VictoryTitle: "Victoria",
		VictoryText:  "Aldric limpia su espada y observa el horizonte. El bosque de Varnwood queda atrás, y con él las criaturas que lo acechaban.\n\nLos mercaderes podrán cruzar el puente. El bosque sanará con el tiempo. Y Aldric tendrá monedas para otra semana.\n\nPero en Valethra, la calma nunca dura. Su leyenda apenas comienza.",
		DefeatTitle: "Fin del Camino",
		DefeatText:  "Aldric cae de rodillas. La oscuridad lo envuelve mientras el frío del acero enemigo se desvanece.\n\nSu historia termina aquí, en un camino olvidado de Valethra. Pero las tierras no olvidan a quienes luchan — aunque caigan.",
	}

	g, err := game.NewGame(cfg)
	if err != nil {
		log.Fatalf("Error initializing game: %v", err)
	}

	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("Valethra")

	if err := ebiten.RunGame(g); err != nil {
		fmt.Printf("Game ended: %v\n", err)
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
				Actions: []entities.ActionEvent{
					&entityactions.CombatAction{
						Enemies: []entities.Character{
							{
								ID:          "enemy-1",
								Name:        "Grukh",
								Description: "Un troll exiliado de su clan que se instaló bajo el puente de Varnock, cobrando peaje a los viajeros.",
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
						Desc: "Enfrentar a Grukh el Troll",
					},
				},
			},
		},
	}
}

func createForestQuest() entities.Quest {
	return entities.Quest{
		Title:        "La Senda de Varnwood",
		Introduction: "Tras cruzar el puente, Aldric se adentra en el bosque de Varnwood. Los mercaderes advirtieron sobre criaturas que acechan entre los árboles desde que la patrulla del barón dejó de recorrer el camino.",
		Events: []entities.Event{
			{
				Description: "Una sombra se desliza entre los troncos. Un lobo con los ojos inyectados en una luz violeta antinatural bloquea el sendero.",
				Actions: []entities.ActionEvent{
					&entityactions.CombatAction{
						Enemies: []entities.Character{
							{
								ID:          "enemy-2",
								Name:        "Lobo Corrompido",
								Description: "Un lobo del bosque de Varnwood corrompido por magia oscura.",
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
				},
			},
			{
				Description: "El bosque se abre en un claro donde los árboles están muertos. Una figura encapuchada murmura sobre un círculo de runas. Un desterrado del Cónclave de Reth.",
				Actions: []entities.ActionEvent{
					&entityactions.CombatAction{
						Enemies: []entities.Character{
							{
								ID:          "enemy-3",
								Name:        "Hereje de Reth",
								Description: "Un mago desterrado del Cónclave de Reth por practicar artes prohibidas.",
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
			},
		},
	}
}
