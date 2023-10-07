package component

import (
	"image"

	imgColor "image/color"

	"github.com/fglo/chopstiqs/pkg/event"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

// Container is a component that contains other components
type Container struct {
	image *ebiten.Image

	Rect image.Rectangle

	backgroundColor imgColor.RGBA

	components []Component

	disabled bool

	width  int
	height int

	posX float64
	posY float64
}

// NewContainer creates a new container
func NewContainer(posX, posY float64, width, height int, backgroundColor imgColor.RGBA) *Container {
	c := &Container{
		image:           ebiten.NewImage(width, height),
		width:           width,
		height:          height,
		posX:            posX,
		posY:            posY,
		Rect:            image.Rectangle{Min: image.Point{int(posX), int(posY)}, Max: image.Point{int(posX) + width, int(posY) + height}},
		backgroundColor: backgroundColor,
	}

	return c
}

// Disable disables the container
func (c *Container) Disable() {
	c.disabled = true
}

// Enable enables the container
func (c *Container) Enable() {
	c.disabled = false
}

// SetDisabled sets the container to disabled or enabled
func (c *Container) SetDisabled(disabled bool) {
	c.disabled = disabled
}

// SetBackgroundColor sets the container's background color
func (c *Container) SetBackgroundColor(backgroundColor imgColor.RGBA) {
	c.backgroundColor = backgroundColor
}

// Position returns the container's position
func (c *Container) Position() (float64, float64) {
	return c.posX, c.posY
}

// Size returns the container's size
func (c *Container) Size() (int, int) {
	return c.width, c.height
}

// AddComponent adds a component to the container
func (c *Container) AddComponent(posX, posY float64, component Component) {
	component.setContainer(c)
	component.SetPosistion(posX, posY)
	c.components = append(c.components, component)
}

// Update updates registered in the container components.
func (c *Container) Update() {
	for _, component := range c.components {
		component.FireEvents()
	}
}

// Draw draws the container's components, executes deferred events and returns the image.
func (c *Container) Draw() *ebiten.Image {
	event.HandleFired()

	c.image.Fill(c.backgroundColor)

	for _, component := range c.components {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(component.Position())
		c.image.DrawImage(component.Draw(), op)
	}

	return c.image
}
