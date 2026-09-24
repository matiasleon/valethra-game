extends SceneTree

const MainScene = preload("res://scenes/main.tscn")

var failures: Array[String] = []


func _init() -> void:
	call_deferred("run_tests")


func run_tests() -> void:
	var main := MainScene.instantiate()
	root.add_child(main)
	await process_frame

	expect_equal(main.phase, &"intro", "The game should start on the introduction")
	expect_true(ResourceLoader.exists("res://assets/art/grukh-bridge.png"), "Bridge art must exist")
	expect_true(ResourceLoader.exists("res://assets/art/aldric.png"), "Aldric art must exist")
	expect_true(ResourceLoader.exists("res://assets/art/grukh.png"), "Grukh art must exist")

	var continue_button := find_button(main, "CONTINUAR")
	expect_true(continue_button != null, "Introduction should expose a continue button")
	continue_button.emit_signal("pressed")
	await process_frame
	expect_equal(main.phase, &"decision", "Continue should open the decisions")

	var decision_text := collect_text(main)
	for forbidden in ["Poder", "Stamina", "Agilidad", "Inteligencia", "Voluntad", "probable", "ventaja", "derrota", "victoria"]:
		expect_true(
			forbidden.to_lower() not in decision_text.to_lower(),
			"Decision screen must not reveal '%s'" % forbidden,
		)

	var attack_button := find_button(main, "1     ATACAR")
	expect_true(attack_button != null, "Decision screen should expose the attack button")
	attack_button.emit_signal("pressed")
	await process_frame
	expect_equal(main.phase, &"result", "A pressed choice should open a result without freeing its active signal")

	for choice in [&"observe", &"provoke", &"negotiate"]:
		main.show_decisions()
		main.resolve_choice(choice)
		expect_equal(main.phase, &"result", "Choice %s should open a result" % choice)

	main.resolve_choice(&"observe")
	var result_text := collect_text(main)
	expect_true("PODER  70" in result_text, "Observation result should reveal Aldric's final power")
	expect_true("PODER  65" in result_text, "Observation result should reveal Grukh's power")

	main.queue_free()
	if failures.is_empty():
		print("PASS: all scene flow tests")
		quit(0)
		return

	for failure in failures:
		push_error(failure)
	quit(1)


func collect_text(node: Node) -> String:
	var values: Array[String] = []
	collect_text_recursive(node, values)
	return "\n".join(values)


func collect_text_recursive(node: Node, values: Array[String]) -> void:
	if node is Label:
		values.append((node as Label).text)
	elif node is Button:
		values.append((node as Button).text)

	for child in node.get_children():
		collect_text_recursive(child, values)


func find_button(node: Node, button_text: String) -> Button:
	if node is Button and (node as Button).text == button_text:
		return node as Button
	for child in node.get_children():
		var match_button := find_button(child, button_text)
		if match_button != null:
			return match_button
	return null


func expect_equal(actual: Variant, expected: Variant, message: String) -> void:
	if actual != expected:
		failures.append("%s: expected %s, got %s" % [message, expected, actual])


func expect_true(value: bool, message: String) -> void:
	if not value:
		failures.append(message)
