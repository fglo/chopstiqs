package component

import (
	"image/color"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type TextInputCursorDrawer interface {
	Draw(*textInputCursor) *ebiten.Image
}

type DefaultTextInputCursorDrawer struct {
	Color color.RGBA
}

func (d *DefaultTextInputCursorDrawer) Draw(cursor *textInputCursor) *ebiten.Image {
	cursor.image.WritePixels(d.draw(cursor))
	return cursor.image
}

func (d *DefaultTextInputCursorDrawer) draw(cursor *textInputCursor) []byte {
	arr := make([]byte, cursor.pixelRows*cursor.pixelCols)

	if cursor.frameCount >= 50 {
		return arr
	}

	for rowId := cursor.firstPixelRowId; rowId <= cursor.lastPixelRowId; rowId++ {
		rowNumber := cursor.pixelCols * rowId

		for colId := cursor.firstPixelColId; colId <= cursor.lastPixelColId; colId += 4 {
			arr[colId+rowNumber] = d.Color.R
			arr[colId+1+rowNumber] = d.Color.G
			arr[colId+2+rowNumber] = d.Color.B
			arr[colId+3+rowNumber] = d.Color.A
		}
	}

	return arr
}
