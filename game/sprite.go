package game

import "github.com/hajimehoshi/ebiten/v2"

// AnimState representa los estados visuales de un personaje.
type AnimState int

const (
	StateIdle   AnimState = iota // Quieto, respirando
	StateAttack                  // Atacando
	StateHurt                    // Recibiendo daño
	StateDead                    // Muerto
)

// Sprite contiene los frames de una animación específica.
// Por ejemplo, el sprite de "idle" tiene 6 frames.
type Sprite struct {
	frames []*ebiten.Image
}

// NewSprite crea un sprite desde un slice de frames.
func NewSprite(frames []*ebiten.Image) *Sprite {
	return &Sprite{frames: frames}
}

// Frame retorna el frame en el índice dado.
// Si el índice está fuera de rango, retorna el primer frame.
func (s *Sprite) Frame(index int) *ebiten.Image {
	if index < 0 || index >= len(s.frames) {
		return s.frames[0]
	}
	return s.frames[index]
}

// FrameCount retorna la cantidad de frames.
func (s *Sprite) FrameCount() int {
	return len(s.frames)
}

// Animator controla qué frame mostrar y cuándo avanzar la animación.
// Cada personaje tiene su propio Animator.
type Animator struct {
	sprites       map[AnimState]*Sprite // Sprites por estado
	state         AnimState             // Estado actual
	currentFrame  int                   // Frame actual dentro del sprite
	ticksPerFrame int                   // Cuántos Update() antes de cambiar frame
	tickCounter   int                   // Contador de ticks
	loop          bool                  // Si la animación hace loop
	finished      bool                  // Si una animación no-loop terminó
}

// NewAnimator crea un animator con los sprites dados.
// ticksPerFrame controla la velocidad: 8 significa ~7.5 FPS de animación (60/8).
func NewAnimator(sprites map[AnimState]*Sprite, ticksPerFrame int) *Animator {
	return &Animator{
		sprites:       sprites,
		state:         StateIdle,
		ticksPerFrame: ticksPerFrame,
		loop:          true, // Idle hace loop por defecto
	}
}

// SetState cambia el estado de animación.
// Reinicia el frame al principio y configura si hace loop.
func (a *Animator) SetState(state AnimState) {
	if a.state == state && !a.finished {
		return // Ya estamos en ese estado y no terminó
	}

	a.state = state
	a.currentFrame = 0
	a.tickCounter = 0
	a.finished = false

	// Solo idle hace loop. Attack, hurt y death se reproducen una vez.
	a.loop = (state == StateIdle)
}

// Update avanza la animación. Llamar una vez por frame (60 veces/segundo).
func (a *Animator) Update() {
	if a.finished {
		return // Animación terminó, no avanzar
	}

	sprite := a.sprites[a.state]
	if sprite == nil {
		return // No hay sprite para este estado
	}

	a.tickCounter++
	if a.tickCounter >= a.ticksPerFrame {
		a.tickCounter = 0
		a.currentFrame++

		// ¿Llegamos al final?
		if a.currentFrame >= sprite.FrameCount() {
			if a.loop {
				a.currentFrame = 0 // Volver al principio
			} else {
				a.currentFrame = sprite.FrameCount() - 1 // Quedarse en el último
				a.finished = true
			}
		}
	}
}

// CurrentFrame retorna la imagen actual para dibujar.
func (a *Animator) CurrentFrame() *ebiten.Image {
	sprite := a.sprites[a.state]
	if sprite == nil {
		return nil
	}
	return sprite.Frame(a.currentFrame)
}

// IsFinished retorna true si una animación no-loop terminó.
// Útil para saber cuándo volver a idle después de attack/hurt.
func (a *Animator) IsFinished() bool {
	return a.finished
}

// State retorna el estado actual.
func (a *Animator) State() AnimState {
	return a.state
}

// CharacterSprites agrupa todos los sprites y el animator de un personaje.
type CharacterSprites struct {
	Animator *Animator
	FlipX    bool // Si true, el sprite se dibuja espejado (para enemigos)
}

// LoadCharacterSprites carga todos los sprites de un personaje.
// basePath es "assets/soldier" o "assets/orc".
func LoadCharacterSprites(basePath string, flipX bool) (*CharacterSprites, error) {
	const frameSize = 100 // Cada frame es 100x100

	// Cargar cada sprite sheet
	idleFrames, err := LoadSpriteSheet(basePath+"/idle.png", frameSize)
	if err != nil {
		return nil, err
	}

	attackFrames, err := LoadSpriteSheet(basePath+"/attack.png", frameSize)
	if err != nil {
		return nil, err
	}

	hurtFrames, err := LoadSpriteSheet(basePath+"/hurt.png", frameSize)
	if err != nil {
		return nil, err
	}

	deathFrames, err := LoadSpriteSheet(basePath+"/death.png", frameSize)
	if err != nil {
		return nil, err
	}

	// Crear mapa de sprites por estado
	sprites := map[AnimState]*Sprite{
		StateIdle:   NewSprite(idleFrames),
		StateAttack: NewSprite(attackFrames),
		StateHurt:   NewSprite(hurtFrames),
		StateDead:   NewSprite(deathFrames),
	}

	// Crear animator con velocidad de 8 ticks por frame (~7.5 FPS)
	animator := NewAnimator(sprites, 8)

	return &CharacterSprites{
		Animator: animator,
		FlipX:    flipX,
	}, nil
}
