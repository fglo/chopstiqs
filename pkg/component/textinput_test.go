package component

import (
	"reflect"
	"runtime"
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
	is.Equal(ti.cursorPosition, 1)

	keyRelease(t, ebiten.KeyLeft)
	time.Sleep(textInputActionKeyRepeatDelay)
	handleState(t, ti)
	is.Equal(ti.cursorPosition, 1)

	keyPress(t, ebiten.KeyLeft)
	time.Sleep(textInputActionKeyRepeatDelay)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(ti.cursorPosition, 0)

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
	is.Equal(ti.cursorPosition, 1)

	keyRelease(t, ebiten.KeyRight)
	time.Sleep(textInputActionKeyRepeatDelay)
	handleState(t, ti)
	is.Equal(ti.cursorPosition, 1)

	keyPress(t, ebiten.KeyRight)
	time.Sleep(textInputActionKeyRepeatDelay)
	handleState(t, ti)
	handleState(t, ti)
	is.Equal(ti.cursorPosition, 2)

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

func TestTextInput_pressedKeyHandler(t *testing.T) {
	resetInput(t)

	type args struct {
		key ebiten.Key
	}

	tests := []struct {
		name                string
		args                args
		pressedModifierKeys []ebiten.Key
		before              func()
		want                func()
	}{
		{
			name: "CursorLeft",
			args: args{
				key: ebiten.KeyLeft,
			},
			want: (&TextInput{}).CursorLeft,
		},
		{
			name: "CursorRight",
			args: args{
				key: ebiten.KeyRight,
			},
			want: (&TextInput{}).CursorRight,
		},
		{
			name: "Home",
			args: args{
				key: ebiten.KeyHome,
			},
			want: (&TextInput{}).Home,
		},
		{
			name: "End",
			args: args{
				key: ebiten.KeyEnd,
			},
			want: (&TextInput{}).End,
		},
		{
			name: "Delete",
			args: args{
				key: ebiten.KeyDelete,
			},
			want: (&TextInput{}).Delete,
		},
		{
			name: "Backspace",
			args: args{
				key: ebiten.KeyBackspace,
			},
			want: (&TextInput{}).Backspace,
		},
		{
			name: "Left with ctrl on windows",
			args: args{
				key: ebiten.KeyLeft,
			},
			pressedModifierKeys: []ebiten.Key{ebiten.KeyControl},
			before:              func() { input.OS = input.Windows },
			want:                (&TextInput{}).Home,
		},
		{
			name: "Right with ctrl on windows",
			args: args{
				key: ebiten.KeyRight,
			},
			pressedModifierKeys: []ebiten.Key{ebiten.KeyControl},
			before:              func() { input.OS = input.Windows },
			want:                (&TextInput{}).End,
		},
		{
			name: "Left with meta on windows",
			args: args{
				key: ebiten.KeyLeft,
			},
			pressedModifierKeys: []ebiten.Key{ebiten.KeyMeta},
			before:              func() { input.OS = input.Windows },
			want:                (&TextInput{}).CursorLeft,
		},
		{
			name: "Right with meta on windows",
			args: args{
				key: ebiten.KeyRight,
			},
			pressedModifierKeys: []ebiten.Key{ebiten.KeyMeta},
			before:              func() { input.OS = input.Windows },
			want:                (&TextInput{}).CursorRight,
		},
		{
			name: "Left with meta on macos",
			args: args{
				key: ebiten.KeyLeft,
			},
			pressedModifierKeys: []ebiten.Key{ebiten.KeyMeta},
			before:              func() { input.OS = input.MacOS },
			want:                (&TextInput{}).Home,
		},
		{
			name: "Right with meta on macos",
			args: args{
				key: ebiten.KeyRight,
			},
			pressedModifierKeys: []ebiten.Key{ebiten.KeyMeta},
			before:              func() { input.OS = input.MacOS },
			want:                (&TextInput{}).End,
		},
		{
			name: "Left with ctrl on macos",
			args: args{
				key: ebiten.KeyLeft,
			},
			pressedModifierKeys: []ebiten.Key{ebiten.KeyControl},
			before:              func() { input.OS = input.MacOS },
			want:                (&TextInput{}).CursorLeft,
		},
		{
			name: "Right with ctrl on macos",
			args: args{
				key: ebiten.KeyRight,
			},
			pressedModifierKeys: []ebiten.Key{ebiten.KeyControl},
			before:              func() { input.OS = input.MacOS },
			want:                (&TextInput{}).CursorRight,
		},
		{
			name: "WordLeft",
			args: args{
				key: ebiten.KeyLeft,
			},
			pressedModifierKeys: []ebiten.Key{ebiten.KeyAlt},
			want:                (&TextInput{}).WordLeft,
		},
		{
			name: "WordRight",
			args: args{
				key: ebiten.KeyRight,
			},
			pressedModifierKeys: []ebiten.Key{ebiten.KeyAlt},
			want:                (&TextInput{}).WordRight,
		},
		{
			name: "DeleteWord",
			args: args{
				key: ebiten.KeyDelete,
			},
			pressedModifierKeys: []ebiten.Key{ebiten.KeyAlt},
			want:                (&TextInput{}).DeleteWord,
		},
		{
			name: "BackspaceWord",
			args: args{
				key: ebiten.KeyBackspace,
			},
			pressedModifierKeys: []ebiten.Key{ebiten.KeyAlt},
			want:                (&TextInput{}).BackspaceWord,
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
				ti.modifierKeys[key] = true
			}

			got := ti.pressedKeyHandler(tt.args.key)

			wantName := runtime.FuncForPC(reflect.ValueOf(tt.want).Pointer()).Name()
			gotName := runtime.FuncForPC(reflect.ValueOf(got).Pointer()).Name()

			if gotName != wantName {
				t.Errorf("TextInput.pressedKeyHandler() = %s, want %s", gotName, wantName)
			}
		})
	}
}
