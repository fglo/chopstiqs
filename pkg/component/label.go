package component

import (
	"image"
	"image/color"

	colorutils "github.com/fglo/chopstiqs/internal/color"
	fontutils "github.com/fglo/chopstiqs/internal/font"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type HorizontalAlignment int

const (
	AlignmentLeft HorizontalAlignment = iota
	AlignmentCenteredHorizontally
	AlignmentRight
)

type VerticalAlignment int

const (
	AlignmentTop VerticalAlignment = iota
	AlignmentCenteredVertically
	AlignmentBottom
)

type Label struct {
	component
	text  string
	color color.RGBA

	font    font.Face
	metrics fontutils.Metrics

	horizontalAlignment HorizontalAlignment
	verticalAlignment   VerticalAlignment

	alignHorizontally func()
	alignVertically   func()
	textOriginX       int
	textOriginY       int

	bounds image.Rectangle

	Inverted bool
}

type LabelOptions struct {
	Color color.Color
	Font  font.Face

	HorizontalAlignment HorizontalAlignment
	VerticalAlignment   VerticalAlignment

	Inverted bool

	Padding *Padding
}

func NewLabel(text string, opt *LabelOptions) *Label {
	l := &Label{
		color:       color.RGBA{230, 230, 230, 255},
		font:        fontutils.DefaultFontFace,
		metrics:     fontutils.NewMetrics(fontutils.DefaultFontFace.Metrics()),
		textOriginX: 0,
		Inverted:    false,
	}

	l.SetText(text)

	if opt != nil {
		if opt.Color != nil {
			l.color = colorutils.ToRGBA(opt.Color)
		}

		if opt.Font != nil {
			l.font = opt.Font
			l.metrics = fontutils.NewMetrics(l.font.Metrics())
		}

		l.horizontalAlignment = opt.HorizontalAlignment
		l.verticalAlignment = opt.VerticalAlignment
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

	l.setUpComponent(opt)

	l.align()

	return l
}

func (l *Label) setUpComponent(opt *LabelOptions) {
	var componentOptions ComponentOptions

	if opt != nil {
		componentOptions = ComponentOptions{
			Padding: opt.Padding,
		}
	}

	l.component.setUpComponent(&componentOptions)
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
	l.posX = l.posX - float64(l.bounds.Dx())/2
}

func (l *Label) centerVertically() {
	l.posY = l.posY - float64(l.bounds.Dy())/2
}

func (l *Label) alignToLeft() {
	// l.posX = l.posX
}

func (l *Label) alignToRight() {
	l.posX = l.posX - float64(l.bounds.Dx())
}

func (l *Label) alignToTop() {
	// l.posY = l.posY
}

func (l *Label) alignToBottom() {
	l.posY = l.posY - float64(l.bounds.Dy())
}

func (l *Label) SetText(labelText string) {
	// TODO: change deprecated function
	l.bounds = text.BoundString(fontutils.DefaultFontFace, labelText) // nolint
	l.text = labelText
	l.textOriginY = -l.bounds.Min.Y

	l.SetDimensions(l.bounds.Dx(), l.bounds.Dy())

	if l.container != nil {
		l.container.SetWidth(l.container.Width() + l.component.width)
	}
}

func (l *Label) InvertColor() {
	l.Inverted = !l.Inverted
}

func (l *Label) Draw() *ebiten.Image {
	if l.hidden {
		return l.image
	}

	l.image = ebiten.NewImage(l.widthWithPadding, l.heightWithPadding)

	if l.Inverted {
		text.Draw(l.image, l.text, l.font, l.textOriginX+l.padding.Left, l.textOriginY+l.padding.Top, colorutils.Invert(l.color))
	} else {
		text.Draw(l.image, l.text, l.font, l.textOriginX+l.padding.Left, l.textOriginY+l.padding.Top, l.color)
	}

	l.component.Draw()

	return l.image
}

func (l *Label) SetPosX(posX float64) {
	l.component.SetPosX(posX)
	l.align()
}

func (l *Label) SetPosY(posY float64) {
	l.component.SetPosY(posY)
	l.align()
}

func (l *Label) SetPosistion(posX, posY float64) {
	l.component.SetPosision(posX, posY)
	l.align()
}
