package component

import (
	imgColor "image/color"

	"github.com/fglo/chopstiqs/pkg/event"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

// SimpleContainer is a component that contains other components
type SimpleContainer struct {
	component

	components      []Component
	backgroundColor imgColor.RGBA
}

type SimpleContainerOptions struct {
	LeftPadding   *int
	RightPadding  *int
	TopPadding    *int
	BottomPadding *int
}

// NewSimpleContainer creates a new simple container
func NewSimpleContainer(width, height int, options *SimpleContainerOptions) *SimpleContainer {
	c := &SimpleContainer{
		components: make([]Component, 0),
	}

	c.component = c.createComponent(width, height, options)

	return c
}

func (c *SimpleContainer) createComponent(width, height int, options *SimpleContainerOptions) component {
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

// AddComponent adds a component to the container
func (c *SimpleContainer) AddComponent(component Component) {
	component.setContainer(c)
	c.components = append(c.components, component)
}

func (c *SimpleContainer) SetBackgroundColor(color imgColor.RGBA) {
	c.backgroundColor = color
}

func (c *SimpleContainer) GetBackgroundColor() imgColor.RGBA {
	return c.backgroundColor
}

// FireEvents fires the container's components deferred events
func (c *SimpleContainer) FireEvents() {
	for _, component := range c.components {
		component.FireEvents()
	}
}

// Draw draws the container's components, executes deferred events and returns the image.
func (c *SimpleContainer) Draw() *ebiten.Image {
	event.HandleFired()

	c.image.Fill(c.backgroundColor)

	for _, component := range c.components {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(component.Position())
		c.image.DrawImage(component.Draw(), op)
	}

	c.component.Draw()

	return c.image
}
