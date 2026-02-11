package game

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TransitionScene displays a title, message, and waits for user confirmation.
// Used for intros, quest introductions, victory, and defeat screens.
type TransitionScene struct {
	title       string
	message     string
	accentColor color.RGBA
	bgColor     color.RGBA
}

// NewTransitionScene creates a transition scene with a custom accent color.
func NewTransitionScene(title, message string, accent color.RGBA) *TransitionScene {
	return &TransitionScene{
		title:       title,
		message:     message,
		accentColor: accent,
		bgColor:     color.RGBA{25, 25, 35, 255},
	}
}

// Preset constructors for common transition types.

func NewIntroScene(title, message string) *TransitionScene {
	return NewTransitionScene(title, message, color.RGBA{220, 180, 80, 255}) // Gold
}

func NewQuestIntroScene(title, message string) *TransitionScene {
	return NewTransitionScene(title, message, color.RGBA{140, 160, 220, 255}) // Blue
}

func NewVictoryScene(title, message string) *TransitionScene {
	return NewTransitionScene(title, message, color.RGBA{80, 220, 100, 255}) // Green
}

func NewDefeatScene(title, message string) *TransitionScene {
	return NewTransitionScene(title, message, color.RGBA{220, 80, 80, 255}) // Red
}

func (t *TransitionScene) Update() SceneResult {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return SceneResult{Done: true, NextTag: "confirm"}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) ||
		inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return SceneResult{Done: true, NextTag: "exit"}
	}

	return SceneResult{}
}

func (t *TransitionScene) Draw(screen *ebiten.Image) {
	screen.Fill(t.bgColor)

	// Accent bar at top
	drawRect(screen, 0, 0, ScreenWidth, 4, t.accentColor)

	// Title
	titleX := (ScreenWidth / 2) - (len(t.title) * 3)
	if titleX < 20 {
		titleX = 20
	}
	ebitenutil.DebugPrintAt(screen, "=== "+t.title+" ===", titleX, 60)

	// Accent line below title
	lineWidth := len(t.title)*6 + 48
	lineX := titleX
	drawRect(screen, lineX, 80, lineWidth, 2, t.accentColor)

	// Message (wrapped)
	wrapped := wrapText(t.message, 80)
	lines := strings.Split(wrapped, "\n")
	y := 110
	for _, line := range lines {
		if line == "" {
			y += 10 // extra spacing for blank lines
		} else {
			ebitenutil.DebugPrintAt(screen, line, 60, y)
			y += 18
		}
	}

	// Help text at bottom
	help := "[SPACE/ENTER] Continuar"
	helpX := (ScreenWidth / 2) - (len(help) * 3)
	drawRect(screen, 0, ScreenHeight-40, ScreenWidth, 40, color.RGBA{35, 35, 45, 255})
	ebitenutil.DebugPrintAt(screen, help, helpX, ScreenHeight-28)

	// Accent bar at bottom
	drawRect(screen, 0, ScreenHeight-4, ScreenWidth, 4, t.accentColor)
}
