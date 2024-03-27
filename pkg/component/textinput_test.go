package component

import (
	"testing"
	"time"

	"github.com/fglo/chopstiqs/pkg/event"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/matryer/is"
)

func handleState(t *testing.T, ti *TextInput) {
	t.Helper()

	ti.state = ti.state(ti)
	ti.eventManager.HandleFired()
}

func TestTextInput_PressedEnter(t *testing.T) {
	is := is.New(t)

	firedEventsCounter := 0

	eventManager := event.NewManager()

	ti := NewTextInput(&TextInputOptions{})
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
}
