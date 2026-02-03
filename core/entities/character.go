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
// If attack exceeds armor, character dies. Otherwise, armor is reduced.
// If armor is 0, damage goes directly to health.
func (c *Character) Attack(target *Character) string {
	// Check if target is already defeated
	if target.Attributes.Health <= 0 {
		return fmt.Sprintf("%s cannot attack %s - target is already defeated\n", c.Name, target.Name)
	}

	attackPower := c.Attributes.AttackPower
	armor := target.Attributes.Armor

	// Si el ataque es mayor a la armadura, el personaje muere
	if attackPower > armor {
		target.Attributes.Armor = 0
		target.Attributes.Health = 0
		return fmt.Sprintf("%s attacks %s. Attack exceeds armor (%d > %d). %s dies!\n",
			c.Name, target.Name, attackPower, armor, target.Name)
	}

	// Si el ataque es menor o igual a la armadura, solo se reduce la armadura
	if armor > 0 {
		initialArmor := target.Attributes.Armor
		target.Attributes.Armor -= attackPower
		if target.Attributes.Armor < 0 {
			target.Attributes.Armor = 0
		}
		armorReduced := initialArmor - target.Attributes.Armor
		return fmt.Sprintf("%s attacks %s. Armor reduced by %d. %s has %d health and %d armor remaining.\n",
			c.Name, target.Name, armorReduced, target.Name, target.Attributes.Health, target.Attributes.Armor)
	}

	// Si no tiene armadura (Armor = 0), el ataque va directamente a Health
	healthDamage := attackPower
	target.Attributes.Health -= healthDamage
	if target.Attributes.Health < 0 {
		target.Attributes.Health = 0
	}

	return fmt.Sprintf("%s attacks %s. No armor! Health damaged by %d. %s has %d health remaining.\n",
		c.Name, target.Name, healthDamage, target.Name, target.Attributes.Health)
}

func (c *Character) Spell(target *Character) string {
	return fmt.Sprintf("spell %s", target.Name)
}
