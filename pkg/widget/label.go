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
	left     position = 0
	centered position = 1
	right    position = 2
)

type Label struct {
	widget
	text  string
	color color.RGBA
	font  font.Face

	position position
	textPosX int
	textPosY int

	textBounds image.Rectangle

	Inverted bool
}

type LabelOpt func(b *Label)
type LabelOptions struct {
	opts []LabelOpt
}

func NewLabel(posX, posY float64, labelText string, color color.RGBA, options *LabelOptions) *Label {
	// TODO: change deprecated function
	bounds := text.BoundString(fontutils.DefaultFontFace, labelText) // nolint

	width := bounds.Dx()
	height := bounds.Dy()

	lbl := &Label{
		text:       labelText,
		color:      color,
		font:       fontutils.DefaultFontFace,
		position:   left,
		textPosX:   0,
		textPosY:   -bounds.Min.Y,
		textBounds: bounds,
		Inverted:   false,
	}

	lbl.widget = lbl.createWidget(posX, posY, width, height)

	if options != nil {
		for _, o := range options.opts {
			o(lbl)
		}
	}

	return lbl
}

func (o *LabelOptions) Centered() *LabelOptions {
	o.opts = append(o.opts, func(l *Label) {
		l.position = centered
		l.posX = l.posX - float64(l.textBounds.Dx())/2
		l.posY = l.posY - float64(l.textBounds.Dy())/2
	})

	return o
}

func (o *LabelOptions) CenteredHorizontally() *LabelOptions {
	o.opts = append(o.opts, func(l *Label) {
		l.position = centered
		l.posX = l.posX - float64(l.textBounds.Dx())/2
	})

	return o
}

func (o *LabelOptions) CenteredVertically() *LabelOptions {
	o.opts = append(o.opts, func(l *Label) {
		l.position = centered
		l.posY = l.posY - float64(l.textBounds.Dy())/2
	})

	return o
}

func (o *LabelOptions) Left() *LabelOptions {
	o.opts = append(o.opts, func(l *Label) {
		l.position = left
	})

	return o
}

func (o *LabelOptions) Right() *LabelOptions {
	o.opts = append(o.opts, func(l *Label) {
		l.position = right
		l.posX = l.posX - float64(l.textBounds.Dx())
	})

	return o
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

func (l *Label) createWidget(posX, posY float64, width, height int) widget {
	widgetOptions := &WidgetOptions{}

	return *NewWidget(posX, posY, width, height, widgetOptions)
}
