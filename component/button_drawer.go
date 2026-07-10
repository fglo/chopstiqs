package component

import (
	"image/color"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type ButtonDrawer interface {
	Draw(*Button) *ebiten.Image
}

type DefaultButtonDrawer struct {
	Color         color.RGBA
	ColorPressed  color.RGBA
	ColorHovered  color.RGBA
	ColorDisabled color.RGBA

	pixelBuf []byte
	bgRow    []byte // background
}

func (d DefaultButtonDrawer) Draw(bttn *Button) *ebiten.Image {
	borderColor := d.Color
	insidesColor := &d.Color

	switch {
	case bttn.pressed:
		borderColor = d.ColorPressed
		insidesColor = nil
	case bttn.hovering:
		borderColor = d.ColorHovered
		insidesColor = &d.ColorHovered
	case bttn.disabled:
		borderColor = d.ColorDisabled
		insidesColor = &d.ColorDisabled
	}

	bttn.image.WritePixels(d.drawWithColor(bttn, borderColor, insidesColor))

	return bttn.image
}

func (d *DefaultButtonDrawer) getBuffer(bttn *Button) []byte {
	size := bttn.pixelRows * bttn.pixelCols

	if len(d.pixelBuf) != size {
		d.pixelBuf = make([]byte, size)
	}

	// zero it out for reuse
	for i := range d.pixelBuf {
		d.pixelBuf[i] = 0
	}
	return d.pixelBuf
}

func (d *DefaultButtonDrawer) getBgRow(bttn *Button) []byte {
	if len(d.bgRow) != bttn.pixelCols {
		d.bgRow = make([]byte, bttn.pixelCols)
	}

	bgColor := bttn.container.GetBackgroundColor()
	if d.bgRow[0] == bgColor.R &&
		d.bgRow[1] == bgColor.G &&
		d.bgRow[2] == bgColor.B &&
		d.bgRow[3] == bgColor.A {
		return d.bgRow
	}

	// fill it with background Color
	for i := range d.bgRow {
		d.bgRow[i], d.bgRow[i+1], d.bgRow[i+2], d.bgRow[i+3] = bgColor.R, bgColor.G, bgColor.B, bgColor.A
	}
	return d.bgRow
}

func (d *DefaultButtonDrawer) drawWithColor(bttn *Button, borderColor color.RGBA, insidesColor *color.RGBA) []byte {
	arr := d.getBuffer(bttn)
	bgRow := d.getBgRow(bttn)

	for rowId := bttn.firstPixelRowId; rowId <= bttn.lastPixelRowId; rowId++ {
		copy(arr[bttn.pixelCols*rowId:], bgRow)

		rowNumber := bttn.pixelCols * rowId
		isFirstOrLastRow := rowId == bttn.firstPixelRowId || rowId == bttn.lastPixelRowId
		rowInsideColoredSection := rowId > bttn.secondPixelRowId && rowId < bttn.penultimatePixelRowId

		for colId := bttn.firstPixelColId; colId <= bttn.lastPixelColId; colId += 4 {
			isFirstOrLastCol := colId == bttn.firstPixelColId || colId == bttn.lastPixelColId
			colInsideColoredSection := colId > bttn.secondPixelColId && colId < bttn.penultimatePixelColId

			if (isFirstOrLastRow || isFirstOrLastCol) && !(isFirstOrLastRow && isFirstOrLastCol) { // border
				arr[colId+rowNumber] = borderColor.R
				arr[colId+1+rowNumber] = borderColor.G
				arr[colId+2+rowNumber] = borderColor.B
				arr[colId+3+rowNumber] = borderColor.A
			} else if rowInsideColoredSection && colInsideColoredSection && insidesColor != nil { // insides
				arr[colId+rowNumber] = insidesColor.R
				arr[colId+1+rowNumber] = insidesColor.G
				arr[colId+2+rowNumber] = insidesColor.B
				arr[colId+3+rowNumber] = insidesColor.A
			}
		}
	}

	return arr
}
