package widget

import (
	"image"

	"github.com/fglo/chopstiqs/pkg/event"
	"github.com/fglo/chopstiqs/pkg/input"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type Widget interface {
	Draw() *ebiten.Image
	Position() (float64, float64)
	Size() (int, int)
	FireEvents(mouse *input.Mouse)
	Disable()
	Enable()
	setContainer(*Container)
	SetWidth(width int)
	SetHeight(height int)
	SetDimensions(width, height int)
}

type widget struct {
	container *Container

	image *ebiten.Image

	Rect image.Rectangle

	disabled bool

	width  int
	height int

	pixelCols int
	pixelRows int

	lastPixelRowId        int
	penultimatePixelRowId int
	lastPixelColId        int
	penultimatePixelColId int

	posX float64
	posY float64

	lastUpdateMouseLeftButtonPressed  bool
	lastUpdateMouseRightButtonPressed bool
	lastUpdateCursorEntered           bool

	MouseButtonPressedEvent  event.Event
	MouseButtonReleasedEvent event.Event
	CursorEnterEvent         event.Event
	CursorExitEvent          event.Event
}

type WidgetOpt func(w *widget)
type WidgetOptions struct {
	opts []WidgetOpt
}

func NewWidget(posX, posY float64, width, height int, options *WidgetOptions) *widget {
	w := &widget{
		image:  ebiten.NewImage(width, height),
		width:  width,
		height: height,
		posX:   posX,
		posY:   posY,
		Rect:   image.Rectangle{Min: image.Point{int(posX), int(posY)}, Max: image.Point{int(posX) + width, int(posY) + height}},

		pixelCols:             width * 4,
		lastPixelColId:        width*4 - 4,
		penultimatePixelColId: width*4 - 8,

		pixelRows:             height,
		lastPixelRowId:        height - 1,
		penultimatePixelRowId: height - 2,
	}

	for _, o := range options.opts {
		o(w)
	}

	return w
}

func (o *WidgetOptions) Disabled() *WidgetOptions {
	o.opts = append(o.opts, func(w *widget) {
		w.disabled = true
	})

	return o
}

func (w *widget) setContainer(container *Container) {
	w.container = container
}

func (w *widget) SetWidth(width int) {
	w.width = width
	w.pixelCols = width * 4
	w.lastPixelColId = width*4 - 4
	w.penultimatePixelColId = width*4 - 8

	w.image = ebiten.NewImage(width, w.height)
	w.Rect = image.Rectangle{Min: image.Point{int(w.posX), int(w.posY)}, Max: image.Point{int(w.posX) + w.width, int(w.posY) + w.height}}
}

func (w *widget) SetHeight(height int) {
	w.height = height
	w.pixelRows = height * 4
	w.lastPixelRowId = height*4 - 4
	w.penultimatePixelRowId = height*4 - 8

	w.image = ebiten.NewImage(w.width, height)
	w.Rect = image.Rectangle{Min: image.Point{int(w.posX), int(w.posY)}, Max: image.Point{int(w.posX) + w.width, int(w.posY) + w.height}}
}

func (w *widget) SetDimensions(width, height int) {
	w.width = width
	w.pixelCols = width * 4
	w.lastPixelColId = width*4 - 4
	w.penultimatePixelColId = width*4 - 8

	w.height = height
	w.pixelRows = height * 4
	w.lastPixelRowId = height*4 - 4
	w.penultimatePixelRowId = height*4 - 8

	w.image = ebiten.NewImage(width, height)
	w.Rect = image.Rectangle{Min: image.Point{int(w.posX), int(w.posY)}, Max: image.Point{int(w.posX) + w.width, int(w.posY) + w.height}}
}

func (w *widget) Disable() {
	w.disabled = true
}

func (w *widget) Enable() {
	w.disabled = false
}

func (w *widget) SetDisabled(disabled bool) {
	w.disabled = disabled
}

func (w *widget) Position() (float64, float64) {
	return w.posX, w.posY
}

func (w *widget) Size() (int, int) {
	return w.width, w.height
}

func (w *widget) FireEvents(mouse *input.Mouse) {
	p := image.Point{mouse.CursorPosX, mouse.CursorPosY}
	mouseEntered := p.In(w.Rect)

	if mouseEntered {
		w.lastUpdateCursorEntered = true

		if mouse.LeftButtonJustPressed {
			w.lastUpdateMouseLeftButtonPressed = true
			w.MouseButtonPressedEvent.Fire(&WidgetMouseButtonPressedEventArgs{
				Widget: w,
				Button: ebiten.MouseButtonLeft,
			})
		} else {
			w.CursorEnterEvent.Fire(&WidgetCursorEnterEventArgs{
				Widget: w,
			})
		}

		if mouse.RightButtonJustPressed {
			w.lastUpdateMouseRightButtonPressed = true
			w.MouseButtonPressedEvent.Fire(&WidgetMouseButtonPressedEventArgs{
				Widget: w,
				Button: ebiten.MouseButtonRight,
			})
		}
	} else {
		w.lastUpdateCursorEntered = false
		w.CursorExitEvent.Fire(&WidgetCursorExitEventArgs{
			Widget: w,
		})
	}

	if !mouse.LeftButtonPressed && w.lastUpdateMouseLeftButtonPressed {
		w.lastUpdateMouseLeftButtonPressed = false
		w.MouseButtonReleasedEvent.Fire(&WidgetMouseButtonReleasedEventArgs{
			Widget: w,
			Inside: mouseEntered,
			Button: ebiten.MouseButtonLeft,
		})
	}

	if !mouse.RightButtonPressed && w.lastUpdateMouseRightButtonPressed {
		w.lastUpdateMouseRightButtonPressed = false
		w.MouseButtonReleasedEvent.Fire(&WidgetMouseButtonReleasedEventArgs{
			Widget: w,
			Inside: mouseEntered,
			Button: ebiten.MouseButtonRight,
		})
	}
}

// WidgetMouseButtonPressedHandlerFunc is a function that handles mouse button press events.
type WidgetMouseButtonPressedHandlerFunc func(args *WidgetMouseButtonPressedEventArgs) //nolint:golint
// WidgetMouseButtonPressedEventArgs are the arguments for mouse button press events.
type WidgetMouseButtonPressedEventArgs struct { //nolint:golint
	Widget *widget
	Button ebiten.MouseButton
}

func (o *WidgetOptions) MouseButtonPressedHandler(f WidgetMouseButtonPressedHandlerFunc) *WidgetOptions {
	o.opts = append(o.opts, func(w *widget) {
		w.MouseButtonPressedEvent.AddHandler(func(args interface{}) {
			f(args.(*WidgetMouseButtonPressedEventArgs))
		})
	})

	return o
}

// WidgetMouseButtonReleasedHandlerFunc is a function that handles mouse button release events.
type WidgetMouseButtonReleasedHandlerFunc func(args *WidgetMouseButtonReleasedEventArgs) //nolint:golint
// WidgetMouseButtonReleasedEventArgs are the arguments for mouse button release events.
type WidgetMouseButtonReleasedEventArgs struct { //nolint:golint
	Widget *widget
	Button ebiten.MouseButton
	Inside bool
}

func (o *WidgetOptions) MouseButtonReleasedHandler(f WidgetMouseButtonReleasedHandlerFunc) *WidgetOptions {
	o.opts = append(o.opts, func(w *widget) {
		w.MouseButtonReleasedEvent.AddHandler(func(args interface{}) {
			f(args.(*WidgetMouseButtonReleasedEventArgs))
		})
	})

	return o
}

// WidgetCursorEnterHandlerFunc is a function that handles cursor enter events.
type WidgetCursorEnterHandlerFunc func(args *WidgetCursorEnterEventArgs) //nolint:golint
// WidgetCursorEnterEventArgs are the arguments for cursor enter events.
type WidgetCursorEnterEventArgs struct { //nolint:golint
	Widget *widget
}

func (o *WidgetOptions) CursorEnterHandler(f WidgetCursorEnterHandlerFunc) *WidgetOptions {
	o.opts = append(o.opts, func(w *widget) {
		w.CursorEnterEvent.AddHandler(func(args interface{}) {
			f(args.(*WidgetCursorEnterEventArgs))
		})
	})

	return o
}

// WidgetCursorExitHandlerFunc is a function that handles cursor exit events.
type WidgetCursorExitHandlerFunc func(args *WidgetCursorExitEventArgs) //nolint:golint
// WidgetCursorExitEventArgs are the arguments for cursor exit events.
type WidgetCursorExitEventArgs struct { //nolint:golint
	Widget *widget
}

func (o *WidgetOptions) CursorExitHandler(f WidgetCursorExitHandlerFunc) *WidgetOptions {
	o.opts = append(o.opts, func(w *widget) {
		w.CursorExitEvent.AddHandler(func(args interface{}) {
			f(args.(*WidgetCursorExitEventArgs))
		})
	})

	return o
}
