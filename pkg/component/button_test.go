package component

import (
	"testing"

	"github.com/matryer/is"
)

func TestButton_Clicked(t *testing.T) {
	is := is.New(t)

	eventFiredCounter := 0

	cb := NewButton(&ButtonOptions{})
	cb.AddClickedHandler(func(args *ButtonClickedEventArgs) {
		eventFiredCounter++
	})

	leftMouseButtonClick(t, &cb.component)
	is.Equal(eventFiredCounter, 1)

	leftMouseButtonPress(t, &cb.component)
	is.Equal(eventFiredCounter, 1)

	leftMouseButtonRelease(t, &cb.component)
	is.Equal(eventFiredCounter, 2)
}

func TestButton_Pressed(t *testing.T) {
	is := is.New(t)

	eventFiredCounter := 0

	cb := NewButton(&ButtonOptions{})
	cb.AddPressedHandler(func(args *ButtonPressedEventArgs) {
		eventFiredCounter++
	})

	leftMouseButtonClick(t, &cb.component)
	is.Equal(eventFiredCounter, 1)

	leftMouseButtonPress(t, &cb.component)
	is.Equal(eventFiredCounter, 2)

	leftMouseButtonRelease(t, &cb.component)
	is.Equal(eventFiredCounter, 2)
}

func TestButton_Released(t *testing.T) {
	is := is.New(t)

	eventFiredCounter := 0

	cb := NewButton(&ButtonOptions{})
	cb.AddReleasedHandler(func(args *ButtonReleasedEventArgs) {
		eventFiredCounter++
	})

	leftMouseButtonClick(t, &cb.component)
	is.Equal(eventFiredCounter, 1)

	leftMouseButtonPress(t, &cb.component)
	is.Equal(eventFiredCounter, 1)

	leftMouseButtonRelease(t, &cb.component)
	is.Equal(eventFiredCounter, 2)
}
