package component

import (
	"testing"

	"github.com/fglo/chopstiqs/pkg/event"
	"github.com/matryer/is"
)

func TestCheckbox_Click(t *testing.T) {
	is := is.New(t)

	firedEventsCounter := 0

	cb := NewCheckBox(&CheckBoxOptions{})
	cb.AddToggledHandler(func(args *CheckBoxToggledEventArgs) {
		firedEventsCounter++
	})

	leftMouseButtonClick(t, &cb.component)
	is.Equal(firedEventsCounter, 1)

	leftMouseButtonPress(t, &cb.component)
	is.Equal(firedEventsCounter, 1)

	leftMouseButtonRelease(t, &cb.component)
	is.Equal(firedEventsCounter, 2)
}

func TestCheckbox_Toggle(t *testing.T) {
	is := is.New(t)

	cb := NewCheckBox(&CheckBoxOptions{})

	leftMouseButtonClick(t, &cb.component)
	is.Equal(cb.Checked(), true)

	leftMouseButtonClick(t, &cb.component)
	is.Equal(cb.Checked(), false)

	leftMouseButtonClick(t, &cb.component)
	is.Equal(cb.Checked(), true)
}

func TestCheckbox_SetChecked(t *testing.T) {
	is := is.New(t)

	firedEventsCounter := 0

	cb := NewCheckBox(&CheckBoxOptions{})
	cb.AddToggledHandler(func(args *CheckBoxToggledEventArgs) {
		firedEventsCounter++
	})

	cb.Set(false)
	event.HandleFired()

	is.Equal(cb.Checked(), false)
	is.Equal(firedEventsCounter, 0)

	cb.Set(true)
	event.HandleFired()

	is.Equal(cb.Checked(), true)
	is.Equal(firedEventsCounter, 1)

	cb.Set(false)
	event.HandleFired()

	is.Equal(cb.Checked(), false)
	is.Equal(firedEventsCounter, 2)
}
