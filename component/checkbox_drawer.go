package component

import (
	"image/color"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type CheckBoxDrawer interface {
	Draw(checkbox *CheckBox) *ebiten.Image
}

type DefaultCheckBoxDrawer struct {
	Color color.RGBA

	pixelBuf []byte
	bgRow    []byte
}

func (d DefaultCheckBoxDrawer) Draw(cb *CheckBox) *ebiten.Image {
	if cb.Checked() {
		cb.image.WritePixels(d.draw(cb, true))
	} else {
		cb.image.WritePixels(d.draw(cb, false))
	}

	return cb.image
}

func (d *DefaultCheckBoxDrawer) getBuffer(cb *CheckBox) []byte {
	size := cb.pixelRows * cb.pixelCols

	if len(d.pixelBuf) != size {
		d.pixelBuf = make([]byte, size)
	}

	for i := range d.pixelBuf {
		d.pixelBuf[i] = 0
	}

	return d.pixelBuf
}

func (d *DefaultCheckBoxDrawer) getBgRow(cb *CheckBox) []byte {
	if len(d.bgRow) != cb.pixelCols {
		d.bgRow = make([]byte, cb.pixelCols)
	}

	bgColor := cb.container.GetBackgroundColor()
	if len(d.bgRow) >= 4 && d.bgRow[0] == bgColor.R && d.bgRow[1] == bgColor.G && d.bgRow[2] == bgColor.B && d.bgRow[3] == bgColor.A {
		return d.bgRow
	}

	for i := 0; i < len(d.bgRow); i += 4 {
		d.bgRow[i], d.bgRow[i+1], d.bgRow[i+2], d.bgRow[i+3] = bgColor.R, bgColor.G, bgColor.B, bgColor.A
	}

	return d.bgRow
}

func (d *DefaultCheckBoxDrawer) draw(cb *CheckBox, checked bool) []byte {
	arr := d.getBuffer(cb)
	bgRow := d.getBgRow(cb)

	for rowId := cb.firstPixelRowId; rowId <= cb.lastPixelRowId; rowId++ {
		copy(arr[cb.pixelCols*rowId:], bgRow)
		isFirstOrLastRow := rowId == cb.firstPixelRowId || rowId == cb.lastPixelRowId
		rowInsideColoredSection := rowId > cb.secondPixelRowId && rowId < cb.penultimatePixelRowId

		rowNumber := cb.pixelCols * rowId
		for colId := cb.firstPixelColId; colId <= cb.lastPixelColId; colId += 4 {
			isFirstOrLastCol := colId == cb.firstPixelColId || colId == cb.lastPixelColId
			colInsideColoredSection := colId > cb.secondPixelColId && colId < cb.penultimatePixelColId

			if isFirstOrLastRow || isFirstOrLastCol { // border
				d.setPixel(arr, rowNumber, colId, d.Color)
			} else if checked && colInsideColoredSection && rowInsideColoredSection { // insides
				d.setPixel(arr, rowNumber, colId, d.Color)
			}
		}
	}

	return arr
}

func (d *DefaultCheckBoxDrawer) setPixel(arr []byte, rowNumber, colId int, color color.RGBA) {
	arr[colId+rowNumber] = color.R
	arr[colId+1+rowNumber] = color.G
	arr[colId+2+rowNumber] = color.B
	arr[colId+3+rowNumber] = color.A
}
