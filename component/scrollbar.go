package component

import (
	"image/color"
	"math"

	"github.com/fglo/chopstiqs/event"
	"github.com/fglo/chopstiqs/input"
	"github.com/fglo/chopstiqs/option"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type ScrollBar struct {
	component

	step       float64
	stepPixels float64

	value float64

	pressed  bool
	hovering bool

	handle *Button

	scrolling bool

	ScrolledEvent *event.Event

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

	drawer       ScrollBarDrawer
	handleDrawer ButtonDrawer
}

type ScrollBarOptions struct {
	Container container

	Step         option.OptFloat
	DefaultValue option.OptFloat

	Width option.OptInt

	Padding *Padding

	Drawer       ScrollBarDrawer
	HandleDrawer ButtonDrawer
}

type ScrollBarScrolledEventArgs struct {
	ScrollBar *ScrollBar
	Value     float64
	Change    float64
}

type ScrollBarPressedEventArgs struct {
	ScrollBar *ScrollBar
}

type ScrollBarReleasedEventArgs struct {
	ScrollBar *ScrollBar
	Inside    bool
}

type ScrollBarClickedEventArgs struct {
	ScrollBar *ScrollBar
}

type ScrollBarScrolledHandlerFunc func(args *ScrollBarScrolledEventArgs)

type ScrollBarPressedHandlerFunc func(args *ScrollBarPressedEventArgs)

type ScrollBarReleasedHandlerFunc func(args *ScrollBarReleasedEventArgs)

type ScrollBarClickedHandlerFunc func(args *ScrollBarClickedEventArgs)

func NewScrollBar(opt *ScrollBarOptions) *ScrollBar {
	s := &ScrollBar{
		ScrolledEvent: &event.Event{},
		PressedEvent:  &event.Event{},
		ReleasedEvent: &event.Event{},
		ClickedEvent:  &event.Event{},

		drawer: DefaultScrollBarDrawer{
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

	width := 15
	height := 0

	s.step = 1

	if opt != nil {
		if opt.Container != nil {
			height = opt.Container.Height()
		} else {
			_ = 0 // lint
			// TODO: errors from component constructors
		}

		if opt.Width.IsSet() {
			width = opt.Width.Val()
		}

		s.SetDimensions(width, height)

		if opt.Step.IsSet() {
			s.step = opt.Step.Val()
		}

		if opt.DefaultValue.IsSet() {
			s.value = opt.DefaultValue.Val()
		} else {
			s.value = 0
		}

		if opt.Drawer != nil {
			s.drawer = opt.Drawer
		}

		if opt.HandleDrawer != nil {
			s.handleDrawer = opt.HandleDrawer
		}
	}

	steps := math.Round(float64(s.component.height)/s.step) + 1
	s.stepPixels = float64(s.component.height-4) / steps

	s.handle = NewButton(&ButtonOptions{Width: option.Int(7), Height: option.Int(s.component.height), Drawer: s.handleDrawer})
	s.handle.setContainer(s)
	s.handle.SetPosition(s.calcHandlePosition(), 0)

	s.handle.AddPressedHandler(func(args *ButtonPressedEventArgs) {
		s.scrolling = true
	})

	s.handle.AddReleasedHandler(func(args *ButtonReleasedEventArgs) {
		s.scrolling = false
	})

	s.setUpComponent(opt)

	s.setDrawingDimensions()

	return s
}

func (s *ScrollBar) setUpComponent(opt *ScrollBarOptions) {
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
			s.scrolling = true

			if s.handle.posX >= 0 && s.handle.posX <= float64(s.width) {
				s.updateHandlePosition()
			}

			s.eventManager.Fire(s.PressedEvent, &ScrollBarPressedEventArgs{
				ScrollBar: s,
			})
		}
	})

	s.component.AddMouseButtonReleasedHandler(func(args *ComponentMouseButtonReleasedEventArgs) {
		if s.pressed && args.Button == ebiten.MouseButtonLeft {
			s.pressed = false
			s.scrolling = false

			s.eventManager.Fire(s.ReleasedEvent, &ScrollBarReleasedEventArgs{
				ScrollBar: s,
				Inside:    args.Inside,
			})

			if !s.disabled {
				s.eventManager.Fire(s.ClickedEvent, &ScrollBarClickedEventArgs{
					ScrollBar: s,
				})
			}
		}
	})
}

func (s *ScrollBar) setContainer(container container) {
	s.component.setContainer(container)
	s.SetHeight(container.Height())
}

func (s *ScrollBar) calcHandlePosition() float64 {
	return (s.value / s.step) * s.stepPixels
}

func (s *ScrollBar) setDrawingDimensions() {
	s.firstPixelColId = s.padding.Left * 4
	s.secondPixelColId = s.firstPixelColId + 4

	s.lastPixelColId = (s.width+s.padding.Left)*4 - 4
	s.penultimatePixelColId = s.lastPixelColId - 4

	s.firstPixelRowId = s.padding.Top + 3
	s.secondPixelRowId = s.firstPixelRowId + 1

	s.lastPixelRowId = s.height + s.padding.Top - 4
	s.penultimatePixelRowId = s.lastPixelRowId - 1
}

func (s *ScrollBar) SetBackgroundColor(color color.RGBA) {
	s.container.SetBackgroundColor(color)
}

func (s *ScrollBar) GetBackgroundColor() color.RGBA {
	return s.container.GetBackgroundColor()
}

func (s *ScrollBar) SetPosition(posX, posY float64) {
	s.component.SetPosition(posX, posY)
	if s.handle != nil {
		s.handle.RecalculateAbsPosition()
	}
}

func (s *ScrollBar) RecalculateAbsPosition() {
	s.component.RecalculateAbsPosition()
	if s.handle != nil {
		s.handle.RecalculateAbsPosition()
	}
}

func (s *ScrollBar) SetDisabled(disabled bool) {
	s.handle.SetDisabled(disabled)
	s.component.SetDisabled(disabled)
}

func (s *ScrollBar) AddScrolledHandler(f ScrollBarScrolledHandlerFunc) *ScrollBar {
	s.ScrolledEvent.AddHandler(func(args any) {
		f(args.(*ScrollBarScrolledEventArgs))
	})

	return s
}

func (s *ScrollBar) GetValue() float64 {
	return s.value
}

// FireEvents checks if the mouse cursor is inside the component and fires events accordingly.
func (s *ScrollBar) FireEvents() {
	s.component.FireEvents()
	s.handle.FireEvents()
}

func (s *ScrollBar) Draw() *ebiten.Image {
	if s.hidden {
		return s.emptyImage
	}

	s.drawer.Draw(s)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(s.handle.Position())
	handleImg := s.handle.Draw()
	s.image.DrawImage(handleImg, op)

	s.component.Draw()

	return s.image
}

func (s *ScrollBar) Set(value float64) {
	prevValue := s.value
	s.value = value
	s.handle.SetPosition(s.calcHandlePosition(), 0)
	s.fireEventOnChange(prevValue)
}

func (s *ScrollBar) SetToMin() {
	prevValue := s.value
	s.value = 0
	s.handle.SetPosX(2)
	s.fireEventOnChange(prevValue)
}

func (s *ScrollBar) SetToMax() {
	prevValue := s.value
	s.value = float64(s.height - s.handle.height) // TODO: it might be wrong
	s.handle.SetPosX(float64(s.width-s.handle.width) - 2)
	s.fireEventOnChange(prevValue)
}

func (s *ScrollBar) fireEventOnChange(prevValue float64) {
	change := math.Round((s.value - prevValue) / s.step)
	if math.Abs(change) < s.step {
		change = 0
	}

	if change != 0 {
		s.eventManager.Fire(s.ScrolledEvent, &ScrollBarScrolledEventArgs{
			ScrollBar: s,
			Change:    change,
			Value:     s.value,
		})
	}
}

func (s *ScrollBar) updateHandlePosition() {
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
		case value >= float64(s.height) || newHandlePosX > float64(s.rect.Max.X-s.handle.width)-s.absPosX:
			s.SetToMax()
		case value <= 0 || newHandlePosX < float64(s.rect.Min.X)-s.absPosX:
			s.SetToMin()
		default:
			s.Set(value)
		}
	}
}
