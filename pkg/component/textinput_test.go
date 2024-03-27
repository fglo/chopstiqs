package component

import (
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/fglo/chopstiqs/pkg/event"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/matryer/is"
)

func handleState(t *testing.T, ti *TextInput) {
	t.Helper()

	ti.state = ti.state(ti)
	ti.eventManager.HandleFired()
}

func TestTextInput_PressedLeft(t *testing.T) {
	is := is.New(t)

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
	ti := NewTextInput(nil)

	type args struct {
		key ebiten.Key
	}

	tests := []struct {
		name string
		args args
		want func()
	}{
		{
			name: "CursorLeft",
			args: args{
				key: ebiten.KeyLeft,
			},
			want: ti.CursorLeft,
		},
		{
			name: "CursorRight",
			args: args{
				key: ebiten.KeyRight,
			},
			want: ti.CursorRight,
		},
		{
			name: "Home",
			args: args{
				key: ebiten.KeyHome,
			},
			want: ti.Home,
		},
		{
			name: "Emd",
			args: args{
				key: ebiten.KeyEnd,
			},
			want: ti.End,
		},
		{
			name: "Delete",
			args: args{
				key: ebiten.KeyDelete,
			},
			want: ti.Delete,
		},
		{
			name: "Backspace",
			args: args{
				key: ebiten.KeyBackspace,
			},
			want: ti.Backspace,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ti.pressedKeyHandler(tt.args.key)

			wantName := runtime.FuncForPC(reflect.ValueOf(tt.want).Pointer()).Name()
			gotName := runtime.FuncForPC(reflect.ValueOf(got).Pointer()).Name()

			if gotName != wantName {
				t.Errorf("TextInput.pressedKeyHandler() = %T, want %T", got, tt.want)
			}
		})
	}
}
