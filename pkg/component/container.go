package component

import (
	imgColor "image/color"

	"github.com/fglo/chopstiqs/pkg/event"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type Container interface {
	Component
	AddComponent(Component)
	SetBackgroundColor(imgColor.RGBA)
	GetBackgroundColor() imgColor.RGBA
}

type container struct {
	component

	layout Layout

	components      []Component
	backgroundColor imgColor.RGBA

	lastComponentPosX int
	lastComponentPosY int
}

type ContainerOptions struct {
	Layout Layout

	Width  *int
	Height *int

	Padding *Padding
}

// Newcontainer creates a new simple container
func NewContainer(options *ContainerOptions) Container {
	c := &container{
		components: make([]Component, 0),
	}

	c.component = c.createComponent(1, 1, options)

	if options != nil {
		if options.Layout != nil {
			if gl, ok := options.Layout.(*GridLayout); ok {
				gl.Setup()
			}
			c.layout = options.Layout
		}

		if options.Width != nil {
			c.SetWidth(*options.Width)
		}

		if options.Height != nil {
			c.SetHeight(*options.Height)
		}
	}

	return c
}

func (c *container) createComponent(width, height int, options *ContainerOptions) component {
	var componentOptions ComponentOptions

	if options != nil {
		componentOptions = ComponentOptions{
			Padding: options.Padding,
		}
	}

	component := NewComponent(width, height, &componentOptions)

	return *component
}

// setContainer sets the component's container.
func (c *container) setContainer(container Container) {
	if c.layout != nil {
		c.layout.Rearrange(c)
	}
	c.component.setContainer(container)
}

// AddComponent adds a component to the container
func (c *container) AddComponent(component Component) {
	c.components = append(c.components, component)
	if c.layout != nil {
		c.layout.Arrange(c, component)
	}
	component.setContainer(c)
}

func (c *container) SetBackgroundColor(color imgColor.RGBA) {
	c.backgroundColor = color
}

func (c *container) GetBackgroundColor() imgColor.RGBA {
	return c.backgroundColor
}

// FireEvents fires the container's components deferred events
func (c *container) FireEvents() {
	for _, component := range c.components {
		component.FireEvents()
	}
}

// Draw draws the container's components, executes deferred events and returns the image.
func (c *container) Draw() *ebiten.Image {
	event.HandleFired()

	c.image.Fill(c.backgroundColor)

	for _, component := range c.components {
		if !component.Hidden() {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(component.Position())
			c.image.DrawImage(component.Draw(), op)
		}
	}

	c.component.Draw()

	return c.image
}
