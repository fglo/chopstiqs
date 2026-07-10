package component

import (
	"image/color"

	"github.com/fglo/chopstiqs/option"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type textInputCursor struct {
	component
	drawer     TextInputCursorDrawer
	frameCount int
	textInput  *TextInput
}

type TextInputCursorOptions struct {
	Width  option.OptInt
	Height option.OptInt

	Drawer TextInputCursorDrawer
}

func newTextInputCursor(textInput *TextInput, options *TextInputCursorOptions) *textInputCursor {
	tic := &textInputCursor{
		drawer: &DefaultTextInputCursorDrawer{
			Color: color.RGBA{230, 230, 230, 255},
		},
		textInput: textInput,
	}

	width := 1
	height := 11

	tic.SetDimensions(width, height)

	if options != nil {
		if options.Width.IsSet() {
			tic.width = options.Width.Val()
		}

		if options.Height.IsSet() {
			tic.height = options.Height.Val()
		}

		tic.SetDimensions(width, height)

		if options.Drawer != nil {
			tic.drawer = options.Drawer
		}
	}

	tic.setUpComponent()

	return tic
}

func (tic *textInputCursor) setUpComponent() {
	var componentOptions ComponentOptions
	tic.component.setUpComponent(&componentOptions)
}

func (tic *textInputCursor) ResetBlink() {
	tic.frameCount = 0
}

func (tic *textInputCursor) incFrameCount() {
	tic.frameCount = (tic.frameCount + 1) % 70
}

func (tic *textInputCursor) Draw() *ebiten.Image {
	tic.incFrameCount()
	tic.drawer.Draw(tic)
	return tic.image
}
