class_name EncounterRules
extends RefCounted

const ALDRIC := {
	"attack": 25,
	"stamina": 30,
	"agility": 15,
	"intelligence": 10,
	"willpower": 12,
}

const GRUKH := {
	"attack": 35,
	"stamina": 30,
}

const CHOICES: Array[StringName] = [
	&"attack",
	&"observe",
	&"provoke",
	&"negotiate",
]


static func resolve(choice: StringName) -> Dictionary:
	assert(choice in CHOICES, "Unknown encounter choice: %s" % choice)

	var aldric_attack: int = ALDRIC["attack"]
	var aldric_stamina: int = ALDRIC["stamina"]
	var grukh_attack: int = GRUKH["attack"]
	var grukh_stamina: int = GRUKH["stamina"]
	var modifier := "Sin preparación"
	var explanation := "Aldric enfrenta a Grukh sin alterar las condiciones del encuentro."
	var outcome := &"defeat"
	var title := "DERROTA"
	var consequence := "Grukh mantiene el control del puente."

	match choice:
		&"observe":
			if ALDRIC["intelligence"] >= 10:
				aldric_attack += 15
				modifier = "Pierna herida: +15 Ataque"
				explanation = "Aldric observa a Grukh y descubre que protege una vieja herida."
		&"provoke":
			if ALDRIC["agility"] >= 15:
				grukh_stamina -= 15
				modifier = "Grukh agotado: -15 Stamina"
				explanation = "Aldric evita sus embestidas hasta que el troll pierde el aliento."
		&"negotiate":
			if ALDRIC["willpower"] >= 12:
				return {
					"choice": choice,
					"outcome": &"agreement",
					"title": "ACUERDO",
					"modifier": "Voluntad 12",
					"explanation": "Aldric sostiene la mirada de Grukh y propone un peaje aceptable.",
					"consequence": "El puente queda abierto y Grukh permanece con vida.",
					"aldric_power": ALDRIC["attack"] + ALDRIC["stamina"],
					"grukh_power": GRUKH["attack"] + GRUKH["stamina"],
				}

	var aldric_power: int = aldric_attack + aldric_stamina
	var grukh_power: int = grukh_attack + grukh_stamina
	if aldric_power > grukh_power:
		outcome = &"victory"
		title = "VICTORIA"
		consequence = "Aldric derrota a Grukh y el puente queda abierto."

	return {
		"choice": choice,
		"outcome": outcome,
		"title": title,
		"modifier": modifier,
		"explanation": explanation,
		"consequence": consequence,
		"aldric_power": aldric_power,
		"grukh_power": grukh_power,
	}
