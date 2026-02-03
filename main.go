package main

import (
	"fmt"

	"saturday-chill/core/entities"
)

func main() {

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
