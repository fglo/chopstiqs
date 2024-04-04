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
}

func (d DefaultCheckBoxDrawer) Draw(cb *CheckBox) *ebiten.Image {
	if cb.Checked() {
		cb.image.WritePixels(d.drawChecked(cb))
	} else {
		cb.image.WritePixels(d.drawUnchecked(cb))
	}

	return cb.image
}

func (d DefaultCheckBoxDrawer) isBorder(cb *CheckBox, rowId, colId int) bool {
	return rowId == cb.firstPixelRowId || rowId == cb.lastPixelRowId || colId == cb.firstPixelColId || colId == cb.lastPixelColId
}

func (d DefaultCheckBoxDrawer) isColored(cb *CheckBox, rowId, colId int) bool {
	return colId > cb.secondPixelColId && colId < cb.penultimatePixelColId && rowId > cb.secondPixelRowId && rowId < cb.penultimatePixelRowId
}

func (d DefaultCheckBoxDrawer) drawUnchecked(cb *CheckBox) []byte {
	arr := make([]byte, cb.component.pixelRows*cb.component.pixelCols)
	backgroundColor := cb.container.GetBackgroundColor()

	for rowId := cb.firstPixelRowId; rowId <= cb.lastPixelRowId; rowId++ {
		rowNumber := cb.component.pixelCols * rowId

		for colId := cb.firstPixelColId; colId <= cb.lastPixelColId; colId += 4 {
			if d.isBorder(cb, rowId, colId) {
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

func (d DefaultCheckBoxDrawer) drawChecked(cb *CheckBox) []byte {
	arr := make([]byte, cb.component.pixelRows*cb.component.pixelCols)
	backgroundColor := cb.container.GetBackgroundColor()

	for rowId := cb.firstPixelRowId; rowId <= cb.lastPixelRowId; rowId++ {
		rowNumber := cb.component.pixelCols * rowId

		for colId := cb.firstPixelColId; colId <= cb.lastPixelColId; colId += 4 {
			if d.isBorder(cb, rowId, colId) || d.isColored(cb, rowId, colId) {
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
