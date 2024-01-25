package component

import (
	"image/color"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type TextInputDrawer interface {
	Draw(textInput *TextInput) *ebiten.Image
}

type DefaultTextInputDrawer struct {
	BorderColor     color.RGBA
	BackgroundColor color.RGBA
}

func (d *DefaultTextInputDrawer) isCorner(textInput *TextInput, rowId, colId int) bool {
	return (rowId == textInput.firstPixelRowId || rowId == textInput.lastPixelRowId) && (colId == textInput.firstPixelColId || colId == textInput.lastPixelColId)
}

func (d *DefaultTextInputDrawer) isBorder(textInput *TextInput, rowId, colId int) bool {
	return rowId == textInput.firstPixelRowId || rowId == textInput.lastPixelRowId || colId == textInput.firstPixelColId || colId == textInput.lastPixelColId
}

func (d DefaultTextInputDrawer) Draw(textInput *TextInput) *ebiten.Image {
	img := ebiten.NewImage(textInput.width, textInput.height)
	img.Fill(d.BackgroundColor)

	textInput.image.WritePixels(d.draw(textInput))

	return img
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

	return arr
}
