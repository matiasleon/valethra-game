package entities

import (
	"strings"
	"testing"
)

// CASO 1: Ataque mayor a armadura → Muerte instantánea
func TestAttack_ExceedsArmor_CharacterDies(t *testing.T) {
	attacker := &Character{
		Name: "Attacker",
		Attributes: Attributes{
			AttackPower: 100,
		},
	}

	defender := &Character{
		Name: "Defender",
		Attributes: Attributes{
			Health: 100,
			Armor:  60,
		},
	}

	initialHealth := defender.Attributes.Health

	attacker.Attack(defender)

	if defender.Attributes.Health != 0 {
		t.Errorf("Expected Health = 0, got %d", defender.Attributes.Health)
	}

	if defender.Attributes.Armor != 0 {
		t.Errorf("Expected Armor = 0, got %d", defender.Attributes.Armor)
	}

	if initialHealth == 0 {
		t.Error("Initial health should not be 0 for this test")
	}
}

// CASO 2: Ataque igual a armadura → Muerte instantánea (límite)
func TestAttack_EqualsArmor_CharacterDies(t *testing.T) {
	attacker := &Character{
		Name: "Attacker",
		Attributes: Attributes{
			AttackPower: 60,
		},
	}

	defender := &Character{
		Name: "Defender",
		Attributes: Attributes{
			Health: 100,
			Armor:  60,
		},
	}

	attacker.Attack(defender)

	if defender.Attributes.Health != 0 {
		t.Errorf("Expected Health = 0 when attack equals armor, got %d", defender.Attributes.Health)
	}

	if defender.Attributes.Armor != 0 {
		t.Errorf("Expected Armor = 0, got %d", defender.Attributes.Armor)
	}
}

// CASO 3: Ataque menor a armadura → Solo reduce armadura
func TestAttack_LessThanArmor_ReducesArmorOnly(t *testing.T) {
	attacker := &Character{
		Name: "Attacker",
		Attributes: Attributes{
			AttackPower: 30,
		},
	}

	defender := &Character{
		Name: "Defender",
		Attributes: Attributes{
			Health: 100,
			Armor:  60,
		},
	}

	initialHealth := defender.Attributes.Health
	expectedArmor := 30 // 60 - 30 = 30

	attacker.Attack(defender)

	if defender.Attributes.Health != initialHealth {
		t.Errorf("Expected Health to remain %d, got %d", initialHealth, defender.Attributes.Health)
	}

	if defender.Attributes.Armor != expectedArmor {
		t.Errorf("Expected Armor = %d, got %d", expectedArmor, defender.Attributes.Armor)
	}
}

// CASO 4: Ataque agota completamente la armadura pero no supera → Armor = 0, Health intacto
func TestAttack_ExhaustsArmorButNotExceeds_ArmorZeroHealthIntact(t *testing.T) {
	attacker := &Character{
		Name: "Attacker",
		Attributes: Attributes{
			AttackPower: 60,
		},
	}

	defender := &Character{
		Name: "Defender",
		Attributes: Attributes{
			Health: 100,
			Armor:  60,
		},
	}

	attacker.Attack(defender)

	// Cuando ataque = armadura, el personaje muere (límite)
	if defender.Attributes.Health != 0 {
		t.Errorf("Expected Health = 0 (attack equals armor = death), got %d", defender.Attributes.Health)
	}

	if defender.Attributes.Armor != 0 {
		t.Errorf("Expected Armor = 0, got %d", defender.Attributes.Armor)
	}
}

// CASO 5: Ataque a personaje ya derrotado
func TestAttack_TargetAlreadyDefeated_ReturnsError(t *testing.T) {
	attacker := &Character{
		Name: "Attacker",
		Attributes: Attributes{
			AttackPower: 50,
		},
	}

	defender := &Character{
		Name: "Defender",
		Attributes: Attributes{
			Health: 0,
			Armor:  60,
		},
	}

	result := attacker.Attack(defender)

	if !strings.Contains(result, "no puede atacar") || !strings.Contains(result, "derrotado") {
		t.Errorf("Expected error message about defeated target, got: %s", result)
	}
}

// CASO 6: Armor se reduce correctamente en múltiples ataques
func TestAttack_MultipleAttacks_ArmorReducesProgressively(t *testing.T) {
	attacker := &Character{
		Name: "Attacker",
		Attributes: Attributes{
			AttackPower: 20,
		},
	}

	defender := &Character{
		Name: "Defender",
		Attributes: Attributes{
			Health: 100,
			Armor:  60,
		},
	}

	// Ataque 1: 60 - 20 = 40
	attacker.Attack(defender)
	if defender.Attributes.Armor != 40 {
		t.Errorf("After first attack: Expected Armor = 40, got %d", defender.Attributes.Armor)
	}

	// Ataque 2: 40 - 20 = 20
	attacker.Attack(defender)
	if defender.Attributes.Armor != 20 {
		t.Errorf("After second attack: Expected Armor = 20, got %d", defender.Attributes.Armor)
	}

	// Ataque 3: 20 - 20 = 0 (pero como es igual, muere)
	attacker.Attack(defender)
	if defender.Attributes.Armor != 0 {
		t.Errorf("After third attack: Expected Armor = 0, got %d", defender.Attributes.Armor)
	}
	if defender.Attributes.Health != 0 {
		t.Errorf("After third attack: Expected Health = 0 (death), got %d", defender.Attributes.Health)
	}
}

// CASO 7: Health solo se afecta cuando hay muerte o cuando Armor = 0
func TestAttack_HealthOnlyAffectedOnDeathOrNoArmor(t *testing.T) {
	attacker := &Character{
		Name: "Attacker",
		Attributes: Attributes{
			AttackPower: 20,
		},
	}

	defender := &Character{
		Name: "Defender",
		Attributes: Attributes{
			Health: 100,
			Armor:  100,
		},
	}

	initialHealth := defender.Attributes.Health

	// Realizar múltiples ataques que no superan la armadura
	for i := 0; i < 4; i++ {
		attacker.Attack(defender)
		if defender.Attributes.Health != initialHealth {
			t.Errorf("After attack %d: Expected Health to remain %d (has armor), got %d", i+1, initialHealth, defender.Attributes.Health)
		}
	}

	// Ahora un ataque que supera la armadura
	strongAttacker := &Character{
		Name: "StrongAttacker",
		Attributes: Attributes{
			AttackPower: 200,
		},
	}

	strongAttacker.Attack(defender)
	if defender.Attributes.Health != 0 {
		t.Errorf("After fatal attack: Expected Health = 0, got %d", defender.Attributes.Health)
	}
}

// CASO 8: Daño progresivo cuando Armor = 0
func TestAttack_NoArmor_HealthDamagedProgressively(t *testing.T) {
	attacker := &Character{
		Name: "Attacker",
		Attributes: Attributes{
			AttackPower: 30,
		},
	}

	defender := &Character{
		Name: "Defender",
		Attributes: Attributes{
			Health: 100,
			Armor:  0, // Sin armadura
		},
	}

	initialHealth := defender.Attributes.Health
	expectedHealthDamage := 30
	expectedHealth := initialHealth - expectedHealthDamage

	attacker.Attack(defender)

	if defender.Attributes.Health != expectedHealth {
		t.Errorf("Expected Health = %d (damaged by %d), got %d", expectedHealth, expectedHealthDamage, defender.Attributes.Health)
	}

	if defender.Attributes.Armor != 0 {
		t.Errorf("Expected Armor to remain 0, got %d", defender.Attributes.Armor)
	}
}

// CASO 9: Armor se agota y el siguiente ataque va a Health
func TestAttack_ArmorExhausted_NextAttackGoesToHealth(t *testing.T) {
	attacker := &Character{
		Name: "Attacker",
		Attributes: Attributes{
			AttackPower: 30,
		},
	}

	defender := &Character{
		Name: "Defender",
		Attributes: Attributes{
			Health: 100,
			Armor:  30, // Exactamente igual al ataque
		},
	}

	// Primer ataque: agota la armadura (30 = 30 → muerte)
	attacker.Attack(defender)
	if defender.Attributes.Health != 0 {
		t.Errorf("Expected Health = 0 (armor exhausted = death), got %d", defender.Attributes.Health)
	}
	if defender.Attributes.Armor != 0 {
		t.Errorf("Expected Armor = 0, got %d", defender.Attributes.Armor)
	}
}

// CASO 10: Múltiples ataques progresivos cuando Armor = 0
func TestAttack_MultipleAttacksNoArmor_HealthReducesProgressively(t *testing.T) {
	attacker := &Character{
		Name: "Attacker",
		Attributes: Attributes{
			AttackPower: 20,
		},
	}

	defender := &Character{
		Name: "Defender",
		Attributes: Attributes{
			Health: 100,
			Armor:  0,
		},
	}

	// Ataque 1: 100 - 20 = 80
	attacker.Attack(defender)
	if defender.Attributes.Health != 80 {
		t.Errorf("After first attack: Expected Health = 80, got %d", defender.Attributes.Health)
	}

	// Ataque 2: 80 - 20 = 60
	attacker.Attack(defender)
	if defender.Attributes.Health != 60 {
		t.Errorf("After second attack: Expected Health = 60, got %d", defender.Attributes.Health)
	}

	// Ataque 3: 60 - 20 = 40
	attacker.Attack(defender)
	if defender.Attributes.Health != 40 {
		t.Errorf("After third attack: Expected Health = 40, got %d", defender.Attributes.Health)
	}
}

// CASO 11: Ataque que reduce armadura parcialmente
func TestAttack_PartialArmorReduction_HealthIntact(t *testing.T) {
	attacker := &Character{
		Name: "Attacker",
		Attributes: Attributes{
			AttackPower: 25,
		},
	}

	defender := &Character{
		Name: "Defender",
		Attributes: Attributes{
			Health: 100,
			Armor:  50,
		},
	}

	initialHealth := defender.Attributes.Health
	expectedArmor := 25 // 50 - 25 = 25

	attacker.Attack(defender)

	if defender.Attributes.Health != initialHealth {
		t.Errorf("Expected Health to remain %d, got %d", initialHealth, defender.Attributes.Health)
	}

	if defender.Attributes.Armor != expectedArmor {
		t.Errorf("Expected Armor = %d, got %d", expectedArmor, defender.Attributes.Armor)
	}
}
