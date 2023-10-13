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

	LeftPadding   *int
	RightPadding  *int
	TopPadding    *int
	BottomPadding *int
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

	b.component = b.createComponent(width, height, options)

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

func (b *Button) createComponent(width, height int, options *ButtonOptions) component {
	var componentOptions ComponentOptions

	if options != nil {
		componentOptions = ComponentOptions{
			LeftPadding:   options.LeftPadding,
			RightPadding:  options.RightPadding,
			TopPadding:    options.TopPadding,
			BottomPadding: options.BottomPadding,
		}
	}

	component := NewComponent(width, height, &componentOptions)

	component.AddCursorEnterHandler(func(args *ComponentCursorEnterEventArgs) {
		if !b.disabled {
			b.hovering = true
		}
	})

	component.AddCursorExitHandler(func(args *ComponentCursorExitEventArgs) {
		b.hovering = false
	})

	component.AddMouseButtonPressedHandler(func(args *ComponentMouseButtonPressedEventArgs) {
		if !b.disabled && args.Button == ebiten.MouseButtonLeft {
			b.pressed = true
			b.PressedEvent.Fire(&ButtonPressedEventArgs{
				Button: b,
			})
		}
	})

	component.AddMouseButtonReleasedHandler(func(args *ComponentMouseButtonReleasedEventArgs) {
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

	return *component
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

	b.component.Draw()

	return b.image
}

func (b *Button) isCorner(rowId, colId int) bool {
	return (rowId == b.firstPixelRowId || rowId == b.lastPixelRowId) && (colId == b.firstPixelColId || colId == b.lastPixelColId)
}

func (b *Button) isBorder(rowId, colId int) bool {
	return rowId == b.firstPixelRowId || rowId == b.lastPixelRowId || colId == b.firstPixelColId || colId == b.lastPixelColId
}

func (b *Button) isColored(rowId, colId int) bool {
	return colId > b.secondPixelColId && colId < b.penultimatePixelColId && rowId > b.secondPixelRowId && rowId < b.penultimatePixelRowId
}

func (b *Button) draw() []byte {
	arr := make([]byte, b.pixelRows*b.pixelCols)
	backgroundColor := b.container.GetBackgroundColor()

	for rowId := b.firstPixelRowId; rowId <= b.lastPixelRowId; rowId++ {
		rowNumber := b.pixelCols * rowId

		for colId := b.firstPixelColId; colId <= b.lastPixelColId; colId += 4 {
			if b.isCorner(rowId, colId) {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			} else if b.isBorder(rowId, colId) || b.isColored(rowId, colId) {
				arr[colId+rowNumber] = b.color.R
				arr[colId+1+rowNumber] = b.color.G
				arr[colId+2+rowNumber] = b.color.B
				arr[colId+3+rowNumber] = b.color.A
			} else {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			}
		}
	}

	return arr
}

func (b *Button) drawPressed() []byte {
	arr := make([]byte, b.pixelRows*b.pixelCols)
	backgroundColor := b.container.GetBackgroundColor()

	for rowId := b.firstPixelRowId; rowId <= b.lastPixelRowId; rowId++ {
		rowNumber := b.pixelCols * rowId

		for colId := b.firstPixelColId; colId <= b.lastPixelColId; colId += 4 {
			if b.isCorner(rowId, colId) {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			} else if b.isBorder(rowId, colId) {
				arr[colId+rowNumber] = b.colorPressed.R
				arr[colId+1+rowNumber] = b.colorPressed.G
				arr[colId+2+rowNumber] = b.colorPressed.B
				arr[colId+3+rowNumber] = b.colorPressed.A
			} else {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			}
		}
	}

	return arr
}

func (b *Button) drawHovered() []byte {
	arr := make([]byte, b.pixelRows*b.pixelCols)
	backgroundColor := b.container.GetBackgroundColor()

	for rowId := b.firstPixelRowId; rowId <= b.lastPixelRowId; rowId++ {
		rowNumber := b.pixelCols * rowId

		for colId := b.firstPixelColId; colId <= b.lastPixelColId; colId += 4 {
			if b.isCorner(rowId, colId) {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			} else if b.isBorder(rowId, colId) || b.isColored(rowId, colId) {
				arr[colId+rowNumber] = b.colorHovered.R
				arr[colId+1+rowNumber] = b.colorHovered.G
				arr[colId+2+rowNumber] = b.colorHovered.B
				arr[colId+3+rowNumber] = b.colorHovered.A
			} else {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			}
		}
	}

	return arr
}

func (b *Button) drawDisabled() []byte {
	arr := make([]byte, b.pixelRows*b.pixelCols)
	backgroundColor := b.container.GetBackgroundColor()

	for rowId := b.firstPixelRowId; rowId <= b.lastPixelRowId; rowId++ {
		rowNumber := b.pixelCols * rowId

		for colId := b.firstPixelColId; colId <= b.lastPixelColId; colId += 4 {
			if b.isCorner(rowId, colId) {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			} else if b.isBorder(rowId, colId) || b.isColored(rowId, colId) {
				arr[colId+rowNumber] = b.colorDisabled.R
				arr[colId+1+rowNumber] = b.colorDisabled.G
				arr[colId+2+rowNumber] = b.colorDisabled.B
				arr[colId+3+rowNumber] = b.colorDisabled.A
			} else {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			}
		}
	}

	return arr
}
