package component

import (
	imgColor "image/color"

	"github.com/fglo/chopstiqs/option"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type container interface {
	Component
	// SetBackgroundColor sets the container's background color
	SetBackgroundColor(color imgColor.RGBA)
	// GetBackgroundColor gets the container's background color
	GetBackgroundColor() imgColor.RGBA
	// FireEvents fires the container's components deferred events
	FireEvents()
}

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
func (c *Container) setContainer(container container) {
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
	for _, component := range components {
		c.AddComponent(component)
	}
}

// func (c *Container) AddFocusedHandler(f ComponentFocusedHandlerFunc) Component {
// 	c.component.AddFocusedHandler(f)

// 	for _, component := range c.components {
// 		component.AddFocusedHandler(func(args *ComponentFocusedEventArgs) {
// 			c.eventManager.Fire(c.FocusedEvent, &ComponentFocusedEventArgs{
// 				Focused:   component.Focused(),
// 				Component: component,
// 			})
// 		})
// 	}

// 	return c
// }

// SetPosX sets the container's position X.
func (c *Container) SetPosX(posX float64) {
	c.component.SetPosX(posX)
	for _, component := range c.components {
		component.RecalculateAbsPosition()
	}
}

// SetPosY sets the container's position Y.
func (c *Container) SetPosY(posY float64) {
	c.component.SetPosY(posY)
	for _, component := range c.components {
		component.RecalculateAbsPosition()
	}
}

// SetPosition sets the container's position (x and y).
func (c *Container) SetPosition(posX, posY float64) {
	c.component.SetPosition(posX, posY)
	for _, component := range c.components {
		component.RecalculateAbsPosition()
	}
}

func (c *Container) RecalculateAbsPosition() {
	c.component.RecalculateAbsPosition()
	for _, component := range c.components {
		component.RecalculateAbsPosition()
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
