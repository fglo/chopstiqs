package component

import (
	"testing"

	"github.com/fglo/chopstiqs/event"
	"github.com/matryer/is"
)

func TestButton_Clicked(t *testing.T) {
	is := is.New(t)

	firedEventsCounter := 0

	eventManager := event.NewManager()

	b := NewButton(&ButtonOptions{})
	b.SetEventManager(eventManager)
	b.AddClickedHandler(func(args *ButtonClickedEventArgs) {
		firedEventsCounter++
	})

	leftMouseButtonClick(t, &b.component)
	is.Equal(firedEventsCounter, 1)

	leftMouseButtonPress(t, &b.component)
	is.Equal(firedEventsCounter, 1)

	leftMouseButtonRelease(t, &b.component)
	is.Equal(firedEventsCounter, 2)
}

func TestButton_Pressed(t *testing.T) {
	is := is.New(t)

	firedEventsCounter := 0

	eventManager := event.NewManager()

	b := NewButton(&ButtonOptions{})
	b.SetEventManager(eventManager)
	b.AddPressedHandler(func(args *ButtonPressedEventArgs) {
		firedEventsCounter++
	})

	leftMouseButtonClick(t, &b.component)
	is.Equal(firedEventsCounter, 1)

	leftMouseButtonPress(t, &b.component)
	is.Equal(firedEventsCounter, 2)

	leftMouseButtonRelease(t, &b.component)
	is.Equal(firedEventsCounter, 2)
}

func TestButton_Released(t *testing.T) {
	is := is.New(t)

	firedEventsCounter := 0

	eventManager := event.NewManager()

	b := NewButton(&ButtonOptions{})
	b.SetEventManager(eventManager)
	b.AddReleasedHandler(func(args *ButtonReleasedEventArgs) {
		firedEventsCounter++
	})

	leftMouseButtonClick(t, &b.component)
	is.Equal(firedEventsCounter, 1)

	leftMouseButtonPress(t, &b.component)
	is.Equal(firedEventsCounter, 1)

	leftMouseButtonRelease(t, &b.component)
	is.Equal(firedEventsCounter, 2)
}
