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
}

func (d DefaultButtonDrawer) Draw(bttn *Button) *ebiten.Image {
	if bttn.pressed {
		bttn.image.WritePixels(d.drawPressed(bttn))
	} else if bttn.hovering {
		bttn.image.WritePixels(d.drawHovered(bttn))
	} else if bttn.disabled {
		bttn.image.WritePixels(d.drawDisabled(bttn))
	} else {
		bttn.image.WritePixels(d.draw(bttn))
	}

	return bttn.image
}

func (d *DefaultButtonDrawer) isCorner(bttn *Button, rowId, colId int) bool {
	return (rowId == bttn.firstPixelRowId || rowId == bttn.lastPixelRowId) && (colId == bttn.firstPixelColId || colId == bttn.lastPixelColId)
}

func (d *DefaultButtonDrawer) isBorder(bttn *Button, rowId, colId int) bool {
	return rowId == bttn.firstPixelRowId || rowId == bttn.lastPixelRowId || colId == bttn.firstPixelColId || colId == bttn.lastPixelColId
}

func (d *DefaultButtonDrawer) isColored(bttn *Button, rowId, colId int) bool {
	return colId > bttn.secondPixelColId && colId < bttn.penultimatePixelColId && rowId > bttn.secondPixelRowId && rowId < bttn.penultimatePixelRowId
}

func (d *DefaultButtonDrawer) draw(bttn *Button) []byte {
	arr := make([]byte, bttn.pixelRows*bttn.pixelCols)
	backgroundColor := bttn.container.GetBackgroundColor()

	for rowId := bttn.firstPixelRowId; rowId <= bttn.lastPixelRowId; rowId++ {
		rowNumber := bttn.pixelCols * rowId

		for colId := bttn.firstPixelColId; colId <= bttn.lastPixelColId; colId += 4 {
			if d.isCorner(bttn, rowId, colId) {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			} else if d.isBorder(bttn, rowId, colId) || d.isColored(bttn, rowId, colId) {
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

func (d *DefaultButtonDrawer) drawPressed(bttn *Button) []byte {
	arr := make([]byte, bttn.pixelRows*bttn.pixelCols)
	backgroundColor := bttn.container.GetBackgroundColor()

	for rowId := bttn.firstPixelRowId; rowId <= bttn.lastPixelRowId; rowId++ {
		rowNumber := bttn.pixelCols * rowId

		for colId := bttn.firstPixelColId; colId <= bttn.lastPixelColId; colId += 4 {
			if d.isCorner(bttn, rowId, colId) {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			} else if d.isBorder(bttn, rowId, colId) {
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

func (d *DefaultButtonDrawer) drawHovered(bttn *Button) []byte {
	arr := make([]byte, bttn.pixelRows*bttn.pixelCols)
	backgroundColor := bttn.container.GetBackgroundColor()

	for rowId := bttn.firstPixelRowId; rowId <= bttn.lastPixelRowId; rowId++ {
		rowNumber := bttn.pixelCols * rowId

		for colId := bttn.firstPixelColId; colId <= bttn.lastPixelColId; colId += 4 {
			if d.isCorner(bttn, rowId, colId) {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			} else if d.isBorder(bttn, rowId, colId) || d.isColored(bttn, rowId, colId) {
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

func (d *DefaultButtonDrawer) drawDisabled(bttn *Button) []byte {
	arr := make([]byte, bttn.pixelRows*bttn.pixelCols)
	backgroundColor := bttn.container.GetBackgroundColor()

	for rowId := bttn.firstPixelRowId; rowId <= bttn.lastPixelRowId; rowId++ {
		rowNumber := bttn.pixelCols * rowId

		for colId := bttn.firstPixelColId; colId <= bttn.lastPixelColId; colId += 4 {
			if d.isCorner(bttn, rowId, colId) {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			} else if d.isBorder(bttn, rowId, colId) || d.isColored(bttn, rowId, colId) {
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
