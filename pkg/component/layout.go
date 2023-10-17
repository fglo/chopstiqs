package component

type Layout interface {
	Rearrange(*container)
	Arrange(*container, Component)
}

var HorizontalListLayout horizontalListLayout

type horizontalListLayout struct {
}

func (hl horizontalListLayout) Rearrange(c *container) {
	width := 0
	height := 0

	c.lastComponentDim = 0

	for _, component := range c.components {
		component.SetPosision(float64(c.leftPadding)+c.lastComponentDim, float64(c.topPadding))
		c.lastComponentDim += float64(component.WidthWithPadding())

		width += component.WidthWithPadding()

		if component.HeightWithPadding() > height {
			height = component.HeightWithPadding()
		}
	}

	c.SetDimensions(width, height)
}

func (hl horizontalListLayout) Arrange(c *container, component Component) {
	width := c.width
	height := c.height

	component.SetPosision(float64(c.leftPadding)+c.lastComponentDim, float64(c.topPadding))
	c.lastComponentDim += float64(component.WidthWithPadding())

	width += component.WidthWithPadding()

	if component.HeightWithPadding() > height {
		height = component.HeightWithPadding()
	}

	c.SetDimensions(width, height)
}

var VerticalListLayout verticalListLayout

type verticalListLayout struct {
}

func (vl verticalListLayout) Rearrange(c *container) {
	width := 0
	height := 0

	c.lastComponentDim = 0.

	for _, component := range c.components {
		component.SetPosision(float64(c.leftPadding), c.lastComponentDim+float64(c.topPadding))
		c.lastComponentDim += float64(component.HeightWithPadding())

		if component.WidthWithPadding() > width {
			width = component.WidthWithPadding()
		}

		height += component.HeightWithPadding()
	}

	c.SetDimensions(width, height)
}

func (vl verticalListLayout) Arrange(c *container, component Component) {
	width := c.width
	height := c.height

	component.SetPosision(float64(c.leftPadding), c.lastComponentDim+float64(c.topPadding))
	c.lastComponentDim += float64(component.HeightWithPadding())

	if component.WidthWithPadding() > width {
		width = component.WidthWithPadding()
	}

	height += component.HeightWithPadding()

	c.SetDimensions(width, height)
}
