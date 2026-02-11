package game

// CharacterLayout holds computed layout positions for a single character.
type CharacterLayout struct {
	X, Y         int // sprite position
	SpriteW      int // sprite width (scaled)
	SpriteH      int // sprite height (scaled)
	BarWidth     int // width of health/armor bars
	NameX        int // X position for name text
	HealthBarY   int // Y position for health bar
	HealthTextY  int // Y position for health text
	ArmorBarY    int // Y position for armor bar
	ArmorTextY   int // Y position for armor text
	AttackTextY  int // Y position for attack text
}

// CombatLayout holds all precomputed layout positions for one frame.
type CombatLayout struct {
	TitleX, TitleY         int
	EventDescX, EventDescY int
	Hero                   CharacterLayout
	Enemy                  CharacterLayout
	LogX, LogY, LogMaxLines int
	StatusBarY, StatusTextY int
}

// LayoutManager computes all layout positions based on configuration.
type LayoutManager struct {
	cfg LayoutConfig
}

// NewLayoutManager creates a layout manager with the given configuration.
func NewLayoutManager(cfg LayoutConfig) *LayoutManager {
	return &LayoutManager{cfg: cfg}
}

// Compute calculates all layout positions for the current frame.
// It needs sprite dimensions to properly position HUD elements.
func (lm *LayoutManager) Compute(titleLen int, heroSprites, enemySprites *CharacterSprites) CombatLayout {
	layout := CombatLayout{}

	// Title (centered)
	layout.TitleX = (lm.cfg.ScreenWidth / 2) - (titleLen * lm.cfg.CharWidthPX / 2)
	layout.TitleY = lm.cfg.TitleY

	// Event description
	layout.EventDescX = lm.cfg.LogLeftMargin
	layout.EventDescY = lm.cfg.EventDescY

	// Character layouts
	heroW, heroH := lm.spriteDimensions(heroSprites)
	enemyW, enemyH := lm.spriteDimensions(enemySprites)
	maxSpriteHeight := heroH
	if enemyH > maxSpriteHeight {
		maxSpriteHeight = enemyH
	}

	layout.Hero = lm.characterLayout(lm.cfg.HeroX, lm.cfg.SpritesY, heroW, heroH)
	layout.Enemy = lm.characterLayout(lm.cfg.EnemyX, lm.cfg.SpritesY, enemyW, enemyH)

	// Combat log
	hudBottom := lm.cfg.SpritesY + maxSpriteHeight + lm.cfg.AttackTextOffset + lm.cfg.CharHUDTail
	layout.LogX = lm.cfg.LogLeftMargin
	layout.LogY = hudBottom + lm.cfg.HUDPadding

	// Calculate available space for log
	statusTop := lm.cfg.ScreenHeight - lm.cfg.StatusBarHeight
	availableLogHeight := statusTop - layout.LogY - lm.cfg.HUDPadding
	layout.LogMaxLines = (availableLogHeight - lm.cfg.LogHeaderHeight) / lm.cfg.LogLineHeight
	if layout.LogMaxLines < lm.cfg.MinLogLines {
		layout.LogMaxLines = lm.cfg.MinLogLines
	}

	// Status bar
	layout.StatusBarY = lm.cfg.ScreenHeight - lm.cfg.StatusBarHeight
	layout.StatusTextY = lm.cfg.ScreenHeight - lm.cfg.StatusTextOffset

	return layout
}

// characterLayout computes layout for a single character.
func (lm *LayoutManager) characterLayout(x, y, spriteW, spriteH int) CharacterLayout {
	cl := CharacterLayout{
		X:       x,
		Y:       y,
		SpriteW: spriteW,
		SpriteH: spriteH,
		BarWidth: spriteW,
	}

	// Name (centered above sprite)
	// Note: actual nameX depends on character name length, computed at draw time

	// HUD offsets (relative to sprite bottom)
	cl.HealthBarY = y + spriteH + lm.cfg.HealthBarOffset
	cl.HealthTextY = y + spriteH + lm.cfg.HealthTextOffset
	cl.ArmorBarY = y + spriteH + lm.cfg.ArmorBarOffset
	cl.ArmorTextY = y + spriteH + lm.cfg.ArmorTextOffset
	cl.AttackTextY = y + spriteH + lm.cfg.AttackTextOffset

	return cl
}

// spriteDimensions calculates scaled sprite dimensions.
func (lm *LayoutManager) spriteDimensions(sprites *CharacterSprites) (int, int) {
	if sprites == nil {
		size := int(100 * lm.cfg.SpriteScale)
		return size, size
	}
	frame := sprites.Animator.CurrentFrame()
	if frame == nil {
		size := int(100 * lm.cfg.SpriteScale)
		return size, size
	}
	w := int(float64(frame.Bounds().Dx()) * lm.cfg.SpriteScale)
	h := int(float64(frame.Bounds().Dy()) * lm.cfg.SpriteScale)
	return w, h
}
