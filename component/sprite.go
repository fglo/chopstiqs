package component

import (
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
