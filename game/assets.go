// Package game provides the Ebitengine runtime and scenes.
package game

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/png" // Registra el decoder PNG

	"github.com/hajimehoshi/ebiten/v2"
)

// Embebe todos los PNGs de soldier y orc en el binario.
// Esto significa que al compilar, las imágenes quedan DENTRO del ejecutable.
// No necesitás distribuir los archivos .png por separado.
//
//go:embed assets/soldier/*.png assets/orc/*.png
var assetsFS embed.FS

// LoadImage carga una imagen PNG del filesystem embebido.
// path debe ser relativo a game/, ej: "assets/soldier/idle.png"
func LoadImage(path string) (*ebiten.Image, error) {
	// 1. Lee los bytes del archivo embebido
	data, err := assetsFS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading embedded file %s: %w", path, err)
	}

	// 2. Decodifica los bytes como imagen PNG
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decoding image %s: %w", path, err)
	}

	// 3. Convierte a imagen de Ebitengine (que puede dibujarse en pantalla)
	return ebiten.NewImageFromImage(img), nil
}

// LoadSpriteSheet carga un sprite sheet y lo corta en frames de frameSize x frameSize.
// Retorna un slice de imágenes, una por cada frame.
func LoadSpriteSheet(path string, frameSize int) ([]*ebiten.Image, error) {
	// 1. Carga la imagen completa
	data, err := assetsFS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading sprite sheet %s: %w", path, err)
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decoding sprite sheet %s: %w", path, err)
	}

	// 2. Calcula cuántos frames tiene (ancho total / ancho de frame)
	bounds := img.Bounds()
	frameCount := bounds.Dx() / frameSize

	// 3. Corta cada frame
	frames := make([]*ebiten.Image, frameCount)
	for i := 0; i < frameCount; i++ {
		// SubImage corta una porción de la imagen original
		// x va de i*100 a (i+1)*100, y va de 0 a 100
		x := i * frameSize
		subImg := img.(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(image.Rect(x, 0, x+frameSize, frameSize))

		frames[i] = ebiten.NewImageFromImage(subImg)
	}

	return frames, nil
}
