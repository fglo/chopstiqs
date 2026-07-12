package component

import (
	"github.com/fglo/chopstiqs/option"
)

type Layout interface {
	Rearrange(*Container)
	Arrange(*Container, Component)
}

type HorizontalListLayout struct {
	ColumnGap int
}

func (hl *HorizontalListLayout) Rearrange(c *Container) {
	width := 0
	height := 0

	c.lastComponentPosX = 0

	widthSet := 0
	strechedHorizontally := 0

	for _, component := range c.components {
		if component.HeightWithPadding() > height {
			height = component.HeightWithPadding()
		}

		if component.HorizontalAlignment() != option.StretchedHorizontally {
			widthSet += component.Width()
		} else {
			strechedHorizontally += 1
		}
	}

	stretchedWidth := 0
	if strechedHorizontally > 0 {
		stretchedWidth = (c.width - widthSet) / strechedHorizontally
	}

	for _, component := range c.components {
		if component.HorizontalAlignment() == option.StretchedHorizontally && stretchedWidth > 0 {
			component.SetWidth(stretchedWidth)
		}
		if component.VerticalAlignment() == option.StretchedVertically {
			component.SetHeight(height)
		}
		width, height = hl.arrange(c, component, height)
	}

	c.SetDimensions(width, height)
}

func (hl *HorizontalListLayout) Arrange(c *Container, component Component) {
	width, height := hl.arrange(c, component, c.height)
	c.SetDimensions(width, height)
}

func (hl *HorizontalListLayout) arrange(c *Container, component Component, height int) (int, int) {
	if component.HeightWithPadding() > height {
		height = component.HeightWithPadding()
	}

	dy := 0.

	switch component.VerticalAlignment() {
	case option.AlignmentBottom:
		dy = float64(height - component.HeightWithPadding())
	case option.AlignmentCenteredVertically:
		dy = float64(height-component.HeightWithPadding()) / 2
	}

	component.SetPosition(float64(c.padding.Left+c.lastComponentPosX), float64(c.padding.Top)+dy)

	c.lastComponentPosX += component.WidthWithPadding() + hl.ColumnGap

	return c.lastComponentPosX, height
}

type VerticalListLayout struct {
	RowGap int
}

func (vl *VerticalListLayout) Rearrange(c *Container) {
	width := 0
	height := 0

	c.lastComponentPosY = 0.

	heightSet := 0
	strechedVertically := 0

	for _, component := range c.components {
		if component.WidthWithPadding() > width {
			width = component.WidthWithPadding()
		}

		if component.VerticalAlignment() != option.StretchedVertically {
			heightSet += component.HeightWithPadding()
		} else {
			strechedVertically += 1
		}
	}

	stretchedHeight := 0
	if strechedVertically > 0 {
		stretchedHeight = (c.height - heightSet) / strechedVertically
	}

	for _, component := range c.components {
		if component.VerticalAlignment() == option.StretchedVertically && stretchedHeight > 0 {
			component.SetHeight(stretchedHeight)
		}
		if component.HorizontalAlignment() == option.StretchedHorizontally {
			component.SetWidth(width)
		}
		width, height = vl.arrange(c, component, width)
	}

	c.SetDimensions(width, height)
}

func (vl *VerticalListLayout) Arrange(c *Container, component Component) {
	width, height := vl.arrange(c, component, c.width)
	c.SetDimensions(width, height)
}

func (vl *VerticalListLayout) arrange(c *Container, component Component, width int) (int, int) {
	if component.WidthWithPadding() > width {
		width = component.WidthWithPadding()
	}

	dx := 0.

	switch component.HorizontalAlignment() {
	case option.AlignmentRight:
		dx = float64(width - component.WidthWithPadding())
	case option.AlignmentCenteredHorizontally:
		dx = float64(width-component.WidthWithPadding()) / 2
	}

	component.SetPosition(float64(c.padding.Left)+dx, float64(c.lastComponentPosY+c.padding.Top))

	c.lastComponentPosY += component.HeightWithPadding() + vl.RowGap

	return width, c.lastComponentPosY
}

type GridLayout struct {
	Columns       int
	ColumnsWidths []int
	ColumnGap     int

	Rows        int
	RowsHeights []int
	RowGap      int

	fixedColumnWidths bool
	fixedRowHeights   bool
}

func (gl *GridLayout) Setup() {
	if len(gl.ColumnsWidths) == 0 {
		gl.ColumnsWidths = make([]int, gl.Columns)
	} else {
		gl.Columns = len(gl.ColumnsWidths)
		gl.fixedColumnWidths = true
	}

	if len(gl.RowsHeights) == 0 {
		gl.RowsHeights = make([]int, gl.Rows)
	} else {
		gl.Rows = len(gl.RowsHeights)
		gl.fixedRowHeights = true
	}
}

func (gl *GridLayout) Rearrange(c *Container) {
	width := 0
	height := 0

	c.lastComponentPosX = 0.
	c.lastComponentPosY = 0.

	currColId := 0
	currRowId := 0

	for _, component := range c.components {
		if currRowId >= gl.Rows {
			component.SetHidden(true)
			continue
		}

		if !gl.fixedColumnWidths && (gl.ColumnsWidths[currColId] == 0 || gl.ColumnsWidths[currColId] < component.WidthWithPadding()) {
			gl.ColumnsWidths[currColId] = component.WidthWithPadding()
		}

		if !gl.fixedRowHeights && (gl.RowsHeights[currRowId] == 0 || gl.RowsHeights[currRowId] < component.HeightWithPadding()) {
			gl.RowsHeights[currRowId] = component.HeightWithPadding()
		}

		currColId++
		if currColId == gl.Columns {
			currColId = 0
			currRowId++
		}
	}

	for _, colwidth := range gl.ColumnsWidths {
		width += colwidth + c.padding.Left + c.padding.Right
	}

	for _, rowheight := range gl.RowsHeights {
		height += rowheight + c.padding.Top + c.padding.Bottom
	}

	currColId = 0
	currRowId = 0

	for _, component := range c.components {
		if currRowId >= gl.Rows {
			break
		}

		component.SetPosition(float64(c.padding.Left+c.lastComponentPosX), float64(c.padding.Top+c.lastComponentPosY))

		c.lastComponentPosX += gl.ColumnsWidths[currColId] + gl.ColumnGap

		currColId++
		if currColId == gl.Columns {
			c.lastComponentPosX = 0
			c.lastComponentPosY += gl.RowsHeights[currRowId] + gl.RowGap
			currColId = 0
			currRowId++
		}

		if c.lastComponentPosY > height {
			height = c.lastComponentPosY
		}
	}

	c.SetDimensions(width, height)
}

func (gl *GridLayout) Arrange(c *Container, component Component) {
	gl.Rearrange(c)
}
