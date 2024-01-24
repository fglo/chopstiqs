package component

import (
	imgColor "image/color"

	"github.com/fglo/chopstiqs/pkg/event"
	"github.com/fglo/chopstiqs/pkg/option"
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

	Width  option.OptInt
	Height option.OptInt

	Padding *Padding
}

// Newcontainer creates a new simple container
func NewContainer(options *ContainerOptions) Container {
	c := &container{
		components: make([]Component, 0),
	}

	c.component.width = 1
	c.component.height = 1

	if options != nil {
		if options.Layout != nil {
			if gl, ok := options.Layout.(*GridLayout); ok {
				gl.Setup()
			}
			c.layout = options.Layout
		}

		if options.Width.HasVal() {
			c.SetWidth(options.Width.Val())
		}

		if options.Height.HasVal() {
			c.SetHeight(options.Height.Val())
		}
	}

	c.setUpComponent(options)

	return c
}

func (c *container) setUpComponent(options *ContainerOptions) {
	var componentOptions ComponentOptions

	if options != nil {
		componentOptions = ComponentOptions{
			Padding: options.Padding,
		}
	}

	c.component.setUpComponent(&componentOptions)
}

// setContainer sets the component's container.
func (c *container) setContainer(container Container) {
	if c.layout != nil {
		c.layout.Rearrange(c)
	}
	c.component.setContainer(container)
}

// SetDisabled sets the container's and its component disabled states
func (c *container) SetDisabled(disabled bool) {
	for _, component := range c.components {
		component.SetDisabled(disabled)
	}
	c.component.SetDisabled(disabled)
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
