package input

import (
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

	AnyKeyPressed     bool
	AnyJustKeyPressed bool
	KeyPressed        [ebiten.KeyMax + 1]bool
	KeyJustPressed    [ebiten.KeyMax + 1]bool
)

func init() {
	DetectSystem()
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
	MouseLeftButtonJustPressed = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	MouseLastUpdateLeftButtonPressed = MouseLeftButtonPressed

	MouseRightButtonPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)
	MouseRightButtonJustPressed = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)
	MouseLastUpdateRightButtonPressed = MouseRightButtonPressed

	InputChars = ebiten.AppendInputChars(InputChars)
	AnyKeyPressed = false

	DetectPressedKeys()
}

func Draw() {

}

func AfterDraw() {
	InputChars = InputChars[:0]
}
