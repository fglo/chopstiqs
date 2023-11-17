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
	if slider.pressed {
		slider.image.WritePixels(d.drawPressed(slider))
	} else if slider.hovering {
		slider.image.WritePixels(d.drawHovered(slider))
	} else if slider.disabled {
		slider.image.WritePixels(d.drawDisabled(slider))
	} else {
		slider.image.WritePixels(d.draw(slider))
	}

	return slider.image
}

func (d *DefaultSliderDrawer) draw(slider *Slider) []byte {
	arr := make([]byte, slider.pixelRows*slider.pixelCols)
	backgroundColor := slider.container.GetBackgroundColor()

	for rowId := slider.firstPixelRowId; rowId <= slider.lastPixelRowId; rowId++ {
		rowNumber := slider.pixelCols * rowId

		for colId := slider.firstPixelColId; colId <= slider.lastPixelColId; colId += 4 {
			if d.isCorner(slider, rowId, colId) {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			} else if d.isBorder(slider, rowId, colId) || (d.isColored(slider, rowId, colId) && colId <= int(slider.handle.posX)*4) {
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

	return arr
}

func (d *DefaultSliderDrawer) drawPressed(slider *Slider) []byte {
	arr := make([]byte, slider.pixelRows*slider.pixelCols)
	backgroundColor := slider.container.GetBackgroundColor()

	for rowId := slider.firstPixelRowId; rowId <= slider.lastPixelRowId; rowId++ {
		rowNumber := slider.pixelCols * rowId

		for colId := slider.firstPixelColId; colId <= slider.lastPixelColId; colId += 4 {
			if d.isCorner(slider, rowId, colId) {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			} else if d.isBorder(slider, rowId, colId) || (d.isColored(slider, rowId, colId) && colId <= int(slider.handle.posX)*4) {
				arr[colId+rowNumber] = d.ColorPressed.R
				arr[colId+1+rowNumber] = d.ColorPressed.G
				arr[colId+2+rowNumber] = d.ColorPressed.B
				arr[colId+3+rowNumber] = d.ColorPressed.A
			} else {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			}
		}
	}

	return arr
}

func (d *DefaultSliderDrawer) drawHovered(slider *Slider) []byte {
	arr := make([]byte, slider.pixelRows*slider.pixelCols)
	backgroundColor := slider.container.GetBackgroundColor()

	for rowId := slider.firstPixelRowId; rowId <= slider.lastPixelRowId; rowId++ {
		rowNumber := slider.pixelCols * rowId

		for colId := slider.firstPixelColId; colId <= slider.lastPixelColId; colId += 4 {
			if d.isCorner(slider, rowId, colId) {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			} else if d.isBorder(slider, rowId, colId) || (d.isColored(slider, rowId, colId) && colId <= int(slider.handle.posX)*4) {
				arr[colId+rowNumber] = d.ColorHovered.R
				arr[colId+1+rowNumber] = d.ColorHovered.G
				arr[colId+2+rowNumber] = d.ColorHovered.B
				arr[colId+3+rowNumber] = d.ColorHovered.A
			} else {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			}
		}
	}

	return arr
}

func (d *DefaultSliderDrawer) drawDisabled(slider *Slider) []byte {
	arr := make([]byte, slider.pixelRows*slider.pixelCols)
	backgroundColor := slider.container.GetBackgroundColor()

	for rowId := slider.firstPixelRowId; rowId <= slider.lastPixelRowId; rowId++ {
		rowNumber := slider.pixelCols * rowId

		for colId := slider.firstPixelColId; colId <= slider.lastPixelColId; colId += 4 {
			if d.isCorner(slider, rowId, colId) {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			} else if d.isBorder(slider, rowId, colId) || (d.isColored(slider, rowId, colId) && colId <= int(slider.handle.posX)*4) {
				arr[colId+rowNumber] = d.ColorDisabled.R
				arr[colId+1+rowNumber] = d.ColorDisabled.G
				arr[colId+2+rowNumber] = d.ColorDisabled.B
				arr[colId+3+rowNumber] = d.ColorDisabled.A
			} else {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			}
		}
	}

	return arr
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
