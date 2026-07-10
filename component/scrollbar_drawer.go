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
}

func (d DefaultScrollBarDrawer) Draw(ScrollBar *ScrollBar) *ebiten.Image {
	if ScrollBar.pressed {
		ScrollBar.image.WritePixels(d.drawPressed(ScrollBar))
	} else if ScrollBar.hovering {
		ScrollBar.image.WritePixels(d.drawHovered(ScrollBar))
	} else if ScrollBar.disabled {
		ScrollBar.image.WritePixels(d.drawDisabled(ScrollBar))
	} else {
		ScrollBar.image.WritePixels(d.draw(ScrollBar))
	}

	return ScrollBar.image
}

func (d *DefaultScrollBarDrawer) draw(ScrollBar *ScrollBar) []byte {
	arr := make([]byte, ScrollBar.pixelRows*ScrollBar.pixelCols)
	backgroundColor := ScrollBar.container.GetBackgroundColor()

	for rowId := ScrollBar.firstPixelRowId; rowId <= ScrollBar.lastPixelRowId; rowId++ {
		rowNumber := ScrollBar.pixelCols * rowId

		for colId := ScrollBar.firstPixelColId; colId <= ScrollBar.lastPixelColId; colId += 4 {
			if d.isCorner(ScrollBar, rowId, colId) {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			} else if d.isBorder(ScrollBar, rowId, colId) {
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

func (d *DefaultScrollBarDrawer) drawPressed(ScrollBar *ScrollBar) []byte {
	arr := make([]byte, ScrollBar.pixelRows*ScrollBar.pixelCols)
	backgroundColor := ScrollBar.container.GetBackgroundColor()

	for rowId := ScrollBar.firstPixelRowId; rowId <= ScrollBar.lastPixelRowId; rowId++ {
		rowNumber := ScrollBar.pixelCols * rowId

		for colId := ScrollBar.firstPixelColId; colId <= ScrollBar.lastPixelColId; colId += 4 {
			if d.isCorner(ScrollBar, rowId, colId) {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			} else if d.isBorder(ScrollBar, rowId, colId) {
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

func (d *DefaultScrollBarDrawer) drawHovered(ScrollBar *ScrollBar) []byte {
	arr := make([]byte, ScrollBar.pixelRows*ScrollBar.pixelCols)
	backgroundColor := ScrollBar.container.GetBackgroundColor()

	for rowId := ScrollBar.firstPixelRowId; rowId <= ScrollBar.lastPixelRowId; rowId++ {
		rowNumber := ScrollBar.pixelCols * rowId

		for colId := ScrollBar.firstPixelColId; colId <= ScrollBar.lastPixelColId; colId += 4 {
			if d.isCorner(ScrollBar, rowId, colId) {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			} else if d.isBorder(ScrollBar, rowId, colId) {
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

func (d *DefaultScrollBarDrawer) drawDisabled(ScrollBar *ScrollBar) []byte {
	arr := make([]byte, ScrollBar.pixelRows*ScrollBar.pixelCols)
	backgroundColor := ScrollBar.container.GetBackgroundColor()

	for rowId := ScrollBar.firstPixelRowId; rowId <= ScrollBar.lastPixelRowId; rowId++ {
		rowNumber := ScrollBar.pixelCols * rowId

		for colId := ScrollBar.firstPixelColId; colId <= ScrollBar.lastPixelColId; colId += 4 {
			if d.isCorner(ScrollBar, rowId, colId) {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			} else if d.isBorder(ScrollBar, rowId, colId) {
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

func (d DefaultScrollBarDrawer) isCorner(ScrollBar *ScrollBar, rowId, colId int) bool {
	return (rowId == ScrollBar.firstPixelRowId || rowId == ScrollBar.lastPixelRowId) && (colId == ScrollBar.firstPixelColId || colId == ScrollBar.lastPixelColId)
}

func (d DefaultScrollBarDrawer) isBorder(ScrollBar *ScrollBar, rowId, colId int) bool {
	return rowId == ScrollBar.firstPixelRowId || rowId == ScrollBar.lastPixelRowId || colId == ScrollBar.firstPixelColId || colId == ScrollBar.lastPixelColId
}
