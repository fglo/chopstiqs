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

var OS OperatingSystem

func init() {
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

var (
	CursorPosX int
	CursorPosY int

	MouseLeftButtonPressed           bool
	MouseLeftButtonJustPressed       bool
	MouseLastUpdateLeftButtonPressed bool

	MouseRightButtonPressed           bool
	MouseRightButtonJustPressed       bool
	MouseLastUpdateRightButtonPressed bool

	InputChars []rune

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
}

func Draw() {

}

func AfterDraw() {
	InputChars = InputChars[:0]
}
