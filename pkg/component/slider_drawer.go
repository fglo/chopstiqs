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
}

func (d DefaultSliderDrawer) Draw(slider *Slider) *ebiten.Image {
	arr := make([]byte, slider.pixelRows*slider.pixelCols)
	backgroundColor := slider.container.GetBackgroundColor()

	for rowId := slider.firstPixelRowId; rowId <= slider.lastPixelRowId; rowId++ {
		rowNumber := slider.pixelCols * rowId

		for colId := slider.firstPixelColId; colId <= slider.lastPixelColId; colId += 4 {
			if slider.isCorner(rowId, colId) {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			} else if slider.isBorder(rowId, colId) || (slider.isColored(rowId, colId) && colId <= int(slider.handle.posX)*4) {
				arr[colId+rowNumber] = d.Color.R
				arr[colId+1+rowNumber] = d.Color.G
				arr[colId+2+rowNumber] = d.Color.B
				arr[colId+3+rowNumber] = d.Color.A
			} else {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			}
		}
	}

	slider.image.WritePixels(arr)

	return slider.image
}

func (d DefaultSliderDrawer) isCorner(slider *Slider, rowId, colId int) bool {
	return (rowId == slider.firstPixelRowId || rowId == slider.lastPixelRowId) && (colId == slider.firstPixelColId || colId == slider.lastPixelColId)
}

func (d DefaultSliderDrawer) isBorder(slider *Slider, rowId, colId int) bool {
	return rowId == slider.firstPixelRowId || rowId == slider.lastPixelRowId || colId == slider.firstPixelColId || colId == slider.lastPixelColId
}

func (d DefaultSliderDrawer) isColored(slider *Slider, rowId, colId int) bool {
	return colId > slider.secondPixelColId && colId < slider.penultimatePixelColId && rowId > slider.secondPixelRowId && rowId < slider.penultimatePixelRowId
}
