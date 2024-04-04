package component

import (
	"image/color"
	"math"

	"github.com/fglo/chopstiqs/pkg/event"
	"github.com/fglo/chopstiqs/pkg/input"
	"github.com/fglo/chopstiqs/pkg/option"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type Slider struct {
	component
	min float64
	max float64

	step       float64
	stepPixels float64

	value float64

	pressed  bool
	hovering bool

	handle *Button

	sliding bool

	SlidedEvent *event.Event

	PressedEvent  *event.Event
	ReleasedEvent *event.Event
	ClickedEvent  *event.Event

	firstPixelRowId       int
	secondPixelRowId      int
	lastPixelRowId        int
	penultimatePixelRowId int

	firstPixelColId       int
	secondPixelColId      int
	lastPixelColId        int
	penultimatePixelColId int

	drawer       SliderDrawer
	handleDrawer ButtonDrawer
}

type SliderOptions struct {
	Min          option.OptFloat
	Max          option.OptFloat
	Step         option.OptFloat
	DefaultValue option.OptFloat

	Width  option.OptInt
	Height option.OptInt

	Padding *Padding

	Drawer       SliderDrawer
	HandleDrawer ButtonDrawer
}

type SliderSlidedEventArgs struct {
	Slider *Slider
	Value  float64
	Change float64
}

type SliderPressedEventArgs struct {
	Slider *Slider
}

type SliderReleasedEventArgs struct {
	Slider *Slider
	Inside bool
}

type SliderClickedEventArgs struct {
	Slider *Slider
}

type SliderSlidedHandlerFunc func(args *SliderSlidedEventArgs)

type SliderPressedHandlerFunc func(args *SliderPressedEventArgs)

type SliderReleasedHandlerFunc func(args *SliderReleasedEventArgs)

type SliderClickedHandlerFunc func(args *SliderClickedEventArgs)

func NewSlider(opt *SliderOptions) *Slider {
	s := &Slider{
		SlidedEvent:   &event.Event{},
		PressedEvent:  &event.Event{},
		ReleasedEvent: &event.Event{},
		ClickedEvent:  &event.Event{},

		drawer: DefaultSliderDrawer{
			Color:         color.RGBA{230, 230, 230, 255},
			ColorPressed:  color.RGBA{230, 230, 230, 255},
			ColorHovered:  color.RGBA{230, 230, 230, 255},
			ColorDisabled: color.RGBA{150, 150, 150, 255},
		},
		handleDrawer: &DefaultButtonDrawer{
			Color:         color.RGBA{230, 230, 230, 255},
			ColorPressed:  color.RGBA{200, 200, 200, 255},
			ColorHovered:  color.RGBA{250, 250, 250, 255},
			ColorDisabled: color.RGBA{150, 150, 150, 255},
		},
	}

	s.component.width = 45
	s.component.height = 15

	if opt != nil {
		if opt.Min.IsSet() {
			s.min = opt.Min.Val()
		}

		if opt.Max.IsSet() {
			s.max = opt.Max.Val()
		}

		if opt.Step.IsSet() {
			s.step = opt.Step.Val()
		}

		if opt.DefaultValue.IsSet() {
			s.value = opt.DefaultValue.Val()
		} else {
			s.value = s.min
		}

		if opt.Width.IsSet() {
			s.component.width = opt.Width.Val()
		}

		if opt.Height.IsSet() {
			s.component.height = opt.Height.Val()
		}

		if opt.Drawer != nil {
			s.drawer = opt.Drawer
		}

		if opt.HandleDrawer != nil {
			s.handleDrawer = opt.HandleDrawer
		}
	}

	steps := math.Round((s.max-s.min)/s.step) + 1
	s.stepPixels = float64(s.component.width-4) / steps

	s.handle = NewButton(&ButtonOptions{Width: option.Int(7), Height: option.Int(s.component.height), Drawer: s.handleDrawer})
	s.handle.SetPosition(s.calcHandlePosition(), 0)

	s.handle.AddPressedHandler(func(args *ButtonPressedEventArgs) {
		s.sliding = true
	})

	s.handle.AddReleasedHandler(func(args *ButtonReleasedEventArgs) {
		s.sliding = false
	})

	s.setUpComponent(opt)

	s.setDrawingDimensions()

	return s
}

func (s *Slider) setUpComponent(opt *SliderOptions) {
	var componentOptions ComponentOptions

	if opt != nil {
		componentOptions = ComponentOptions{
			Padding: opt.Padding,
		}
	}

	s.component.setUpComponent(&componentOptions)

	s.component.AddCursorEnterHandler(func(args *ComponentCursorEnterEventArgs) {
		if !s.disabled {
			s.hovering = true
		}
	})

	s.component.AddCursorExitHandler(func(args *ComponentCursorExitEventArgs) {
		s.hovering = false
	})

	s.component.AddMouseButtonPressedHandler(func(args *ComponentMouseButtonPressedEventArgs) {
		if !s.disabled && args.Button == ebiten.MouseButtonLeft {
			s.pressed = true
			s.sliding = true

			if s.handle.posX >= 0 && s.handle.posX <= float64(s.width) {
				s.updateHandlePosition()
			}

			s.eventManager.Fire(s.PressedEvent, &SliderPressedEventArgs{
				Slider: s,
			})
		}
	})

	s.component.AddMouseButtonReleasedHandler(func(args *ComponentMouseButtonReleasedEventArgs) {
		if s.pressed && args.Button == ebiten.MouseButtonLeft {
			s.pressed = false
			s.sliding = false

			s.eventManager.Fire(s.ReleasedEvent, &SliderReleasedEventArgs{
				Slider: s,
				Inside: args.Inside,
			})

			if !s.disabled {
				s.eventManager.Fire(s.ClickedEvent, &SliderClickedEventArgs{
					Slider: s,
				})
			}
		}
	})
}

func (s *Slider) calcHandlePosition() float64 {
	return (s.value / s.step) * s.stepPixels
}

func (s *Slider) setDrawingDimensions() {
	s.firstPixelColId = s.padding.Left * 4
	s.secondPixelColId = s.firstPixelColId + 4

	s.lastPixelColId = (s.width+s.padding.Left)*4 - 4
	s.penultimatePixelColId = s.lastPixelColId - 4

	s.firstPixelRowId = s.padding.Top + 3
	s.secondPixelRowId = s.firstPixelRowId + 1

	s.lastPixelRowId = s.height + s.padding.Top - 4
	s.penultimatePixelRowId = s.lastPixelRowId - 1
}

func (s *Slider) AddComponent(Component) {
	panic("Slider can't have children")
}

func (s *Slider) SetBackgroundColor(color color.RGBA) {
	s.container.SetBackgroundColor(color)
}

func (s *Slider) GetBackgroundColor() color.RGBA {
	return s.container.GetBackgroundColor()
}

func (s *Slider) SetPosition(posX, posY float64) {
	s.component.SetPosition(posX, posY)
	s.handle.SetPosition(s.handle.posX, s.handle.posY)
}

func (s *Slider) setContainer(container *Container) {
	s.component.setContainer(container)
	s.handle.setContainer(s.container)
}

func (s *Slider) SetDisabled(disabled bool) {
	s.handle.SetDisabled(disabled)
	s.component.SetDisabled(disabled)
}

func (s *Slider) AddSlidedHandler(f SliderSlidedHandlerFunc) *Slider {
	s.SlidedEvent.AddHandler(func(args interface{}) {
		f(args.(*SliderSlidedEventArgs))
	})

	return s
}

func (s *Slider) GetValue() float64 {
	return s.value
}

// FireEvents checks if the mouse cursor is inside the component and fires events accordingly.
func (s *Slider) FireEvents() {
	s.component.FireEvents()
	s.handle.FireEvents()
}

func (s *Slider) Draw() *ebiten.Image {
	if s.hidden {
		return s.image
	}

	s.drawer.Draw(s)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(s.handle.Position())
	handleImg := s.handle.Draw()
	s.image.DrawImage(handleImg, op)

	s.component.Draw()

	return s.image
}

func (s *Slider) Set(value float64) {
	prevValue := s.value
	s.value = value
	s.handle.SetPosition(s.calcHandlePosition(), 0)
	s.fireEventOnChange(prevValue)
}

func (s *Slider) SetToMin() {
	prevValue := s.value
	s.value = s.min
	s.handle.SetPosX(2)
	s.fireEventOnChange(prevValue)
}

func (s *Slider) SetToMax() {
	prevValue := s.value
	s.value = s.max
	s.handle.SetPosX(float64(s.width-s.handle.width) - 2)
	s.fireEventOnChange(prevValue)
}

func (s *Slider) fireEventOnChange(prevValue float64) {
	change := math.Round((s.value - prevValue) / s.step)
	if math.Abs(change) < s.step {
		change = 0
	}

	if change != 0 {
		s.eventManager.Fire(s.SlidedEvent, &SliderSlidedEventArgs{
			Slider: s,
			Change: change,
			Value:  s.value,
		})
	}
}

func (s *Slider) updateHandlePosition() {
	currCursorPosX := input.CursorPosX

	switch {
	case currCursorPosX >= s.rect.Max.X:
		s.SetToMax()
	case currCursorPosX <= s.rect.Min.X:
		s.SetToMin()
	default:
		diff := float64(currCursorPosX) - s.absPosX
		steps := math.Floor(diff / s.stepPixels)
		value := float64(steps) * s.step
		newHandlePosX := s.calcHandlePosition()

		switch {
		case value >= s.max || newHandlePosX > float64(s.rect.Max.X-s.handle.width)-s.absPosX:
			s.SetToMax()
		case value <= s.min || newHandlePosX < float64(s.rect.Min.X)-s.absPosX:
			s.SetToMin()
		default:
			s.Set(value)
		}
	}
}
