package component

import (
	"testing"

	"github.com/fglo/chopstiqs/event"
	"github.com/fglo/chopstiqs/option"
	"github.com/matryer/is"
)

func TestCheckbox_Click(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()

	firedEventsCounter := 0

	cb := NewCheckBox(&CheckBoxOptions{})
	cb.SetEventManager(eventManager)
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

	eventManager := event.NewManager()

	cb := NewCheckBox(&CheckBoxOptions{})
	cb.SetEventManager(eventManager)
	leftMouseButtonClick(t, &cb.component)
	is.Equal(cb.Checked(), true)

	leftMouseButtonClick(t, &cb.component)
	is.Equal(cb.Checked(), false)

	leftMouseButtonClick(t, &cb.component)
	is.Equal(cb.Checked(), true)
}

func TestCheckbox_SetChecked(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()

	firedEventsCounter := 0

	cb := NewCheckBox(&CheckBoxOptions{})
	cb.SetEventManager(eventManager)
	cb.AddToggledHandler(func(args *CheckBoxToggledEventArgs) {
		firedEventsCounter++
	})

	cb.Set(false)
	eventManager.HandleFired()

	is.Equal(cb.Checked(), false)
	is.Equal(firedEventsCounter, 0)

	cb.Set(true)
	eventManager.HandleFired()

	is.Equal(cb.Checked(), true)
	is.Equal(firedEventsCounter, 1)

	cb.Set(false)
	eventManager.HandleFired()

	is.Equal(cb.Checked(), false)
	is.Equal(firedEventsCounter, 2)
}

func TestCheckbox_InitialState(t *testing.T) {
	is := is.New(t)

	cb := NewCheckBox(&CheckBoxOptions{})

	is.Equal(cb.Checked(), false)
	is.Equal(cb.Disable(), false)
	is.Equal(cb.Hidden(), false)
}

func TestCheckbox_Disabled_NoToggle(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()

	cb := NewCheckBox(&CheckBoxOptions{})
	cb.SetEventManager(eventManager)
	cb.SetDisabled(true)

	leftMouseButtonClick(t, &cb.component)
	is.Equal(cb.Checked(), false)
}

func TestCheckbox_Disabled_NoToggledEvent(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()

	cb := NewCheckBox(&CheckBoxOptions{})
	cb.SetEventManager(eventManager)
	cb.SetDisabled(true)

	fired := 0
	cb.AddToggledHandler(func(args *CheckBoxToggledEventArgs) {
		fired++
	})

	leftMouseButtonClick(t, &cb.component)
	is.Equal(fired, 0)
}

func TestCheckbox_SetChecked_AlreadyChecked_NoEvent(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()

	cb := NewCheckBox(&CheckBoxOptions{})
	cb.SetEventManager(eventManager)

	fired := 0
	cb.AddToggledHandler(func(args *CheckBoxToggledEventArgs) {
		fired++
	})

	cb.Set(true)
	eventManager.HandleFired()
	is.Equal(cb.Checked(), true)
	is.Equal(fired, 1)

	cb.Set(true)
	eventManager.HandleFired()
	is.Equal(cb.Checked(), true)
	is.Equal(fired, 1)
}

func TestCheckbox_WithOptions(t *testing.T) {
	is := is.New(t)

	cb := NewCheckBox(&CheckBoxOptions{
		Width:  option.Int(100),
		Height: option.Int(50),
	})

	is.Equal(cb.Width(), 100)
	is.Equal(cb.Height(), 50)
}
