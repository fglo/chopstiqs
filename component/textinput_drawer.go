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
	backgroundColor color.RGBA
	invertedBgColor color.RGBA

	pixelBuf []byte
	bgRow    []byte
}

func (d *DefaultTextInputDrawer) Draw(textInput *TextInput) *ebiten.Image {
	d.backgroundColor = textInput.container.GetBackgroundColor()
	d.invertedBgColor = colorutils.Invert(d.backgroundColor)

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

func (d *DefaultTextInputDrawer) getBuffer(textInput *TextInput) []byte {
	size := textInput.pixelRows * textInput.pixelCols

	if len(d.pixelBuf) != size {
		d.pixelBuf = make([]byte, size)
	}

	for i := range d.pixelBuf {
		d.pixelBuf[i] = 0
	}

	return d.pixelBuf
}

func (d *DefaultTextInputDrawer) getBgRow(textInput *TextInput) []byte {
	if len(d.bgRow) != textInput.pixelCols {
		d.bgRow = make([]byte, textInput.pixelCols)
	}

	bgColor := textInput.container.GetBackgroundColor()
	if len(d.bgRow) >= 4 && d.bgRow[0] == bgColor.R && d.bgRow[1] == bgColor.G && d.bgRow[2] == bgColor.B && d.bgRow[3] == bgColor.A {
		return d.bgRow
	}

	for i := 0; i < len(d.bgRow); i += 4 {
		d.bgRow[i], d.bgRow[i+1], d.bgRow[i+2], d.bgRow[i+3] = bgColor.R, bgColor.G, bgColor.B, bgColor.A
	}

	return d.bgRow
}

func (d *DefaultTextInputDrawer) draw(textInput *TextInput, borderColor color.RGBA) []byte {
	arr := d.getBuffer(textInput)
	bgRow := d.getBgRow(textInput)

	selectingFromColId := -1
	selectingToColId := -1

	if textInput.HasSelectedText() {
		selectingFromColId = (textInput.possibleCursorPosXs[textInput.selectionStart] - textInput.scrollOffset + textInput.padding.Left + textInput.cursor.Width()) * 4
		selectingToColId = (textInput.possibleCursorPosXs[textInput.selectionEnd] - textInput.scrollOffset + textInput.padding.Left + textInput.cursor.Width()) * 4
	}

	for rowId := textInput.firstPixelRowId; rowId <= textInput.lastPixelRowId; rowId++ {
		copy(arr[textInput.pixelCols*rowId:], bgRow)
		isFirstOrLastRow := rowId == textInput.firstPixelRowId || rowId == textInput.lastPixelRowId
		insideTextInput := textInput.firstPixelRowId+1 < rowId && rowId < textInput.lastPixelRowId-1

		rowNumber := textInput.pixelCols * rowId
		for colId := textInput.firstPixelColId; colId <= textInput.lastPixelColId; colId += 4 {
			isFirstOrLastCol := colId == textInput.firstPixelColId || colId == textInput.lastPixelColId
			insideSelectedText := selectingFromColId+4 < colId && colId <= selectingToColId

			if isFirstOrLastRow != isFirstOrLastCol { // border
				d.setPixel(arr, rowNumber, colId, borderColor)
			} else if insideSelectedText && insideTextInput { // selected text background
				d.setPixel(arr, rowNumber, colId, d.invertedBgColor)
			}
		}
	}

	return arr
}

func (d *DefaultTextInputDrawer) setPixel(arr []byte, rowNumber, colId int, color color.RGBA) {
	arr[colId+rowNumber] = color.R
	arr[colId+1+rowNumber] = color.G
	arr[colId+2+rowNumber] = color.B
	arr[colId+3+rowNumber] = color.A
}
