// Package models provides data structures for game characters and their attributes.
package entities

import (
	"fmt"
)

type Attributes struct {
	Health       int `json:"health"`
	AttackPower  int `json:"attack_power"`
	Agility      int `json:"agility"`
	Intelligence int `json:"intelligence"`
	Willpower    int `json:"willpower"`
	Speed        int `json:"speed"`
	Level        int `json:"level"`
	Armor        int `json:"armor"`
}

type Character struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Attributes  Attributes `json:"attributes"`
}

// Attack performs an attack on the target character.
// If armor > 0 and attack >= armor, character dies instantly.
// If armor > 0 and attack < armor, only armor is reduced.
// If armor == 0, damage goes directly to health progressively.
func (c *Character) Attack(target *Character) string {
	if target.Attributes.Health <= 0 {
		return fmt.Sprintf("%s no puede atacar a %s — ya ha sido derrotado\n", c.Name, target.Name)
	}

	attackPower := c.Attributes.AttackPower
	armor := target.Attributes.Armor

	// Si tiene armadura y el ataque es mayor o igual, muerte instantánea
	if armor > 0 && attackPower >= armor {
		target.Attributes.Armor = 0
		target.Attributes.Health = 0
		return fmt.Sprintf("%s ataca a %s. El ataque destroza la armadura (%d >= %d). %s cae derrotado!\n",
			c.Name, target.Name, attackPower, armor, target.Name)
	}

	// Si tiene armadura y el ataque es menor, solo reduce armadura
	if armor > 0 {
		target.Attributes.Armor -= attackPower
		return fmt.Sprintf("%s ataca a %s. Armadura reducida en %d. %s tiene %d de vida y %d de armadura.\n",
			c.Name, target.Name, attackPower, target.Name, target.Attributes.Health, target.Attributes.Armor)
	}

	// Sin armadura: daño progresivo a health
	target.Attributes.Health -= attackPower
	if target.Attributes.Health < 0 {
		target.Attributes.Health = 0
	}

	return fmt.Sprintf("%s ataca a %s. Sin armadura! Daño directo: %d. %s tiene %d de vida.\n",
		c.Name, target.Name, attackPower, target.Name, target.Attributes.Health)
}

func (c *Character) Spell(target *Character) string {
	return fmt.Sprintf("spell %s", target.Name)
}
