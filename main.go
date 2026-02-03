package main

import (
	"fmt"

	"saturday-chill/models"
)

func main() {

}

func testQuest() {
	batman := &models.Character{
		ID:          "1",
		Name:        "Batman",
		Description: "Batman is a superhero who fights crime in Gotham City.",
		Attributes: models.Attributes{
			Health:      100,
			AttackPower: 70,
			Armor:       100,
		},
	}

	superman := &models.Character{
		ID:          "2",
		Name:        "Superman",
		Description: "Superman is a superhero who fights crime in Metropolis.",
		Attributes: models.Attributes{
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
