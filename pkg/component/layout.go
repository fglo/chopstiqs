package component

type Layout interface {
	Rearrange(*container)
	Arrange(*container, Component)
}

type HorizontalListLayout struct {
	ColumnGap int
}

func (hl *HorizontalListLayout) Rearrange(c *container) {
	width := 0
	height := 0

	c.lastComponentPosX = 0

	for _, component := range c.components {
		width, height = hl.arrange(c, component, width, height)
	}

	c.SetDimensions(width, height)
}

func (hl *HorizontalListLayout) Arrange(c *container, component Component) {
	width, height := hl.arrange(c, component, c.width, c.height)
	c.SetDimensions(width, height)
}

func (hl *HorizontalListLayout) arrange(c *container, component Component, width, height int) (int, int) {
	component.SetPosision(float64(c.padding.Left+c.lastComponentPosX), float64(c.padding.Top))
	c.lastComponentPosX += component.WidthWithPadding() + hl.ColumnGap

	if component.HeightWithPadding() > height {
		height = component.HeightWithPadding()
	}

	return c.lastComponentPosX, height
}

type VerticalListLayout struct {
	RowGap int
}

func (vl *VerticalListLayout) Rearrange(c *container) {
	width := 0
	height := 0

	c.lastComponentPosY = 0.

	for _, component := range c.components {
		width, height = vl.arrange(c, component, width, height)
	}

	c.SetDimensions(width, height)
}

func (vl *VerticalListLayout) Arrange(c *container, component Component) {
	width, height := vl.arrange(c, component, c.width, c.height)
	c.SetDimensions(width, height)
}

func (vl *VerticalListLayout) arrange(c *container, component Component, width, height int) (int, int) {
	component.SetPosision(float64(c.padding.Left), float64(c.lastComponentPosY+c.padding.Top))
	c.lastComponentPosY += component.HeightWithPadding() + vl.RowGap

	if component.WidthWithPadding() > width {
		width = component.WidthWithPadding()
	}

	return width, c.lastComponentPosY
}

type GridLayout struct {
	Columns       int
	ColumnsWidths []int
	Rows          int
	RowsHeights   []int

	Width int

	currColId int
	currRowId int
}

func (gl *GridLayout) Rearrange(c *container) {
	width := 0
	height := 0

	c.lastComponentPosX = 0.
	c.lastComponentPosY = 0.

	gl.currColId = 0
	gl.currRowId = 0

	for _, component := range c.components {
		width, height = gl.arrange(c, component, c.width, c.height)
	}

	c.SetDimensions(width, height)
}

func (gl *GridLayout) Arrange(c *container, component Component) {
	width, height := gl.arrange(c, component, c.width, c.height)
	c.SetDimensions(width, height)
}

func (gl *GridLayout) arrange(c *container, component Component, width, height int) (int, int) {
	// if gl.currRowId >= gl.Rows {
	// 	// component.SetVisible(false)
	// 	return width, height
	// }

	// component.SetPosision(float64(c.padding.Left+gl.columnWidth), float64(c.padding.Top+c.lastComponentPosY))
	// if gl.currColId < gl.Columns {
	// 	if component.HeightWithPadding() > gl.maxHeightInRow {
	// 		gl.maxHeightInRow = component.HeightWithPadding()
	// 	}
	// }

	// gl.currColId++
	// if gl.currColId == gl.Columns {
	// 	gl.currColId = 0
	// 	gl.currRowId++
	// 	c.lastComponentPosX = 0 // reset last component dim for next row.
	// 	c.lastComponentPosY += gl.maxHeightInRow
	// 	gl.maxHeightInRow = 0 // reset max height for next row.
	// }

	// if c.lastComponentPosY > height {
	// 	height = c.lastComponentPosY
	// }

	return width, height
}
