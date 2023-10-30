package input

import ebiten "github.com/hajimehoshi/ebiten/v2"

var (
	CursorPosX int
	CursorPosY int

	MouseLeftButtonPressed           bool
	MouseLeftButtonJustPressed       bool
	MouseLastUpdateLeftButtonPressed bool

	MouseRightButtonPressed           bool
	MouseRightButtonJustPressed       bool
	MouseLastUpdateRightButtonPressed bool
)

func Update() {
	CursorPosX, CursorPosY = ebiten.CursorPosition()

	MouseLeftButtonPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	MouseLeftButtonJustPressed = MouseLeftButtonPressed != MouseLastUpdateLeftButtonPressed
	MouseLastUpdateLeftButtonPressed = MouseLeftButtonPressed

	MouseRightButtonPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)
	MouseRightButtonJustPressed = MouseRightButtonPressed != MouseLastUpdateRightButtonPressed
	MouseLastUpdateRightButtonPressed = MouseRightButtonPressed
}
