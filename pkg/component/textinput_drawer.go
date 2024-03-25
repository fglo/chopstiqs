package component

import (
	"image/color"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type TextInputDrawer interface {
	Draw(textInput *TextInput) *ebiten.Image
	ResetCursorBlink()
}

type DefaultTextInputDrawer struct {
	BorderColor     color.RGBA
	BackgroundColor color.RGBA

	frameCount int
}

func (d *DefaultTextInputDrawer) isCorner(textInput *TextInput, rowId, colId int) bool {
	return (rowId == textInput.firstPixelRowId || rowId == textInput.lastPixelRowId) && (colId == textInput.firstPixelColId || colId == textInput.lastPixelColId)
}

func (d *DefaultTextInputDrawer) isBorder(textInput *TextInput, rowId, colId int) bool {
	return rowId == textInput.firstPixelRowId || rowId == textInput.lastPixelRowId || colId == textInput.firstPixelColId || colId == textInput.lastPixelColId
}

func (d *DefaultTextInputDrawer) ResetCursorBlink() {
	d.frameCount = 0
}

func (d *DefaultTextInputDrawer) Draw(textInput *TextInput) *ebiten.Image {
	// img := ebiten.NewImage(textInput.width, textInput.height)
	// img.Fill(d.BackgroundColor)
	textInput.image.WritePixels(d.draw(textInput))
	return textInput.image
}

func (d *DefaultTextInputDrawer) draw(textInput *TextInput) []byte {
	arr := make([]byte, textInput.pixelRows*textInput.pixelCols)
	backgroundColor := textInput.container.GetBackgroundColor()

	for rowId := textInput.firstPixelRowId; rowId <= textInput.lastPixelRowId; rowId++ {
		rowNumber := textInput.pixelCols * rowId

		for colId := textInput.firstPixelColId; colId <= textInput.lastPixelColId; colId += 4 {
			if d.isCorner(textInput, rowId, colId) {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			} else if d.isBorder(textInput, rowId, colId) {
				arr[colId+rowNumber] = d.BorderColor.R
				arr[colId+1+rowNumber] = d.BorderColor.G
				arr[colId+2+rowNumber] = d.BorderColor.B
				arr[colId+3+rowNumber] = d.BorderColor.A
			} else {
				arr[colId+rowNumber] = d.BackgroundColor.R
				arr[colId+1+rowNumber] = d.BackgroundColor.G
				arr[colId+2+rowNumber] = d.BackgroundColor.B
				arr[colId+3+rowNumber] = d.BackgroundColor.A
			}
		}
	}

	colNumber := (textInput.cursorPosX() + textInput.scrollOffset) * 4

	if !textInput.focused || colNumber >= textInput.pixelCols {
		return arr
	}

	if d.frameCount = (d.frameCount + 1) % 90; d.frameCount < 50 {
		lineHeight := textInput.metrics.Ascent - textInput.metrics.Descent - 1
		lineTop := textInput.textPosY - textInput.metrics.Ascent

		for rowId := lineTop; rowId > 0 && rowId < lineTop+lineHeight && rowId < textInput.pixelRows; rowId++ {
			rowNumber := textInput.pixelCols * rowId

			arr[colNumber+rowNumber] = d.BorderColor.R
			arr[colNumber+1+rowNumber] = d.BorderColor.G
			arr[colNumber+2+rowNumber] = d.BorderColor.B
			arr[colNumber+3+rowNumber] = d.BorderColor.A
		}
	}

	return arr
}
