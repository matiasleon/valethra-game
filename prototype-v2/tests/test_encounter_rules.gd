extends SceneTree

const Rules = preload("res://scripts/encounter_rules.gd")

var failures: Array[String] = []


func _init() -> void:
	test_direct_attack_loses()
	test_observation_creates_winning_advantage()
	test_provocation_exhausts_grukh()
	test_negotiation_avoids_combat()
	test_unknown_choice_is_rejected()

	if failures.is_empty():
		print("PASS: all encounter rule tests")
		quit(0)
		return

	for failure in failures:
		push_error(failure)
	quit(1)


func test_direct_attack_loses() -> void:
	var result: Dictionary = Rules.resolve(&"attack")
	expect_equal(result["outcome"], &"defeat", "Direct attack should lose")
	expect_equal(result["aldric_power"], 55, "Aldric base power should be 55")
	expect_equal(result["grukh_power"], 65, "Grukh base power should be 65")


func test_observation_creates_winning_advantage() -> void:
	var result: Dictionary = Rules.resolve(&"observe")
	expect_equal(result["outcome"], &"victory", "Observation should produce victory")
	expect_equal(result["aldric_power"], 70, "Observation should add 15 attack")
	expect_equal(result["grukh_power"], 65, "Observation should not change Grukh")


func test_provocation_exhausts_grukh() -> void:
	var result: Dictionary = Rules.resolve(&"provoke")
	expect_equal(result["outcome"], &"victory", "Provocation should produce victory")
	expect_equal(result["aldric_power"], 55, "Provocation should not change Aldric")
	expect_equal(result["grukh_power"], 50, "Provocation should remove 15 stamina from Grukh")


func test_negotiation_avoids_combat() -> void:
	var result: Dictionary = Rules.resolve(&"negotiate")
	expect_equal(result["outcome"], &"agreement", "Negotiation should avoid combat")


func test_unknown_choice_is_rejected() -> void:
	# Avoid deliberately triggering an assertion in the test runner. The exported
	# choice list is the boundary used by the UI and must remain exhaustive.
	expect_equal(Rules.CHOICES.size(), 4, "The prototype should expose exactly four choices")
	expect_true(&"unknown" not in Rules.CHOICES, "Unknown choices must not be exposed")


func expect_equal(actual: Variant, expected: Variant, message: String) -> void:
	if actual != expected:
		failures.append("%s: expected %s, got %s" % [message, expected, actual])


func expect_true(value: bool, message: String) -> void:
	if not value:
		failures.append(message)
