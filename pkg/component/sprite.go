package component

import (
	"encoding/xml"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	component
	image *ebiten.Image
}

type SpriteOptions struct {
	Padding *Padding
}

func NewSprite(image *ebiten.Image, options *SpriteOptions) *Sprite {
	s := &Sprite{
		image: image,
	}

	bounds := image.Bounds()
	s.SetDimensions(bounds.Dx(), bounds.Dy())

	s.setUpComponent(options)

	return s
}

func (s *Sprite) setUpComponent(options *SpriteOptions) {
	var componentOptions ComponentOptions

	if options != nil {
		componentOptions = ComponentOptions{
			Padding: options.Padding,
		}
	}

	s.component.setUpComponent(&componentOptions)
}

func (s *Sprite) SetImage(image *ebiten.Image) {
	s.image = image
}

func (s *Sprite) Draw() *ebiten.Image {
	return s.image
}

func (s *Sprite) MarshalYAML() (interface{}, error) {
	return struct {
		Sprite SpriteOptions
	}{
		Sprite: SpriteOptions{
			Padding: &s.padding,
		},
	}, nil
}

type SpriteXML struct {
	XMLName xml.Name `xml:"sprite"`
	Padding *Padding `xml:"padding,attr,omitempty"`
}

func (s *Sprite) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "sprite"

	return e.EncodeElement(SpriteXML{
		Padding: &s.padding,
	}, start)
}

func (s *Sprite) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return nil
}
