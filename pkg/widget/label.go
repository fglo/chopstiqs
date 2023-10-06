package widget

import (
	"image"
	"image/color"

	"github.com/fglo/chopstiqs/pkg/colorutils"
	"github.com/fglo/chopstiqs/pkg/fontutils"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type HorizontalAlignment int

const (
	AlignmentLeft                 HorizontalAlignment = iota
	AlignmentCenteredHorizontally HorizontalAlignment = iota
	AlignmentRight                HorizontalAlignment = iota
)

type VerticalAlignment int

const (
	AlignmentTop                VerticalAlignment = iota
	AlignmentCenteredVertically VerticalAlignment = iota
	AlignmentBottom             VerticalAlignment = iota
)

type Label struct {
	widget
	text  string
	color color.RGBA
	font  font.Face

	horizontalAlignment HorizontalAlignment
	verticalAlignment   VerticalAlignment

	alignHorizontally func()
	alignVertically   func()
	textPosX          int
	textPosY          int

	textBounds image.Rectangle

	Inverted bool
}

type LabelOptions struct {
	Color color.Color
	Font  font.Face

	HorizontalAlignment HorizontalAlignment
	VerticalAlignment   VerticalAlignment

	Inverted bool
}

func NewLabel(labelText string, options *LabelOptions) *Label {
	// TODO: change deprecated function
	bounds := text.BoundString(fontutils.DefaultFontFace, labelText) // nolint

	width := bounds.Dx()
	height := bounds.Dy()

	l := &Label{
		text:       labelText,
		color:      color.RGBA{230, 230, 230, 255},
		font:       fontutils.DefaultFontFace,
		textPosX:   0,
		textPosY:   -bounds.Min.Y,
		textBounds: bounds,
		Inverted:   false,
	}

	if options != nil {
		if options.Color != nil {
			l.color = colorutils.ColorToRGBA(options.Color)
		}

		l.horizontalAlignment = options.HorizontalAlignment
		l.verticalAlignment = options.VerticalAlignment
	}

	switch l.horizontalAlignment {
	case AlignmentLeft:
		l.alignHorizontally = l.alignToLeft
	case AlignmentCenteredHorizontally:
		l.alignHorizontally = l.centerHorizontally
	case AlignmentRight:
		l.alignHorizontally = l.alignToRight
	}

	switch l.verticalAlignment {
	case AlignmentTop:
		l.alignVertically = l.alignToTop
	case AlignmentCenteredVertically:
		l.alignVertically = l.centerVertically
	case AlignmentBottom:
		l.alignVertically = l.alignToBottom
	}

	l.widget = l.createWidget(width, height)

	l.align()

	return l
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
