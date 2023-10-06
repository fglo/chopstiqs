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

	color         color.RGBA
	colorPressed  color.RGBA
	colorHovered  color.RGBA
	colorDisabled color.RGBA
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

func NewButton(options *ButtonOptions) *Button {
	b := &Button{
		PressedEvent:  &event.Event{},
		ReleasedEvent: &event.Event{},
		ClickedEvent:  &event.Event{},

		color:         color.RGBA{230, 230, 230, 255},
		colorPressed:  color.RGBA{200, 200, 200, 255},
		colorHovered:  color.RGBA{220, 220, 220, 255},
		colorDisabled: color.RGBA{150, 150, 150, 255},
	}

	b.widget = b.createWidget(45, 15)

	if options != nil {
		for _, o := range options.opts {
			o(b)
		}
	}

	return b
}

func (o *ButtonOptions) Label(labelText string, color color.RGBA) *ButtonOptions {
	// TODO: change deprecated function
	bounds := text.BoundString(fontutils.DefaultFontFace, labelText) // nolint

	lblOpts := &LabelOptions{}
	label := NewLabel(labelText, lblOpts.Color(color).Centered())
	label.SetPosistion(float64(bounds.Dx())/2+5, 7.5)

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

func (o *ButtonOptions) Color(color, colorPressed, colorHovered, colorDisabled color.RGBA) *ButtonOptions {
	o.opts = append(o.opts, func(b *Button) {
		b.color = color
		b.colorPressed = colorPressed
		b.colorHovered = colorHovered
		b.colorDisabled = colorDisabled
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
		rowNumber := b.pixelCols * i
		for j := 0; j < b.pixelCols; j += 4 {
			if i == 0 && (j == 0 || j == b.lastPixelColId) || i == b.lastPixelRowId && (j == 0 || j == b.lastPixelColId) {
				continue
			} else if (i == 0 || i == b.lastPixelRowId || j == 0 || j == b.lastPixelColId) ||
				(j > 4 && j < b.penultimatePixelColId && i > 1 && i < b.penultimatePixelRowId) {
				arr[j+rowNumber] = b.color.R
				arr[j+1+rowNumber] = b.color.G
				arr[j+2+rowNumber] = b.color.B
				arr[j+3+rowNumber] = b.color.A
			} else {
				arr[j+rowNumber] = b.container.backgroundColor.R
				arr[j+1+rowNumber] = b.container.backgroundColor.G
				arr[j+2+rowNumber] = b.container.backgroundColor.B
				arr[j+3+rowNumber] = b.container.backgroundColor.A
			}
		}
	}

	return arr
}

func (b *Button) drawPressed() []byte {
	arr := make([]byte, b.pixelRows*b.pixelCols)

	for i := 0; i < b.pixelRows; i++ {
		rowNumber := b.pixelCols * i

		for j := 0; j < b.pixelCols; j += 4 {
			if i == 0 && (j == 0 || j == b.lastPixelColId) || i == b.lastPixelRowId && (j == 0 || j == b.lastPixelColId) {
				continue
			} else if i == 0 || i == b.lastPixelRowId || j == 0 || j == b.lastPixelColId {
				arr[j+rowNumber] = b.colorPressed.R
				arr[j+1+rowNumber] = b.colorPressed.G
				arr[j+2+rowNumber] = b.colorPressed.B
				arr[j+3+rowNumber] = b.colorPressed.A
			} else {
				arr[j+rowNumber] = b.container.backgroundColor.R
				arr[j+1+rowNumber] = b.container.backgroundColor.G
				arr[j+2+rowNumber] = b.container.backgroundColor.B
				arr[j+3+rowNumber] = b.container.backgroundColor.A
			}
		}
	}

	return arr
}

func (b *Button) drawHovered() []byte {
	arr := make([]byte, b.pixelRows*b.pixelCols)

	for i := 0; i < b.pixelRows; i++ {
		rowNumber := b.pixelCols * i

		for j := 0; j < b.pixelCols; j += 4 {
			if (i == 0 || i == b.lastPixelRowId) && (j == 0 || j == b.lastPixelColId) {
				continue
			} else if (i == 0 || i == b.lastPixelRowId || j == 0 || j == b.lastPixelColId) ||
				j > 4 && j < b.penultimatePixelColId && i > 1 && i < b.penultimatePixelRowId {
				arr[j+rowNumber] = b.colorHovered.R
				arr[j+1+rowNumber] = b.colorHovered.G
				arr[j+2+rowNumber] = b.colorHovered.B
				arr[j+3+rowNumber] = b.colorHovered.A
			} else {
				arr[j+rowNumber] = b.container.backgroundColor.R
				arr[j+1+rowNumber] = b.container.backgroundColor.G
				arr[j+2+rowNumber] = b.container.backgroundColor.B
				arr[j+3+rowNumber] = b.container.backgroundColor.A
			}
		}
	}

	return arr
}

func (b *Button) drawDisabled() []byte {
	arr := make([]byte, b.pixelRows*b.pixelCols)

	for i := 0; i < b.pixelRows; i++ {
		rowNumber := b.pixelCols * i

		for j := 0; j < b.pixelCols; j += 4 {
			if (i == 0 || i == b.lastPixelRowId) && (j == 0 || j == b.lastPixelColId) {
				continue
			} else if i == 0 || i == b.lastPixelRowId || j == 0 || j == b.lastPixelColId ||
				(j > 4 && j < b.penultimatePixelColId && i > 1 && i < b.penultimatePixelRowId) {
				arr[j+rowNumber] = b.colorDisabled.R
				arr[j+1+rowNumber] = b.colorDisabled.G
				arr[j+2+rowNumber] = b.colorDisabled.B
				arr[j+3+rowNumber] = b.colorDisabled.A
			} else {
				arr[j+rowNumber] = b.container.backgroundColor.R
				arr[j+1+rowNumber] = b.container.backgroundColor.G
				arr[j+2+rowNumber] = b.container.backgroundColor.B
				arr[j+3+rowNumber] = b.container.backgroundColor.A
			}
		}
	}

	return arr
}

func (b *Button) createWidget(width, height int) widget {
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

	return *NewWidget(width, height, widgetOptions)
}
