package component

import (
	"testing"
	"time"

	"github.com/fglo/chopstiqs/pkg/event"
	"github.com/fglo/chopstiqs/pkg/input"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/matryer/is"
)

func handleState(t *testing.T, ti *TextInput) {
	t.Helper()

	ti.state = ti.state(ti)
	ti.eventManager.HandleFired()
}

func resetInput(t *testing.T) {
	t.Helper()

	input.SetSystem()
	input.InputChars = ebiten.AppendInputChars(input.InputChars)
	input.AnyKeyPressed = false
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		input.KeyPressed[k] = false
	}
}

func TestTextInput_PressedLeft(t *testing.T) {
	is := is.New(t)
	resetInput(t)

	eventManager := event.NewManager()

	ti := NewTextInput(&TextInputOptions{})
	ti.SetValue("qwerty")
	ti.SetEventManager(eventManager)

	ti.focused = true
	ti.cursorPosition = 2

	keyPress(t, ebiten.KeyLeft)
	time.Sleep(textInputActionKeyRepeatDelay)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(int(ti.cursorPosition), 1)

	keyRelease(t, ebiten.KeyLeft)
	time.Sleep(textInputActionKeyRepeatDelay)
	handleState(t, ti)
	is.Equal(int(ti.cursorPosition), 1)

	keyPress(t, ebiten.KeyLeft)
	time.Sleep(textInputActionKeyRepeatDelay)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(int(ti.cursorPosition), 0)

	keyRelease(t, ebiten.KeyLeft)
}

func TestTextInput_PressedRight(t *testing.T) {
	is := is.New(t)
	resetInput(t)

	eventManager := event.NewManager()

	ti := NewTextInput(&TextInputOptions{})
	ti.SetValue("qwerty")
	ti.SetEventManager(eventManager)

	ti.focused = true
	ti.cursorPosition = 0

	keyPress(t, ebiten.KeyRight)
	time.Sleep(textInputActionKeyRepeatDelay)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(int(ti.cursorPosition), 1)

	keyRelease(t, ebiten.KeyRight)
	time.Sleep(textInputActionKeyRepeatDelay)
	handleState(t, ti)
	is.Equal(int(ti.cursorPosition), 1)

	keyPress(t, ebiten.KeyRight)
	time.Sleep(textInputActionKeyRepeatDelay)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(int(ti.cursorPosition), 2)

	keyRelease(t, ebiten.KeyRight)
}

func TestTextInput_PressedEnter(t *testing.T) {
	is := is.New(t)
	resetInput(t)

	firedEventsCounter := 0

	eventManager := event.NewManager()

	ti := NewTextInput(&TextInputOptions{})
	ti.SetValue("qwerty")
	ti.SetEventManager(eventManager)
	ti.AddSubmittedHandler(func(args *TextInputSubmittedEventArgs) {
		firedEventsCounter++
	})

	ti.focused = true

	keyPress(t, ebiten.KeyEnter)
	time.Sleep(textInputActionKeyRepeatDelay)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(firedEventsCounter, 1)

	keyRelease(t, ebiten.KeyEnter)
	time.Sleep(textInputActionKeyRepeatDelay)
	handleState(t, ti)
	is.Equal(firedEventsCounter, 1)

	keyPress(t, ebiten.KeyEnter)
	time.Sleep(textInputActionKeyRepeatDelay)
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
			name:                "Left with ctrl on windows",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyControl},
			before:              func() { input.OS = input.Windows },
			want:                textInputHome,
		},
		{
			name:                "Left with meta on windows",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyMeta},
			before:              func() { input.OS = input.Windows },
			want:                textInputCursorLeft,
		},
		{
			name:                "Left with meta on macos",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyMeta},
			before:              func() { input.OS = input.MacOS },
			want:                textInputHome,
		},
		{
			name:                "Left with ctrl on macos",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyControl},
			before:              func() { input.OS = input.MacOS },
			want:                textInputCursorLeft,
		},
		{
			name:                "WordLeft",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyAlt},
			want:                textInputWordLeft,
		},
		{
			name:                "Left + Alt + Control on windows",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyAlt, ebiten.KeyControl},
			before:              func() { input.OS = input.Windows },
			want:                textInputIdle,
		},
		{
			name:                "Left + Alt + Meta on MacOS",
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
				t.Errorf("TextInput.pressedKeyHandler() = %v, want %v", got, tt.want)
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
			name:                "Right with ctrl on windows",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyControl},
			before:              func() { input.OS = input.Windows },
			want:                textInputEnd,
		},
		{
			name:                "Right with meta on windows",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyMeta},
			before:              func() { input.OS = input.Windows },
			want:                textInputCursorRight,
		},
		{
			name:                "Right with meta on macos",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyMeta},
			before:              func() { input.OS = input.MacOS },
			want:                textInputEnd,
		},
		{
			name:                "Right with ctrl on macos",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyControl},
			before:              func() { input.OS = input.MacOS },
			want:                textInputCursorRight,
		},
		{
			name:                "WordRight",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyAlt},
			want:                textInputWordRight,
		},
		{
			name:                "Right + Alt + Control on windows",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyAlt, ebiten.KeyControl},
			before:              func() { input.OS = input.Windows },
			want:                textInputIdle,
		},
		{
			name:                "Right + Alt + Meta on MacOS",
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
				t.Errorf("TextInput.pressedKeyHandler() = %v, want %v", got, tt.want)
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
				t.Errorf("TextInput.pressedKeyHandler() = %v, want %v", got, tt.want)
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
				t.Errorf("TextInput.pressedKeyHandler() = %v, want %v", got, tt.want)
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
			name:                "DeleteWord",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyAlt},
			want:                textInputDeleteWord,
		},
		{
			name:                "RemoveSelection",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyShift},
			before: func(ti *TextInput) {
				ti.selectingFrom = 0
				ti.cursorPosition = 2
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
				t.Errorf("TextInput.pressedKeyHandler() = %v, want %v", got, tt.want)
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
			name:                "BackspaceWord",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyAlt},
			want:                textInputBackspaceWord,
		},
		{
			name:                "RemoveSelection",
			pressedModifierKeys: []ebiten.Key{ebiten.KeyShift},
			before: func(ti *TextInput) {
				ti.selectingFrom = 0
				ti.cursorPosition = 2
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
				t.Errorf("TextInput.pressedKeyHandler() = %v, want %v", got, tt.want)
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
				t.Errorf("TextInput.pressedKeyHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
