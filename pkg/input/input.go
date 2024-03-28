package input

import (
	"runtime"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type OperatingSystem string

const (
	Windows OperatingSystem = "windows"
	Linux   OperatingSystem = "linux"
	MacOS   OperatingSystem = "macos"
)

var (
	OS OperatingSystem

	CursorPosX int
	CursorPosY int

	MouseLeftButtonPressed           bool
	MouseLeftButtonJustPressed       bool
	MouseLastUpdateLeftButtonPressed bool

	MouseRightButtonPressed           bool
	MouseRightButtonJustPressed       bool
	MouseLastUpdateRightButtonPressed bool

	InputChars []rune

	AnyKeyPressed bool
	KeyPressed    map[ebiten.Key]bool = make(map[ebiten.Key]bool)
)

func init() {
	SetSystem()
}

func SetSystem() {
	switch runtime.GOOS {
	case "windows":
		OS = Windows
	case "linux":
		OS = Linux
	case "darwin":
		OS = MacOS
	}
}

func OSWindows() bool {
	return OS == Windows
}

func OSLinux() bool {
	return OS == Linux
}

func OSMacOS() bool {
	return OS == MacOS
}

const (
	KeyNone ebiten.Key = -1
)

func Update() {
	CursorPosX, CursorPosY = ebiten.CursorPosition()

	MouseLeftButtonPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	MouseLeftButtonJustPressed = MouseLeftButtonPressed != MouseLastUpdateLeftButtonPressed
	MouseLastUpdateLeftButtonPressed = MouseLeftButtonPressed

	MouseRightButtonPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)
	MouseRightButtonJustPressed = MouseRightButtonPressed != MouseLastUpdateRightButtonPressed
	MouseLastUpdateRightButtonPressed = MouseRightButtonPressed

	InputChars = ebiten.AppendInputChars(InputChars)
	AnyKeyPressed = false
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		p := ebiten.IsKeyPressed(k)
		KeyPressed[k] = p
		if p {
			AnyKeyPressed = true
		}
	}
}

func Draw() {

}

func AfterDraw() {
	InputChars = InputChars[:0]
}
