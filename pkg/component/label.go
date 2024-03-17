package component

import (
	"encoding/xml"
	"image"
	"image/color"

	colorutils "github.com/fglo/chopstiqs/pkg/color"
	font1 "github.com/fglo/chopstiqs/pkg/font"
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

	Padding *Padding
}

func NewLabel(text string, opt *LabelOptions) *Label {
	l := &Label{
		color:    color.RGBA{230, 230, 230, 255},
		font:     font1.DefaultFontFace,
		textPosX: 0,
		Inverted: false,
	}

	l.SetText(text)

	if opt != nil {
		if opt.Color != nil {
			l.color = colorutils.ColorToRGBA(opt.Color)
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

func (l *Label) SetText(labelText string) {
	// TODO: change deprecated function
	bounds := text.BoundString(font1.DefaultFontFace, labelText) // nolint

	l.text = labelText

	l.textPosY = -bounds.Min.Y
	l.textBounds = bounds

	l.SetDimensions(bounds.Dx(), bounds.Dy())

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
		text.Draw(l.image, l.text, l.font, l.textPosX+l.padding.Left, l.textPosY+l.padding.Top, color.RGBA{255 - l.color.R, 255 - l.color.G, 255 - l.color.B, l.color.A})
	} else {
		text.Draw(l.image, l.text, l.font, l.textPosX+l.padding.Left, l.textPosY+l.padding.Top, l.color)
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

type LabelXML struct {
	XMLName             xml.Name                    `xml:"label"`
	Color               colorutils.RGBASerializable `xml:"color,attr,omitempty"`
	Font                font.Face
	HorizontalAlignment HorizontalAlignment `xml:"horizontalAlignment,attr,omitempty"`
	VerticalAlignment   VerticalAlignment   `xml:"verticalAlignment,attr,omitempty"`
	Inverted            bool                `xml:"inverted,attr,omitempty"`
	Padding             *Padding            `xml:"padding,attr,omitempty"`
}

func (l *Label) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "label"

	return e.EncodeElement(LabelXML{
		Color:               colorutils.RGBASerializable(l.color),
		Font:                l.font,
		HorizontalAlignment: l.horizontalAlignment,
		VerticalAlignment:   l.verticalAlignment,
		Inverted:            l.Inverted,
		Padding:             &l.padding,
	}, start)
}

func (l *Label) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return nil
}
