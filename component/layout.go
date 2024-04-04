package component

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

	for _, component := range c.components {
		width, height = hl.arrange(c, component, width, height)
	}

	c.SetDimensions(width, height)
}

func (hl *HorizontalListLayout) Arrange(c *Container, component Component) {
	width, height := hl.arrange(c, component, c.width, c.height)
	c.SetDimensions(width, height)
}

func (hl *HorizontalListLayout) arrange(c *Container, component Component, width, height int) (int, int) {
	component.SetPosition(float64(c.padding.Left+c.lastComponentPosX), float64(c.padding.Top))
	c.lastComponentPosX += component.WidthWithPadding() + hl.ColumnGap

	if component.HeightWithPadding() > height {
		height = component.HeightWithPadding()
	}

	return c.lastComponentPosX, height
}

type VerticalListLayout struct {
	RowGap int
}

func (vl *VerticalListLayout) Rearrange(c *Container) {
	width := 0
	height := 0

	c.lastComponentPosY = 0.

	for _, component := range c.components {
		width, height = vl.arrange(c, component, width, height)
	}

	c.SetDimensions(width, height)
}

func (vl *VerticalListLayout) Arrange(c *Container, component Component) {
	width, height := vl.arrange(c, component, c.width, c.height)
	c.SetDimensions(width, height)
}

func (vl *VerticalListLayout) arrange(c *Container, component Component, width, height int) (int, int) {
	component.SetPosition(float64(c.padding.Left), float64(c.lastComponentPosY+c.padding.Top))
	c.lastComponentPosY += component.HeightWithPadding() + vl.RowGap

	if component.WidthWithPadding() > width {
		width = component.WidthWithPadding()
	}

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
