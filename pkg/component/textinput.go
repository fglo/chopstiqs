package component

import (
	"image"
	"image/color"

	"github.com/fglo/chopstiqs/internal/colorutils"
	"github.com/fglo/chopstiqs/internal/fontutils"
	"github.com/fglo/chopstiqs/pkg/event"
	"github.com/fglo/chopstiqs/pkg/option"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type TextInput struct {
	component

	value string

	color color.RGBA
	font  font.Face

	textPosX   int
	textPosY   int
	textBounds image.Rectangle

	ClickedEvent *event.Event

	active bool

	drawer TextInputDrawer
}

type TextInputOptions struct {
	Width  option.OptInt
	Height option.OptInt

	Drawer TextInputDrawer

	Color color.Color
	Font  font.Face

	Padding *Padding
}

type TextInputClickedEventArgs struct {
	TextInput *TextInput
}

type TextInputClickedHandlerFunc func(args *TextInputClickedEventArgs)

func NewTextInput(options *TextInputOptions) *TextInput {
	ti := &TextInput{
		ClickedEvent: &event.Event{},

		color: color.RGBA{230, 230, 230, 255},
		font:  fontutils.DefaultFontFace,

		textPosX: 3,

		drawer: &DefaultTextInputDrawer{
			BorderColor: color.RGBA{230, 230, 230, 255},
		},
	}

	ti.width = 60
	ti.height = 15

	if options != nil {
		if options.Width.IsSet() {
			ti.width = options.Width.Val()
		}

		if options.Height.IsSet() {
			ti.height = options.Height.Val()
		}

		if options.Color != nil {
			ti.color = colorutils.ColorToRGBA(options.Color)
		}

		if options.Font != nil {
			ti.font = options.Font
		}

		if options.Drawer != nil {
			ti.drawer = options.Drawer
		}
	}

	ti.setUpComponent(options)

	return ti
}

func (ti *TextInput) setUpComponent(options *TextInputOptions) {
	var componentOptions ComponentOptions

	if options != nil {
		componentOptions = ComponentOptions{
			Padding: options.Padding,
		}
	}

	ti.component.setUpComponent(&componentOptions)

	// b.component.AddCursorEnterHandler(func(args *ComponentCursorEnterEventArgs) {
	// 	if !b.disabled {
	// 		b.hovering = true
	// 	}
	// })

	// b.component.AddCursorExitHandler(func(args *ComponentCursorExitEventArgs) {
	// 	b.hovering = false
	// })

	ti.component.AddMouseButtonReleasedHandler(func(args *ComponentMouseButtonReleasedEventArgs) {
		if !ti.disabled && args.Button == ebiten.MouseButtonLeft {
			ti.ClickedEvent.Fire(&TextInputClickedEventArgs{
				TextInput: ti,
			})
		}
	})
}

func (ti *TextInput) AddClickedHandler(f TextInputClickedHandlerFunc) *TextInput {
	ti.ClickedEvent.AddHandler(func(args interface{}) { f(args.(*TextInputClickedEventArgs)) })

	return ti
}

// SetValue sets the value of the text input.
func (ti *TextInput) SetValue(value string) {
	// TODO: change deprecated function
	bounds := text.BoundString(fontutils.DefaultFontFace, value) // nolint

	ti.value = value

	ti.textPosY = -bounds.Min.Y + ti.heightWithPadding/2 - 3
	ti.textBounds = bounds
}

func (ti *TextInput) Draw() *ebiten.Image {
	if ti.hidden {
		return ti.image
	}

	ti.drawer.Draw(ti)

	text.Draw(ti.image, ti.value, ti.font, ti.textPosX+ti.padding.Left, ti.textPosY+ti.padding.Top, ti.color)

	ti.component.Draw()

	return ti.image
}
