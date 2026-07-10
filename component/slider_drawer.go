package component

import (
	"image/color"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type SliderDrawer interface {
	Draw(slider *Slider) *ebiten.Image
}

type DefaultSliderDrawer struct {
	Color         color.RGBA
	ColorPressed  color.RGBA
	ColorHovered  color.RGBA
	ColorDisabled color.RGBA

	pixelBuf []byte
	bgRow    []byte
}

func (d DefaultSliderDrawer) Draw(slider *Slider) *ebiten.Image {
	borderColor := d.Color

	switch {
	case slider.pressed:
		borderColor = d.ColorPressed
	case slider.hovering:
		borderColor = d.ColorHovered
	case slider.disabled:
		borderColor = d.ColorDisabled
	}

	slider.image.WritePixels(d.drawWithColor(slider, borderColor))
	return slider.image
}

func (d *DefaultSliderDrawer) getBuffer(slider *Slider) []byte {
	size := slider.pixelRows * slider.pixelCols

	if len(d.pixelBuf) != size {
		d.pixelBuf = make([]byte, size)
	}

	for i := range d.pixelBuf {
		d.pixelBuf[i] = 0
	}

	return d.pixelBuf
}

func (d *DefaultSliderDrawer) getBgRow(slider *Slider) []byte {
	if len(d.bgRow) != slider.pixelCols {
		d.bgRow = make([]byte, slider.pixelCols)
	}

	bgColor := slider.container.GetBackgroundColor()
	if len(d.bgRow) >= 4 && d.bgRow[0] == bgColor.R && d.bgRow[1] == bgColor.G && d.bgRow[2] == bgColor.B && d.bgRow[3] == bgColor.A {
		return d.bgRow
	}

	for i := 0; i < len(d.bgRow); i += 4 {
		d.bgRow[i], d.bgRow[i+1], d.bgRow[i+2], d.bgRow[i+3] = bgColor.R, bgColor.G, bgColor.B, bgColor.A
	}

	return d.bgRow
}

func (d *DefaultSliderDrawer) drawWithColor(slider *Slider, borderColor color.RGBA) []byte {
	arr := d.getBuffer(slider)
	bgRow := d.getBgRow(slider)

	for rowId := slider.firstPixelRowId; rowId <= slider.lastPixelRowId; rowId++ {
		copy(arr[slider.pixelCols*rowId:], bgRow)
		isFirstOrLastRow := rowId == slider.firstPixelRowId || rowId == slider.lastPixelRowId
		rowInsideColoredSection := rowId > slider.secondPixelRowId && rowId < slider.penultimatePixelRowId

		rowNumber := slider.pixelCols * rowId
		for colId := slider.firstPixelColId; colId <= slider.lastPixelColId; colId += 4 {
			isFirstOrLastCol := colId == slider.firstPixelColId || colId == slider.lastPixelColId
			colInsideColoredSection := colId > slider.secondPixelColId && colId < slider.penultimatePixelColId

			if isFirstOrLastRow != isFirstOrLastCol { // border without corners
				d.setPixel(arr, rowNumber, colId, borderColor)
			} else if rowInsideColoredSection && colInsideColoredSection && colId <= int(slider.handle.posX)*4 { // colored section
				d.setPixel(arr, rowNumber, colId, borderColor)
			}
		}
	}

	return arr
}

func (d *DefaultSliderDrawer) setPixel(arr []byte, rowNumber, colId int, color color.RGBA) {
	arr[colId+rowNumber] = color.R
	arr[colId+1+rowNumber] = color.G
	arr[colId+2+rowNumber] = color.B
	arr[colId+3+rowNumber] = color.A
}
