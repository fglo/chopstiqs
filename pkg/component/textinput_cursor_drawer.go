package component

import (
	"image/color"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type TextInputCursorDrawer interface {
	Draw(*textInputCursor) *ebiten.Image
	ResetBlink()
}

type DefaultTextInputCursorDrawer struct {
	Color color.RGBA

	frameCount int
}

func (d *DefaultTextInputCursorDrawer) ResetBlink() {
	d.frameCount = 0
}

func (d *DefaultTextInputCursorDrawer) Draw(textInputCursor *textInputCursor) *ebiten.Image {
	textInputCursor.image.WritePixels(d.draw(textInputCursor))
	return textInputCursor.image
}

func (d *DefaultTextInputCursorDrawer) draw(textInput *textInputCursor) []byte {
	arr := make([]byte, textInput.pixelRows*textInput.pixelCols)

	if d.frameCount = (d.frameCount + 1) % 90; d.frameCount >= 50 {
		return arr
	}

	for rowId := textInput.firstPixelRowId; rowId <= textInput.lastPixelRowId; rowId++ {
		rowNumber := textInput.pixelCols * rowId

		for colId := textInput.firstPixelColId; colId <= textInput.lastPixelColId; colId += 4 {
			arr[colId+rowNumber] = d.Color.R
			arr[colId+1+rowNumber] = d.Color.G
			arr[colId+2+rowNumber] = d.Color.B
			arr[colId+3+rowNumber] = d.Color.A
		}
	}

	return arr
}
