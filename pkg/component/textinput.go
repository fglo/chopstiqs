package component

import (
	"image/color"
	"math"
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

// textInputCursorPosition is a type indicating that value is one of the possible cursor positions, not coordinate in the X axis
type textInputCursorPosition int

// textInputAction is a type refering to a action triggered by a action key with (or without) modifier keys
type textInputAction int

const (
	textInputIdle textInputAction = iota
	textInputCursorLeft
	textInputWordLeft
	textInputCursorRight
	textInputWordRight
	textInputHome
	textInputEnd
	textInputDelete
	textInputDeleteWord
	textInputBackspace
	textInputBackspaceWord
	textInputRemoveSelection
	textInputSubmit
)

type TextInput struct {
	component

	hovering bool

	value string

	color         color.RGBA
	colorDisabled color.RGBA
	colorHovered  color.RGBA
	font          font.Face
	metrics       fontutils.Metrics

	textPosX int
	textPosY int

	scrollOffset int

	cursor              textInputCursor
	cursorPosition      textInputCursorPosition
	possibleCursorPosXs []int

	pressed          bool
	pressedPosition  textInputCursorPosition
	releasedPosition textInputCursorPosition

	selecting      bool
	selectingFrom  textInputCursorPosition
	selectionStart textInputCursorPosition
	selectionEnd   textInputCursorPosition

	ClickedEvent   *event.Event
	PressedEvent   *event.Event
	ReleasedEvent  *event.Event
	ChangedEvent   *event.Event
	SubmittedEvent *event.Event

	submitOnUnfocus bool

	drawer TextInputDrawer

	state textInputState

	lastActionKeyPressed ebiten.Key
	readyForActionRepeat *atomic.Int32
	readyForNewAction    *atomic.Bool

	actionKeyHandlers   map[ebiten.Key]func() textInputAction
	actionHandlers      map[textInputAction]func()
	modifierKeysPressed map[ebiten.Key]bool

	onSubmitFunc        TextInputOnSubmitFunc
	inputValidationFunc TextInputValidationFunc
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

	OnSubmitFunc        TextInputOnSubmitFunc
	InputValidationFunc TextInputValidationFunc

	SubmitOnUnfocus bool

	CursorOptions *TextInputCursorOptions
}

type TextInputClickedEventArgs struct {
	TextInput *TextInput
}

type TextInputClickedHandlerFunc func(args *TextInputClickedEventArgs)

type TextInputPressedEventArgs struct {
	TextInput *TextInput
}

type TextInputPressedHandlerFunc func(args *TextInputClickedEventArgs)

type TextInputReleasedEventArgs struct {
	TextInput *TextInput
	Inside    bool
}

type TextInputReleasedHandlerFunc func(args *TextInputReleasedEventArgs)

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
		ClickedEvent:   &event.Event{},
		PressedEvent:   &event.Event{},
		ReleasedEvent:  &event.Event{},
		ChangedEvent:   &event.Event{},
		SubmittedEvent: &event.Event{},

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

		onSubmitFunc:        func(s string) string { return s },
		inputValidationFunc: func(s string) (bool, string) { return true, s },

		selectingFrom: -1,
	}

	ti.state = ti.idleStateFactory()
	ti.readyForActionRepeat.Store(0)
	ti.readyForNewAction.Store(true)

	ti.actionKeyHandlers = map[ebiten.Key]func() textInputAction{
		ebiten.KeyLeft:      ti.handleKeyLeft,
		ebiten.KeyRight:     ti.handleKeyRight,
		ebiten.KeyHome:      ti.handleKeyHome,
		ebiten.KeyEnd:       ti.handleKeyEnd,
		ebiten.KeyBackspace: ti.handleKeyBackspace,
		ebiten.KeyDelete:    ti.handleKeyDelete,
		ebiten.KeyEnter:     ti.handleKeyEnter,
	}

	ti.actionHandlers = map[textInputAction]func(){
		textInputCursorLeft:      ti.CursorLeft,
		textInputWordLeft:        ti.WordLeft,
		textInputCursorRight:     ti.CursorRight,
		textInputWordRight:       ti.WordRight,
		textInputHome:            ti.Home,
		textInputEnd:             ti.End,
		textInputBackspace:       ti.Backspace,
		textInputBackspaceWord:   ti.BackspaceWord,
		textInputDelete:          ti.Delete,
		textInputDeleteWord:      ti.DeleteWord,
		textInputRemoveSelection: ti.RemoveSelection,
		textInputSubmit:          ti.Submit,
	}

	ti.modifierKeysPressed = map[ebiten.Key]bool{
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
		ti.submitOnUnfocus = options.SubmitOnUnfocus

		if options.Width.IsSet() {
			ti.width = options.Width.Val()
		}

		if options.Height.IsSet() {
			ti.height = options.Height.Val()
			ti.cursor.height = ti.height - 4
		}

		if options.Color != nil {
			ti.color = colorutils.ToRGBA(options.Color)
		}

		if options.ColorDisabled != nil {
			ti.colorDisabled = colorutils.ToRGBA(options.ColorDisabled)
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

		if options.InputValidationFunc != nil {
			ti.inputValidationFunc = options.InputValidationFunc
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
			ti.cursor.ResetBlink()

			if !args.Focused {
				ti.Deselect()

				if ti.submitOnUnfocus {
					ti.Submit()
				}
			}
		}
	})

	ti.component.AddMouseButtonPressedHandler(func(args *ComponentMouseButtonPressedEventArgs) {
		if ti.disabled || args.Button != ebiten.MouseButtonLeft {
			return
		}

		ti.checkForShift()

		ti.moveCursor(ti.findClosestPossibleCursorPosition())

		if !ti.pressed {
			ti.pressedPosition = ti.cursorPosition
			ti.releasedPosition = -1
		}

		ti.pressed = true

		if !ti.selecting {
			ti.Deselect()
			ti.selectingFrom = ti.pressedPosition
		}

		ti.eventManager.Fire(ti.PressedEvent, &TextInputPressedEventArgs{
			TextInput: ti,
		})
	})

	ti.component.AddMouseButtonReleasedHandler(func(args *ComponentMouseButtonReleasedEventArgs) {
		if !ti.pressed || args.Button != ebiten.MouseButtonLeft {
			return
		}

		ti.moveCursor(ti.findClosestPossibleCursorPosition())
		ti.pressed = false
		ti.releasedPosition = ti.cursorPosition

		ti.eventManager.Fire(ti.ReleasedEvent, &TextInputReleasedEventArgs{
			TextInput: ti,
			Inside:    args.Inside,
		})

		if !ti.disabled && ti.pressedPosition == ti.releasedPosition {
			ti.eventManager.Fire(ti.ClickedEvent, &TextInputClickedEventArgs{
				TextInput: ti,
			})
		}

		ti.pressedPosition = -1
	})
}

func (ti *TextInput) AddClickedHandler(f TextInputClickedHandlerFunc) *TextInput {
	ti.ClickedEvent.AddHandler(func(args interface{}) { f(args.(*TextInputClickedEventArgs)) })

	return ti
}

func (ti *TextInput) AddChangedHandler(f TextInputChangedHandlerFunc) *TextInput {
	ti.ChangedEvent.AddHandler(func(args interface{}) { f(args.(*TextInputChangedEventArgs)) })

	return ti
}

func (ti *TextInput) AddSubmittedHandler(f TextInputSubmittedHandlerFunc) *TextInput {
	ti.SubmittedEvent.AddHandler(func(args interface{}) { f(args.(*TextInputSubmittedEventArgs)) })

	return ti
}

func (ti *TextInput) Value() string {
	return ti.value
}

// SetValue sets the value of the text input.
func (ti *TextInput) SetValue(value string) {
	if valid, valueAfterValidation := ti.inputValidationFunc(value); valid {
		ti.setValue(valueAfterValidation)
	}
}

func (ti *TextInput) HasSelectedText() bool {
	return ti.selectionStart != -1 && ti.selectionEnd != -1 && ti.selectionStart != ti.selectionEnd
}

func (ti *TextInput) Deselect() {
	ti.selecting = false
	ti.selectingFrom = -1
}

func (ti *TextInput) CursorLeft() {
	if ti.cursorPosition > 0 {
		ti.moveCursor(ti.cursorPosition - 1)
	}
}

func (ti *TextInput) WordLeft() {
	if ti.cursorPosition <= 0 {
		return
	}
	ti.moveCursor(ti.findPositionBeforeWord())
}

func (ti *TextInput) CursorRight() {
	if int(ti.cursorPosition) < len(ti.possibleCursorPosXs)-1 {
		ti.moveCursor(ti.cursorPosition + 1)
	}
}

func (ti *TextInput) WordRight() {
	if int(ti.cursorPosition) >= len(ti.possibleCursorPosXs)-1 {
		return
	}
	ti.moveCursor(ti.findPositionAfterWord())
}

func (ti *TextInput) Home() {
	ti.moveCursor(0)
}

func (ti *TextInput) End() {
	ti.moveCursor(textInputCursorPosition(len(ti.possibleCursorPosXs) - 1))
}

func (ti *TextInput) Insert(chars []rune) {
	newValue := ti.value[0:ti.cursorPosition] + string(chars) + ti.value[ti.cursorPosition:]

	if valid, valueAfterValidation := ti.inputValidationFunc(newValue); valid {
		ti.setValue(valueAfterValidation)
		ti.moveCursor(ti.cursorPosition + textInputCursorPosition(len(chars)))
		ti.fireChangedEvent()
	}

}

func (ti *TextInput) Delete() {
	if ti.cursorPosition < textInputCursorPosition(len(ti.value)) {
		ti.setValue(ti.value[0:ti.cursorPosition] + ti.value[ti.cursorPosition+1:])
		ti.fireChangedEvent()
	}
}

func (ti *TextInput) DeleteWord() {
	if ti.cursorPosition < textInputCursorPosition(len(ti.value)) {
		spaceToTheRightPosition := ti.findPositionAfterWord()
		ti.setValue(ti.value[0:ti.cursorPosition] + ti.value[spaceToTheRightPosition:])
		ti.fireChangedEvent()
	}
}

func (ti *TextInput) Backspace() {
	if ti.cursorPosition > 0 {
		ti.setValue(ti.value[0:ti.cursorPosition-1] + ti.value[ti.cursorPosition:])
		ti.fireChangedEvent()
		ti.CursorLeft()
	}
}

func (ti *TextInput) BackspaceWord() {
	if ti.cursorPosition > 0 {
		spaceToTheLeftPosition := ti.findPositionBeforeWord()
		ti.setValue(ti.value[0:spaceToTheLeftPosition] + ti.value[ti.cursorPosition:])
		ti.fireChangedEvent()
		ti.moveCursor(spaceToTheLeftPosition)
	}
}

func (ti *TextInput) RemoveSelection() {
	if ti.HasSelectedText() {
		ti.setValue(ti.value[0:ti.selectionStart] + ti.value[ti.selectionEnd:])
		ti.fireChangedEvent()
		ti.moveCursor(ti.selectionStart)
		ti.Deselect()
	}
}

func (ti *TextInput) Submit() {
	ti.setValue(ti.onSubmitFunc(ti.value))
	ti.eventManager.Fire(ti.SubmittedEvent, &TextInputSubmittedEventArgs{
		TextInput: ti,
		Text:      ti.value,
	})
}

func (ti *TextInput) fireChangedEvent() {
	ti.eventManager.Fire(ti.ChangedEvent, &TextInputChangedEventArgs{
		TextInput: ti,
		Text:      ti.value,
	})
}

func (ti *TextInput) findPositionBeforeWord() textInputCursorPosition {
	if ti.cursorPosition <= 0 {
		return 0
	}

	var tmpCursorPosition textInputCursorPosition

	for i := int(ti.cursorPosition) - 1; i >= 0; i-- {
		if i >= len(ti.value) {
			continue
		}

		if i == 0 {
			return 0
		}

		if !wordSeparatorRegex.MatchString(string(ti.value[i])) {
			tmpCursorPosition = textInputCursorPosition(i)
			break
		}
	}

	for i := int(tmpCursorPosition); i >= 0; i-- {
		if i >= len(ti.value) {
			continue
		}
		if wordSeparatorRegex.MatchString(string(ti.value[i])) {
			return textInputCursorPosition(i + 1)
		}
	}

	return 0
}

func (ti *TextInput) findPositionAfterWord() textInputCursorPosition {
	if int(ti.cursorPosition) >= len(ti.value) {
		return textInputCursorPosition(len(ti.value))
	}

	var tmpCursorPosition textInputCursorPosition

	for i := int(ti.cursorPosition); i < len(ti.possibleCursorPosXs); i++ {
		if i == len(ti.value) {
			return textInputCursorPosition(len(ti.value))
		}

		if !wordSeparatorRegex.MatchString(string(ti.value[i])) {
			tmpCursorPosition = textInputCursorPosition(i)
			break
		}
	}

	for i := int(tmpCursorPosition); i < len(ti.possibleCursorPosXs); i++ {
		if i == len(ti.possibleCursorPosXs)-1 || wordSeparatorRegex.MatchString(string(ti.value[i])) {
			return textInputCursorPosition(i)
		}
	}

	return textInputCursorPosition(len(ti.value))
}

func (ti *TextInput) cursorPosX() int {
	return ti.possibleCursorPosXs[ti.cursorPosition] + ti.textPosX + ti.padding.Left - 1
}

func (ti *TextInput) findClosestPossibleCursorPosition() textInputCursorPosition {
	cursorPosX := input.CursorPosX - int(ti.absPosX) - ti.textPosX - ti.padding.Left + 1

	if cursorPosX <= ti.possibleCursorPosXs[0] {
		return 0
	}

	var lastId textInputCursorPosition = textInputCursorPosition(len(ti.possibleCursorPosXs) - 1)

	if cursorPosX >= ti.possibleCursorPosXs[lastId] {
		return lastId
	}

	var min textInputCursorPosition = 0
	var max textInputCursorPosition = lastId

	getClosest := func(a, b textInputCursorPosition, target int) textInputCursorPosition {
		if target-ti.possibleCursorPosXs[a] >= ti.possibleCursorPosXs[b]-target {
			return b
		} else {
			return a
		}
	}

	for min <= max {
		mid := (min + max) / 2
		switch {
		case cursorPosX < ti.possibleCursorPosXs[mid]:
			if mid > 0 && cursorPosX > ti.possibleCursorPosXs[mid-1] {
				return getClosest(mid-1, mid, cursorPosX)
			}

			max = mid - 1
		case cursorPosX > ti.possibleCursorPosXs[mid]:
			if mid < lastId && cursorPosX < ti.possibleCursorPosXs[mid+1] {
				return getClosest(mid, mid+1, cursorPosX)
			}

			min = mid + 1
		default:
			return mid
		}
	}

	return getClosest(max, min, cursorPosX)
}

func (ti *TextInput) calcScrollOffset() int {
	cursorPosX := ti.cursorPosX()
	scrollOffsetLowerBound := 0
	scrollOffsetUpperBound := fontutils.MeasureString(ti.value, ti.font) - (ti.width - ti.textPosX - ti.cursor.width - 2)
	if scrollOffsetUpperBound < 0 {
		scrollOffsetUpperBound = 0
	}

	applyBoundsToScrollOffset := func(offset int) int {
		switch {
		case offset < 0:
			return 0
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

func (ti *TextInput) moveCursor(position textInputCursorPosition) {
	ti.cursorPosition = position
	ti.cursor.ResetBlink()
}

func (ti *TextInput) updateSelectionBounds() {
	ti.selectionStart = textInputCursorPosition(math.Min(float64(ti.selectingFrom), float64(ti.cursorPosition)))
	ti.selectionEnd = textInputCursorPosition(math.Max(float64(ti.selectingFrom), float64(ti.cursorPosition)))
}

func (ti *TextInput) setValue(value string) {
	ti.value = value
	ti.afterChange()
}

func (ti *TextInput) afterChange() {
	ti.textPosY = ti.metrics.Ascent - ti.metrics.Descent - 1
	ti.possibleCursorPosXs = make([]int, len(ti.value)+1)
	ti.possibleCursorPosXs[0] = 0

	for i, c := range ti.value {
		ti.possibleCursorPosXs[i+1] = ti.possibleCursorPosXs[i] + fontutils.MeasureString(string(c), ti.font)
	}
}

func (ti *TextInput) actionKeyPressed() (bool, ebiten.Key) {
	for key := range ti.modifierKeysPressed {
		ti.modifierKeysPressed[key] = input.KeyPressed[key]
	}

	for key := range ti.actionKeyHandlers {
		if input.KeyPressed[key] {
			return true, key
		}
	}

	ti.lastActionKeyPressed = input.KeyNone

	return false, input.KeyNone
}

func (ti *TextInput) handleKeyLeft() textInputAction {
	ti.checkForShift()

	switch {
	case ti.modifierKeysPressed[ebiten.KeyAlt] && (ti.modifierKeysPressed[ebiten.KeyControl] || ti.modifierKeysPressed[ebiten.KeyMeta]):
		return textInputIdle
	case ti.modifierKeysPressed[ebiten.KeyAlt]:
		return textInputWordLeft
	case input.OSMacOS() && ti.modifierKeysPressed[ebiten.KeyMeta] || !input.OSMacOS() && ti.modifierKeysPressed[ebiten.KeyControl]:
		return textInputHome
	default:
		return textInputCursorLeft
	}
}

func (ti *TextInput) handleKeyRight() textInputAction {
	ti.checkForShift()

	switch {
	case ti.modifierKeysPressed[ebiten.KeyAlt] && (ti.modifierKeysPressed[ebiten.KeyControl] || ti.modifierKeysPressed[ebiten.KeyMeta]):
		return textInputIdle
	case ti.modifierKeysPressed[ebiten.KeyAlt]:
		return textInputWordRight
	case input.OSMacOS() && ti.modifierKeysPressed[ebiten.KeyMeta] || !input.OSMacOS() && ti.modifierKeysPressed[ebiten.KeyControl]:
		return textInputEnd
	default:
		return textInputCursorRight
	}
}

func (ti *TextInput) handleKeyHome() textInputAction {
	ti.checkForShift()

	switch {
	case ti.modifierKeysPressed[ebiten.KeyAlt]:
		return textInputIdle
	default:
		return textInputHome
	}
}

func (ti *TextInput) handleKeyEnd() textInputAction {
	ti.checkForShift()

	switch {
	case ti.modifierKeysPressed[ebiten.KeyAlt]:
		return textInputIdle
	default:
		return textInputEnd
	}
}

func (ti *TextInput) handleKeyDelete() textInputAction {
	ti.checkForShift()

	switch {
	case ti.HasSelectedText():
		return textInputRemoveSelection
	case ti.modifierKeysPressed[ebiten.KeyShift]:
		return textInputIdle
	case ti.modifierKeysPressed[ebiten.KeyAlt]:
		return textInputDeleteWord
	default:
		return textInputDelete
	}
}

func (ti *TextInput) handleKeyBackspace() textInputAction {
	ti.checkForShift()

	switch {
	case ti.HasSelectedText():
		return textInputRemoveSelection
	case ti.modifierKeysPressed[ebiten.KeyShift]:
		return textInputIdle
	case ti.modifierKeysPressed[ebiten.KeyAlt]:
		return textInputBackspaceWord
	default:
		return textInputBackspace
	}
}

func (ti *TextInput) handleKeyEnter() textInputAction {
	for _, pressed := range ti.modifierKeysPressed {
		if pressed {
			return textInputIdle
		}
	}

	return textInputSubmit
}

func (ti *TextInput) checkForShift() {
	if !ti.modifierKeysPressed[ebiten.KeyShift] {
		ti.selecting = false
		ti.Deselect()
	} else if !ti.selecting {
		ti.selecting = true
		if !ti.HasSelectedText() {
			ti.selectingFrom = ti.cursorPosition
		}
	}
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

		handler, found := ti.actionKeyHandlers[pressedKey]
		if !found {
			return ti.idleStateFactory()
		}

		action := handler()
		if action == textInputIdle {
			return ti.idleStateFactory()
		}

		ti.actionHandlers[action]()

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

func (ti *TextInput) drawText(clr color.RGBA) {
	textStartPosX := ti.textPosX - ti.scrollOffset + ti.padding.Left

	if !ti.HasSelectedText() {
		text.Draw(ti.image, ti.value, ti.font, textStartPosX, ti.textPosY+ti.padding.Top, clr)
		return
	}

	if ti.selectionStart > 0 {
		text.Draw(ti.image, ti.value[0:ti.selectionStart], ti.font, textStartPosX, ti.textPosY+ti.padding.Top, clr)
	}

	text.Draw(ti.image, ti.value[ti.selectionStart:ti.selectionEnd], ti.font, textStartPosX+ti.possibleCursorPosXs[ti.selectionStart], ti.textPosY+ti.padding.Top, colorutils.Invert(clr))

	if int(ti.selectionEnd) <= len(ti.value)-1 {
		text.Draw(ti.image, ti.value[ti.selectionEnd:], ti.font, textStartPosX+ti.possibleCursorPosXs[ti.selectionEnd], ti.textPosY+ti.padding.Top, clr)
	}
}

func (ti *TextInput) Draw() *ebiten.Image {
	if ti.hidden {
		return ti.image
	}

	if !ti.disabled {
		ti.state = ti.state(ti)
	}

	ti.updateSelectionBounds()

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
		ti.drawText(ti.colorDisabled)
	case ti.hovering:
		ti.drawText(ti.colorHovered)
	default:
		ti.drawText(ti.color)
	}

	ti.component.Draw()

	return ti.image
}
