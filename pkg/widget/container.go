package widget

import (
	"image"

	imgColor "image/color"

	"github.com/fglo/chopstiqs/pkg/event"
	"github.com/fglo/chopstiqs/pkg/input"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type Container struct {
	image *ebiten.Image

	Rect image.Rectangle

	backgroundColor imgColor.RGBA

	components []Widget

	disabled bool

	width  int
	height int

	posX float64
	posY float64
}

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

func (c *Container) Disable() {
	c.disabled = true
}

func (c *Container) Enable() {
	c.disabled = false
}

func (c *Container) SetDisabled(disabled bool) {
	c.disabled = disabled
}

func (c *Container) SetBackgroundColor(backgroundColor imgColor.RGBA) {
	c.backgroundColor = backgroundColor
}

func (c *Container) Position() (float64, float64) {
	return c.posX, c.posY
}

func (c *Container) Size() (int, int) {
	return c.width, c.height
}

func (c *Container) AddComponent(component Widget) {
	component.setContainer(c)
	c.components = append(c.components, component)
}

func (c *Container) Update(mouse *input.Mouse) {
	for _, component := range c.components {
		component.FireEvents(mouse)
	}
}

func (c *Container) Draw(mouse *input.Mouse) *ebiten.Image {
	event.ExecuteDeferred()

	c.image.Fill(c.backgroundColor)

	for _, component := range c.components {
		component.FireEvents(mouse)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(component.Position())
		c.image.DrawImage(component.Draw(), op)
	}

	return c.image
}
