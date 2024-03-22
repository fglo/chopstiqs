package component

import (
	"image"
	"image/color"

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

	cursorPos         int
	possibleCursorPos []int

	ClickedEvent *event.Event
	KeyPressed   *event.Event
	// active bool

	drawer TextInputDrawer
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

func NewTextInput(options *TextInputOptions) *TextInput {
	ti := &TextInput{
		ClickedEvent: &event.Event{},

		color:   color.RGBA{230, 230, 230, 255},
		font:    fontutils.DefaultFontFace,
		metrics: fontutils.NewMetrics(fontutils.DefaultFontFace.Metrics()),

		textPosX: 3,

		drawer: &DefaultTextInputDrawer{
			BorderColor: color.RGBA{230, 230, 230, 255},
		},
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

	ti.component.AddMouseButtonReleasedHandler(func(args *ComponentMouseButtonReleasedEventArgs) {
		if !ti.disabled && args.Button == ebiten.MouseButtonLeft {

			ti.cursorPos = ti.findClosestPossibleCursorPosition() + ti.textPosX + ti.padding.Left - 1

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

func (ti *TextInput) findClosestPossibleCursorPosition() int {
	cursorPos := (input.CursorPosX - int(ti.absPosX))

	if cursorPos <= ti.possibleCursorPos[0] {
		return ti.possibleCursorPos[0]
	}

	lastId := len(ti.possibleCursorPos) - 1

	if cursorPos >= ti.possibleCursorPos[lastId] {
		return ti.possibleCursorPos[lastId]
	}

	min := 0
	max := lastId

	for min <= max {
		mid := (min + max) / 2
		if cursorPos < ti.possibleCursorPos[mid] {
			if (ti.possibleCursorPos[mid] - cursorPos) < (cursorPos - ti.possibleCursorPos[mid-1]) {
				return ti.possibleCursorPos[mid]
			}

			max = mid - 1
		} else if cursorPos > ti.possibleCursorPos[mid] {
			if (ti.possibleCursorPos[mid+1] - cursorPos) > (cursorPos - ti.possibleCursorPos[mid]) {
				return ti.possibleCursorPos[mid]
			}

			min = mid + 1
		} else {
			return ti.possibleCursorPos[mid]
		}
	}

	if (ti.possibleCursorPos[min] - cursorPos) < (cursorPos - ti.possibleCursorPos[max]) {
		return ti.possibleCursorPos[min]
	} else {
		return ti.possibleCursorPos[max]
	}
}

func (ti *TextInput) Draw() *ebiten.Image {
	if ti.hidden {
		return ti.image
	}

	ti.drawer.Draw(ti)

	text.Draw(ti.image, ti.value, ti.font, ti.textPosX+ti.padding.Left, ti.textPosY+ti.padding.Top, ti.color)

	ti.component.Draw()

	return ti.image
}
