package component

import (
	"image"

	"github.com/fglo/chopstiqs/pkg/event"
	"github.com/fglo/chopstiqs/pkg/input"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

// Component is an abstraction of a user interface component, like a button or checkbox.
type Component interface {
	// Draw draws the component to it's image.
	Draw() *ebiten.Image
	// Position returns the component's position.
	Position() (posX float64, posY float64)
	// Size returns the component's size (width and height).
	Size() (width int, height int)
	// SetDisabled sets the component's disabled state.
	SetDisabled(disabled bool)
	// FireEvents fires the component's events.
	FireEvents()
	// SetWidth sets the component's width.
	SetWidth(width int)
	// SetHeight sets the component's height.
	SetHeight(height int)
	// SetDimensions sets the component's dimensions.
	SetDimensions(width, height int)
	// SetPosX sets the component's position X.
	SetPosX(posX float64)
	// SetPosY sets the component's position Y.
	SetPosY(posY float64)
	// SetPosistion sets the component's position (x and y)
	SetPosistion(posX, posY float64)

	setContainer(Container)
}

// component is an abstraction of a user interface component, like a button or checkbox.
type component struct {
	container Container

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

// ComponentOption is a function that configures a component.
type ComponentOption func(c *component)

// ComponentOptions is a struct that holds component options.
type ComponentOptions struct {
	opts []ComponentOption
}

// NewComponent creates a new component.
func NewComponent(width, height int, options *ComponentOptions) *component {
	w := &component{
		image:  ebiten.NewImage(width, height),
		width:  width,
		height: height,
		posX:   0,
		posY:   0,
		Rect:   image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{width, height}},

		pixelCols:             width * 4,
		lastPixelColId:        width*4 - 4,
		penultimatePixelColId: width*4 - 8,

		pixelRows:             height,
		lastPixelRowId:        height - 1,
		penultimatePixelRowId: height - 2,
	}

	if options != nil {
		for _, o := range options.opts {
			o(w)
		}
	}

	return w
}

// setContainer sets the component's container.
func (c *component) setContainer(container Container) {
	c.container = container
}

// SetPosX sets the component's position X.
func (c *component) SetPosX(posX float64) {
	c.posX = posX
	c.Rect = image.Rectangle{Min: image.Point{int(c.posX), int(c.posY)}, Max: image.Point{int(c.posX) + c.width, int(c.posY) + c.height}}
}

// SetPosY sets the component's position Y.
func (c *component) SetPosY(posY float64) {
	c.posY = posY
	c.Rect = image.Rectangle{Min: image.Point{int(c.posX), int(c.posY)}, Max: image.Point{int(c.posX) + c.width, int(c.posY) + c.height}}
}

// SetPosistion sets the component's position (x and y).
func (c *component) SetPosistion(posX, posY float64) {
	c.posX = posX
	c.posY = posY
	c.Rect = image.Rectangle{Min: image.Point{int(c.posX), int(c.posY)}, Max: image.Point{int(c.posX) + c.width, int(c.posY) + c.height}}
}

// SetWidth sets the component's width.
func (c *component) SetWidth(width int) {
	c.width = width
	c.pixelCols = width * 4
	c.lastPixelColId = width*4 - 4
	c.penultimatePixelColId = width*4 - 8

	c.image = ebiten.NewImage(width, c.height)
	c.Rect = image.Rectangle{Min: image.Point{int(c.posX), int(c.posY)}, Max: image.Point{int(c.posX) + c.width, int(c.posY) + c.height}}
}

// SetHeight sets the component's height.
func (c *component) SetHeight(height int) {
	c.height = height
	c.pixelRows = height
	c.lastPixelRowId = height - 1
	c.penultimatePixelRowId = height - 2

	c.image = ebiten.NewImage(c.width, height)
	c.Rect = image.Rectangle{Min: image.Point{int(c.posX), int(c.posY)}, Max: image.Point{int(c.posX) + c.width, int(c.posY) + c.height}}
}

// SetDimensions sets the component's dimensions (width and height).
func (c *component) SetDimensions(width, height int) {
	c.width = width
	c.pixelCols = width * 4
	c.lastPixelColId = width*4 - 4
	c.penultimatePixelColId = width*4 - 8

	c.height = height
	c.pixelRows = height
	c.lastPixelRowId = height - 1
	c.penultimatePixelRowId = height - 2

	c.image = ebiten.NewImage(width, height)
	c.Rect = image.Rectangle{Min: image.Point{int(c.posX), int(c.posY)}, Max: image.Point{int(c.posX) + c.width, int(c.posY) + c.height}}
}

// SetDisabled sets the component's disabled state.
func (c *component) SetDisabled(disabled bool) {
	c.disabled = disabled
}

// Position returns the component's position (x and y).
func (c *component) Position() (float64, float64) {
	return c.posX, c.posY
}

// Size returns the component's size (width and height).
func (c *component) Size() (int, int) {
	return c.width, c.height
}

// FireEvents checks if the mouse cursor is inside the component and fires events accordingly.
func (c *component) FireEvents() {
	p := image.Point{input.CursorPosX, input.CursorPosY}
	mouseEntered := p.In(c.Rect)

	if mouseEntered {
		c.lastUpdateCursorEntered = true

		if input.MouseLeftButtonJustPressed {
			c.lastUpdateMouseLeftButtonPressed = true
			c.MouseButtonPressedEvent.Fire(&ComponentMouseButtonPressedEventArgs{
				Component: c,
				Button:    ebiten.MouseButtonLeft,
			})
		} else {
			c.CursorEnterEvent.Fire(&ComponentCursorEnterEventArgs{
				Component: c,
			})
		}

		if input.MouseRightButtonJustPressed {
			c.lastUpdateMouseRightButtonPressed = true
			c.MouseButtonPressedEvent.Fire(&ComponentMouseButtonPressedEventArgs{
				Component: c,
				Button:    ebiten.MouseButtonRight,
			})
		}
	} else {
		c.lastUpdateCursorEntered = false
		c.CursorExitEvent.Fire(&ComponentCursorExitEventArgs{
			Component: c,
		})
	}

	if !input.MouseLeftButtonPressed && c.lastUpdateMouseLeftButtonPressed {
		c.lastUpdateMouseLeftButtonPressed = false
		c.MouseButtonReleasedEvent.Fire(&ComponentMouseButtonReleasedEventArgs{
			Component: c,
			Inside:    mouseEntered,
			Button:    ebiten.MouseButtonLeft,
		})
	}

	if !input.MouseRightButtonPressed && c.lastUpdateMouseRightButtonPressed {
		c.lastUpdateMouseRightButtonPressed = false
		c.MouseButtonReleasedEvent.Fire(&ComponentMouseButtonReleasedEventArgs{
			Component: c,
			Inside:    mouseEntered,
			Button:    ebiten.MouseButtonRight,
		})
	}
}

// ComponentMouseButtonPressedHandlerFunc is a function that handles mouse button press events.
type ComponentMouseButtonPressedHandlerFunc func(args *ComponentMouseButtonPressedEventArgs) //nolint:golint
// ComponentMouseButtonPressedEventArgs are the arguments for mouse button press events.
type ComponentMouseButtonPressedEventArgs struct { //nolint:golint
	Component *component
	Button    ebiten.MouseButton
}

func (o *ComponentOptions) AddMouseButtonPressedHandler(f ComponentMouseButtonPressedHandlerFunc) *ComponentOptions {
	o.opts = append(o.opts, func(c *component) {
		c.MouseButtonPressedEvent.AddHandler(func(args interface{}) {
			f(args.(*ComponentMouseButtonPressedEventArgs))
		})
	})

	return o
}

// ComponentMouseButtonReleasedHandlerFunc is a function that handles mouse button release events.
type ComponentMouseButtonReleasedHandlerFunc func(args *ComponentMouseButtonReleasedEventArgs) //nolint:golint
// ComponentMouseButtonReleasedEventArgs are the arguments for mouse button release events.
type ComponentMouseButtonReleasedEventArgs struct { //nolint:golint
	Component *component
	Button    ebiten.MouseButton
	Inside    bool
}

func (o *ComponentOptions) AddMouseButtonReleasedHandler(f ComponentMouseButtonReleasedHandlerFunc) *ComponentOptions {
	o.opts = append(o.opts, func(c *component) {
		c.MouseButtonReleasedEvent.AddHandler(func(args interface{}) {
			f(args.(*ComponentMouseButtonReleasedEventArgs))
		})
	})

	return o
}

// ComponentCursorEnterHandlerFunc is a function that handles cursor enter events.
type ComponentCursorEnterHandlerFunc func(args *ComponentCursorEnterEventArgs) //nolint:golint
// ComponentCursorEnterEventArgs are the arguments for cursor enter events.
type ComponentCursorEnterEventArgs struct { //nolint:golint
	Component *component
}

func (o *ComponentOptions) AddCursorEnterHandler(f ComponentCursorEnterHandlerFunc) *ComponentOptions {
	o.opts = append(o.opts, func(c *component) {
		c.CursorEnterEvent.AddHandler(func(args interface{}) {
			f(args.(*ComponentCursorEnterEventArgs))
		})
	})

	return o
}

// ComponentCursorExitHandlerFunc is a function that handles cursor exit events.
type ComponentCursorExitHandlerFunc func(args *ComponentCursorExitEventArgs) //nolint:golint
// ComponentCursorExitEventArgs are the arguments for cursor exit events.
type ComponentCursorExitEventArgs struct { //nolint:golint
	Component *component
}

func (o *ComponentOptions) AddCursorExitHandler(f ComponentCursorExitHandlerFunc) *ComponentOptions {
	o.opts = append(o.opts, func(c *component) {
		c.CursorExitEvent.AddHandler(func(args interface{}) {
			f(args.(*ComponentCursorExitEventArgs))
		})
	})

	return o
}
