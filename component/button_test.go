package component

import (
	"testing"

	"github.com/fglo/chopstiqs/event"
	"github.com/fglo/chopstiqs/option"
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

func TestButton_InitialState(t *testing.T) {
	is := is.New(t)

	b := NewButton(&ButtonOptions{})

	is.Equal(b.Disable(), false)
	is.Equal(b.Hidden(), false)
	is.Equal(b.Focused(), false)
}

func TestButton_Disabled_NoClickedEvent(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()

	b := NewButton(&ButtonOptions{})
	b.SetEventManager(eventManager)
	b.SetDisabled(true)

	fired := 0
	b.AddClickedHandler(func(args *ButtonClickedEventArgs) {
		fired++
	})

	leftMouseButtonClick(t, &b.component)
	is.Equal(fired, 0)
}

func TestButton_Disabled_NoPressedEvent(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()

	b := NewButton(&ButtonOptions{})
	b.SetEventManager(eventManager)
	b.SetDisabled(true)

	fired := 0
	b.AddPressedHandler(func(args *ButtonPressedEventArgs) {
		fired++
	})

	leftMouseButtonClick(t, &b.component)
	is.Equal(fired, 0)
}

func TestButton_Disabled_NoReleasedEvent(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()

	b := NewButton(&ButtonOptions{})
	b.SetEventManager(eventManager)
	b.SetDisabled(true)

	fired := 0
	b.AddReleasedHandler(func(args *ButtonReleasedEventArgs) {
		fired++
	})

	leftMouseButtonClick(t, &b.component)
	is.Equal(fired, 0)
}

func TestButton_WithOptions(t *testing.T) {
	is := is.New(t)

	b := NewButton(&ButtonOptions{
		Width:               option.Int(100),
		Height:              option.Int(50),
		HorizontalAlignment: option.AlignmentCenteredHorizontally,
		VerticalAlignment:   option.AlignmentCenteredVertically,
	})

	is.Equal(b.Width(), 100)
	is.Equal(b.Height(), 50)
}
