package component

import (
	"image/color"
	"math"

	"github.com/fglo/chopstiqs/pkg/event"
	"github.com/fglo/chopstiqs/pkg/input"
	"github.com/fglo/chopstiqs/pkg/to"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type Slider struct {
	component
	min float64
	max float64

	step       float64
	stepPixels int

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
	Min          *float64
	Max          *float64
	Step         *float64
	DefaultValue *float64

	Width  *int
	Height *int

	Padding *Padding

	Drawer       SliderDrawer
	HandleDrawer ButtonDrawer
}

type SliderSlidedEventArgs struct {
	Slider *Slider
	Value  float64
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

func NewSlider(options *SliderOptions) *Slider {
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

	if options != nil {
		if options.Min != nil {
			s.min = *options.Min
		}

		if options.Max != nil {
			s.max = *options.Max
		}

		if options.Step != nil {
			s.step = *options.Step
		}

		if options.DefaultValue != nil {
			s.Set(*options.DefaultValue)
		} else {
			s.Set(s.min)
		}

		if options.Width != nil {
			s.component.width = *options.Width
		}

		if options.Height != nil {
			s.component.height = *options.Height
		}

		if options.Drawer != nil {
			s.drawer = options.Drawer
		}

		if options.HandleDrawer != nil {
			s.handleDrawer = options.HandleDrawer
		}
	}

	steps := math.Round((s.max-s.min)/s.step) + 1
	s.stepPixels = int(math.Round(float64(s.component.width-4) / steps))

	s.handle = NewButton(&ButtonOptions{Width: to.Ptr(7), Height: &s.component.height, Drawer: s.handleDrawer})
	s.handle.SetPosision(s.value/s.step*float64(s.stepPixels), 0)

	s.setUpComponent(options)

	s.setDrawingDimensions()

	return s
}

func (s *Slider) setUpComponent(options *SliderOptions) {
	var componentOptions ComponentOptions

	if options != nil {
		componentOptions = ComponentOptions{
			Padding: options.Padding,
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
			s.PressedEvent.Fire(&SliderPressedEventArgs{
				Slider: s,
			})
		}
	})

	s.component.AddMouseButtonReleasedHandler(func(args *ComponentMouseButtonReleasedEventArgs) {
		if !s.disabled && args.Button == ebiten.MouseButtonLeft {
			s.pressed = false
			s.ReleasedEvent.Fire(&SliderReleasedEventArgs{
				Slider: s,
				Inside: args.Inside,
			})

			s.ClickedEvent.Fire(&SliderClickedEventArgs{
				Slider: s,
			})
		}
	})

	s.PressedEvent.AddHandler(func(args interface{}) {
		s.sliding = true
	})

	s.ReleasedEvent.AddHandler(func(args interface{}) {
		s.sliding = false
	})

	s.handle.AddPressedHandler(func(args *ButtonPressedEventArgs) {
		s.sliding = true
	})

	s.handle.AddReleasedHandler(func(args *ButtonReleasedEventArgs) {
		s.sliding = false
	})
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

func (s *Slider) SetPosision(posX, posY float64) {
	s.component.SetPosision(posX, posY)
	s.handle.SetPosision(s.handle.posX, s.handle.posY)
}

func (s *Slider) setContainer(container Container) {
	s.component.setContainer(container)
	s.handle.setContainer(s)
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

func (s *Slider) Set(value float64) {
	s.value = value
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

	if s.sliding && s.handle.posX >= 0 && s.handle.posX <= float64(s.width) {
		s.updateHandlePosition()
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(s.handle.Position())
	handleImg := s.handle.Draw()
	s.image.DrawImage(handleImg, op)

	s.component.Draw()

	return s.image
}

func (s *Slider) setToMin() {
	s.handle.SetPosX(2)
	s.value = s.min
	s.SlidedEvent.Fire(&SliderSlidedEventArgs{
		Slider: s,
		Value:  s.value,
	})
}

func (s *Slider) setToMax() {
	s.handle.SetPosX(float64(s.width-s.handle.width) - 2)
	s.value = s.max
	s.SlidedEvent.Fire(&SliderSlidedEventArgs{
		Slider: s,
		Value:  s.value,
	})
}

func (s *Slider) updateHandlePosition() {
	currCursorPosX := input.CursorPosX

	switch {
	case currCursorPosX > s.rect.Max.X:
		s.setToMax()
	case currCursorPosX < s.rect.Min.X:
		s.setToMin()
	default:
		diff := currCursorPosX - int(s.absPosX)
		steps := math.Round(float64(diff) / float64(s.stepPixels))
		newHandlePosX := steps * float64(s.stepPixels)

		s.value = float64(steps) * s.step

		switch {
		case s.value >= s.max || newHandlePosX >= float64(s.rect.Max.X-s.handle.width)-s.absPosX:
			s.setToMax()
		case s.value <= s.min || newHandlePosX <= float64(s.rect.Min.X)-s.absPosX:
			s.setToMin()
		default:
			s.handle.SetPosX(newHandlePosX)
			s.SlidedEvent.Fire(&SliderSlidedEventArgs{
				Slider: s,
				Value:  s.value,
			})
		}
	}
}
