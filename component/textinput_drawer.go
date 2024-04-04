package component

import (
	"image/color"

	colorutils "github.com/fglo/chopstiqs/color"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type TextInputDrawer interface {
	Draw(textInput *TextInput) *ebiten.Image
}

type DefaultTextInputDrawer struct {
	Color           color.RGBA
	ColorDisabled   color.RGBA
	ColorHovered    color.RGBA
	BackgroundColor color.Color
	backgroundColor color.RGBA
	cornerColor     color.RGBA
}

func (d *DefaultTextInputDrawer) isCorner(textInput *TextInput, rowId, colId int) bool {
	return (rowId == textInput.firstPixelRowId || rowId == textInput.lastPixelRowId) && (colId == textInput.firstPixelColId || colId == textInput.lastPixelColId)
}

func (d *DefaultTextInputDrawer) isBorder(textInput *TextInput, rowId, colId int) bool {
	return rowId == textInput.firstPixelRowId || rowId == textInput.lastPixelRowId || colId == textInput.firstPixelColId || colId == textInput.lastPixelColId
}

func (d *DefaultTextInputDrawer) Draw(textInput *TextInput) *ebiten.Image {
	d.cornerColor = textInput.container.GetBackgroundColor()

	if d.BackgroundColor == nil {
		d.backgroundColor = d.cornerColor
	} else {
		r, g, b, a := d.BackgroundColor.RGBA()
		d.backgroundColor = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	}

	switch {
	case textInput.disabled:
		textInput.image.WritePixels(d.draw(textInput, d.ColorDisabled))
	case textInput.hovering:
		textInput.image.WritePixels(d.draw(textInput, d.ColorHovered))
	default:
		textInput.image.WritePixels(d.draw(textInput, d.Color))
	}

	return textInput.image
}

func (d *DefaultTextInputDrawer) draw(textInput *TextInput, borderColor color.RGBA) []byte {
	arr := make([]byte, textInput.pixelRows*textInput.pixelCols)

	selectingFromColId := -1
	selectingToColId := -1

	if textInput.HasSelectedText() {
		selectingFromColId = (textInput.possibleCursorPosXs[textInput.selectionStart] - textInput.scrollOffset + 3) * 4
		selectingToColId = (textInput.possibleCursorPosXs[textInput.selectionEnd] - textInput.scrollOffset + 1) * 4
	}

	for rowId := textInput.firstPixelRowId; rowId <= textInput.lastPixelRowId; rowId++ {
		rowNumber := textInput.pixelCols * rowId

		for colId := textInput.firstPixelColId; colId <= textInput.lastPixelColId; colId += 4 {
			bgColor := d.backgroundColor
			if selectingFromColId <= colId && colId <= selectingToColId &&
				textInput.firstPixelRowId+1 < rowId && rowId < textInput.lastPixelRowId-1 {
				bgColor = colorutils.Invert(bgColor)
			}

			if d.isCorner(textInput, rowId, colId) {
				arr[colId+rowNumber] = d.cornerColor.R
				arr[colId+1+rowNumber] = d.cornerColor.G
				arr[colId+2+rowNumber] = d.cornerColor.B
				arr[colId+3+rowNumber] = d.cornerColor.A
			} else if d.isBorder(textInput, rowId, colId) {
				arr[colId+rowNumber] = borderColor.R
				arr[colId+1+rowNumber] = borderColor.G
				arr[colId+2+rowNumber] = borderColor.B
				arr[colId+3+rowNumber] = borderColor.A
			} else {
				arr[colId+rowNumber] = bgColor.R
				arr[colId+1+rowNumber] = bgColor.G
				arr[colId+2+rowNumber] = bgColor.B
				arr[colId+3+rowNumber] = bgColor.A
			}
		}
	}

	return arr
}
