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

// TODO: selecting
// TODO: modifier keys
// TODO: input validation
// TODO: on submit
// TODO: on change

type TextInput struct {
	component

	hovering bool

	value string

	color         color.RGBA
	colorDisabled color.RGBA
	colorHovered  color.RGBA
	font          font.Face
	metrics       fontutils.Metrics

	textPosX   int
	textPosY   int
	textBounds image.Rectangle

	scrollOffset int

	cursor              textInputCursor
	cursorPosition      int
	possibleCursorPosXs []int

	ClickedEvent *event.Event
	Changed      *event.Event
	Submitted    *event.Event

	drawer TextInputDrawer

	state textInputState

	lastActionKeyPressed ebiten.Key
	readyForAction       *atomic.Bool
	stateLock            sync.Mutex // TODO: replace this with something that actually works (•_•)

	pressedKeysHandlers map[ebiten.Key]func()

	onSubmitFunc   TextInputOnSubmitFunc
	validationFunc TextInputValidationFunc
}

type TextInputOnSubmitFunc func(string) string
type TextInputValidationFunc func(string) (bool, string)

type TextInputOptions struct {
	Width  option.OptInt
	Height option.OptInt

	Drawer TextInputDrawer

	Color         color.Color
	ColorDisabled color.Color
	ColorHovered  color.Color
	Font          font.Face

	Padding *Padding

	OnSubmitFunc   TextInputOnSubmitFunc
	ValidationFunc TextInputValidationFunc

	CursorOptions *TextInputCursorOptions
}

type TextInputClickedEventArgs struct {
	TextInput *TextInput
}

type TextInputClickedHandlerFunc func(args *TextInputClickedEventArgs)

type TextInputChangedEventArgs struct {
	TextInput *TextInput
	Text      string
}

type TextInputChangedHandlerFunc func(args *TextInputChangedEventArgs)

type TextInputSubmittedEventArgs struct {
	TextInput *TextInput
	Text      string
}

type TextInputSubmittedHandlerFunc func(args *TextInputSubmittedEventArgs)

func NewTextInput(options *TextInputOptions) *TextInput {
	ti := &TextInput{
		ClickedEvent: &event.Event{},
		Changed:      &event.Event{},
		Submitted:    &event.Event{},

		color:         color.RGBA{230, 230, 230, 255},
		colorDisabled: color.RGBA{150, 150, 150, 255},
		colorHovered:  color.RGBA{250, 250, 250, 255},
		font:          fontutils.DefaultFontFace,
		metrics:       fontutils.NewMetrics(fontutils.DefaultFontFace.Metrics()),

		textPosX: 3,

		drawer: &DefaultTextInputDrawer{
			Color:         color.RGBA{230, 230, 230, 255},
			ColorDisabled: color.RGBA{150, 150, 150, 255},
			ColorHovered:  color.RGBA{250, 250, 250, 255},
		},

		lastActionKeyPressed: input.KeyNone,
		readyForAction:       &atomic.Bool{},

		onSubmitFunc:   func(s string) string { return s },
		validationFunc: func(s string) (bool, string) { return true, s },
	}

	ti.state = ti.idleStateFactory()
	ti.readyForAction.Store(true)

	ti.pressedKeysHandlers = map[ebiten.Key]func(){
		ebiten.KeyLeft:      ti.CursorLeft,
		ebiten.KeyRight:     ti.CursorRight,
		ebiten.KeyHome:      ti.Home,
		ebiten.KeyEnd:       ti.End,
		ebiten.KeyBackspace: ti.Backspace,
		ebiten.KeyDelete:    ti.Delete,
		ebiten.KeyEnter:     ti.Submit,
	}

	ti.width = 60
	ti.height = 15

	ti.cursor = *newTextInputCursor(&TextInputCursorOptions{
		Width:  option.Int(1),
		Height: option.Int(ti.height - 4),
	})

	if options != nil {
		if options.Width.IsSet() {
			ti.width = options.Width.Val()
		}

		if options.Height.IsSet() {
			ti.height = options.Height.Val()
			ti.cursor.height = ti.height - 4
		}

		if options.Color != nil {
			ti.color = colorutils.ColorToRGBA(options.Color)
		}

		if options.ColorDisabled != nil {
			ti.colorDisabled = colorutils.ColorToRGBA(options.ColorDisabled)
		}

		if options.Font != nil {
			ti.font = options.Font
			ti.metrics = fontutils.NewMetrics(ti.font.Metrics())
		}

		if options.Drawer != nil {
			ti.drawer = options.Drawer
		}

		if options.OnSubmitFunc != nil {
			ti.onSubmitFunc = options.OnSubmitFunc
		}

		if options.ValidationFunc != nil {
			ti.validationFunc = options.ValidationFunc
		}

		if options.CursorOptions != nil {
			ti.cursor = *newTextInputCursor(options.CursorOptions)
		}
	}

	ti.setUpComponent(options)

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

	ti.component.AddCursorEnterHandler(func(args *ComponentCursorEnterEventArgs) {
		if !ti.disabled {
			ti.hovering = true
		}
	})

	ti.component.AddCursorExitHandler(func(args *ComponentCursorExitEventArgs) {
		ti.hovering = false
	})

	ti.component.AddFocusedHandler(func(args *ComponentFocusedEventArgs) {
		if !ti.disabled {
			ti.focused = args.Focused
			ti.cursor.ResetBlink()

			if !ti.focused {
				ti.Submit()
			}
		}
	})

	ti.component.AddMouseButtonReleasedHandler(func(args *ComponentMouseButtonReleasedEventArgs) {
		if !ti.disabled && ti.focused {
			ti.cursorPosition = ti.findClosestPossibleCursorPosition()
			ti.cursor.ResetBlink()

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

func (ti *TextInput) AddChangedHandler(f TextInputChangedHandlerFunc) *TextInput {
	ti.Changed.AddHandler(func(args interface{}) { f(args.(*TextInputChangedEventArgs)) })

	return ti
}

func (ti *TextInput) AddSubmittedHandler(f TextInputSubmittedHandlerFunc) *TextInput {
	ti.Submitted.AddHandler(func(args interface{}) { f(args.(*TextInputSubmittedEventArgs)) })

	return ti
}

func (ti *TextInput) Value() string {
	return ti.value
}

// SetValue sets the value of the text input.
func (ti *TextInput) SetValue(value string) {
	if valid, valueAfterValidation := ti.validationFunc(value); valid {
		ti.value = valueAfterValidation
		ti.calcTextBounds()
	}
}

func (ti *TextInput) CursorLeft() {
	if ti.cursorPosition > 0 {
		ti.cursorPosition--
		ti.cursor.ResetBlink()
	}
}

func (ti *TextInput) CursorRight() {
	if ti.cursorPosition < len(ti.possibleCursorPosXs)-1 {
		ti.cursorPosition++
		ti.cursor.ResetBlink()
	}
}

func (ti *TextInput) Home() {
	ti.cursorPosition = 0
	ti.cursor.ResetBlink()
}

func (ti *TextInput) End() {
	ti.cursorPosition = len(ti.possibleCursorPosXs) - 1
	ti.cursor.ResetBlink()
}

func (ti *TextInput) Insert(chars []rune) {
	newValue := ti.value[0:ti.cursorPosition] + string(chars) + ti.value[ti.cursorPosition:]

	if valid, valueAfterValidation := ti.validationFunc(newValue); valid {
		ti.value = valueAfterValidation
		ti.calcTextBounds()
		ti.cursorPosition += len(chars)
		ti.cursor.ResetBlink()
		ti.fireChangedEvent()
	}

}

func (ti *TextInput) Delete() {
	if ti.cursorPosition < len(ti.value) {
		ti.value = ti.value[0:ti.cursorPosition] + ti.value[ti.cursorPosition+1:]
		ti.calcTextBounds()
		ti.fireChangedEvent()
	}
}

func (ti *TextInput) Backspace() {
	if ti.cursorPosition > 0 {
		ti.value = ti.value[0:ti.cursorPosition-1] + ti.value[ti.cursorPosition:]
		ti.calcTextBounds()
		ti.CursorLeft()
		ti.fireChangedEvent()
	}
}

func (ti *TextInput) Submit() {
	ti.value = ti.onSubmitFunc(ti.value)
	ti.eventManager.Fire(ti.Submitted, &TextInputSubmittedEventArgs{
		TextInput: ti,
		Text:      ti.value,
	})
}

func (ti *TextInput) fireChangedEvent() {
	ti.eventManager.Fire(ti.Changed, &TextInputChangedEventArgs{
		TextInput: ti,
		Text:      ti.value,
	})
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
	scrollOffsetUpperBound := ti.textBounds.Dx() - (ti.width - ti.textPosX - ti.cursor.width - 2)

	applyBoundsToScrollOffset := func(offset int) int {
		if offset < (ti.width-ti.textPosX)/2 {
			return 0
		}

		if offset > scrollOffsetUpperBound {
			return scrollOffsetUpperBound
		}

		return offset
	}

	ti.scrollOffset = applyBoundsToScrollOffset(ti.scrollOffset)

	cursorPosX := ti.cursorPosX()

	if cursorPosX > ti.scrollOffset && cursorPosX < ti.width+ti.scrollOffset {
		return ti.scrollOffset
	}

	ti.scrollOffset = applyBoundsToScrollOffset(cursorPosX - ti.width/2)

	return ti.scrollOffset
}

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

func (ti *TextInput) actionKeyPressed() (bool, ebiten.Key) {
	actionKeys := []ebiten.Key{
		ebiten.KeyLeft,
		ebiten.KeyRight,
		ebiten.KeyHome,
		ebiten.KeyEnd,
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

func (ti *TextInput) Draw() *ebiten.Image {
	ti.state = ti.state(ti)

	if ti.hidden {
		return ti.image
	}

	ti.drawer.Draw(ti)

	if ti.focused && !ti.disabled {
		ti.scrollOffset = ti.calcScrollOffset()

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(ti.cursorPosX()-ti.scrollOffset), float64(2+ti.padding.Top))
		ti.image.DrawImage(ti.cursor.Draw(), op)
	} else {
		ti.scrollOffset = 0
	}

	switch {
	case ti.disabled:
		text.Draw(ti.image, ti.value, ti.font, ti.textPosX-ti.scrollOffset+ti.padding.Left, ti.textPosY+ti.padding.Top, ti.colorDisabled)
	case ti.hovering:
		text.Draw(ti.image, ti.value, ti.font, ti.textPosX-ti.scrollOffset+ti.padding.Left, ti.textPosY+ti.padding.Top, ti.colorHovered)
	default:
		text.Draw(ti.image, ti.value, ti.font, ti.textPosX-ti.scrollOffset+ti.padding.Left, ti.textPosY+ti.padding.Top, ti.color)
	}

	ti.component.Draw()

	return ti.image
}
