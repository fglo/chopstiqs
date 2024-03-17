package component

import (
	"encoding/xml"
	"fmt"
	imgColor "image/color"
	"strconv"

	"github.com/fglo/chopstiqs/pkg/option"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

// type Container interface {
// 	Component
// 	// SetDisabled sets the container's and its component disabled states
// 	SetDisabled(disabled bool)
// 	// AddComponent adds a component to the container
// 	AddComponent(Component)
// 	// FireEvents fires the container's components deferred events
// 	FireEvents()
// 	// Draw draws the container's components, executes deferred events and returns the image.
// 	Draw() *ebiten.Image
// 	// SetBackgroundColor sets the container's background color
// 	SetBackgroundColor(imgColor.RGBA)
// 	// GetBackgroundColor gets the container's background color
// 	GetBackgroundColor() imgColor.RGBA
// }

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

type ContainerXML struct {
	XMLName    xml.Name      `xml:"container"`
	Layout     Layout        `xml:"layout,attr"`
	Width      option.OptInt `xml:"width,attr"`
	Height     option.OptInt `xml:"height,attr"`
	Padding    *Padding      `xml:"padding,attr"`
	Components []Component
}

func (c *Container) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "container"

	return e.EncodeElement(ContainerXML{
		Layout:     c.layout,
		Width:      option.Int(c.width),
		Height:     option.Int(c.height),
		Padding:    &c.padding,
		Components: c.components,
	}, start)
}

func (c *Container) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	opt := ContainerOptions{
		Padding: &Padding{},
	}

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "layout":
			l, err := UnmarshalLayout(attr)
			if err != nil {
				return fmt.Errorf("error unmarshaling layout attribute: %w", err)
			}

			opt.Layout = l
		case "width":
			if len(attr.Value) > 0 {
				width, err := strconv.Atoi(attr.Value)
				if err != nil {
					return fmt.Errorf("error unmarshaling width attribute: %w", err)
				}
				opt.Width = option.Int(width)
			}
		case "height":
			if len(attr.Value) > 0 {
				height, err := strconv.Atoi(attr.Value)
				if err != nil {
					return fmt.Errorf("error unmarshaling height attribute: %w", err)
				}
				opt.Height = option.Int(height)
			}
		case "padding":
			err := opt.Padding.UnmarshalXMLAttr(attr)
			if err != nil {
				return fmt.Errorf("error unmarshaling padding attribute: %w", err)
			}
		}
	}

	*c = *NewContainer(&opt)

	if err := d.Skip(); err != nil {
		return fmt.Errorf("failed to skip element: %w", err)
	}

	return nil
}
