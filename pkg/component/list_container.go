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

type ListContainerOptions struct {
	Direction Direction
}

// ListContainer is a component that contains other components
type ListContainer struct {
	component

	direction Direction

	components      []Component
	backgroundColor imgColor.RGBA

	lastComponentPosition float64
}

// NewListContainer creates a new List container
func NewListContainer(width, height int, options ListContainerOptions) *ListContainer {
	c := &ListContainer{
		component:  *NewComponent(width, height, nil),
		components: make([]Component, 0),
		direction:  options.Direction,
	}

	return c
}

// AddComponent adds a component to the container
func (c *ListContainer) AddComponent(component Component) {
	component.setContainer(c)

	switch c.direction {
	case Horizontal:
		component.SetPosision(c.PosX()+c.lastComponentPosition, c.PosY())
		c.lastComponentPosition += float64(component.Width())
	case Vertical:
		component.SetPosision(c.PosX(), c.PosY()+c.lastComponentPosition)
		c.lastComponentPosition += float64(component.Height())
	}

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
			op.GeoM.Translate(lastComponentDim, 0)
			lastComponentDim += float64(component.Width())
		case Vertical:
			op.GeoM.Translate(0, lastComponentDim)
			lastComponentDim += float64(component.Height())
		}

		c.image.DrawImage(component.Draw(), op)
	}

	return c.image
}
