package component

import (
	"image"
	"image/color"
	"regexp"

	"github.com/fglo/chopstiqs/pkg/debug"
	"github.com/fglo/chopstiqs/pkg/event"
	"github.com/fglo/chopstiqs/pkg/input"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

var wordSeparatorRegex *regexp.Regexp

func init() {
	var wordSeparator = `[^a-zA-Z0-9_]`
	wordSeparatorRegex = regexp.MustCompile(wordSeparator)
}

// Component is an abstraction of a user interface component, like a button or checkbox.
type Component interface {
	// Draw draws the component to it's image during ebiten.Draw().
	Draw() *ebiten.Image
	// Dimensions returns the component's dimensions (width and height).
	Dimensions() (width int, height int)
	// Width returns the component's width.
	Width() int
	// WidthWithPadding returns the component's width with left and right paddings.
	WidthWithPadding() int
	// Height returns the component's height.
	Height() int
	// HeightWithPadding returns the component's height with top and bottom paddings.
	HeightWithPadding() int
	// Position returns the component's position.
	Position() (posX float64, posY float64)
	// PosX returns the component's position X.
	PosX() float64
	// PosY returns the component's position Y.
	PosY() float64
	// AbsPosition return the component's absolute position.
	AbsPosition() (posX float64, posY float64)
	// AbsPosX returns the component's absolute position X.
	AbsPosX() float64
	// AbsPosY returns the component's absolute position Y.
	AbsPosY() float64
	// Disable returns the component's disabled state.
	Disable() bool
	// SetDisabled sets the component's disabled state.
	SetDisabled(disabled bool)
	// Hidden returns the component's hidden state.
	Hidden() bool
	// SetHidden sets the component's hidden state.
	SetHidden(hidden bool)
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
	// SetPosision sets the component's position (x and y)
	SetPosision(posX, posY float64)
	// SetPadding sets the component's padding.
	SetPadding(padding Padding)
	// SetPaddingTop sets the component's padding top.
	SetPaddingTop(padding int)
	// SetPaddingBottom sets the component's padding bottom.
	SetPaddingBottom(padding int)
	// SetPaddingLeft sets the component's padding left.
	SetPaddingLeft(padding int)
	// SetPaddingRight sets the component's padding right.
	SetPaddingRight(padding int)

	Focused() bool
	SetFocused(bool)

	setContainer(*Container)

	EventManager() *event.Manager
	SetEventManager(*event.Manager)

	AddFocusedHandler(f ComponentFocusedHandlerFunc) Component
}

// component is an abstraction of a user interface component, like a button or checkbox.
type component struct {
	container *Container

	eventManager *event.Manager

	image *ebiten.Image

	rect image.Rectangle

	disabled bool
	hidden   bool

	width             int
	widthWithPadding  int
	height            int
	heightWithPadding int

	padding Padding

	pixelCols int
	pixelRows int

	firstPixelRowId       int
	secondPixelRowId      int
	lastPixelRowId        int
	penultimatePixelRowId int

	firstPixelColId       int
	secondPixelColId      int
	lastPixelColId        int
	penultimatePixelColId int

	absPosX float64
	absPosY float64

	posX float64
	posY float64

	focused bool

	lastUpdateMouseLeftButtonPressed  bool
	lastUpdateMouseRightButtonPressed bool
	lastUpdateCursorEntered           bool

	MouseButtonPressedEvent  *event.Event
	MouseButtonReleasedEvent *event.Event
	CursorEnterEvent         *event.Event
	CursorExitEvent          *event.Event
	FocusedEvent             *event.Event
}

// ComponentOptions is a struct that holds component options.
type ComponentOptions struct {
	Padding  *Padding
	Disabled bool
	Hidden   bool
}

// SetupComponent sets up the component.
func (c *component) setUpComponent(opt *ComponentOptions) {
	c.MouseButtonPressedEvent = &event.Event{}
	c.MouseButtonReleasedEvent = &event.Event{}
	c.CursorEnterEvent = &event.Event{}
	c.CursorExitEvent = &event.Event{}
	c.FocusedEvent = &event.Event{}

	c.padding = DefaultPadding

	if opt != nil {
		if opt.Padding != nil {
			opt.Padding.Validate()
			c.padding = *opt.Padding
		}

		c.disabled = opt.Disabled
		c.hidden = opt.Hidden
	}

	c.SetDimensions(c.width, c.height)
}

// setContainer sets the component's container.
func (c *component) setContainer(container *Container) {
	c.container = container
	c.absPosX = c.posX + c.container.AbsPosX()
	c.absPosY = c.posY + c.container.AbsPosY()
	c.setRect()
	c.eventManager = container.EventManager()
}

func (c *component) EventManager() *event.Manager {
	return c.eventManager
}

func (c *component) SetEventManager(eventManager *event.Manager) {
	c.eventManager = eventManager
}

func (c *component) drawBorders(arr []byte) []byte {
	borderColor := color.RGBA{249, 192, 46, 255}

	firstRowNumber := c.pixelCols * c.padding.Top
	lastRowNumber := c.pixelCols * (c.pixelRows - c.padding.Bottom - 1)
	for colId := c.padding.Left * 4; colId < c.pixelCols-c.padding.Right*4; colId += 4 {
		arr[colId+firstRowNumber] = borderColor.R
		arr[colId+1+firstRowNumber] = borderColor.G
		arr[colId+2+firstRowNumber] = borderColor.B
		arr[colId+3+firstRowNumber] = borderColor.A

		arr[colId+lastRowNumber] = borderColor.R
		arr[colId+1+lastRowNumber] = borderColor.G
		arr[colId+2+lastRowNumber] = borderColor.B
		arr[colId+3+lastRowNumber] = borderColor.A
	}

	firstColumnNumber := c.padding.Left * 4
	lastColumnNumber := c.pixelCols - 4 - c.padding.Right*4
	for rowId := c.padding.Top; rowId < c.pixelRows-c.padding.Bottom; rowId++ {
		rowNumber := c.pixelCols * rowId

		arr[firstColumnNumber+rowNumber] = borderColor.R
		arr[firstColumnNumber+1+rowNumber] = borderColor.G
		arr[firstColumnNumber+2+rowNumber] = borderColor.B
		arr[firstColumnNumber+3+rowNumber] = borderColor.A

		arr[lastColumnNumber+rowNumber] = borderColor.R
		arr[lastColumnNumber+1+rowNumber] = borderColor.G
		arr[lastColumnNumber+2+rowNumber] = borderColor.B
		arr[lastColumnNumber+3+rowNumber] = borderColor.A
	}

	return arr
}

func (c *component) drawPadding(arr []byte) []byte {
	paddingBorderColor := color.RGBA{255, 100, 100, 255}

	lastRowNumber := c.pixelCols * (c.pixelRows - 1)
	for colId := 0; colId < c.pixelCols; colId += 4 {
		arr[colId] = paddingBorderColor.R
		arr[colId+1] = paddingBorderColor.G
		arr[colId+2] = paddingBorderColor.B
		arr[colId+3] = paddingBorderColor.A

		arr[colId+lastRowNumber] = paddingBorderColor.R
		arr[colId+1+lastRowNumber] = paddingBorderColor.G
		arr[colId+2+lastRowNumber] = paddingBorderColor.B
		arr[colId+3+lastRowNumber] = paddingBorderColor.A
	}

	for rowId := 0; rowId < c.pixelRows; rowId++ {
		rowNumber := c.pixelCols * rowId

		arr[rowNumber] = paddingBorderColor.R
		arr[1+rowNumber] = paddingBorderColor.G
		arr[2+rowNumber] = paddingBorderColor.B
		arr[3+rowNumber] = paddingBorderColor.A

		arr[c.pixelCols-4+rowNumber] = paddingBorderColor.R
		arr[c.pixelCols-4+1+rowNumber] = paddingBorderColor.G
		arr[c.pixelCols-4+2+rowNumber] = paddingBorderColor.B
		arr[c.pixelCols-4+3+rowNumber] = paddingBorderColor.A
	}

	return arr
}

func (c *component) Draw() *ebiten.Image {
	if debug.Debug {
		debugImage := ebiten.NewImage(c.widthWithPadding, c.heightWithPadding)

		arr := make([]byte, c.pixelRows*c.pixelCols)

		if debug.ShowComponentBorders {
			arr = c.drawBorders(arr)
		}

		if debug.ShowComponentPadding {
			arr = c.drawPadding(arr)
		}

		debugImage.WritePixels(arr)

		c.image.DrawImage(debugImage, &ebiten.DrawImageOptions{})
	}

	return c.image
}

// SetPosX sets the component's position X.
func (c *component) SetPosX(posX float64) {
	c.posX = posX
	if c.container != nil {
		c.absPosX = posX + c.container.AbsPosX()
	} else {
		c.absPosX = posX
	}

	c.setRect()
}

// SetPosY sets the component's position Y.
func (c *component) SetPosY(posY float64) {
	c.posY = posY
	if c.container != nil {
		c.absPosY = posY + c.container.AbsPosY()
	} else {
		c.absPosY = posY
	}

	c.setRect()
}

// SetPosision sets the component's position (x and y).
func (c *component) SetPosision(posX, posY float64) {
	c.posX = posX
	if c.container != nil {
		c.absPosX = posX + c.container.AbsPosX()
	} else {
		c.absPosX = posX
	}

	c.posY = posY
	if c.container != nil {
		c.absPosY = posY + c.container.AbsPosY()
	} else {
		c.absPosY = posY
	}

	c.setRect()
}

// SetPadding sets the component's padding.
func (c *component) SetPadding(padding Padding) {
	padding.Validate()
	c.padding = padding
	c.SetDimensions(c.width, c.height)
}

// SetPaddingTop sets the component's padding top.
func (c *component) SetPaddingTop(padding int) {
	c.padding.Top = padding
	c.SetDimensions(c.width, c.height)
}

// SetPaddingBottom sets the component's padding bottom.
func (c *component) SetPaddingBottom(padding int) {
	c.padding.Bottom = padding
	c.SetDimensions(c.width, c.height)
}

// SetPaddingLeft sets the component's padding left.
func (c *component) SetPaddingLeft(padding int) {
	c.padding.Left = padding
	c.SetDimensions(c.width, c.height)
}

// SetPaddingRight sets the component's padding right.
func (c *component) SetPaddingRight(padding int) {
	c.padding.Right = padding
	c.SetDimensions(c.width, c.height)
}

// SetWidth sets the component's width.
func (c *component) SetWidth(width int) {
	c.width = width
	c.widthWithPadding = width + c.padding.Left + c.padding.Right
	c.pixelCols = c.widthWithPadding * 4

	c.calcPixelColIds()

	c.setImage()
	c.setRect()
}

// SetHeight sets the component's height.
func (c *component) SetHeight(height int) {
	c.height = height
	c.heightWithPadding = height + c.padding.Top + c.padding.Bottom
	c.pixelRows = c.heightWithPadding

	c.calcPixelRowIds()

	c.setImage()
	c.setRect()
}

// SetDimensions sets the component's dimensions (width and height).
func (c *component) SetDimensions(width, height int) {
	if width != 0 && height != 0 {
		c.width = width
		c.widthWithPadding = width + c.padding.Left + c.padding.Right
		c.pixelCols = c.widthWithPadding * 4

		c.height = height
		c.heightWithPadding = height + c.padding.Top + c.padding.Bottom
		c.pixelRows = c.heightWithPadding

		c.calcPixelColIds()
		c.calcPixelRowIds()

		c.setImage()
		c.setRect()
	}
}

func (c *component) Focused() bool {
	return c.focused
}

func (c *component) SetFocused(focused bool) {
	if c.focused != focused {
		c.focused = focused
		c.eventManager.Fire(c.FocusedEvent, &ComponentFocusedEventArgs{
			Component: c,
			Focused:   focused,
		})
	}
}

func (c *component) calcPixelColIds() {
	c.firstPixelColId = c.padding.Left * 4
	c.secondPixelColId = c.firstPixelColId + 4

	c.lastPixelColId = c.pixelCols - c.padding.Right*4 - 4
	c.penultimatePixelColId = c.lastPixelColId - 4
}

func (c *component) calcPixelRowIds() {
	c.firstPixelRowId = c.padding.Top
	c.secondPixelRowId = c.firstPixelRowId + 1

	c.lastPixelRowId = c.pixelRows - c.padding.Bottom - 1
	c.penultimatePixelRowId = c.lastPixelRowId - 1
}

func (c *component) setImage() {
	c.image = ebiten.NewImage(c.widthWithPadding, c.heightWithPadding)
}

func (c *component) setRect() {
	c.rect = image.Rectangle{Min: image.Point{int(c.absPosX), int(c.absPosY)}, Max: image.Point{int(c.absPosX) + c.widthWithPadding, int(c.absPosY) + c.heightWithPadding}}
}

// Disable returns the component's disabled state.
func (c *component) Disable() bool {
	return c.disabled
}

// SetDisabled sets the component's disabled state.
func (c *component) SetDisabled(disabled bool) {
	c.disabled = disabled
}

// Hidden returns the component's hidden state.
func (c *component) Hidden() bool {
	return c.hidden
}

// SetHidden sets the component's hidden state.
func (c *component) SetHidden(hidden bool) {
	c.hidden = hidden
}

// Position returns the component's position (x and y).
func (c *component) Position() (float64, float64) {
	return c.posX, c.posY
}

// PosX returns the component's position X.
func (c *component) PosX() float64 {
	return c.posX
}

// PosY returns the component's position Y.
func (c *component) PosY() float64 {
	return c.posY
}

// AbsPosition return the component's absolute position.
func (c *component) AbsPosition() (posX float64, posY float64) {
	return c.absPosX, c.absPosY
}

// AbsPosX returns the component's absolute position X.
func (c *component) AbsPosX() float64 {
	return c.absPosX
}

// AbsPosY returns the component's absolute position Y.
func (c *component) AbsPosY() float64 {
	return c.absPosY
}

// Dimensions returns the component's size (width and height).
func (c *component) Dimensions() (int, int) {
	return c.width, c.height
}

// Width returns the component's width.
func (c *component) Width() int {
	return c.width
}

// WidthWithPadding returns the component's width with left and right paddings.
func (c *component) WidthWithPadding() int {
	return c.widthWithPadding
}

// Height returns the component's height.
func (c *component) Height() int {
	return c.height
}

// HeightWithPadding returns the component's height with top and bottom paddings.
func (c *component) HeightWithPadding() int {
	return c.heightWithPadding
}

// FireEvents checks if the mouse cursor is inside the component and fires events accordingly.
func (c *component) FireEvents() {
	p := image.Point{input.CursorPosX, input.CursorPosY}
	mouseEntered := p.In(c.rect)

	if mouseEntered {
		c.lastUpdateCursorEntered = true

		if !input.MouseLeftButtonPressed && !input.MouseRightButtonPressed {
			c.eventManager.Fire(c.CursorEnterEvent, &ComponentCursorEnterEventArgs{
				Component: c,
			})
		}

		if input.MouseLeftButtonJustPressed {
			c.lastUpdateMouseLeftButtonPressed = true
			c.SetFocused(true)
		}

		if input.MouseLeftButtonPressed {
			if c.focused {
				c.eventManager.Fire(c.MouseButtonPressedEvent, &ComponentMouseButtonPressedEventArgs{
					Component: c,
					Button:    ebiten.MouseButtonLeft,
				})
			}
		}

		if input.MouseRightButtonJustPressed {
			c.lastUpdateMouseRightButtonPressed = true
			c.SetFocused(true)
		}

		if input.MouseRightButtonPressed {
			if c.focused {
				c.eventManager.Fire(c.MouseButtonPressedEvent, &ComponentMouseButtonPressedEventArgs{
					Component: c,
					Button:    ebiten.MouseButtonRight,
					Inside:    mouseEntered,
				})
			}
		}
	} else {
		c.lastUpdateCursorEntered = false

		c.eventManager.Fire(c.CursorExitEvent, &ComponentCursorExitEventArgs{
			Component: c,
		})

		if input.MouseLeftButtonPressed && c.lastUpdateMouseLeftButtonPressed {
			c.eventManager.Fire(c.MouseButtonPressedEvent, &ComponentMouseButtonPressedEventArgs{
				Component: c,
				Button:    ebiten.MouseButtonLeft,
				Inside:    mouseEntered,
			})
		}

		if input.MouseLeftButtonJustPressed && !c.lastUpdateMouseLeftButtonPressed || input.MouseRightButtonJustPressed && !c.lastUpdateMouseRightButtonPressed {
			c.SetFocused(false)
		}
	}

	if !input.MouseLeftButtonPressed && c.lastUpdateMouseLeftButtonPressed {
		c.lastUpdateMouseLeftButtonPressed = false
		c.eventManager.Fire(c.MouseButtonReleasedEvent, &ComponentMouseButtonReleasedEventArgs{
			Component: c,
			Inside:    mouseEntered,
			Button:    ebiten.MouseButtonLeft,
		})
	}

	if !input.MouseRightButtonPressed && c.lastUpdateMouseRightButtonPressed {
		c.lastUpdateMouseRightButtonPressed = false
		c.eventManager.Fire(c.MouseButtonReleasedEvent, &ComponentMouseButtonReleasedEventArgs{
			Component: c,
			Inside:    mouseEntered,
			Button:    ebiten.MouseButtonRight,
		})
	}
}

// ComponentMouseButtonJustPressedHandlerFunc is a function that handles mouse button press events.
type ComponentMouseButtonJustPressedHandlerFunc func(args *ComponentMouseButtonJustPressedEventArgs) //nolint:golint
// ComponentMouseButtonPressedEventArgs are the arguments for mouse button press events.
type ComponentMouseButtonJustPressedEventArgs struct { //nolint:golint
	Component Component
	Button    ebiten.MouseButton
}

func (c *component) AddMouseButtonJustPressedHandler(f ComponentMouseButtonJustPressedHandlerFunc) Component {
	c.MouseButtonPressedEvent.AddHandler(func(args interface{}) {
		f(args.(*ComponentMouseButtonJustPressedEventArgs))
	})

	return c
}

// ComponentMouseButtonPressedHandlerFunc is a function that handles mouse button press events.
type ComponentMouseButtonPressedHandlerFunc func(args *ComponentMouseButtonPressedEventArgs) //nolint:golint
// ComponentMouseButtonPressedEventArgs are the arguments for mouse button press events.
type ComponentMouseButtonPressedEventArgs struct { //nolint:golint
	Component Component
	Button    ebiten.MouseButton
	Inside    bool
}

func (c *component) AddMouseButtonPressedHandler(f ComponentMouseButtonPressedHandlerFunc) Component {
	c.MouseButtonPressedEvent.AddHandler(func(args interface{}) {
		f(args.(*ComponentMouseButtonPressedEventArgs))
	})

	return c
}

// ComponentMouseButtonReleasedHandlerFunc is a function that handles mouse button release events.
type ComponentMouseButtonReleasedHandlerFunc func(args *ComponentMouseButtonReleasedEventArgs) //nolint:golint
// ComponentMouseButtonReleasedEventArgs are the arguments for mouse button release events.
type ComponentMouseButtonReleasedEventArgs struct { //nolint:golint
	Component Component
	Button    ebiten.MouseButton
	Inside    bool
}

func (c *component) AddMouseButtonReleasedHandler(f ComponentMouseButtonReleasedHandlerFunc) Component {
	c.MouseButtonReleasedEvent.AddHandler(func(args interface{}) {
		f(args.(*ComponentMouseButtonReleasedEventArgs))
	})

	return c
}

// ComponentCursorEnterHandlerFunc is a function that handles cursor enter events.
type ComponentCursorEnterHandlerFunc func(args *ComponentCursorEnterEventArgs) //nolint:golint
// ComponentCursorEnterEventArgs are the arguments for cursor enter events.
type ComponentCursorEnterEventArgs struct { //nolint:golint
	Component Component
}

func (c *component) AddCursorEnterHandler(f ComponentCursorEnterHandlerFunc) Component {
	c.CursorEnterEvent.AddHandler(func(args interface{}) {
		f(args.(*ComponentCursorEnterEventArgs))
	})

	return c
}

// ComponentCursorExitHandlerFunc is a function that handles cursor exit events.
type ComponentCursorExitHandlerFunc func(args *ComponentCursorExitEventArgs) //nolint:golint
// ComponentCursorExitEventArgs are the arguments for cursor exit events.
type ComponentCursorExitEventArgs struct { //nolint:golint
	Component Component
}

func (c *component) AddCursorExitHandler(f ComponentCursorExitHandlerFunc) Component {
	c.CursorExitEvent.AddHandler(func(args interface{}) {
		f(args.(*ComponentCursorExitEventArgs))
	})

	return c
}

type ComponentFocusedHandlerFunc func(args *ComponentFocusedEventArgs) //nolint:golint
type ComponentFocusedEventArgs struct {
	Component Component
	Focused   bool
}

func (c *component) AddFocusedHandler(f ComponentFocusedHandlerFunc) Component {
	c.FocusedEvent.AddHandler(func(args interface{}) {
		f(args.(*ComponentFocusedEventArgs))
	})

	return c
}
