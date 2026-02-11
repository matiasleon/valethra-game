package game

import "github.com/hajimehoshi/ebiten/v2"

// SceneResult is returned by Scene.Update() to signal scene transitions.
type SceneResult struct {
	Done    bool   // true when the scene has finished
	NextTag string // hint for the scene manager: "victory", "defeat", "next_quest", "exit"
}

// Scene is the interface that all game scenes implement.
type Scene interface {
	Update() SceneResult
	Draw(screen *ebiten.Image)
}
