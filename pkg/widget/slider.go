package widget

import (
	"github.com/fglo/chopstiqs/pkg/event"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type Slider struct {
	widget
	min   float64
	max   float64
	step  float64
	value float64

	SlidedEvent *event.Event
}

type SliderOpt func(b *Slider)
type SliderOptions struct {
	opts []SliderOpt
}

type SliderSlidedEventArgs struct {
	Slider *Slider
}

type SliderSlidedHandlerFunc func(args *SliderSlidedEventArgs)

func NewSlider(posX, posY, min, max, step float64, options *SliderOptions) *Slider {
	s := &Slider{
		min:         min,
		max:         max,
		step:        step,
		value:       min,
		SlidedEvent: &event.Event{},
	}

	s.widget = s.createWidget(posX, posY, 45, 15)

	if options != nil {
		for _, o := range options.opts {
			o(s)
		}
	}

	return s
}

func (o *SliderOptions) Default(value float64) *SliderOptions {
	o.opts = append(o.opts, func(l *Slider) {
		l.Set(value)
	})

	return o
}

func (o *SliderOptions) SlidedHandler(f SliderSlidedHandlerFunc) *SliderOptions {
	o.opts = append(o.opts, func(l *Slider) {
		l.SlidedEvent.AddHandler(func(args interface{}) {
			f(args.(*SliderSlidedEventArgs))
		})
	})

	return o
}

func (s *Slider) Set(value float64) {
	s.value = value
}

func (s *Slider) Draw() *ebiten.Image {
	arr := make([]byte, s.pixelRows*s.pixelCols)

	half := s.lastPixelColId / 2

	for i := 0; i < s.pixelRows; i++ {
		for j := 0; j < s.pixelCols; j += 4 {
			if (i <= 4 || i >= s.lastPixelRowId-4) && (j == 0 || j == s.lastPixelColId) {
				continue
			} else if i == 4 || i == s.lastPixelRowId-4 || j == 0 || j == s.lastPixelColId {
				arr[j+s.pixelCols*i] = 230
				arr[j+1+s.pixelCols*i] = 230
				arr[j+2+s.pixelCols*i] = 230
				arr[j+3+s.pixelCols*i] = 255
			} else {
				arr[j+3+s.pixelCols*i] = 255
			}

			if j == half-12 || j == half+12 || ((j == half-8 || j == half+8) && (i == 2 || i == s.lastPixelRowId-2)) {
				arr[j+3+s.pixelCols*i] = 0
			} else if j > half-12 && j < half+12 && i > 1 && i < s.lastPixelRowId-1 {
				arr[j+s.pixelCols*i] = 230
				arr[j+1+s.pixelCols*i] = 230
				arr[j+2+s.pixelCols*i] = 230
				arr[j+3+s.pixelCols*i] = 255
			}
		}
	}

	s.image.WritePixels(arr)

	return s.image
}

func (s *Slider) createWidget(posX, posY float64, width, height int) widget {
	widgetOptions := &WidgetOptions{}

	return *NewWidget(posX, posY, width, height, widgetOptions)
}
