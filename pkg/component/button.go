package component

import (
	"image/color"

	"github.com/fglo/chopstiqs/internal/colorutils"
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

	color         color.RGBA
	colorPressed  color.RGBA
	colorHovered  color.RGBA
	colorDisabled color.RGBA
}

type ButtonOptions struct {
	Color         color.Color
	ColorPressed  color.Color
	ColorHovered  color.Color
	ColorDisabled color.Color

	Label *Label
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
	width := 45
	height := 15

	b := &Button{
		PressedEvent:  &event.Event{},
		ReleasedEvent: &event.Event{},
		ClickedEvent:  &event.Event{},

		color:         color.RGBA{230, 230, 230, 255},
		colorPressed:  color.RGBA{200, 200, 200, 255},
		colorHovered:  color.RGBA{250, 250, 250, 255},
		colorDisabled: color.RGBA{150, 150, 150, 255},
	}

	b.component = b.createComponent(width, height)

	if options != nil {
		if options.Label != nil {
			b.SetLabel(options.Label)

			b.PressedEvent.AddHandler(func(args interface{}) {
				b.label.Inverted = true
			})

			b.ReleasedEvent.AddHandler(func(args interface{}) {
				b.label.Inverted = false
			})
		}

		if options.Color != nil {
			b.color = colorutils.ColorToRGBA(options.Color)
		}

		if options.ColorPressed != nil {
			b.colorPressed = colorutils.ColorToRGBA(options.ColorPressed)
		}

		if options.ColorHovered != nil {
			b.colorHovered = colorutils.ColorToRGBA(options.ColorHovered)
		}

		if options.ColorDisabled != nil {
			b.colorDisabled = colorutils.ColorToRGBA(options.ColorDisabled)
		}
	}

	return b
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
	backgroundColor := b.container.GetBackgroundColor()

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
				arr[j+rowNumber] = backgroundColor.R
				arr[j+1+rowNumber] = backgroundColor.G
				arr[j+2+rowNumber] = backgroundColor.B
				arr[j+3+rowNumber] = backgroundColor.A
			}
		}
	}

	return arr
}

func (b *Button) drawPressed() []byte {
	arr := make([]byte, b.pixelRows*b.pixelCols)
	backgroundColor := b.container.GetBackgroundColor()

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
				arr[j+rowNumber] = backgroundColor.R
				arr[j+1+rowNumber] = backgroundColor.G
				arr[j+2+rowNumber] = backgroundColor.B
				arr[j+3+rowNumber] = backgroundColor.A
			}
		}
	}

	return arr
}

func (b *Button) drawHovered() []byte {
	arr := make([]byte, b.pixelRows*b.pixelCols)
	backgroundColor := b.container.GetBackgroundColor()

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
				arr[j+rowNumber] = backgroundColor.R
				arr[j+1+rowNumber] = backgroundColor.G
				arr[j+2+rowNumber] = backgroundColor.B
				arr[j+3+rowNumber] = backgroundColor.A
			}
		}
	}

	return arr
}

func (b *Button) drawDisabled() []byte {
	arr := make([]byte, b.pixelRows*b.pixelCols)
	backgroundColor := b.container.GetBackgroundColor()

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
				arr[j+rowNumber] = backgroundColor.R
				arr[j+1+rowNumber] = backgroundColor.G
				arr[j+2+rowNumber] = backgroundColor.B
				arr[j+3+rowNumber] = backgroundColor.A
			}
		}
	}

	return arr
}

func (b *Button) createComponent(width, height int) component {
	componentOptions := &ComponentOptions{}

	componentOptions.AddCursorEnterHandler(func(args *ComponentCursorEnterEventArgs) {
		if !b.disabled {
			b.hovering = true
		}
	})

	componentOptions.AddCursorExitHandler(func(args *ComponentCursorExitEventArgs) {
		b.hovering = false
	})

	componentOptions.AddMouseButtonPressedHandler(func(args *ComponentMouseButtonPressedEventArgs) {
		if !b.disabled && args.Button == ebiten.MouseButtonLeft {
			b.pressed = true
			b.PressedEvent.Fire(&ButtonPressedEventArgs{
				Button: b,
			})
		}
	})

	componentOptions.AddMouseButtonReleasedHandler(func(args *ComponentMouseButtonReleasedEventArgs) {
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

	return *NewComponent(width, height, componentOptions)
}
