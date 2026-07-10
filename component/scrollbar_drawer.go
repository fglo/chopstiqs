package component

import (
	"image/color"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type ScrollBarDrawer interface {
	Draw(ScrollBar *ScrollBar) *ebiten.Image
}

type DefaultScrollBarDrawer struct {
	Color         color.RGBA
	ColorPressed  color.RGBA
	ColorHovered  color.RGBA
	ColorDisabled color.RGBA

	pixelBuf []byte
	bgRow    []byte
}

func (d DefaultScrollBarDrawer) Draw(ScrollBar *ScrollBar) *ebiten.Image {
	borderColor := d.Color

	switch {
	case ScrollBar.pressed:
		borderColor = d.ColorPressed
	case ScrollBar.hovering:
		borderColor = d.ColorHovered
	case ScrollBar.disabled:
		borderColor = d.ColorDisabled
	}

	ScrollBar.image.WritePixels(d.drawWithColor(ScrollBar, borderColor))
	return ScrollBar.image
}

func (d *DefaultScrollBarDrawer) getBuffer(ScrollBar *ScrollBar) []byte {
	size := ScrollBar.pixelRows * ScrollBar.pixelCols

	if len(d.pixelBuf) != size {
		d.pixelBuf = make([]byte, size)
	}

	for i := range d.pixelBuf {
		d.pixelBuf[i] = 0
	}

	return d.pixelBuf
}

func (d *DefaultScrollBarDrawer) getBgRow(ScrollBar *ScrollBar) []byte {
	if len(d.bgRow) != ScrollBar.pixelCols {
		d.bgRow = make([]byte, ScrollBar.pixelCols)
	}

	bgColor := ScrollBar.container.GetBackgroundColor()
	if len(d.bgRow) >= 4 && d.bgRow[0] == bgColor.R && d.bgRow[1] == bgColor.G && d.bgRow[2] == bgColor.B && d.bgRow[3] == bgColor.A {
		return d.bgRow
	}

	for i := 0; i < len(d.bgRow); i += 4 {
		d.bgRow[i], d.bgRow[i+1], d.bgRow[i+2], d.bgRow[i+3] = bgColor.R, bgColor.G, bgColor.B, bgColor.A
	}

	return d.bgRow
}

func (d *DefaultScrollBarDrawer) drawWithColor(ScrollBar *ScrollBar, borderColor color.RGBA) []byte {
	arr := d.getBuffer(ScrollBar)
	bgRow := d.getBgRow(ScrollBar)

	for rowId := ScrollBar.firstPixelRowId; rowId <= ScrollBar.lastPixelRowId; rowId++ {
		copy(arr[ScrollBar.pixelCols*rowId:], bgRow)
		isFirstOrLastRow := rowId == ScrollBar.firstPixelRowId || rowId == ScrollBar.lastPixelRowId

		rowNumber := ScrollBar.pixelCols * rowId
		for colId := ScrollBar.firstPixelColId; colId <= ScrollBar.lastPixelColId; colId += 4 {
			isFirstOrLastCol := colId == ScrollBar.firstPixelColId || colId == ScrollBar.lastPixelColId
			if isFirstOrLastRow != isFirstOrLastCol { // border without corners
				d.setPixel(arr, rowNumber, colId, borderColor)
			}
		}
	}

	return arr
}

func (d *DefaultScrollBarDrawer) setPixel(arr []byte, rowNumber, colId int, color color.RGBA) {
	arr[colId+rowNumber] = color.R
	arr[colId+1+rowNumber] = color.G
	arr[colId+2+rowNumber] = color.B
	arr[colId+3+rowNumber] = color.A
}
