package component

import (
	"image/color"

	"github.com/fglo/chopstiqs/pkg/event"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type Button struct {
	component

	pressed  bool
	hovering bool

	PressedEvent  *event.Event
	ReleasedEvent *event.Event
	ClickedEvent  *event.Event

	label *Label

	drawer ButtonDrawer
}

type ButtonOptions struct {
	Width  *int
	Height *int

	Drawer ButtonDrawer

	Label *Label

	Padding *Padding
}

type ButtonPressedEventArgs struct {
	Button *Button
}

type ButtonReleasedEventArgs struct {
	Button *Button
	Inside bool
}

type ButtonClickedEventArgs struct {
	Button *Button
}

type ButtonPressedHandlerFunc func(args *ButtonPressedEventArgs)

type ButtonReleasedHandlerFunc func(args *ButtonReleasedEventArgs)

type ButtonClickedHandlerFunc func(args *ButtonClickedEventArgs)

func NewButton(options *ButtonOptions) *Button {
	b := &Button{
		PressedEvent:  &event.Event{},
		ReleasedEvent: &event.Event{},
		ClickedEvent:  &event.Event{},

		drawer: &DefaultButtonDrawer{
			Color:         color.RGBA{230, 230, 230, 255},
			ColorPressed:  color.RGBA{200, 200, 200, 255},
			ColorHovered:  color.RGBA{250, 250, 250, 255},
			ColorDisabled: color.RGBA{150, 150, 150, 255},
		},
	}

	b.width = 45
	b.height = 15

	if options != nil {
		if options.Width != nil {
			b.width = *options.Width
		}

		if options.Height != nil {
			b.height = *options.Height
		}

		if options.Label != nil {
			b.SetLabel(options.Label)

			b.PressedEvent.AddHandler(func(args interface{}) {
				b.label.Inverted = true
			})

			b.ReleasedEvent.AddHandler(func(args interface{}) {
				b.label.Inverted = false
			})
		}

		if options.Drawer != nil {
			b.drawer = options.Drawer
		}
	}

	b.setUpComponent(options)

	return b
}

func (b *Button) setUpComponent(options *ButtonOptions) {
	var componentOptions ComponentOptions

	if options != nil {
		componentOptions = ComponentOptions{
			Padding: options.Padding,
		}
	}

	b.component.setUpComponent(&componentOptions)

	b.component.AddCursorEnterHandler(func(args *ComponentCursorEnterEventArgs) {
		if !b.disabled {
			b.hovering = true
		}
	})

	b.component.AddCursorExitHandler(func(args *ComponentCursorExitEventArgs) {
		b.hovering = false
	})

	b.component.AddMouseButtonPressedHandler(func(args *ComponentMouseButtonPressedEventArgs) {
		if !b.disabled && args.Button == ebiten.MouseButtonLeft {
			b.pressed = true
			b.PressedEvent.Fire(&ButtonPressedEventArgs{
				Button: b,
			})
		}
	})

	b.component.AddMouseButtonReleasedHandler(func(args *ComponentMouseButtonReleasedEventArgs) {
		if !b.disabled && args.Button == ebiten.MouseButtonLeft {
			b.pressed = false
			b.ReleasedEvent.Fire(&ButtonReleasedEventArgs{
				Button: b,
				Inside: args.Inside,
			})

			b.ClickedEvent.Fire(&ButtonClickedEventArgs{
				Button: b,
			})
		}
	})
}

func (b *Button) AddPressedHandler(f ButtonPressedHandlerFunc) *Button {
	b.PressedEvent.AddHandler(func(args interface{}) { f(args.(*ButtonPressedEventArgs)) })

	return b
}

func (b *Button) AddReleasedHandler(f ButtonReleasedHandlerFunc) *Button {
	b.ReleasedEvent.AddHandler(func(args interface{}) { f(args.(*ButtonReleasedEventArgs)) })

	return b
}

func (b *Button) AddClickedHandler(f ButtonClickedHandlerFunc) *Button {
	b.ClickedEvent.AddHandler(func(args interface{}) { f(args.(*ButtonClickedEventArgs)) })

	return b
}

// SetLabel sets the label of the button and sets the dimensions of the button accordingly.
func (b *Button) SetLabel(label *Label) {
	b.label = label
	b.label.alignHorizontally = b.label.centerHorizontally
	b.label.alignVertically = b.label.centerVertically

	b.label.SetPosistion(float64(b.label.textBounds.Dx())/2+5, 7.5)

	b.SetDimensions(b.label.width+10, 15)
}

func (b *Button) Draw() *ebiten.Image {
	if b.hidden {
		return b.image
	}

	b.drawer.Draw(b)

	if b.label != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(b.label.Position())
		b.image.DrawImage(b.label.Draw(), op)
	}

	b.component.Draw()

	return b.image
}
