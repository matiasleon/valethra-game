package game

// LayoutConfig centralizes all layout parameters for the combat scene.
// All positions, offsets, and spacing values are defined here.
type LayoutConfig struct {
	// Screen dimensions
	ScreenWidth  int
	ScreenHeight int

	// Vertical zones (Y positions from top)
	TitleY     int
	EventDescY int
	SpritesY   int // baseline Y for character sprites

	// Character horizontal positions
	HeroX  int
	EnemyX int

	// HUD offsets (relative to sprite bottom)
	NameYOffset      int // negative = above sprite
	HealthBarOffset  int
	HealthTextOffset int
	ArmorBarOffset   int
	ArmorTextOffset  int
	AttackTextOffset int

	// Spacing and padding
	HUDPadding  int // padding between HUD and combat log
	CharHUDTail int // extra padding below attack text

	// Combat log
	LogLeftMargin   int
	LogHeaderHeight int
	LogLineHeight   int
	LogMaxLineLen   int // max characters per line before truncation
	MinLogLines     int

	// Status bar
	StatusBarHeight  int
	StatusTextOffset int // offset from screen bottom

	// Sprite rendering
	SpriteScale float64 // scale multiplier for sprites
	BarHeight   int     // height of health/armor bars

	// Text rendering (debug font characteristics)
	CharWidthPX int // approx width of debug font char (for centering)
	WrapWidth   int // event description wrap width (chars)
}

// DefaultLayoutConfig returns the standard layout configuration for 800x600.
func DefaultLayoutConfig(screenWidth, screenHeight int) LayoutConfig {
	return LayoutConfig{
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,

		// Vertical layout
		TitleY:     12,
		EventDescY: 35,
		SpritesY:   105,

		// Character positions
		HeroX:  120,
		EnemyX: 480,

		// HUD offsets
		NameYOffset:      -20,
		HealthBarOffset:  10,
		HealthTextOffset: 23,
		ArmorBarOffset:   38,
		ArmorTextOffset:  51,
		AttackTextOffset: 66,

		// Spacing
		HUDPadding:  15,
		CharHUDTail: 10,

		// Combat log
		LogLeftMargin:   30,
		LogHeaderHeight: 18,
		LogLineHeight:   16,
		LogMaxLineLen:   85,
		MinLogLines:     4,

		// Status bar
		StatusBarHeight:  30,
		StatusTextOffset: 22,

		// Rendering
		SpriteScale: 2.5,
		BarHeight:   10,
		CharWidthPX: 6,
		WrapWidth:   70,
	}
}
