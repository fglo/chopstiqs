//go:build !wasm
// +build !wasm

package input

import (
	"runtime"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func DetectSystem() {
	switch runtime.GOOS {
	case "windows":
		OS = Windows
	case "linux":
		OS = Linux
	case "darwin":
		OS = MacOS
	}
}

func DetectPressedKeys() {
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		KeyPressed[k] = ebiten.IsKeyPressed(k)
		if KeyPressed[k] {
			AnyKeyPressed = true
		}

		KeyJustPressed[k] = inpututil.IsKeyJustPressed(k)
		if KeyJustPressed[k] {
			AnyJustKeyPressed = true
		}
	}
}
