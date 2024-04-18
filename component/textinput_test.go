package component

import (
	"testing"
	"time"

	"github.com/fglo/chopstiqs/event"
	"github.com/fglo/chopstiqs/input"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/matryer/is"
)

func newTestTextInput() *TextInput {
	eventManager := event.NewManager()

	ti := NewTextInput(nil)
	ti.SetEventManager(eventManager)

	return ti
}

func handleState(t *testing.T, ti *TextInput) {
	t.Helper()

	ti.state = ti.state(ti)
	ti.eventManager.HandleFired()
}

func resetInput(t *testing.T) {
	t.Helper()

	input.DetectSystem()
	input.InputChars = ebiten.AppendInputChars(input.InputChars)
	input.AnyKeyPressed = false
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		input.KeyPressed[k] = false
	}
}

func TestTextInput_PressedLeft(t *testing.T) {
	is := is.New(t)
	resetInput(t)

	ti := newTestTextInput()
	ti.SetValue("qwerty")

	ti.focused = true
	ti.cursorPosition = 2

	keyPress(t, ebiten.KeyLeft)
	time.Sleep(textInputActionRepeatDelay)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(int(ti.cursorPosition), 1)

	keyRelease(t, ebiten.KeyLeft)
	time.Sleep(textInputActionRepeatDelay)
	handleState(t, ti)
	is.Equal(int(ti.cursorPosition), 1)

	keyPress(t, ebiten.KeyLeft)
	time.Sleep(textInputActionRepeatDelay)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(int(ti.cursorPosition), 0)

	keyRelease(t, ebiten.KeyLeft)
}

func TestTextInput_PressedLeftWithControl_onWindows(t *testing.T) {
	is := is.New(t)
	resetInput(t)

	input.OS = input.Windows

	ti := newTestTextInput()
	ti.SetValue("01234 6789")

	ti.focused = true
	ti.cursorPosition = 9

	keyPress(t, ebiten.KeyLeft)
	keyPress(t, ebiten.KeyControl)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(int(ti.cursorPosition), 6)

	time.Sleep(time.Millisecond)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(int(ti.cursorPosition), 6)

	time.Sleep(textInputActionRepeatDelay)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(int(ti.cursorPosition), 0)

	keyRelease(t, ebiten.KeyLeft)
	keyRelease(t, ebiten.KeyControl)
}

func TestTextInput_PressedLeftWithAlt_onMacOS(t *testing.T) {
	is := is.New(t)
	resetInput(t)

	input.OS = input.MacOS

	ti := newTestTextInput()
	ti.SetValue("01234 6789")

	ti.focused = true
	ti.cursorPosition = 9

	keyPress(t, ebiten.KeyLeft)
	keyPress(t, ebiten.KeyAlt)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(int(ti.cursorPosition), 6)

	time.Sleep(time.Millisecond)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(int(ti.cursorPosition), 6)

	time.Sleep(textInputActionRepeatDelay)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(int(ti.cursorPosition), 0)

	keyRelease(t, ebiten.KeyLeft)
	keyRelease(t, ebiten.KeyAlt)
}

func TestTextInput_PressedLeftWithShift(t *testing.T) {
	is := is.New(t)
	resetInput(t)

	ti := newTestTextInput()
	ti.SetValue("qwerty")

	ti.focused = true
	ti.cursorPosition = 2

	keyPress(t, ebiten.KeyLeft)
	keyPress(t, ebiten.KeyShift)
	time.Sleep(textInputActionRepeatDelay)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(int(ti.selectingFrom), 2)
	is.Equal(int(ti.cursorPosition), 1)

	keyRelease(t, ebiten.KeyLeft)
	keyRelease(t, ebiten.KeyShift)
}

func TestTextInput_PressedRight(t *testing.T) {
	is := is.New(t)
	resetInput(t)

	ti := newTestTextInput()
	ti.SetValue("qwerty")

	ti.focused = true
	ti.cursorPosition = 0

	keyPress(t, ebiten.KeyRight)
	time.Sleep(textInputActionRepeatDelay)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(int(ti.cursorPosition), 1)

	keyRelease(t, ebiten.KeyRight)
	time.Sleep(textInputActionRepeatDelay)
	handleState(t, ti)
	is.Equal(int(ti.cursorPosition), 1)

	keyPress(t, ebiten.KeyRight)
	time.Sleep(textInputActionRepeatDelay)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(int(ti.cursorPosition), 2)

	keyRelease(t, ebiten.KeyRight)
}

func TestTextInput_PressedRightWithControl_onWindows(t *testing.T) {
	is := is.New(t)
	resetInput(t)

	input.OS = input.Windows

	ti := newTestTextInput()
	ti.SetValue("01234 6789")

	ti.focused = true
	ti.cursorPosition = 1

	keyPress(t, ebiten.KeyRight)
	keyPress(t, ebiten.KeyControl)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(int(ti.cursorPosition), 5)

	time.Sleep(time.Millisecond)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(int(ti.cursorPosition), 5)

	time.Sleep(textInputActionRepeatDelay)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(int(ti.cursorPosition), 10)

	keyRelease(t, ebiten.KeyRight)
	keyRelease(t, ebiten.KeyControl)
}

func TestTextInput_PressedRightWithAlt_onMacOS(t *testing.T) {
	is := is.New(t)
	resetInput(t)

	input.OS = input.MacOS

	ti := newTestTextInput()
	ti.SetValue("01234 6789")

	ti.focused = true
	ti.cursorPosition = 1

	keyPress(t, ebiten.KeyRight)
	keyPress(t, ebiten.KeyAlt)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(int(ti.cursorPosition), 5)

	time.Sleep(time.Millisecond)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(int(ti.cursorPosition), 5)

	time.Sleep(textInputActionRepeatDelay)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(int(ti.cursorPosition), 10)

	keyRelease(t, ebiten.KeyRight)
	keyRelease(t, ebiten.KeyAlt)
}

func TestTextInput_PressedRightWithShift(t *testing.T) {
	is := is.New(t)
	resetInput(t)

	ti := newTestTextInput()
	ti.SetValue("qwerty")

	ti.focused = true
	ti.cursorPosition = 0

	keyPress(t, ebiten.KeyRight)
	keyPress(t, ebiten.KeyShift)
	time.Sleep(textInputActionRepeatDelay)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(int(ti.selectingFrom), 0)
	is.Equal(int(ti.cursorPosition), 1)

	keyRelease(t, ebiten.KeyRight)
	keyRelease(t, ebiten.KeyShift)
}

func TestTextInput_PressedEnter(t *testing.T) {
	is := is.New(t)
	resetInput(t)

	firedEventsCounter := 0

	ti := newTestTextInput()
	ti.SetValue("qwerty")
	ti.AddSubmittedHandler(func(args *TextInputSubmittedEventArgs) {
		firedEventsCounter++
	})

	ti.focused = true

	keyPress(t, ebiten.KeyEnter)
	time.Sleep(textInputActionRepeatDelay)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(firedEventsCounter, 1)

	keyRelease(t, ebiten.KeyEnter)
	time.Sleep(textInputActionRepeatDelay)
	handleState(t, ti)
	is.Equal(firedEventsCounter, 1)

	keyPress(t, ebiten.KeyEnter)
	time.Sleep(textInputActionRepeatDelay)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(firedEventsCounter, 2)

	keyRelease(t, ebiten.KeyEnter)
}

func TestTextInput_handleKeyLeft(t *testing.T) {
	resetInput(t)

	tests := []struct {
		name                string
		pressedModifierKeys []ebiten.Key
		before              func()
		want                textInputAction
	}{
		{
			name: "CursorLeft",
			want: textInputCursorLeft,
		},
		{
			name:                "Left with ctrl (Windows)",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyControl},
			before:              func() { input.OS = input.Windows },
			want:                textInputWordLeft,
		},
		{
			name:                "Left with meta (Windows)",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyMeta},
			before:              func() { input.OS = input.Windows },
			want:                textInputIdle,
		},
		{
			name:                "Left with meta (MacOS)",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyMeta},
			before:              func() { input.OS = input.MacOS },
			want:                textInputHome,
		},
		{
			name:                "Left with ctrl (MacOS)",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyControl},
			before:              func() { input.OS = input.MacOS },
			want:                textInputCursorLeft,
		},
		{
			name:                "Left + Alt (MacOS)",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyAlt},
			before:              func() { input.OS = input.MacOS },
			want:                textInputWordLeft,
		},
		{
			name:                "Left + Alt + Control (Windows)",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyAlt, ebiten.KeyControl},
			before:              func() { input.OS = input.Windows },
			want:                textInputIdle,
		},
		{
			name:                "Left + Alt + Meta (MacOS)",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyAlt, ebiten.KeyMeta},
			before:              func() { input.OS = input.MacOS },
			want:                textInputIdle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetInput(t)

			if tt.before != nil {
				tt.before()
			}

			ti := NewTextInput(nil)

			for _, key := range tt.pressedModifierKeys {
				ti.modifierKeysPressed[key] = true
			}

			got := ti.handleKeyLeft()

			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestTextInput_handleKeyRight(t *testing.T) {
	resetInput(t)

	tests := []struct {
		name                string
		pressedModifierKeys []ebiten.Key
		before              func()
		want                textInputAction
	}{
		{
			name: "CursorRight",
			want: textInputCursorRight,
		},
		{
			name:                "Right with ctrl (Windows)",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyControl},
			before:              func() { input.OS = input.Windows },
			want:                textInputWordRight,
		},
		{
			name:                "Right with meta (Windows)",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyMeta},
			before:              func() { input.OS = input.Windows },
			want:                textInputIdle,
		},
		{
			name:                "Right with meta (MacOS)",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyMeta},
			before:              func() { input.OS = input.MacOS },
			want:                textInputEnd,
		},
		{
			name:                "Right with ctrl (MacOS)",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyControl},
			before:              func() { input.OS = input.MacOS },
			want:                textInputCursorRight,
		},
		{
			name:                "Right + Alt (MacOS)",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyAlt},
			before:              func() { input.OS = input.MacOS },
			want:                textInputWordRight,
		},
		{
			name:                "Right + Alt + Control (Windows)",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyAlt, ebiten.KeyControl},
			before:              func() { input.OS = input.Windows },
			want:                textInputIdle,
		},
		{
			name:                "Right + Alt + Meta (MacOS)",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyAlt, ebiten.KeyMeta},
			before:              func() { input.OS = input.MacOS },
			want:                textInputIdle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetInput(t)

			if tt.before != nil {
				tt.before()
			}

			ti := NewTextInput(nil)

			for _, key := range tt.pressedModifierKeys {
				ti.modifierKeysPressed[key] = true
			}

			got := ti.handleKeyRight()

			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestTextInput_handleKeyHome(t *testing.T) {
	resetInput(t)

	tests := []struct {
		name                string
		pressedModifierKeys []ebiten.Key
		before              func()
		want                textInputAction
	}{
		{
			name: "Home",
			want: textInputHome,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetInput(t)

			if tt.before != nil {
				tt.before()
			}

			ti := NewTextInput(nil)

			for _, key := range tt.pressedModifierKeys {
				ti.modifierKeysPressed[key] = true
			}

			got := ti.handleKeyHome()

			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestTextInput_handleKeyEnd(t *testing.T) {
	resetInput(t)

	tests := []struct {
		name                string
		pressedModifierKeys []ebiten.Key
		before              func()
		want                textInputAction
	}{
		{
			name: "End",
			want: textInputEnd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetInput(t)

			if tt.before != nil {
				tt.before()
			}

			ti := NewTextInput(nil)

			for _, key := range tt.pressedModifierKeys {
				ti.modifierKeysPressed[key] = true
			}

			got := ti.handleKeyEnd()

			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestTextInput_handleKeyDelete(t *testing.T) {
	resetInput(t)

	tests := []struct {
		name                string
		pressedModifierKeys []ebiten.Key
		before              func(*TextInput)
		want                textInputAction
	}{
		{
			name: "Delete",
			want: textInputDelete,
		},
		{
			name:                "Delete + Alt (MacOS)",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyAlt},
			before:              func(*TextInput) { input.OS = input.MacOS },
			want:                textInputDeleteWord,
		},
		{
			name:                "Delete + CTRL (Windows)",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyControl},
			before:              func(*TextInput) { input.OS = input.Windows },
			want:                textInputDeleteWord,
		},
		{
			name: "RemoveSelection",
			before: func(ti *TextInput) {
				ti.selectionStart = 0
				ti.selectionEnd = 2
			},
			want: textInputRemoveSelection,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetInput(t)

			ti := NewTextInput(nil)

			if tt.before != nil {
				tt.before(ti)
			}

			for _, key := range tt.pressedModifierKeys {
				ti.modifierKeysPressed[key] = true
			}

			got := ti.handleKeyDelete()

			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestTextInput_handleKeyBackspace(t *testing.T) {
	resetInput(t)

	tests := []struct {
		name                string
		pressedModifierKeys []ebiten.Key
		before              func(*TextInput)
		want                textInputAction
	}{
		{
			name: "Backspace",
			want: textInputBackspace,
		},
		{
			name:                "Backspace + Alt (MacOS)",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyAlt},
			before:              func(*TextInput) { input.OS = input.MacOS },
			want:                textInputBackspaceWord,
		},
		{
			name:                "Backspace + CTRL (Windows)",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyControl},
			before:              func(*TextInput) { input.OS = input.Windows },
			want:                textInputBackspaceWord,
		},
		{
			name: "RemoveSelection",
			before: func(ti *TextInput) {
				ti.selectionStart = 0
				ti.selectionEnd = 2
			},
			want: textInputRemoveSelection,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetInput(t)

			ti := NewTextInput(nil)

			if tt.before != nil {
				tt.before(ti)
			}

			for _, key := range tt.pressedModifierKeys {
				ti.modifierKeysPressed[key] = true
			}

			got := ti.handleKeyBackspace()

			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestTextInput_handleKeyEnter(t *testing.T) {
	resetInput(t)

	tests := []struct {
		name                string
		pressedModifierKeys []ebiten.Key
		before              func()
		want                textInputAction
	}{
		{
			name: "Submit",
			want: textInputSubmit,
		},
		{
			name:                "Enter + Alt",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyAlt},
			want:                textInputIdle,
		},
		{
			name:                "Enter + Control",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyControl},
			want:                textInputIdle,
		},
		{
			name:                "Enter + Meta",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyMeta},
			want:                textInputIdle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetInput(t)

			if tt.before != nil {
				tt.before()
			}

			ti := NewTextInput(nil)

			for _, key := range tt.pressedModifierKeys {
				ti.modifierKeysPressed[key] = true
			}

			got := ti.handleKeyEnter()

			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestTextInput_handleKeyEscape(t *testing.T) {
	resetInput(t)

	tests := []struct {
		name                string
		pressedModifierKeys []ebiten.Key
		before              func()
		want                textInputAction
	}{
		{
			name: "Unfocus",
			want: textInputUnfocus,
		},
		{
			name:                "Escape + Alt",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyAlt},
			want:                textInputIdle,
		},
		{
			name:                "Escape + Control",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyControl},
			want:                textInputIdle,
		},
		{
			name:                "Escape + Meta",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyMeta},
			want:                textInputIdle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetInput(t)

			if tt.before != nil {
				tt.before()
			}

			ti := NewTextInput(nil)

			for _, key := range tt.pressedModifierKeys {
				ti.modifierKeysPressed[key] = true
			}

			got := ti.handleKeyEscape()

			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestTextInput_handleKeyCtrl_Windows(t *testing.T) {
	resetInput(t)

	tests := []struct {
		name                  string
		pressedAdditionalKeys []ebiten.Key
		before                func()
		want                  textInputAction
	}{
		{
			name:                  "CTRL + Left",
			pressedAdditionalKeys: []ebiten.Key{ebiten.KeyLeft},
			want:                  textInputWordLeft,
		},
		{
			name:                  "CTRL + Right",
			pressedAdditionalKeys: []ebiten.Key{ebiten.KeyRight},
			want:                  textInputWordRight,
		},
		{
			name:                  "CTRL + C",
			pressedAdditionalKeys: []ebiten.Key{ebiten.KeyC},
			want:                  textInputCopy,
		},
		{
			name:                  "CTRL + V",
			pressedAdditionalKeys: []ebiten.Key{ebiten.KeyV},
			want:                  textInputPaste,
		},
		{
			name:                  "CTRL + X",
			pressedAdditionalKeys: []ebiten.Key{ebiten.KeyX},
			want:                  textInputCut,
		},
		{
			name:                  "CTRL + Z",
			pressedAdditionalKeys: []ebiten.Key{ebiten.KeyZ},
			want:                  textInputUndo,
		},
		{
			name:                  "CTRL + Y",
			pressedAdditionalKeys: []ebiten.Key{ebiten.KeyY},
			want:                  textInputRedo,
		},
		{
			name:                  "CTRL + SHIFT + Z",
			pressedAdditionalKeys: []ebiten.Key{ebiten.KeyShift, ebiten.KeyZ},
			want:                  textInputRedo,
		},
		{
			name:                  "CTRL + A",
			pressedAdditionalKeys: []ebiten.Key{ebiten.KeyA},
			want:                  textInputSelectAll,
		},
		{
			name:                  "CTRL + Backspace",
			pressedAdditionalKeys: []ebiten.Key{ebiten.KeyBackspace},
			want:                  textInputBackspaceWord,
		},
		{
			name:                  "CTRL + Delete",
			pressedAdditionalKeys: []ebiten.Key{ebiten.KeyDelete},
			want:                  textInputDeleteWord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetInput(t)

			input.OS = input.Windows

			if tt.before != nil {
				tt.before()
			}

			ti := NewTextInput(nil)

			ti.modifierKeysPressed[ebiten.KeyControl] = true

			for _, key := range tt.pressedAdditionalKeys {
				input.KeyPressed[key] = true
			}

			got := ti.handleKeyCtrl()

			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestTextInput_handleKeyMeta_MacOS(t *testing.T) {
	resetInput(t)

	tests := []struct {
		name                  string
		pressedAdditionalKeys []ebiten.Key
		before                func()
		want                  textInputAction
	}{
		{
			name:                  "CMD + Left",
			pressedAdditionalKeys: []ebiten.Key{ebiten.KeyLeft},
			want:                  textInputHome,
		},
		{
			name:                  "CMD + Right",
			pressedAdditionalKeys: []ebiten.Key{ebiten.KeyRight},
			want:                  textInputEnd,
		},
		{
			name:                  "CMD + C",
			pressedAdditionalKeys: []ebiten.Key{ebiten.KeyC},
			want:                  textInputCopy,
		},
		{
			name:                  "CMD + V",
			pressedAdditionalKeys: []ebiten.Key{ebiten.KeyV},
			want:                  textInputPaste,
		},
		{
			name:                  "CMD + X",
			pressedAdditionalKeys: []ebiten.Key{ebiten.KeyX},
			want:                  textInputCut,
		},
		{
			name:                  "CMD + Z",
			pressedAdditionalKeys: []ebiten.Key{ebiten.KeyZ},
			want:                  textInputUndo,
		},
		{
			name:                  "CMD + Y",
			pressedAdditionalKeys: []ebiten.Key{ebiten.KeyY},
			want:                  textInputIdle,
		},
		{
			name:                  "CMD + SHIFT + Z",
			pressedAdditionalKeys: []ebiten.Key{ebiten.KeyShift, ebiten.KeyZ},
			want:                  textInputRedo,
		},
		{
			name:                  "CMD + A",
			pressedAdditionalKeys: []ebiten.Key{ebiten.KeyA},
			want:                  textInputSelectAll,
		},
		{
			name:                  "CMD + Backspace",
			pressedAdditionalKeys: []ebiten.Key{ebiten.KeyBackspace},
			want:                  textInputBackspaceToBeginning,
		},
		{
			name:                  "CMD + Delete",
			pressedAdditionalKeys: []ebiten.Key{ebiten.KeyDelete},
			want:                  textInputDeleteToEnd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetInput(t)

			input.OS = input.MacOS

			if tt.before != nil {
				tt.before()
			}

			ti := NewTextInput(nil)

			ti.modifierKeysPressed[ebiten.KeyMeta] = true

			for _, key := range tt.pressedAdditionalKeys {
				input.KeyPressed[key] = true
			}

			got := ti.handleKeyMeta()

			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestTextInput_keyCombinations(t *testing.T) {
	resetInput(t)

	tests := []struct {
		name        string
		pressedKeys []ebiten.Key
		before      func(ti *TextInput)
		want        textInputAction
	}{
		{
			name:        "Left",
			pressedKeys: []ebiten.Key{ebiten.KeyLeft},
			want:        textInputCursorLeft,
		},
		{
			name:        "Left + CTRL (Windows)",
			pressedKeys: []ebiten.Key{ebiten.KeyLeft, ebiten.KeyControl},
			before:      func(ti *TextInput) { input.OS = input.Windows },
			want:        textInputWordLeft,
		},
		{
			name:        "Left + meta (Windows)",
			pressedKeys: []ebiten.Key{ebiten.KeyLeft, ebiten.KeyMeta},
			before:      func(ti *TextInput) { input.OS = input.Windows },
			want:        textInputIdle,
		},
		{
			name:        "Left + meta (MacOS)",
			pressedKeys: []ebiten.Key{ebiten.KeyLeft, ebiten.KeyMeta},
			before:      func(ti *TextInput) { input.OS = input.MacOS },
			want:        textInputHome,
		},
		{
			name:        "Left + CTRL (MacOS)",
			pressedKeys: []ebiten.Key{ebiten.KeyLeft, ebiten.KeyControl},
			before:      func(ti *TextInput) { input.OS = input.MacOS },
			want:        textInputIdle,
		},
		{
			name:        "Left + Alt (MacOS)",
			pressedKeys: []ebiten.Key{ebiten.KeyLeft, ebiten.KeyAlt},
			before:      func(ti *TextInput) { input.OS = input.MacOS },
			want:        textInputWordLeft,
		},
		{
			name:        "Left + Alt + Control (Windows)",
			pressedKeys: []ebiten.Key{ebiten.KeyLeft, ebiten.KeyAlt, ebiten.KeyControl},
			before:      func(ti *TextInput) { input.OS = input.Windows },
			want:        textInputIdle,
		},
		{
			name:        "Left + Alt + Meta (MacOS)",
			pressedKeys: []ebiten.Key{ebiten.KeyLeft, ebiten.KeyAlt, ebiten.KeyMeta},
			before:      func(ti *TextInput) { input.OS = input.MacOS },
			want:        textInputIdle,
		},
		{
			name:        "Right",
			pressedKeys: []ebiten.Key{ebiten.KeyRight},
			want:        textInputCursorRight,
		},
		{
			name:        "Right + CTRL (Windows)",
			pressedKeys: []ebiten.Key{ebiten.KeyRight, ebiten.KeyControl},
			before:      func(ti *TextInput) { input.OS = input.Windows },
			want:        textInputWordRight,
		},
		{
			name:        "Right + meta (Windows)",
			pressedKeys: []ebiten.Key{ebiten.KeyRight, ebiten.KeyMeta},
			before:      func(ti *TextInput) { input.OS = input.Windows },
			want:        textInputIdle,
		},
		{
			name:        "Right + meta (MacOS)",
			pressedKeys: []ebiten.Key{ebiten.KeyRight, ebiten.KeyMeta},
			before:      func(ti *TextInput) { input.OS = input.MacOS },
			want:        textInputEnd,
		},
		{
			name:        "Right + CTRL (MacOS)",
			pressedKeys: []ebiten.Key{ebiten.KeyRight, ebiten.KeyControl},
			before:      func(ti *TextInput) { input.OS = input.MacOS },
			want:        textInputIdle,
		},
		{
			name:        "Right + Alt (MacOS)",
			pressedKeys: []ebiten.Key{ebiten.KeyRight, ebiten.KeyAlt},
			before:      func(ti *TextInput) { input.OS = input.MacOS },
			want:        textInputWordRight,
		},
		{
			name:        "Right + Alt + Control (Windows)",
			pressedKeys: []ebiten.Key{ebiten.KeyRight, ebiten.KeyAlt, ebiten.KeyControl},
			before:      func(ti *TextInput) { input.OS = input.Windows },
			want:        textInputIdle,
		},
		{
			name:        "Right + Alt + Meta (MacOS)",
			pressedKeys: []ebiten.Key{ebiten.KeyRight, ebiten.KeyAlt, ebiten.KeyMeta},
			before:      func(ti *TextInput) { input.OS = input.MacOS },
			want:        textInputIdle,
		},
		{
			name:        "Home",
			pressedKeys: []ebiten.Key{ebiten.KeyHome},
			want:        textInputHome,
		},
		{
			name:        "End",
			pressedKeys: []ebiten.Key{ebiten.KeyEnd},
			want:        textInputEnd,
		},
		{
			name:        "Delete",
			pressedKeys: []ebiten.Key{ebiten.KeyDelete},
			want:        textInputDelete,
		},
		{
			name:        "Delete + Shift (Windows)",
			pressedKeys: []ebiten.Key{ebiten.KeyDelete, ebiten.KeyShift},
			before:      func(ti *TextInput) { input.OS = input.Windows },
			want:        textInputRemoveLine,
		},
		{
			name:        "Delete + Shift (MacOS)",
			pressedKeys: []ebiten.Key{ebiten.KeyDelete, ebiten.KeyShift},
			before:      func(ti *TextInput) { input.OS = input.MacOS },
			want:        textInputIdle,
		},
		{
			name:        "Delete + CTRL (Windows)",
			pressedKeys: []ebiten.Key{ebiten.KeyDelete, ebiten.KeyControl},
			before:      func(ti *TextInput) { input.OS = input.Windows },
			want:        textInputDeleteWord,
		},
		{
			name:        "Delete + Alt (MacOS)",
			pressedKeys: []ebiten.Key{ebiten.KeyDelete, ebiten.KeyAlt},
			before:      func(ti *TextInput) { input.OS = input.MacOS },
			want:        textInputDeleteWord,
		},
		{
			name:        "Delete with selected text",
			pressedKeys: []ebiten.Key{ebiten.KeyDelete},
			before: func(ti *TextInput) {
				ti.selectionStart = 0
				ti.selectionEnd = 2
			},
			want: textInputRemoveSelection,
		},
		{
			name:        "Backspace",
			pressedKeys: []ebiten.Key{ebiten.KeyBackspace},
			want:        textInputBackspace,
		},
		{
			name:        "Backspace + Shift",
			pressedKeys: []ebiten.Key{ebiten.KeyBackspace, ebiten.KeyShift},
			want:        textInputBackspace,
		},
		{
			name:        "Backspace + CTRL (Windows)",
			pressedKeys: []ebiten.Key{ebiten.KeyBackspace, ebiten.KeyControl},
			before:      func(ti *TextInput) { input.OS = input.Windows },
			want:        textInputBackspaceWord,
		},
		{
			name:        "Backspace + Alt (MacOS)",
			pressedKeys: []ebiten.Key{ebiten.KeyBackspace, ebiten.KeyAlt},
			before:      func(ti *TextInput) { input.OS = input.MacOS },
			want:        textInputBackspaceWord,
		},
		{
			name:        "Backspace with selected text",
			pressedKeys: []ebiten.Key{ebiten.KeyBackspace},
			before: func(ti *TextInput) {
				ti.selectionStart = 0
				ti.selectionEnd = 2
			},
			want: textInputRemoveSelection,
		},
		{
			name:        "Enter",
			pressedKeys: []ebiten.Key{ebiten.KeyEnter},
			want:        textInputSubmit,
		},
		{
			name:        "Enter + Alt",
			pressedKeys: []ebiten.Key{ebiten.KeyEnter, ebiten.KeyAlt},
			want:        textInputIdle,
		},
		{
			name:        "Enter + Control",
			pressedKeys: []ebiten.Key{ebiten.KeyEnter, ebiten.KeyControl},
			want:        textInputIdle,
		},
		{
			name:        "Enter + Meta",
			pressedKeys: []ebiten.Key{ebiten.KeyEnter, ebiten.KeyMeta},
			want:        textInputIdle,
		},
		{
			name:        "Escape",
			pressedKeys: []ebiten.Key{ebiten.KeyEscape},
			want:        textInputUnfocus,
		},
		{
			name:        "Escape + Alt",
			pressedKeys: []ebiten.Key{ebiten.KeyEscape, ebiten.KeyAlt},
			want:        textInputIdle,
		},
		{
			name:        "Escape + Control",
			pressedKeys: []ebiten.Key{ebiten.KeyEscape, ebiten.KeyControl},
			want:        textInputIdle,
		},
		{
			name:        "Escape + Meta",
			pressedKeys: []ebiten.Key{ebiten.KeyEscape, ebiten.KeyMeta},
			want:        textInputIdle,
		},
		{
			name:        "CTRL + Left (Windows)",
			pressedKeys: []ebiten.Key{ebiten.KeyControl, ebiten.KeyLeft},
			before:      func(ti *TextInput) { input.OS = input.Windows },
			want:        textInputWordLeft,
		},
		{
			name:        "CTRL + Right (Windows)",
			pressedKeys: []ebiten.Key{ebiten.KeyControl, ebiten.KeyRight},
			before:      func(ti *TextInput) { input.OS = input.Windows },
			want:        textInputWordRight,
		},
		{
			name:        "CTRL + C (Windows)",
			pressedKeys: []ebiten.Key{ebiten.KeyControl, ebiten.KeyC},
			before:      func(ti *TextInput) { input.OS = input.Windows },
			want:        textInputCopy,
		},
		{
			name:        "CTRL + V (Windows)",
			pressedKeys: []ebiten.Key{ebiten.KeyControl, ebiten.KeyV},
			before:      func(ti *TextInput) { input.OS = input.Windows },
			want:        textInputPaste,
		},
		{
			name:        "CTRL + X (Windows)",
			pressedKeys: []ebiten.Key{ebiten.KeyControl, ebiten.KeyX},
			before:      func(ti *TextInput) { input.OS = input.Windows },
			want:        textInputCut,
		},
		{
			name:        "CTRL + Z (Windows)",
			pressedKeys: []ebiten.Key{ebiten.KeyControl, ebiten.KeyZ},
			before:      func(ti *TextInput) { input.OS = input.Windows },
			want:        textInputUndo,
		},
		{
			name:        "CTRL + Y (Windows)",
			pressedKeys: []ebiten.Key{ebiten.KeyControl, ebiten.KeyY},
			before:      func(ti *TextInput) { input.OS = input.Windows },
			want:        textInputRedo,
		},
		{
			name:        "CTRL + SHIFT + Z (Windows)",
			pressedKeys: []ebiten.Key{ebiten.KeyControl, ebiten.KeyShift, ebiten.KeyZ},
			before:      func(ti *TextInput) { input.OS = input.Windows },
			want:        textInputRedo,
		},
		{
			name:        "CTRL + A (Windows)",
			pressedKeys: []ebiten.Key{ebiten.KeyControl, ebiten.KeyA},
			before:      func(ti *TextInput) { input.OS = input.Windows },
			want:        textInputSelectAll,
		},
		{
			name:        "CTRL + Backspace (Windows)",
			pressedKeys: []ebiten.Key{ebiten.KeyControl, ebiten.KeyBackspace},
			before:      func(ti *TextInput) { input.OS = input.Windows },
			want:        textInputBackspaceWord,
		},
		{
			name:        "CTRL + Delete (Windows)",
			pressedKeys: []ebiten.Key{ebiten.KeyControl, ebiten.KeyDelete},
			before:      func(ti *TextInput) { input.OS = input.Windows },
			want:        textInputDeleteWord,
		},
		{
			name:        "CMD + Left (MacOS)",
			pressedKeys: []ebiten.Key{ebiten.KeyMeta, ebiten.KeyLeft},
			before:      func(ti *TextInput) { input.OS = input.MacOS },
			want:        textInputHome,
		},
		{
			name:        "CMD + Right (MacOS)",
			pressedKeys: []ebiten.Key{ebiten.KeyMeta, ebiten.KeyRight},
			before:      func(ti *TextInput) { input.OS = input.MacOS },
			want:        textInputEnd,
		},
		{
			name:        "CMD + C (MacOS)",
			pressedKeys: []ebiten.Key{ebiten.KeyMeta, ebiten.KeyC},
			before:      func(ti *TextInput) { input.OS = input.MacOS },
			want:        textInputCopy,
		},
		{
			name:        "CMD + V (MacOS)",
			pressedKeys: []ebiten.Key{ebiten.KeyMeta, ebiten.KeyV},
			before:      func(ti *TextInput) { input.OS = input.MacOS },
			want:        textInputPaste,
		},
		{
			name:        "CMD + X (MacOS)",
			pressedKeys: []ebiten.Key{ebiten.KeyMeta, ebiten.KeyX},
			before:      func(ti *TextInput) { input.OS = input.MacOS },
			want:        textInputCut,
		},
		{
			name:        "CMD + Z (MacOS)",
			pressedKeys: []ebiten.Key{ebiten.KeyMeta, ebiten.KeyZ},
			before:      func(ti *TextInput) { input.OS = input.MacOS },
			want:        textInputUndo,
		},
		{
			name:        "CMD + Y (MacOS)",
			pressedKeys: []ebiten.Key{ebiten.KeyMeta, ebiten.KeyY},
			before:      func(ti *TextInput) { input.OS = input.MacOS },
			want:        textInputIdle,
		},
		{
			name:        "CMD + SHIFT + Z (MacOS)",
			pressedKeys: []ebiten.Key{ebiten.KeyMeta, ebiten.KeyShift, ebiten.KeyZ},
			before:      func(ti *TextInput) { input.OS = input.MacOS },
			want:        textInputRedo,
		},
		{
			name:        "CMD + A (MacOS)",
			pressedKeys: []ebiten.Key{ebiten.KeyMeta, ebiten.KeyA},
			before:      func(ti *TextInput) { input.OS = input.MacOS },
			want:        textInputSelectAll,
		},
		{
			name:        "CMD + Backspace (MacOS)",
			pressedKeys: []ebiten.Key{ebiten.KeyMeta, ebiten.KeyBackspace},
			before:      func(ti *TextInput) { input.OS = input.MacOS },
			want:        textInputBackspaceToBeginning,
		},
		{
			name:        "CMD + Delete (MacOS)",
			pressedKeys: []ebiten.Key{ebiten.KeyMeta, ebiten.KeyDelete},
			before:      func(ti *TextInput) { input.OS = input.MacOS },
			want:        textInputDeleteToEnd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetInput(t)

			ti := NewTextInput(nil)

			if tt.before != nil {
				tt.before(ti)
			}

			for _, key := range tt.pressedKeys {
				input.KeyPressed[key] = true
			}

			if pressed, key := ti.actionKeyPressed(); pressed {
				got := ti.handleActionKey(key)

				if got != tt.want {
					t.Errorf("got %s, want %s", got, tt.want)
				}
			} else {
				t.Fatal("No action keys were pressed.")
			}
		})
	}
}

func TestTextInput_Insert(t *testing.T) {
	type args struct {
		chars []rune
	}

	tests := []struct {
		name   string
		args   args
		before func(ti *TextInput)
		want   string
	}{
		{
			name: "Insert some characters",
			args: args{
				chars: []rune{'q', 'w', 'e', 'r', 't', 'y'},
			},
			want: "qwerty",
		},
		{
			name: "Insert some characters at the beginning",
			args: args{
				chars: []rune{'q', 'w', 'e'},
			},
			before: func(ti *TextInput) {
				ti.value = "rty"
				ti.cursorPosition = 0
			},
			want: "qwerty",
		},
		{
			name: "Insert some characters in the middle",
			args: args{
				chars: []rune{'y', 'u', 'i'},
			},
			before: func(ti *TextInput) {
				ti.value = "qwertop"
				ti.cursorPosition = 5
			},
			want: "qwertyuiop",
		},
		{
			name: "Insert some characters at the end",
			args: args{
				chars: []rune{'r', 't', 'y'},
			},
			before: func(ti *TextInput) {
				ti.value = "qwe"
				ti.cursorPosition = 3
			},
			want: "qwerty",
		},
		{
			name: "Replace some",
			args: args{
				chars: []rune{'a', 's', 'd', 'f'},
			},
			before: func(ti *TextInput) {
				ti.value = "qwerty"
				ti.selectingFrom = 0
				ti.cursorPosition = 6
				ti.updateSelectionBounds()
			},
			want: "asdf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resetInput(t)

			ti := newTestTextInput()

			if tt.before != nil {
				tt.before(ti)
			}

			ti.Insert(tt.args.chars)

			if ti.value != tt.want {
				t.Errorf("got %s, want %s", ti.value, tt.want)
			}
		})
	}
}
