package component

import (
	"testing"

	"github.com/matryer/is"
)

func TestButton_Clicked(t *testing.T) {
	is := is.New(t)

	firedEventsCounter := 0

	cb := NewButton(&ButtonOptions{})
	cb.AddClickedHandler(func(args *ButtonClickedEventArgs) {
		firedEventsCounter++
	})

	leftMouseButtonClick(t, &cb.component)
	is.Equal(firedEventsCounter, 1)

	leftMouseButtonPress(t, &cb.component)
	is.Equal(firedEventsCounter, 1)

	leftMouseButtonRelease(t, &cb.component)
	is.Equal(firedEventsCounter, 2)
}

func TestButton_Pressed(t *testing.T) {
	is := is.New(t)

	firedEventsCounter := 0

	cb := NewButton(&ButtonOptions{})
	cb.AddPressedHandler(func(args *ButtonPressedEventArgs) {
		firedEventsCounter++
	})

	leftMouseButtonClick(t, &cb.component)
	is.Equal(firedEventsCounter, 1)

	leftMouseButtonPress(t, &cb.component)
	is.Equal(firedEventsCounter, 2)

	leftMouseButtonRelease(t, &cb.component)
	is.Equal(firedEventsCounter, 2)
}

func TestButton_Released(t *testing.T) {
	is := is.New(t)

	firedEventsCounter := 0

	cb := NewButton(&ButtonOptions{})
	cb.AddReleasedHandler(func(args *ButtonReleasedEventArgs) {
		firedEventsCounter++
	})

	leftMouseButtonClick(t, &cb.component)
	is.Equal(firedEventsCounter, 1)

	leftMouseButtonPress(t, &cb.component)
	is.Equal(firedEventsCounter, 1)

	leftMouseButtonRelease(t, &cb.component)
	is.Equal(firedEventsCounter, 2)
}
