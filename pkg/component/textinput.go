package component

import (
	"image"
	"image/color"
	"sync"
	"sync/atomic"
	"time"

	colorutils "github.com/fglo/chopstiqs/internal/color"
	fontutils "github.com/fglo/chopstiqs/internal/font"
	"github.com/fglo/chopstiqs/pkg/event"
	"github.com/fglo/chopstiqs/pkg/input"
	"github.com/fglo/chopstiqs/pkg/option"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type TextInput struct {
	component

	value string

	color   color.RGBA
	font    font.Face
	metrics fontutils.Metrics

	textPosX   int
	textPosY   int
	textBounds image.Rectangle

	scrollOffset int

	cursorPosition      int
	possibleCursorPosXs []int

	ClickedEvent    *event.Event
	KeyPressedEvent *event.Event

	drawer TextInputDrawer

	state textInputState

	lastActionKeyPressed ebiten.Key
	readyForAction       *atomic.Bool
	stateLock            sync.Mutex // TODO: replace this with something that actually works (•_•)

	pressedKeysHandlers map[ebiten.Key]func()
}

type TextInputOptions struct {
	Width  option.OptInt
	Height option.OptInt

	Drawer TextInputDrawer

	Color color.Color
	Font  font.Face

	Padding *Padding
}

type TextInputClickedEventArgs struct {
	TextInput *TextInput
}

type TextInputClickedHandlerFunc func(args *TextInputClickedEventArgs)

type KeyPressedEventArgs struct {
	Key ebiten.Key
}

type KeyPressedHandlerFunc func(args *KeyPressedEventArgs)

func NewTextInput(options *TextInputOptions) *TextInput {
	ti := &TextInput{
		ClickedEvent:    &event.Event{},
		KeyPressedEvent: &event.Event{},

		color:   color.RGBA{230, 230, 230, 255},
		font:    fontutils.DefaultFontFace,
		metrics: fontutils.NewMetrics(fontutils.DefaultFontFace.Metrics()),

		textPosX: 3,

		drawer: &DefaultTextInputDrawer{
			BorderColor: color.RGBA{230, 230, 230, 255},
		},

		lastActionKeyPressed: input.KeyNone,
		readyForAction:       &atomic.Bool{},
	}

	ti.state = ti.idleStateFactory()
	ti.readyForAction.Store(true)

	ti.pressedKeysHandlers = map[ebiten.Key]func(){
		ebiten.KeyLeft:      ti.CursorLeft,
		ebiten.KeyRight:     ti.CursorRight,
		ebiten.KeyBackspace: ti.Backspace,
		ebiten.KeyDelete:    ti.Delete,
		ebiten.KeyEnter:     ti.Submit,
	}

	ti.width = 60
	ti.height = 15

	if options != nil {
		if options.Width.IsSet() {
			ti.width = options.Width.Val()
		}

		if options.Height.IsSet() {
			ti.height = options.Height.Val()
		}

		if options.Color != nil {
			ti.color = colorutils.ColorToRGBA(options.Color)
		}

		if options.Font != nil {
			ti.font = options.Font
			ti.metrics = fontutils.NewMetrics(ti.font.Metrics())
		}

		if options.Drawer != nil {
			ti.drawer = options.Drawer
		}
	}

	ti.setUpComponent(options)

	// ti.KeyPressedEvent.AddHandler(func(args interface{}) {
	// 	keyPressedArgs := args.(*KeyPressedEventArgs)

	// 	switch keyPressedArgs.Key {
	// 	case ebiten.KeyLeft:
	// 		ti.CursorLeft()
	// 	case ebiten.KeyRight:
	// 		ti.CursorRight()
	// 	case ebiten.KeyHome:
	// 		ti.Home()
	// 	case ebiten.KeyEnd:
	// 		ti.End()
	// 	case ebiten.KeyBackspace:
	// 		ti.Backspace()
	// 	case ebiten.KeyDelete:
	// 		ti.Delete()
	// 	default:
	// 		ti.value += ebiten.KeyName(keyPressedArgs.Key)
	// 	}
	// })

	return ti
}

func (ti *TextInput) setUpComponent(options *TextInputOptions) {
	var componentOptions ComponentOptions

	if options != nil {
		componentOptions = ComponentOptions{
			Padding: options.Padding,
		}
	}

	ti.component.setUpComponent(&componentOptions)

	// ti.component.AddCursorEnterHandler(func(args *ComponentCursorEnterEventArgs) {
	// 	if !ti.disabled {
	// 		ti.hovering = true
	// 	}
	// })

	// ti.component.AddCursorExitHandler(func(args *ComponentCursorExitEventArgs) {
	// 	ti.hovering = false
	// })

	ti.component.AddFocusedHandler(func(args *ComponentFocusedEventArgs) {
		ti.focused = args.Focused
		ti.drawer.ResetCursorBlink()
	})

	ti.component.AddMouseButtonReleasedHandler(func(args *ComponentMouseButtonReleasedEventArgs) {
		if !ti.disabled && ti.focused {
			ti.cursorPosition = ti.findClosestPossibleCursorPosition()

			ti.eventManager.Fire(ti.ClickedEvent, &TextInputClickedEventArgs{
				TextInput: ti,
			})
		}
	})
}

func (ti *TextInput) AddClickedHandler(f TextInputClickedHandlerFunc) *TextInput {
	ti.ClickedEvent.AddHandler(func(args interface{}) { f(args.(*TextInputClickedEventArgs)) })

	return ti
}

// SetValue sets the value of the text input.
func (ti *TextInput) SetValue(value string) {
	ti.value = value
	ti.calcTextBounds()
}

func (ti *TextInput) cursorPosX() int {
	return ti.possibleCursorPosXs[ti.cursorPosition] + ti.textPosX + ti.padding.Left - 1
}

func (ti *TextInput) findClosestPossibleCursorPosition() int {
	cursorPos := input.CursorPosX - int(ti.absPosX) - ti.textPosX - ti.padding.Left + 1

	if cursorPos <= ti.possibleCursorPosXs[0] {
		return 0
	}

	lastId := len(ti.possibleCursorPosXs) - 1

	if cursorPos >= ti.possibleCursorPosXs[lastId] {
		return lastId
	}

	min := 0
	max := lastId

	getClosest := func(a, b, target int) int {
		if target-ti.possibleCursorPosXs[a] >= ti.possibleCursorPosXs[b]-target {
			return b
		} else {
			return a
		}
	}

	for min <= max {
		mid := (min + max) / 2
		switch {
		case cursorPos < ti.possibleCursorPosXs[mid]:
			if mid > 0 && cursorPos > ti.possibleCursorPosXs[mid-1] {
				return getClosest(mid-1, mid, cursorPos)
			}

			max = mid - 1
		case cursorPos > ti.possibleCursorPosXs[mid]:
			if mid < lastId && cursorPos < ti.possibleCursorPosXs[mid+1] {
				return getClosest(mid, mid+1, cursorPos)
			}

			min = mid + 1
		default:
			return mid
		}
	}

	return getClosest(max, min, cursorPos)
}

func (ti *TextInput) calcScrollOffset() int {
	scrollOffsetUpperBound := ti.textBounds.Dx() - (ti.width - ti.textPosX - 1 - 2)

	applyBoundsToScrollOffset := func(offset int) int {
		if offset < ti.width/2 {
			return 0
		}

		if offset > scrollOffsetUpperBound {
			return scrollOffsetUpperBound
		}

		return offset
	}

	ti.scrollOffset = applyBoundsToScrollOffset(ti.scrollOffset)

	if ti.cursorPosX() > ti.scrollOffset && ti.cursorPosX() < ti.width+ti.scrollOffset {
		return ti.scrollOffset
	}

	ti.scrollOffset = applyBoundsToScrollOffset(ti.cursorPosX() - ti.width/2)

	return ti.scrollOffset
}

func (ti *TextInput) CursorLeft() {
	if ti.cursorPosition > 0 {
		ti.cursorPosition--
		ti.drawer.ResetCursorBlink()
	}
}

func (ti *TextInput) CursorRight() {
	if ti.cursorPosition < len(ti.possibleCursorPosXs)-1 {
		ti.cursorPosition++
		ti.drawer.ResetCursorBlink()
	}
}

func (ti *TextInput) Home() {
	ti.cursorPosition = 0
	ti.drawer.ResetCursorBlink()
}

func (ti *TextInput) End() {
	ti.cursorPosition = len(ti.possibleCursorPosXs) - 1
	ti.drawer.ResetCursorBlink()
}

func (ti *TextInput) Insert(chars []rune) {
	ti.value = ti.value[0:ti.cursorPosition] + string(chars) + ti.value[ti.cursorPosition:]
	ti.calcTextBounds()
	ti.cursorPosition += len(chars)
	ti.drawer.ResetCursorBlink()
}

func (ti *TextInput) Delete() {
	if ti.cursorPosition < len(ti.value) {
		ti.value = ti.value[0:ti.cursorPosition] + ti.value[ti.cursorPosition+1:]
		ti.calcTextBounds()
	}
}

func (ti *TextInput) Backspace() {
	if ti.cursorPosition > 0 {
		ti.value = ti.value[0:ti.cursorPosition-1] + ti.value[ti.cursorPosition:]
		ti.calcTextBounds()
		ti.CursorLeft()
	}
}

func (ti *TextInput) Submit() {}

func (ti *TextInput) calcTextBounds() {
	// TODO: change deprecated function
	bounds := text.BoundString(ti.font, ti.value) // nolint

	ti.textPosY = ti.metrics.Ascent - ti.metrics.Descent - 1
	ti.textBounds = bounds

	ti.possibleCursorPosXs = make([]int, len(ti.value)+1)
	ti.possibleCursorPosXs[0] = 0

	for i, c := range ti.value {
		ti.possibleCursorPosXs[i+1] = ti.possibleCursorPosXs[i] + fontutils.MeasureString(string(c), ti.font)
	}
}

func (ti *TextInput) Draw() *ebiten.Image {
	ti.state = ti.state(ti)

	if ti.hidden {
		return ti.image
	}

	ti.drawer.Draw(ti)

	if ti.focused && ti.cursorPosition > 0 {
		ti.scrollOffset = ti.calcScrollOffset()
	} else if !ti.focused {
		ti.scrollOffset = 0
	}

	text.Draw(ti.image, ti.value, ti.font, ti.textPosX-ti.scrollOffset+ti.padding.Left, ti.textPosY+ti.padding.Top, ti.color)

	ti.component.Draw()

	return ti.image
}

func (ti *TextInput) actionKeyPressed() (bool, ebiten.Key) {
	actionKeys := []ebiten.Key{
		ebiten.KeyLeft,
		ebiten.KeyRight,
		ebiten.KeyBackspace,
		ebiten.KeyDelete,
		ebiten.KeyEnter,
	}

	for _, key := range actionKeys {
		if ebiten.IsKeyPressed(key) {
			return true, key
		}
	}

	ti.lastActionKeyPressed = input.KeyNone

	return false, input.KeyNone
}

type textInputState func(ti *TextInput) textInputState

var textInputActionKeyRepeatDelay = 350 * time.Millisecond
var textInputActionKeyRepeatInterval = 50 * time.Millisecond

func (ti *TextInput) idleStateFactory() textInputState {
	return func(ti *TextInput) textInputState {
		ti.stateLock.Lock()
		defer ti.stateLock.Unlock()

		if !ti.focused || ti.disabled {
			return ti.idleStateFactory()
		}

		if len(input.InputChars) > 0 {
			return ti.inputStateFactory(input.InputChars)
		}

		if pressed, key := ti.actionKeyPressed(); pressed {
			return ti.actionStateFactory(key)
		}

		return ti.idleStateFactory()
	}
}

func (ti *TextInput) inputStateFactory(chars []rune) textInputState {
	return func(ti *TextInput) textInputState {
		ti.stateLock.Lock()
		defer ti.stateLock.Unlock()

		if !ti.focused || ti.disabled {
			return ti.idleStateFactory()
		}

		ti.Insert(chars)

		if pressed, key := ti.actionKeyPressed(); pressed {
			return ti.actionStateFactory(key)
		}

		return ti.idleStateFactory()
	}
}

func (ti *TextInput) actionStateFactory(pressedKey ebiten.Key) textInputState {
	return func(ti *TextInput) textInputState {
		ti.stateLock.Lock()
		defer ti.stateLock.Unlock()

		if !ti.focused || ti.disabled {
			return ti.idleStateFactory()
		}

		delay := textInputActionKeyRepeatDelay
		if ti.lastActionKeyPressed == pressedKey {
			delay = textInputActionKeyRepeatInterval

			if !ti.readyForAction.Load() {
				return ti.idleStateFactory()
			}
		}

		ti.lastActionKeyPressed = pressedKey

		ti.pressedKeysHandlers[pressedKey]()

		ti.readyForAction.Store(false)

		time.AfterFunc(delay, func() {
			ti.readyForAction.Store(true)
		})

		return ti.idleStateFactory()
	}
}
