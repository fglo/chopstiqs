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

	pixelBuf []byte
	zeroBuf  []byte
}

func (d *DefaultTextInputCursorDrawer) Draw(cursor *textInputCursor) *ebiten.Image {
	cursor.image.WritePixels(d.draw(cursor))
	return cursor.image
}

func (d *DefaultTextInputCursorDrawer) getBuffer(cursor *textInputCursor) []byte {
	size := cursor.pixelRows * cursor.pixelCols

	if len(d.pixelBuf) != size {
		d.pixelBuf = make([]byte, size)
	}

	for i := range d.pixelBuf {
		d.pixelBuf[i] = 0
	}

	return d.pixelBuf
}

func (d *DefaultTextInputCursorDrawer) getZeroBuffer(cursor *textInputCursor) []byte {
	size := cursor.pixelRows * cursor.pixelCols

	if len(d.zeroBuf) != size {
		d.zeroBuf = make([]byte, size)
	}

	for i := range d.zeroBuf {
		d.zeroBuf[i] = 0
	}

	return d.zeroBuf
}

func (d *DefaultTextInputCursorDrawer) draw(cursor *textInputCursor) []byte {
	if cursor.frameCount >= 40 {
		return d.getZeroBuffer(cursor)
	}

	arr := d.getBuffer(cursor)

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
