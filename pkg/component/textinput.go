package component

import (
	"image"
	"image/color"
	"regexp"
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

var wordSeparatorRegex *regexp.Regexp

func init() {
	var wordSeparator = `[^a-zA-Z0-9_]`
	wordSeparatorRegex = regexp.MustCompile(wordSeparator)
}

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
	readyForActionRepeat *atomic.Int32
	readyForNewAction    *atomic.Bool

	pressedKeysHandlers map[ebiten.Key]func()
	pressedModifierKeys map[ebiten.Key]bool

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
		readyForActionRepeat: &atomic.Int32{},
		readyForNewAction:    &atomic.Bool{},

		onSubmitFunc:   func(s string) string { return s },
		validationFunc: func(s string) (bool, string) { return true, s },
	}

	ti.state = ti.idleStateFactory()
	ti.readyForActionRepeat.Store(0)
	ti.readyForNewAction.Store(true)

	ti.pressedKeysHandlers = map[ebiten.Key]func(){
		ebiten.KeyLeft:      ti.CursorLeft,
		ebiten.KeyRight:     ti.CursorRight,
		ebiten.KeyHome:      ti.Home,
		ebiten.KeyEnd:       ti.End,
		ebiten.KeyBackspace: ti.Backspace,
		ebiten.KeyDelete:    ti.Delete,
		ebiten.KeyEnter:     ti.Submit,
	}

	ti.pressedModifierKeys = map[ebiten.Key]bool{
		ebiten.KeyControl: false,
		ebiten.KeyAlt:     false,
		ebiten.KeyShift:   false,
		ebiten.KeyMeta:    false,
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

func (ti *TextInput) WordLeft() {
	if ti.cursorPosition <= 0 {
		return
	}
	ti.cursorPosition = ti.findPositionBeforeWord()
	ti.cursor.ResetBlink()
}

func (ti *TextInput) CursorRight() {
	if ti.cursorPosition < len(ti.possibleCursorPosXs)-1 {
		ti.cursorPosition++
		ti.cursor.ResetBlink()
	}
}

func (ti *TextInput) WordRight() {
	if ti.cursorPosition >= len(ti.possibleCursorPosXs)-1 {
		return
	}
	ti.cursorPosition = ti.findPositionAfterWord()
	ti.cursor.ResetBlink()
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

func (ti *TextInput) DeleteWord() {
	if ti.cursorPosition < len(ti.value) {
		spaceToTheRightPosition := ti.findPositionAfterWord()
		ti.value = ti.value[0:ti.cursorPosition] + ti.value[spaceToTheRightPosition:]
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

func (ti *TextInput) BackspaceWord() {
	if ti.cursorPosition > 0 {
		spaceToTheLeftPosition := ti.findPositionBeforeWord()
		ti.value = ti.value[0:spaceToTheLeftPosition] + ti.value[ti.cursorPosition:]
		ti.cursorPosition = spaceToTheLeftPosition + 1
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

func (ti *TextInput) findPositionBeforeWord() int {
	tmpCursorPosition := ti.cursorPosition
	if tmpCursorPosition <= 0 {
		return 0
	}

	for i := ti.cursorPosition - 1; i >= 0; i-- {
		if i >= len(ti.value) {
			continue
		}
		if !wordSeparatorRegex.MatchString(string(ti.value[i])) {
			tmpCursorPosition = i
			break
		}
	}

	for i := tmpCursorPosition; i >= 0; i-- {
		if i >= len(ti.value) {
			continue
		}
		if wordSeparatorRegex.MatchString(string(ti.value[i])) {
			return i + 1
		}
	}

	return 0
}

func (ti *TextInput) findPositionAfterWord() int {
	tmpCursorPosition := ti.cursorPosition
	if tmpCursorPosition >= len(ti.value) {
		return len(ti.value)
	}

	for i := ti.cursorPosition; i < len(ti.possibleCursorPosXs); i++ {
		if !wordSeparatorRegex.MatchString(string(ti.value[i])) {
			tmpCursorPosition = i
			break
		}
	}

	for i := tmpCursorPosition; i < len(ti.possibleCursorPosXs); i++ {
		if i == len(ti.possibleCursorPosXs)-1 || wordSeparatorRegex.MatchString(string(ti.value[i])) {
			return i
		}
	}

	return len(ti.value)
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
	cursorPosX := ti.cursorPosX()
	scrollOffsetLowerBound := 0
	scrollOffsetUpperBound := ti.textBounds.Dx() - (ti.width - ti.textPosX - ti.cursor.width - 2)
	if scrollOffsetUpperBound < 0 {
		scrollOffsetUpperBound = 0
	}

	applyBoundsToScrollOffset := func(offset int) int {
		switch {
		case offset < 0:
			return 0
		// case offset-cursorPosX > 0 && offset < (ti.width-ti.textPosX)/2:
		// 	return 0
		case offset < scrollOffsetLowerBound:
			return scrollOffsetLowerBound
		case offset > scrollOffsetUpperBound:
			return scrollOffsetUpperBound
		default:
			return offset
		}
	}

	ti.scrollOffset = applyBoundsToScrollOffset(ti.scrollOffset)

	if cursorPosX > ti.scrollOffset && cursorPosX < ti.width+ti.scrollOffset-1 {
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

func (ti *TextInput) getKeyPressedHandler(key ebiten.Key) func() {
	switch key {
	case ebiten.KeyLeft:
		if (input.OSMacOS() && ti.pressedModifierKeys[ebiten.KeyMeta]) || ti.pressedModifierKeys[ebiten.KeyControl] {
			return ti.Home
		} else if ti.pressedModifierKeys[ebiten.KeyAlt] {
			return ti.WordLeft
		}
	case ebiten.KeyRight:
		if (input.OSMacOS() && ti.pressedModifierKeys[ebiten.KeyMeta]) || ti.pressedModifierKeys[ebiten.KeyControl] {
			return ti.End
		} else if ti.pressedModifierKeys[ebiten.KeyAlt] {
			return ti.WordRight
		}
	case ebiten.KeyDelete:
		if ti.pressedModifierKeys[ebiten.KeyAlt] {
			return ti.DeleteWord
		}
	case ebiten.KeyBackspace:
		if ti.pressedModifierKeys[ebiten.KeyAlt] {
			return ti.BackspaceWord
		}
	}

	return ti.pressedKeysHandlers[key]
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

	for key := range ti.pressedModifierKeys {
		ti.pressedModifierKeys[key] = input.KeyPressed[key]
	}

	for _, key := range actionKeys {
		if input.KeyPressed[key] {
			return true, key
		}
	}

	ti.lastActionKeyPressed = input.KeyNone

	return false, input.KeyNone
}

type textInputState func(ti *TextInput) textInputState

func (ti *TextInput) idleStateFactory() textInputState {
	return func(ti *TextInput) textInputState {
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

var textInputActionKeyRepeatDelay = 350 * time.Millisecond
var textInputActionKeyRepeatInterval = 35 * time.Millisecond
var textInputDelayBeforeNewAction = 35 * time.Millisecond

func (ti *TextInput) actionStateFactory(pressedKey ebiten.Key) textInputState {
	return func(ti *TextInput) textInputState {
		if !ti.focused || ti.disabled {
			return ti.idleStateFactory()
		}

		delay := textInputActionKeyRepeatDelay
		if ti.lastActionKeyPressed == pressedKey {
			delay = textInputActionKeyRepeatInterval

			if ti.readyForActionRepeat.Load() > 0 {
				return ti.idleStateFactory()
			}
		} else if !ti.readyForNewAction.Load() {
			return ti.idleStateFactory()
		}

		ti.lastActionKeyPressed = pressedKey

		ti.getKeyPressedHandler(pressedKey)()

		ti.readyForActionRepeat.Add(1)
		time.AfterFunc(delay, func() {
			ti.readyForActionRepeat.Add(-1)
		})

		ti.readyForNewAction.Store(false)
		time.AfterFunc(textInputDelayBeforeNewAction, func() {
			ti.readyForNewAction.Store(true)
		})

		return ti.idleStateFactory()
	}
}

func (ti *TextInput) Draw() *ebiten.Image {
	if ti.hidden {
		return ti.image
	}

	if !ti.disabled {
		ti.state = ti.state(ti)
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
