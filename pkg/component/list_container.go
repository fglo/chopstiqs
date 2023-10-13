package component

import (
	imgColor "image/color"

	"github.com/fglo/chopstiqs/pkg/event"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type Direction int

const (
	Horizontal Direction = iota
	Vertical
)

// ListContainer is a component that contains other components
type ListContainer struct {
	component

	direction Direction

	components      []Component
	backgroundColor imgColor.RGBA

	lastComponentPosition float64
}

type ListContainerOptions struct {
	Direction Direction

	LeftPadding   *int
	RightPadding  *int
	TopPadding    *int
	BottomPadding *int
}

// NewListContainer creates a new List container
func NewListContainer(options *ListContainerOptions) *ListContainer {
	c := &ListContainer{
		components: make([]Component, 0),
	}

	c.component = c.createComponent(0, 0, options)

	if options != nil {
		c.direction = options.Direction
	}

	return c
}

func (c *ListContainer) createComponent(width, height int, options *ListContainerOptions) component {
	var componentOptions ComponentOptions

	if options != nil {
		componentOptions = ComponentOptions{
			LeftPadding:   options.LeftPadding,
			RightPadding:  options.RightPadding,
			TopPadding:    options.TopPadding,
			BottomPadding: options.BottomPadding,
		}
	}

	component := NewComponent(width, height, &componentOptions)

	return *component
}

// setContainer sets the component's container.
func (c *ListContainer) setContainer(container Container) {
	c.component.setContainer(container)

	c.lastComponentPosition = 0

	for _, component := range c.components {
		switch c.direction {
		case Horizontal:
			component.SetPosision(c.PosX()+c.lastComponentPosition+float64(c.leftPadding), c.PosY()+float64(c.topPadding))
			c.lastComponentPosition += float64(component.WidthWithPadding())
		case Vertical:
			component.SetPosision(c.PosX()+float64(c.leftPadding), c.PosY()+c.lastComponentPosition+float64(c.topPadding))
			c.lastComponentPosition += float64(component.HeightWithPadding())
		}
	}
}

// AddComponent adds a component to the container
func (c *ListContainer) AddComponent(component Component) {
	width := c.width
	height := c.height

	switch c.direction {
	case Horizontal:
		component.SetPosision(c.PosX()+c.lastComponentPosition+float64(c.leftPadding), c.PosY()+float64(c.topPadding))
		c.lastComponentPosition += float64(component.WidthWithPadding())

		width += component.WidthWithPadding()

		if component.HeightWithPadding() > height {
			height = component.HeightWithPadding()
		}
	case Vertical:
		component.SetPosision(c.PosX()+float64(c.leftPadding), c.PosY()+c.lastComponentPosition+float64(c.topPadding))
		c.lastComponentPosition += float64(component.HeightWithPadding())

		if component.WidthWithPadding() > width {
			width = component.WidthWithPadding()
		}

		height += component.HeightWithPadding()
	}

	component.setContainer(c)

	c.SetDimensions(width, height)

	c.components = append(c.components, component)
}

func (c *ListContainer) SetBackgroundColor(color imgColor.RGBA) {
	c.backgroundColor = color
}

func (c *ListContainer) GetBackgroundColor() imgColor.RGBA {
	return c.backgroundColor
}

// FireEvents fires the container's components deferred events
func (c *ListContainer) FireEvents() {
	for _, component := range c.components {
		component.FireEvents()
	}
}

// Draw draws the container's components, executes deferred events and returns the image.
func (c *ListContainer) Draw() *ebiten.Image {
	event.HandleFired()

	c.image.Fill(c.backgroundColor)

	lastComponentDim := 0.

	for _, component := range c.components {
		op := &ebiten.DrawImageOptions{}

		switch c.direction {
		case Horizontal:
			op.GeoM.Translate(lastComponentDim+float64(c.leftPadding), float64(c.topPadding))
			lastComponentDim += float64(component.WidthWithPadding())
		case Vertical:
			op.GeoM.Translate(float64(c.leftPadding), lastComponentDim+float64(c.topPadding))
			lastComponentDim += float64(component.HeightWithPadding())
		}

		c.image.DrawImage(component.Draw(), op)
	}

	c.component.Draw()

	return c.image
}
