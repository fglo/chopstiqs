package component

import (
	imgColor "image/color"

	"github.com/fglo/chopstiqs/pkg/option"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type Container struct {
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
func NewContainer(opt *ContainerOptions) *Container {
	c := &Container{
		components: make([]Component, 0),
	}

	c.SetDimensions(1, 1)

	if opt != nil {
		if opt.Layout != nil {
			if gl, ok := opt.Layout.(*GridLayout); ok {
				gl.Setup()
			}
			c.layout = opt.Layout
		}

		if opt.Width.IsSet() && opt.Height.IsSet() {
			c.SetDimensions(opt.Width.Val(), opt.Height.Val())
		} else {
			if opt.Width.IsSet() {
				c.SetWidth(opt.Width.Val())
			}

			if opt.Height.IsSet() {
				c.SetHeight(opt.Height.Val())
			}
		}
	}

	c.setUpComponent(opt)

	return c
}

func (c *Container) setUpComponent(opt *ContainerOptions) {
	var componentOptions ComponentOptions

	if opt != nil {
		componentOptions = ComponentOptions{
			Padding: opt.Padding,
		}
	}

	c.component.setUpComponent(&componentOptions)
}

// setContainer sets the component's container.
func (c *Container) setContainer(container *Container) {
	if c.layout != nil {
		c.layout.Rearrange(c)
	}
	c.component.setContainer(container)
}

// SetDisabled sets the container's and its component disabled states
func (c *Container) SetDisabled(disabled bool) {
	for _, component := range c.components {
		component.SetDisabled(disabled)
	}
	c.component.SetDisabled(disabled)
}

// AddComponent adds a component to the container
func (c *Container) AddComponent(component Component) {
	c.components = append(c.components, component)
	if c.layout != nil {
		c.layout.Arrange(c, component)
	}
	component.setContainer(c)
	component.AddFocusedHandler(func(args *ComponentFocusedEventArgs) {
		c.eventManager.Fire(c.FocusedEvent, &ComponentFocusedEventArgs{
			Focused:   component.Focused(),
			Component: component,
		})
	})
}

// AddComponents adds components to the container
func (c *Container) AddComponents(components ...Component) {
	c.components = append(c.components, components...)

	for _, component := range components {
		if c.layout != nil {
			c.layout.Arrange(c, component)
		}
		component.setContainer(c)
		component.AddFocusedHandler(func(args *ComponentFocusedEventArgs) {
			c.eventManager.Fire(c.FocusedEvent, &ComponentFocusedEventArgs{
				Focused:   component.Focused(),
				Component: component,
			})
		})
	}
}

// SetBackgroundColor sets the container's background color
func (c *Container) SetBackgroundColor(color imgColor.RGBA) {
	c.backgroundColor = color
}

// GetBackgroundColor gets the container's background color
func (c *Container) GetBackgroundColor() imgColor.RGBA {
	return c.backgroundColor
}

// FireEvents fires the container's components deferred events
func (c *Container) FireEvents() {
	for _, component := range c.components {
		component.FireEvents()
	}
}

// Draw draws the container's components, executes deferred events and returns the image.
func (c *Container) Draw() *ebiten.Image {
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
