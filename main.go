package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"saturday-chill/core/entities"
	"saturday-chill/ui"
)

func main() {
	trollQuest()
}

func trollQuest() {
	hero := &entities.Character{
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

	troll := &entities.Character{
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
	}

	quest := entities.Quest{
		Title:        "El Peaje de Grukh",
		Introduction: "Los mercaderes de Varnock llevan semanas sin poder cruzar el puente sur. Un troll se ha asentado debajo y exige un peaje absurdo a todo el que pase. El gremio de comerciantes ofrece una recompensa por eliminarlo. Nadie pregunta si el troll tiene razones para estar ahí.",
		Events: []entities.Event{
			{
				Description: "Aldric llega al puente de Varnock al anochecer. El olor a carne quemada y huesos apilados confirma los rumores. Grukh emerge de las sombras bajo el puente, con un garrote improvisado.",
				Enemies:     []entities.Character{*troll},
			},
		},
	}

	model := ui.NewCombatModel(quest, hero, troll)
	p := tea.NewProgram(model)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func testQuest() {
	batman := &entities.Character{
		ID:          "1",
		Name:        "Batman",
		Description: "Batman is a superhero who fights crime in Gotham City.",
		Attributes: entities.Attributes{
			Health:      100,
			AttackPower: 70,
			Armor:       100,
		},
	}

	superman := &entities.Character{
		ID:          "2",
		Name:        "Superman",
		Description: "Superman is a superhero who fights crime in Metropolis.",
		Attributes: entities.Attributes{
			Health:      10000,
			AttackPower: 10000,
			Armor:       10000,
		},
	}

	var result = batman.Attack(superman)
	fmt.Println(result)

	result = batman.Spell(superman)
	fmt.Println(result)

	result = superman.Attack(batman)
	fmt.Println(result)
}
