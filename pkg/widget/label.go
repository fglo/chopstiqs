package widget

import (
	"image"
	"image/color"

	"github.com/fglo/chopstiqs/pkg/fontutils"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type position int

const (
	left                 position = iota
	right                position = iota
	centered             position = iota
	centeredHorizontally position = iota
	centeredVertically   position = iota
	top                  position = iota
	bottom               position = iota
)

type Label struct {
	widget
	text  string
	color color.RGBA
	font  font.Face

	position          position
	alignHorizontally func()
	alignVertically   func()
	textPosX          int
	textPosY          int

	textBounds image.Rectangle

	Inverted bool
}

type LabelOpt func(b *Label)
type LabelOptions struct {
	opts []LabelOpt
}

func NewLabel(labelText string, options *LabelOptions) *Label {
	// TODO: change deprecated function
	bounds := text.BoundString(fontutils.DefaultFontFace, labelText) // nolint

	width := bounds.Dx()
	height := bounds.Dy()

	lbl := &Label{
		text:       labelText,
		color:      color.RGBA{230, 230, 230, 255},
		font:       fontutils.DefaultFontFace,
		position:   left,
		textPosX:   0,
		textPosY:   -bounds.Min.Y,
		textBounds: bounds,
		Inverted:   false,
	}

	lbl.widget = lbl.createWidget(width, height)

	if options != nil {
		for _, o := range options.opts {
			o(lbl)
		}
	}

	return lbl
}

func (o *LabelOptions) Color(color color.RGBA) *LabelOptions {
	o.opts = append(o.opts, func(l *Label) {
		l.color = color
	})

	return o
}

func (o *LabelOptions) Centered() *LabelOptions {
	o.opts = append(o.opts, func(l *Label) {
		l.position = centered
		l.alignHorizontally = l.centerHorizontally
		l.alignVertically = l.centerVertically
	})

	return o
}

func (o *LabelOptions) CenteredHorizontally() *LabelOptions {
	o.opts = append(o.opts, func(l *Label) {
		l.position = centeredHorizontally
		l.alignHorizontally = l.centerHorizontally
	})

	return o
}

func (o *LabelOptions) CenteredVertically() *LabelOptions {
	o.opts = append(o.opts, func(l *Label) {
		l.position = centeredVertically
		l.alignVertically = l.centerVertically
	})

	return o
}

func (o *LabelOptions) Left() *LabelOptions {
	o.opts = append(o.opts, func(l *Label) {
		l.position = left
		l.alignHorizontally = l.alignToLeft
	})

	return o
}

func (o *LabelOptions) Right() *LabelOptions {
	o.opts = append(o.opts, func(l *Label) {
		l.position = right
		l.alignHorizontally = l.alignToRight
	})

	return o
}

func (o *LabelOptions) Top() *LabelOptions {
	o.opts = append(o.opts, func(l *Label) {
		l.position = top
		l.alignVertically = l.alignToTop
	})

	return o
}

func (o *LabelOptions) Bottom() *LabelOptions {
	o.opts = append(o.opts, func(l *Label) {
		l.position = bottom
		l.alignVertically = l.alignToBottom
	})

	return o
}

func (l *Label) align() {
	if l.alignHorizontally != nil {
		l.alignHorizontally()
	}

	if l.alignVertically != nil {
		l.alignVertically()
	}
}

func (l *Label) centerHorizontally() {
	l.posX = l.posX - float64(l.textBounds.Dx())/2
}

func (l *Label) centerVertically() {
	l.posY = l.posY - float64(l.textBounds.Dy())/2
}

func (l *Label) alignToLeft() {
	// l.posX = l.posX
}

func (l *Label) alignToRight() {
	l.posX = l.posX - float64(l.textBounds.Dx())
}

func (l *Label) alignToTop() {
	// l.posY = l.posY
}

func (l *Label) alignToBottom() {
	l.posY = l.posY - float64(l.textBounds.Dy())
}

func (l *Label) Invert() {
	l.Inverted = !l.Inverted
}

func (l *Label) Draw() *ebiten.Image {
	if l.Inverted {
		text.Draw(l.image, l.text, l.font, 0, l.textPosY, color.RGBA{255 - l.color.R, 255 - l.color.G, 255 - l.color.B, l.color.A})
	} else {
		text.Draw(l.image, l.text, l.font, 0, l.textPosY, l.color)
	}

	return l.image
}

func (l *Label) createWidget(width, height int) widget {
	widgetOptions := &WidgetOptions{}

	return *NewWidget(width, height, widgetOptions)
}

func (l *Label) SetPosX(posX float64) {
	l.widget.SetPosX(posX)
	l.align()
}

func (l *Label) SetPosY(posY float64) {
	l.widget.SetPosY(posY)
	l.align()
}

func (l *Label) SetPosistion(posX, posY float64) {
	l.widget.SetPosistion(posX, posY)
	l.align()
}
