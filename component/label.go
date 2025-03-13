package component

import (
	"image"
	"image/color"

	colorutils "github.com/fglo/chopstiqs/color"
	fontutils "github.com/fglo/chopstiqs/font"
	"github.com/fglo/chopstiqs/option"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type Label struct {
	component
	text  string
	color color.RGBA

	font    font.Face
	metrics fontutils.Metrics

	horizontalAlignment option.HorizontalAlignment
	verticalAlignment   option.VerticalAlignment

	textOriginX int
	textOriginY int

	bounds image.Rectangle

	Inverted bool
}

type LabelOptions struct {
	Color color.Color
	Font  font.Face

	HorizontalAlignment option.HorizontalAlignment
	VerticalAlignment   option.VerticalAlignment

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

func (l *Label) centerHorizontally() {
	l.SetPosX(float64(l.container.WidthWithPadding()-l.bounds.Dx()) / 2)
}

func (l *Label) alignToLeft() {
	l.SetPosX(0)
}

func (l *Label) alignToRight() {
	l.SetPosX(float64(l.container.WidthWithPadding() - l.bounds.Dx()))
}

func (l *Label) centerVertically() {
	l.SetPosY(float64(l.container.HeightWithPadding()-l.bounds.Dy()) / 2)
}

func (l *Label) alignToTop() {
	l.SetPosY(0)
}

func (l *Label) alignToBottom() {
	l.SetPosY(float64(l.container.HeightWithPadding() - l.bounds.Dy()))
}

func (l *Label) align() {
	switch l.horizontalAlignment {
	case option.AlignmentLeft:
		l.alignToLeft()
	case option.AlignmentCenteredHorizontally:
		l.centerHorizontally()
	case option.AlignmentRight:
		l.alignToRight()
	}

	switch l.verticalAlignment {
	case option.AlignmentTop:
		l.alignToTop()
	case option.AlignmentCenteredVertically:
		l.centerVertically()
	case option.AlignmentBottom:
		l.alignToBottom()
	}
}

func (l *Label) SetText(labelText string) {
	// TODO: change deprecated function
	l.bounds = text.BoundString(fontutils.DefaultFontFace, labelText) // nolint
	l.text = labelText
	l.textOriginY = -l.bounds.Min.Y

	l.SetDimensions(l.bounds.Dx(), l.bounds.Dy())

	if l.container != nil && l.container.Width() < l.widthWithPadding {
		l.container.SetWidth(l.widthWithPadding)
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
