package widget

import (
	"image/color"

	"github.com/fglo/chopstiqs/pkg/event"
	"github.com/fglo/chopstiqs/pkg/fontutils"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Button struct {
	widget

	pressed  bool
	hovering bool

	PressedEvent  *event.Event
	ReleasedEvent *event.Event
	ClickedEvent  *event.Event

	label *Label
}

type ButtonOpt func(b *Button)
type ButtonOptions struct {
	opts []ButtonOpt
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

func NewButton(posX, posY float64, options *ButtonOptions) *Button {
	b := &Button{
		PressedEvent:  &event.Event{},
		ReleasedEvent: &event.Event{},
		ClickedEvent:  &event.Event{},
	}

	b.widget = b.createWidget(posX, posY, 45, 15)

	if options != nil {
		for _, o := range options.opts {
			o(b)
		}
	}

	return b
}

func (o *ButtonOptions) Text(labelText string, color color.RGBA) *ButtonOptions {
	// TODO: change deprecated function
	bounds := text.BoundString(fontutils.DefaultFontFace, labelText) // nolint

	lblOpts := &LabelOptions{}
	label := NewLabel(float64(bounds.Dx())/2+5, 7.5, labelText, color, lblOpts.Centered())

	o.PressedHandler(func(args *ButtonPressedEventArgs) {
		label.Inverted = true
	})

	o.ReleasedHandler(func(args *ButtonReleasedEventArgs) {
		label.Inverted = false
	})

	o.opts = append(o.opts, func(b *Button) {
		b.SetLabel(label)
	})

	return o
}

func (o *ButtonOptions) PressedHandler(f ButtonPressedHandlerFunc) *ButtonOptions {
	o.opts = append(o.opts, func(b *Button) {
		b.PressedEvent.AddHandler(func(args interface{}) {
			f(args.(*ButtonPressedEventArgs))
		})
	})

	return o
}

func (o *ButtonOptions) ReleasedHandler(f ButtonReleasedHandlerFunc) *ButtonOptions {
	o.opts = append(o.opts, func(b *Button) {
		b.ReleasedEvent.AddHandler(func(args interface{}) {
			f(args.(*ButtonReleasedEventArgs))
		})
	})

	return o
}

func (o *ButtonOptions) ClickedHandler(f ButtonClickedHandlerFunc) *ButtonOptions {
	o.opts = append(o.opts, func(b *Button) {
		b.ClickedEvent.AddHandler(func(args interface{}) {
			f(args.(*ButtonClickedEventArgs))
		})
	})

	return o
}

func (b *Button) SetLabel(label *Label) {
	b.label = label
	b.SetWidth(label.width + 10)
}

func (b *Button) Draw() *ebiten.Image {
	if b.pressed {
		b.image.WritePixels(b.drawPressed())
	} else if b.hovering {
		b.image.WritePixels(b.drawHovered())
	} else if b.disabled {
		b.image.WritePixels(b.drawDisabled())
	} else {
		b.image.WritePixels(b.draw())
	}

	if b.label != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(b.label.Position())
		b.image.DrawImage(b.label.Draw(), op)
	}

	return b.image
}

func (b *Button) draw() []byte {
	arr := make([]byte, b.pixelRows*b.pixelCols)

	for i := 0; i < b.pixelRows; i++ {
		for j := 0; j < b.pixelCols; j += 4 {
			if i == 0 && (j == 0 || j == b.lastPixelColId) || i == b.lastPixelRowId && (j == 0 || j == b.lastPixelColId) {
				continue
			} else if i == 0 || i == b.lastPixelRowId || j == 0 || j == b.lastPixelColId {
				arr[j+b.pixelCols*i] = 230
				arr[j+1+b.pixelCols*i] = 230
				arr[j+2+b.pixelCols*i] = 230
			} else if j > 4 && j < b.penultimatePixelColId && i > 1 && i < b.penultimatePixelRowId {
				arr[j+b.pixelCols*i] = 230
				arr[j+1+b.pixelCols*i] = 230
				arr[j+2+b.pixelCols*i] = 230
			} else {
				arr[j+b.pixelCols*i] = b.container.backgroundColor.R
				arr[j+1+b.pixelCols*i] = b.container.backgroundColor.G
				arr[j+2+b.pixelCols*i] = b.container.backgroundColor.B
			}
			arr[j+3+b.pixelCols*i] = 255
		}
	}

	return arr
}

func (b *Button) drawPressed() []byte {
	arr := make([]byte, b.pixelRows*b.pixelCols)

	for i := 0; i < b.pixelRows; i++ {
		for j := 0; j < b.pixelCols; j += 4 {
			if i == 0 && (j == 0 || j == b.lastPixelColId) || i == b.lastPixelRowId && (j == 0 || j == b.lastPixelColId) {
				continue
			} else if i == 0 || i == b.lastPixelRowId || j == 0 || j == b.lastPixelColId {
				arr[j+b.pixelCols*i] = 200
				arr[j+1+b.pixelCols*i] = 200
				arr[j+2+b.pixelCols*i] = 200
			} else {
				arr[j+b.pixelCols*i] = b.container.backgroundColor.R
				arr[j+1+b.pixelCols*i] = b.container.backgroundColor.G
				arr[j+2+b.pixelCols*i] = b.container.backgroundColor.B
			}
			arr[j+3+b.pixelCols*i] = 255
		}
	}

	return arr
}

func (b *Button) drawHovered() []byte {
	arr := make([]byte, b.pixelRows*b.pixelCols)

	for i := 0; i < b.pixelRows; i++ {
		for j := 0; j < b.pixelCols; j += 4 {
			if (i == 0 || i == b.lastPixelRowId) && (j == 0 || j == b.lastPixelColId) {
				continue
			} else if i == 0 || i == b.lastPixelRowId || j == 0 || j == b.lastPixelColId {
				arr[j+b.pixelCols*i] = 220
				arr[j+1+b.pixelCols*i] = 220
				arr[j+2+b.pixelCols*i] = 220
			} else if j > 4 && j < b.penultimatePixelColId && i > 1 && i < b.penultimatePixelRowId {
				arr[j+b.pixelCols*i] = 200
				arr[j+1+b.pixelCols*i] = 200
				arr[j+2+b.pixelCols*i] = 200
			} else {
				arr[j+b.pixelCols*i] = b.container.backgroundColor.R
				arr[j+1+b.pixelCols*i] = b.container.backgroundColor.G
				arr[j+2+b.pixelCols*i] = b.container.backgroundColor.B
			}
			arr[j+3+b.pixelCols*i] = 255
		}
	}

	return arr
}

func (b *Button) drawDisabled() []byte {
	arr := make([]byte, b.pixelRows*b.pixelCols)

	for i := 0; i < b.pixelRows; i++ {
		for j := 0; j < b.pixelCols; j += 4 {
			if (i == 0 || i == b.lastPixelRowId) && (j == 0 || j == b.lastPixelColId) {
				continue
			} else if i == 0 || i == b.lastPixelRowId || j == 0 || j == b.lastPixelColId ||
				(j > 4 && j < b.penultimatePixelColId && i > 1 && i < b.penultimatePixelRowId) {
				arr[j+b.pixelCols*i] = 175
				arr[j+1+b.pixelCols*i] = 175
				arr[j+2+b.pixelCols*i] = 175
				arr[j+3+b.pixelCols*i] = 255
			} else {
				arr[j+3+b.pixelCols*i] = 0
			}
			arr[j+3+b.pixelCols*i] = 255
		}
	}

	return arr
}

func (b *Button) createWidget(posX, posY float64, width, height int) widget {
	widgetOptions := &WidgetOptions{}

	widgetOptions.CursorEnterHandler(func(args *WidgetCursorEnterEventArgs) {
		if !b.disabled {
			b.hovering = true
		}
	})

	widgetOptions.CursorExitHandler(func(args *WidgetCursorExitEventArgs) {
		b.hovering = false
	})

	widgetOptions.MouseButtonPressedHandler(func(args *WidgetMouseButtonPressedEventArgs) {
		if !b.disabled && args.Button == ebiten.MouseButtonLeft {
			b.pressed = true
			b.PressedEvent.Fire(&ButtonPressedEventArgs{
				Button: b,
			})
		}
	})

	widgetOptions.MouseButtonReleasedHandler(func(args *WidgetMouseButtonReleasedEventArgs) {
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

	return *NewWidget(posX, posY, width, height, widgetOptions)
}
