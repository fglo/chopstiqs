package component

import (
	"image/color"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type TextInputDrawer interface {
	Draw(textInput *TextInput) *ebiten.Image
}

type DefaultTextInputDrawer struct {
	Color           color.RGBA
	ColorDisabled   color.RGBA
	ColorHovered    color.RGBA
	BackgroundColor color.RGBA
}

func (d *DefaultTextInputDrawer) isCorner(textInput *TextInput, rowId, colId int) bool {
	return (rowId == textInput.firstPixelRowId || rowId == textInput.lastPixelRowId) && (colId == textInput.firstPixelColId || colId == textInput.lastPixelColId)
}

func (d *DefaultTextInputDrawer) isBorder(textInput *TextInput, rowId, colId int) bool {
	return rowId == textInput.firstPixelRowId || rowId == textInput.lastPixelRowId || colId == textInput.firstPixelColId || colId == textInput.lastPixelColId
}

func (d *DefaultTextInputDrawer) Draw(textInput *TextInput) *ebiten.Image {
	switch {
	case textInput.disabled:
		textInput.image.WritePixels(d.drawDisabled(textInput))
	case textInput.hovering:
		textInput.image.WritePixels(d.drawHovered(textInput))
	default:
		textInput.image.WritePixels(d.draw(textInput))
	}

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
				arr[colId+rowNumber] = d.Color.R
				arr[colId+1+rowNumber] = d.Color.G
				arr[colId+2+rowNumber] = d.Color.B
				arr[colId+3+rowNumber] = d.Color.A
			} else {
				arr[colId+rowNumber] = d.BackgroundColor.R
				arr[colId+1+rowNumber] = d.BackgroundColor.G
				arr[colId+2+rowNumber] = d.BackgroundColor.B
				arr[colId+3+rowNumber] = d.BackgroundColor.A
			}
		}
	}

	return arr
}

func (d *DefaultTextInputDrawer) drawDisabled(textInput *TextInput) []byte {
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
				arr[colId+rowNumber] = d.ColorDisabled.R
				arr[colId+1+rowNumber] = d.ColorDisabled.G
				arr[colId+2+rowNumber] = d.ColorDisabled.B
				arr[colId+3+rowNumber] = d.ColorDisabled.A
			} else {
				arr[colId+rowNumber] = d.BackgroundColor.R
				arr[colId+1+rowNumber] = d.BackgroundColor.G
				arr[colId+2+rowNumber] = d.BackgroundColor.B
				arr[colId+3+rowNumber] = d.BackgroundColor.A
			}
		}
	}

	return arr
}

func (d *DefaultTextInputDrawer) drawHovered(textInput *TextInput) []byte {
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
				arr[colId+rowNumber] = d.ColorHovered.R
				arr[colId+1+rowNumber] = d.ColorHovered.G
				arr[colId+2+rowNumber] = d.ColorHovered.B
				arr[colId+3+rowNumber] = d.ColorHovered.A
			} else {
				arr[colId+rowNumber] = d.BackgroundColor.R
				arr[colId+1+rowNumber] = d.BackgroundColor.G
				arr[colId+2+rowNumber] = d.BackgroundColor.B
				arr[colId+3+rowNumber] = d.BackgroundColor.A
			}
		}
	}

	return arr
}
