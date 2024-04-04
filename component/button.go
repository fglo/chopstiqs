package component

import (
	"image/color"

	"github.com/fglo/chopstiqs/event"
	"github.com/fglo/chopstiqs/option"
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
	Width  option.OptInt
	Height option.OptInt

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

func NewButton(opt *ButtonOptions) *Button {
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

	if opt != nil {
		if opt.Width.IsSet() {
			b.width = opt.Width.Val()
		}

		if opt.Height.IsSet() {
			b.height = opt.Height.Val()
		}

		if opt.Label != nil {
			b.SetLabel(opt.Label)

			b.PressedEvent.AddHandler(func(args interface{}) {
				b.label.Inverted = true
			})

			b.ReleasedEvent.AddHandler(func(args interface{}) {
				b.label.Inverted = false
			})
		}

		if opt.Drawer != nil {
			b.drawer = opt.Drawer
		}
	}

	b.setUpComponent(opt)

	return b
}

func (b *Button) setUpComponent(opt *ButtonOptions) {
	var componentOptions ComponentOptions

	if opt != nil {
		componentOptions = ComponentOptions{
			Padding: opt.Padding,
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
			b.eventManager.Fire(b.PressedEvent, &ButtonPressedEventArgs{
				Button: b,
			})
		}
	})

	b.component.AddMouseButtonReleasedHandler(func(args *ComponentMouseButtonReleasedEventArgs) {
		if b.pressed && args.Button == ebiten.MouseButtonLeft {
			b.pressed = false
			b.eventManager.Fire(b.ReleasedEvent, &ButtonReleasedEventArgs{
				Button: b,
				Inside: args.Inside,
			})

			if !b.disabled {
				b.eventManager.Fire(b.ClickedEvent, &ButtonClickedEventArgs{
					Button: b,
				})
			}
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

	b.label.SetPosistion(float64(b.label.bounds.Dx())/2+5, 7.5)

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
