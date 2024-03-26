package component

import (
	"image/color"

	"github.com/fglo/chopstiqs/pkg/option"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type textInputCursor struct {
	component
	drawer TextInputCursorDrawer
}

type TextInputCursorOptions struct {
	Width  option.OptInt
	Height option.OptInt

	Drawer TextInputCursorDrawer
}

func newTextInputCursor(options *TextInputCursorOptions) *textInputCursor {
	tic := &textInputCursor{
		drawer: &DefaultTextInputCursorDrawer{
			Color: color.RGBA{230, 230, 230, 255},
		},
	}

	tic.width = 1
	tic.height = 11

	if options != nil {
		if options.Width.IsSet() {
			tic.width = options.Width.Val()
		}

		if options.Height.IsSet() {
			tic.height = options.Height.Val()
		}

		if options.Drawer != nil {
			tic.drawer = options.Drawer
		}
	}

	tic.setUpComponent(options)

	return tic
}

func (tic *textInputCursor) setUpComponent(options *TextInputCursorOptions) {
	var componentOptions ComponentOptions
	tic.component.setUpComponent(&componentOptions)
}

func (tic *textInputCursor) ResetBlink() {
	tic.drawer.ResetBlink()
}

func (tic *textInputCursor) Draw() *ebiten.Image {
	tic.drawer.Draw(tic)
	return tic.image
}
