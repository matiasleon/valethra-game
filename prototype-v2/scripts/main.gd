extends Control

const EncounterRulesScript = preload("res://scripts/encounter_rules.gd")
const BRIDGE_TEXTURE = preload("res://assets/art/grukh-bridge.png")
const ALDRIC_TEXTURE = preload("res://assets/art/aldric.png")
const GRUKH_TEXTURE = preload("res://assets/art/grukh.png")

const COLOR_PANEL := Color("#090c10ed")
const COLOR_PANEL_LIGHT := Color("#15171bed")
const COLOR_BORDER := Color("#9a7848")
const COLOR_BORDER_BRIGHT := Color("#d0a866")
const COLOR_TEXT := Color("#f0e4cf")
const COLOR_MUTED := Color("#b9aa92")
const COLOR_ACCENT := Color("#c69a5a")
const COLOR_SUCCESS := Color("#a8bd75")
const COLOR_DANGER := Color("#c97870")

var phase := &"intro"
var overlay: VBoxContainer
var stage_veil: ColorRect
var aldric_art: TextureRect
var grukh_art: TextureRect


func _ready() -> void:
	build_stage()
	build_overlay()
	show_intro()


func _unhandled_input(event: InputEvent) -> void:
	if event is InputEventKey and event.pressed and not event.echo:
		if event.keycode == KEY_ESCAPE:
			get_tree().quit()
			return

		if phase == &"intro" and (event.keycode == KEY_ENTER or event.keycode == KEY_SPACE):
			show_decisions()
			return

		if phase == &"decision":
			match event.keycode:
				KEY_1:
					resolve_choice(&"attack")
				KEY_2:
					resolve_choice(&"observe")
				KEY_3:
					resolve_choice(&"provoke")
				KEY_4:
					resolve_choice(&"negotiate")

		if phase == &"result" and event.keycode == KEY_R:
			show_intro()


func build_stage() -> void:
	var background := TextureRect.new()
	background.texture = BRIDGE_TEXTURE
	background.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	background.expand_mode = TextureRect.EXPAND_IGNORE_SIZE
	background.stretch_mode = TextureRect.STRETCH_KEEP_ASPECT_COVERED
	background.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(background)

	var atmosphere := ColorRect.new()
	atmosphere.color = Color("#07101a24")
	atmosphere.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	atmosphere.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(atmosphere)

	aldric_art = make_character_art(ALDRIC_TEXTURE, Vector2(26, 114), Vector2(430, 590))
	add_child(aldric_art)
	grukh_art = make_character_art(GRUKH_TEXTURE, Vector2(706, 112), Vector2(570, 486))
	add_child(grukh_art)

	stage_veil = ColorRect.new()
	stage_veil.color = Color("#02040718")
	stage_veil.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	stage_veil.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(stage_veil)

	add_child(make_nameplate("ALDRIC", Vector2(48, 432), Vector2(250, 52)))
	add_child(make_nameplate("GRUKH", Vector2(982, 432), Vector2(250, 52)))

	var frame := PanelContainer.new()
	frame.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	frame.offset_left = 14
	frame.offset_top = 14
	frame.offset_right = -14
	frame.offset_bottom = -14
	frame.mouse_filter = Control.MOUSE_FILTER_IGNORE
	frame.add_theme_stylebox_override("panel", make_style(Color.TRANSPARENT, COLOR_BORDER, 2, 10, 0))
	add_child(frame)


func build_overlay() -> void:
	var margin := MarginContainer.new()
	margin.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	margin.add_theme_constant_override("margin_left", 72)
	margin.add_theme_constant_override("margin_top", 28)
	margin.add_theme_constant_override("margin_right", 72)
	margin.add_theme_constant_override("margin_bottom", 36)
	add_child(margin)

	overlay = VBoxContainer.new()
	overlay.add_theme_constant_override("separation", 12)
	margin.add_child(overlay)


func show_intro() -> void:
	phase = &"intro"
	clear_overlay()
	set_stage_emphasis(1.0, 1.0, 0.08)

	overlay.add_child(make_chapter_header("CONTRATO I", "EL PEAJE DE GRUKH"))
	overlay.add_spacer(false)

	var dialogue := make_panel(COLOR_PANEL, COLOR_BORDER_BRIGHT, 2, 14, 24)
	dialogue.custom_minimum_size = Vector2(0, 198)
	overlay.add_child(dialogue)

	var content := VBoxContainer.new()
	content.add_theme_constant_override("separation", 9)
	dialogue.add_child(content)
	content.add_child(make_label("GRUKH", 24, COLOR_ACCENT, HORIZONTAL_ALIGNMENT_CENTER))
	content.add_child(make_separator())
	content.add_child(make_label("“Nadie cruza este puente sin pagar.”", 30, COLOR_TEXT, HORIZONTAL_ALIGNMENT_CENTER))
	content.add_child(make_label("El troll sostiene el garrote, pero no avanza.", 17, COLOR_MUTED, HORIZONTAL_ALIGNMENT_CENTER))
	content.add_child(make_button("CONTINUAR", show_decisions, 250))
	animate_overlay()


func show_decisions() -> void:
	phase = &"decision"
	clear_overlay()
	set_stage_emphasis(0.82, 0.82, 0.24)

	overlay.add_child(make_chapter_header("EL PUENTE DE VARNOCK", "¿QUÉ HACE ALDRIC?"))
	overlay.add_spacer(false)

	var decision_panel := make_panel(COLOR_PANEL, COLOR_BORDER, 2, 14, 22)
	decision_panel.custom_minimum_size = Vector2(0, 296)
	overlay.add_child(decision_panel)
	var box := VBoxContainer.new()
	box.add_theme_constant_override("separation", 14)
	decision_panel.add_child(box)
	box.add_child(make_label("Grukh espera su respuesta.", 18, COLOR_MUTED, HORIZONTAL_ALIGNMENT_CENTER))

	var choices := GridContainer.new()
	choices.columns = 2
	choices.add_theme_constant_override("h_separation", 15)
	choices.add_theme_constant_override("v_separation", 15)
	box.add_child(choices)
	choices.add_child(make_choice_button("1", "ATACAR", &"attack"))
	choices.add_child(make_choice_button("2", "OBSERVAR ANTES DE ACTUAR", &"observe"))
	choices.add_child(make_choice_button("3", "PROVOCARLO", &"provoke"))
	choices.add_child(make_choice_button("4", "NEGOCIAR", &"negotiate"))
	box.add_child(make_label("Mouse o teclas 1–4", 14, COLOR_MUTED.darkened(0.15), HORIZONTAL_ALIGNMENT_CENTER))
	animate_overlay()


func resolve_choice(choice: StringName) -> void:
	var result: Dictionary = EncounterRulesScript.resolve(choice)
	show_result(result)


func show_result(result: Dictionary) -> void:
	phase = &"result"
	clear_overlay()

	var outcome_color := COLOR_DANGER
	if result["outcome"] == &"victory":
		outcome_color = COLOR_SUCCESS
	elif result["outcome"] == &"agreement":
		outcome_color = COLOR_ACCENT
	set_stage_emphasis(0.52, 0.52, 0.48)

	overlay.add_child(make_chapter_header("RESOLUCIÓN", result["title"], outcome_color))
	overlay.add_spacer(false)

	var result_panel := make_panel(COLOR_PANEL, outcome_color.darkened(0.18), 2, 14, 22)
	result_panel.custom_minimum_size = Vector2(0, 350)
	overlay.add_child(result_panel)
	var content := VBoxContainer.new()
	content.add_theme_constant_override("separation", 11)
	result_panel.add_child(content)
	content.add_child(make_label(result["explanation"], 22, COLOR_TEXT, HORIZONTAL_ALIGNMENT_CENTER))
	content.add_child(make_label(result["modifier"], 24, outcome_color, HORIZONTAL_ALIGNMENT_CENTER))

	if result["outcome"] != &"agreement":
		var comparison := HBoxContainer.new()
		comparison.alignment = BoxContainer.ALIGNMENT_CENTER
		comparison.add_theme_constant_override("separation", 26)
		comparison.add_child(make_power_card("ALDRIC", result["aldric_power"], outcome_color))
		comparison.add_child(make_label("VS", 22, COLOR_MUTED, HORIZONTAL_ALIGNMENT_CENTER))
		comparison.add_child(make_power_card("GRUKH", result["grukh_power"], COLOR_TEXT))
		content.add_child(comparison)

	content.add_child(make_separator())
	content.add_child(make_label(result["consequence"], 20, COLOR_TEXT, HORIZONTAL_ALIGNMENT_CENTER))
	content.add_child(make_button("VOLVER A INTENTAR", show_intro, 290))
	content.add_child(make_label("R: reiniciar  ·  ESC: salir", 13, COLOR_MUTED, HORIZONTAL_ALIGNMENT_CENTER))
	animate_overlay()


func make_character_art(texture: Texture2D, position: Vector2, dimensions: Vector2) -> TextureRect:
	var art := TextureRect.new()
	art.texture = texture
	art.expand_mode = TextureRect.EXPAND_IGNORE_SIZE
	art.stretch_mode = TextureRect.STRETCH_SCALE
	art.position = position
	art.size = dimensions
	art.mouse_filter = Control.MOUSE_FILTER_IGNORE
	art.texture_filter = CanvasItem.TEXTURE_FILTER_LINEAR_WITH_MIPMAPS
	return art


func make_nameplate(text: String, position: Vector2, dimensions: Vector2) -> PanelContainer:
	var plate := make_panel(Color("#0c0e12e8"), COLOR_BORDER, 1, 7, 7)
	plate.position = position
	plate.size = dimensions
	plate.mouse_filter = Control.MOUSE_FILTER_IGNORE
	plate.add_child(make_label(text, 18, COLOR_TEXT, HORIZONTAL_ALIGNMENT_CENTER))
	return plate


func make_chapter_header(kicker: String, title: String, title_color: Color = COLOR_TEXT) -> PanelContainer:
	var header := make_panel(Color("#090c10dc"), COLOR_BORDER, 1, 9, 10)
	header.custom_minimum_size = Vector2(0, 92)
	var box := VBoxContainer.new()
	box.add_theme_constant_override("separation", 0)
	header.add_child(box)
	box.add_child(make_label(kicker, 13, COLOR_ACCENT, HORIZONTAL_ALIGNMENT_CENTER))
	box.add_child(make_title(title, 36, title_color))
	return header


func make_choice_button(number: String, text: String, choice: StringName) -> Button:
	var button := Button.new()
	button.text = "%s     %s" % [number, text]
	button.alignment = HORIZONTAL_ALIGNMENT_LEFT
	button.custom_minimum_size = Vector2(0, 76)
	button.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	button.add_theme_font_size_override("font_size", 19)
	button.add_theme_color_override("font_color", COLOR_TEXT)
	button.add_theme_color_override("font_hover_color", Color.WHITE)
	button.add_theme_color_override("font_focus_color", Color.WHITE)
	button.add_theme_stylebox_override("normal", make_style(COLOR_PANEL_LIGHT, COLOR_BORDER.darkened(0.15), 1, 8, 18))
	button.add_theme_stylebox_override("hover", make_style(Color("#25221dee"), COLOR_BORDER_BRIGHT, 2, 8, 18))
	button.add_theme_stylebox_override("focus", make_style(Color("#25221dee"), COLOR_BORDER_BRIGHT, 2, 8, 18))
	button.add_theme_stylebox_override("pressed", make_style(Color("#101216f5"), COLOR_ACCENT, 2, 8, 18))
	button.pressed.connect(resolve_choice.bind(choice), Object.CONNECT_DEFERRED)
	return button


func make_button(text: String, callback: Callable, width: float) -> Button:
	var button := Button.new()
	button.text = text
	button.custom_minimum_size = Vector2(width, 49)
	button.size_flags_horizontal = Control.SIZE_SHRINK_CENTER
	button.add_theme_font_size_override("font_size", 17)
	button.add_theme_color_override("font_color", COLOR_TEXT)
	button.add_theme_color_override("font_hover_color", Color.WHITE)
	button.add_theme_stylebox_override("normal", make_style(Color("#4a2925f2"), COLOR_BORDER, 1, 7, 12))
	button.add_theme_stylebox_override("hover", make_style(Color("#673a32f5"), COLOR_BORDER_BRIGHT, 2, 7, 12))
	button.add_theme_stylebox_override("pressed", make_style(Color("#2d1917f5"), COLOR_ACCENT, 2, 7, 12))
	button.pressed.connect(callback, Object.CONNECT_DEFERRED)
	return button


func make_power_card(character_name: String, power: int, value_color: Color) -> PanelContainer:
	var panel := make_panel(COLOR_PANEL_LIGHT, COLOR_BORDER.darkened(0.25), 1, 8, 12)
	panel.custom_minimum_size = Vector2(210, 74)
	var box := VBoxContainer.new()
	box.add_theme_constant_override("separation", 0)
	box.add_child(make_label(character_name, 13, COLOR_MUTED, HORIZONTAL_ALIGNMENT_CENTER))
	box.add_child(make_label("PODER  %d" % power, 26, value_color, HORIZONTAL_ALIGNMENT_CENTER))
	panel.add_child(box)
	return panel


func make_title(text: String, font_size: int, color: Color = COLOR_TEXT) -> Label:
	var label := make_label(text, font_size, color, HORIZONTAL_ALIGNMENT_CENTER)
	label.add_theme_color_override("font_shadow_color", Color("#000000e6"))
	label.add_theme_constant_override("shadow_offset_x", 2)
	label.add_theme_constant_override("shadow_offset_y", 2)
	return label


func make_label(text: String, font_size: int, color: Color, alignment: HorizontalAlignment) -> Label:
	var label := Label.new()
	label.text = text
	label.horizontal_alignment = alignment
	label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	label.add_theme_font_size_override("font_size", font_size)
	label.add_theme_color_override("font_color", color)
	label.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	label.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	return label


func make_separator() -> HSeparator:
	var separator := HSeparator.new()
	separator.add_theme_constant_override("separation", 8)
	separator.add_theme_stylebox_override("separator", make_style(Color.TRANSPARENT, COLOR_BORDER.darkened(0.25), 1, 0, 0))
	return separator


func make_panel(color: Color, border_color: Color = COLOR_BORDER, border_width: int = 2, radius: int = 12, padding: int = 18) -> PanelContainer:
	var panel := PanelContainer.new()
	panel.add_theme_stylebox_override("panel", make_style(color, border_color, border_width, radius, padding))
	return panel


func make_style(color: Color, border_color: Color, border_width: int, radius: int, padding: int) -> StyleBoxFlat:
	var style := StyleBoxFlat.new()
	style.bg_color = color
	style.border_color = border_color
	style.set_border_width_all(border_width)
	style.set_corner_radius_all(radius)
	style.set_content_margin_all(padding)
	return style


func clear_overlay() -> void:
	for child in overlay.get_children():
		overlay.remove_child(child)
		child.free()


func set_stage_emphasis(aldric_alpha: float, grukh_alpha: float, veil_alpha: float) -> void:
	aldric_art.modulate.a = aldric_alpha
	grukh_art.modulate.a = grukh_alpha
	stage_veil.color = Color(0.01, 0.015, 0.025, veil_alpha)


func animate_overlay() -> void:
	overlay.modulate.a = 0.0
	var tween := create_tween()
	tween.set_ease(Tween.EASE_OUT)
	tween.set_trans(Tween.TRANS_QUAD)
	tween.tween_property(overlay, "modulate:a", 1.0, 0.24)
