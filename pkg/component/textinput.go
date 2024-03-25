package component

import (
	"image"
	"image/color"
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

	cursorPos                  int
	currentPossibleCursorPosId int
	possibleCursorPos          []int

	ClickedEvent    *event.Event
	KeyPressedEvent *event.Event

	drawer TextInputDrawer

	state textInputState

	lastActionKeyPressed ebiten.Key
	canRepeatAction      *atomic.Bool

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
		canRepeatAction:      &atomic.Bool{},
	}

	ti.state = ti.idleStateFactory()
	ti.canRepeatAction.Store(true)

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

	ti.KeyPressedEvent.AddHandler(func(args interface{}) {
		keyPressedArgs := args.(*KeyPressedEventArgs)

		switch keyPressedArgs.Key {
		case ebiten.KeyLeft:
			ti.CursorLeft()
		case ebiten.KeyRight:
			ti.CursorRight()
		case ebiten.KeyHome:
			ti.Home()
		case ebiten.KeyEnd:
			ti.End()
		case ebiten.KeyBackspace:
			ti.Backspace()
		case ebiten.KeyDelete:
			ti.Delete()
		default:
			ti.value += ebiten.KeyName(keyPressedArgs.Key)
		}
	})

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
			ti.currentPossibleCursorPosId = ti.findClosestPossibleCursorPosition()
			ti.setCursorPosition(ti.possibleCursorPos[ti.currentPossibleCursorPosId])

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

var actionKeyRepeatDelay = 500 * time.Millisecond
var actionKeyRepeatInterval = 60 * time.Millisecond

func (ti *TextInput) Insert(chars []rune) {
	ti.value += string(chars) // TODO: proper handling
}

type textInputState func(ti *TextInput) textInputState

func (ti *TextInput) idleStateFactory() textInputState {
	return func(ti *TextInput) textInputState {
		if !ti.focused || ti.disabled {
			return ti.idleStateFactory()
		}

		if len(input.InputChars) > 0 {
			return ti.inputStateFactory()
		}

		if pressed, key := ti.actionKeyPressed(); pressed {
			return ti.actionStateFactory(key)
		}

		return ti.idleStateFactory()
	}
}

func (ti *TextInput) inputStateFactory() textInputState {
	return func(ti *TextInput) textInputState {
		if !ti.focused || ti.disabled {
			return ti.idleStateFactory()
		}

		ti.Insert(input.InputChars)

		if pressed, key := ti.actionKeyPressed(); pressed {
			return ti.actionStateFactory(key)
		}

		return ti.idleStateFactory()
	}
}

func (ti *TextInput) actionStateFactory(pressedKey ebiten.Key) textInputState {
	return func(ti *TextInput) textInputState {
		if !ti.focused || ti.disabled {
			return ti.idleStateFactory()
		}

		delay := actionKeyRepeatDelay
		if ti.lastActionKeyPressed == pressedKey {
			delay = actionKeyRepeatInterval

			if !ti.canRepeatAction.Load() {
				return ti.idleStateFactory()
			}
		}

		ti.canRepeatAction.Store(false)

		time.AfterFunc(delay, func() {
			ti.canRepeatAction.Store(true)
		})

		ti.lastActionKeyPressed = pressedKey
		ti.pressedKeysHandlers[pressedKey]()

		return ti.idleStateFactory()
	}
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

	return false, input.KeyNone
}

// SetValue sets the value of the text input.
func (ti *TextInput) SetValue(value string) {
	// TODO: change deprecated function
	bounds := text.BoundString(ti.font, value) // nolint

	ti.value = value

	ti.textPosY = -bounds.Min.Y + ti.heightWithPadding/2 - 3
	ti.textBounds = bounds

	lastPos := 0

	ti.possibleCursorPos = make([]int, len(value)+1)
	ti.possibleCursorPos[0] = lastPos

	for i, c := range value {
		lastPos += int(fontutils.FixedInt26_6ToInt(font.MeasureString(ti.font, string(c))))
		ti.possibleCursorPos[i+1] = lastPos
	}

	_ = ti.possibleCursorPos
}

func (ti *TextInput) setCursorPosition(pos int) {
	ti.cursorPos = pos + ti.textPosX + ti.padding.Left - 1
}

func (ti *TextInput) findClosestPossibleCursorPosition() int {
	cursorPos := input.CursorPosX - int(ti.absPosX) - ti.textPosX - ti.padding.Left + 1

	if cursorPos <= ti.possibleCursorPos[0] {
		return 0
	}

	lastId := len(ti.possibleCursorPos) - 1

	if cursorPos >= ti.possibleCursorPos[lastId] {
		return lastId
	}

	min := 0
	max := lastId

	getClosest := func(a, b, target int) int {
		if target-ti.possibleCursorPos[a] >= ti.possibleCursorPos[b]-target {
			return b
		} else {
			return a
		}
	}

	for min <= max {
		mid := (min + max) / 2
		switch {
		case cursorPos < ti.possibleCursorPos[mid]:
			if mid > 0 && cursorPos > ti.possibleCursorPos[mid-1] {
				return getClosest(mid-1, mid, cursorPos)
			}

			max = mid - 1
		case cursorPos > ti.possibleCursorPos[mid]:
			if mid < lastId && cursorPos < ti.possibleCursorPos[mid+1] {
				return getClosest(mid, mid+1, cursorPos)
			}

			min = mid + 1
		default:
			return mid
		}
	}

	return getClosest(max, min, cursorPos)
}

func (ti *TextInput) CursorLeft() {
	if ti.currentPossibleCursorPosId > 0 {
		ti.currentPossibleCursorPosId--
		ti.setCursorPosition(ti.possibleCursorPos[ti.currentPossibleCursorPosId])
		ti.drawer.ResetCursorBlink()
	}
}

func (ti *TextInput) CursorRight() {
	if ti.currentPossibleCursorPosId < len(ti.possibleCursorPos)-1 {
		ti.currentPossibleCursorPosId++
		ti.setCursorPosition(ti.possibleCursorPos[ti.currentPossibleCursorPosId])
		ti.drawer.ResetCursorBlink()
	}
}

func (ti *TextInput) Home() {
	ti.currentPossibleCursorPosId = 0
	ti.setCursorPosition(ti.possibleCursorPos[ti.currentPossibleCursorPosId])
	ti.drawer.ResetCursorBlink()
}

func (ti *TextInput) End() {
	ti.currentPossibleCursorPosId--
	ti.setCursorPosition(ti.possibleCursorPos[len(ti.possibleCursorPos)-1])
	ti.drawer.ResetCursorBlink()
}

func (ti *TextInput) Delete() {}

func (ti *TextInput) Backspace() {}

func (ti *TextInput) Submit() {}

func (ti *TextInput) Draw() *ebiten.Image {
	// ti.handleState()
	ti.state = ti.state(ti)

	if ti.hidden {
		return ti.image
	}

	ti.drawer.Draw(ti)

	text.Draw(ti.image, ti.value, ti.font, ti.textPosX+ti.padding.Left, ti.textPosY+ti.padding.Top, ti.color)

	ti.component.Draw()

	return ti.image
}
