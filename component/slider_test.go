package component

import (
	"testing"

	"github.com/fglo/chopstiqs/event"
	"github.com/fglo/chopstiqs/option"
	"github.com/matryer/is"
)

func newTestSlider(min, max, step float64) *Slider {
	eventManager := event.NewManager()
	s := NewSlider(&SliderOptions{
		Min:  option.Float(min),
		Max:  option.Float(max),
		Step: option.Float(step),
	})
	s.SetEventManager(eventManager)
	return s
}

func TestSlider_NewDefaults(t *testing.T) {
	is := is.New(t)

	s := NewSlider(nil)

	is.Equal(s.GetValue(), 0.)
}

func TestSlider_SetValue(t *testing.T) {
	is := is.New(t)

	s := newTestSlider(0, 100, 1)

	fired := 0
	s.AddSlidedHandler(func(args *SliderSlidedEventArgs) {
		fired++
	})

	s.SetEventManager(event.NewManager())
	s.Set(50)
	s.eventManager.HandleFired()

	is.Equal(s.GetValue(), 50.)
	is.Equal(fired, 1)
}

func TestSlider_SetToMin(t *testing.T) {
	is := is.New(t)

	s := newTestSlider(0, 100, 1)
	s.SetEventManager(event.NewManager())

	s.Set(50)
	s.eventManager.HandleFired()

	fired := 0
	s.AddSlidedHandler(func(args *SliderSlidedEventArgs) {
		fired++
	})

	s.SetToMin()
	s.eventManager.HandleFired()

	is.Equal(s.GetValue(), 0.)
	is.Equal(fired, 1)
}

func TestSlider_SetToMax(t *testing.T) {
	is := is.New(t)

	s := newTestSlider(0, 100, 1)
	s.SetEventManager(event.NewManager())

	fired := 0
	s.AddSlidedHandler(func(args *SliderSlidedEventArgs) {
		fired++
	})

	s.SetToMax()
	s.eventManager.HandleFired()

	is.Equal(s.GetValue(), 100.)
	is.Equal(fired, 1)
}

func TestSlider_SlidedEvent_NoChange(t *testing.T) {
	is := is.New(t)

	s := newTestSlider(0, 100, 1)
	s.SetEventManager(event.NewManager())

	fired := 0
	s.AddSlidedHandler(func(args *SliderSlidedEventArgs) {
		fired++
	})

	s.Set(0)
	s.eventManager.HandleFired()

	is.Equal(fired, 0)
}

func TestSlider_PressedEvent(t *testing.T) {
	is := is.New(t)

	s := newTestSlider(0, 100, 1)
	s.SetEventManager(event.NewManager())

	fired := 0
	s.PressedEvent.AddHandler(func(args any) {
		fired++
	})

	s.eventManager.Fire(s.PressedEvent, &SliderPressedEventArgs{Slider: s})
	s.eventManager.HandleFired()

	is.Equal(fired, 1)
}

func TestSlider_ReleasedEvent(t *testing.T) {
	is := is.New(t)

	s := newTestSlider(0, 100, 1)
	s.SetEventManager(event.NewManager())

	fired := 0
	s.ReleasedEvent.AddHandler(func(args any) {
		fired++
	})

	s.eventManager.Fire(s.ReleasedEvent, &SliderReleasedEventArgs{Slider: s, Inside: true})
	s.eventManager.HandleFired()

	is.Equal(fired, 1)
}

func TestSlider_ClickedEvent(t *testing.T) {
	is := is.New(t)

	s := newTestSlider(0, 100, 1)
	s.SetEventManager(event.NewManager())

	fired := 0
	s.ClickedEvent.AddHandler(func(args any) {
		fired++
	})

	s.eventManager.Fire(s.ClickedEvent, &SliderClickedEventArgs{Slider: s})
	s.eventManager.HandleFired()

	is.Equal(fired, 1)
}

func TestSlider_Disabled_NoEvents(t *testing.T) {
	is := is.New(t)

	s := newTestSlider(0, 100, 1)
	s.SetEventManager(event.NewManager())
	s.SetDisabled(true)

	fired := 0
	s.AddSlidedHandler(func(args *SliderSlidedEventArgs) {
		fired++
	})

	s.Set(50)
	s.eventManager.HandleFired()

	is.Equal(fired, 1)
}

func TestSlider_SetDisabled_CascadesToHandle(t *testing.T) {
	is := is.New(t)

	s := newTestSlider(0, 100, 1)
	s.SetEventManager(event.NewManager())

	s.SetDisabled(true)
	is.Equal(s.handle.Disable(), true)

	s.SetDisabled(false)
	is.Equal(s.handle.Disable(), false)
}

func TestSlider_GetValue(t *testing.T) {
	is := is.New(t)

	s := newTestSlider(0, 100, 1)
	s.SetEventManager(event.NewManager())

	s.Set(42)
	is.Equal(s.GetValue(), 42.)
}

func TestSlider_SlidedEventArgs(t *testing.T) {
	is := is.New(t)

	s := newTestSlider(0, 100, 1)
	s.SetEventManager(event.NewManager())

	var receivedValue float64
	var receivedChange float64

	s.AddSlidedHandler(func(args *SliderSlidedEventArgs) {
		receivedValue = args.Value
		receivedChange = args.Change
	})

	s.Set(50)
	s.eventManager.HandleFired()

	is.Equal(receivedValue, 50.)
	is.Equal(receivedChange, 50.)
}

func TestSlider_DefaultMin(t *testing.T) {
	is := is.New(t)

	s := newTestSlider(10, 100, 1)

	is.Equal(s.GetValue(), 10.)
}

func TestSlider_SetPosition_UpdatesHandle(t *testing.T) {
	is := is.New(t)

	s := newTestSlider(0, 100, 1)
	s.SetEventManager(event.NewManager())

	s.SetPosition(50, 60)

	absX, absY := s.handle.AbsPosition()
	is.True(absX >= 50)
	is.True(absY >= 60)
}

func TestSlider_FireEvents(t *testing.T) {
	is := is.New(t)

	s := newTestSlider(0, 100, 1)
	s.SetEventManager(event.NewManager())

	fired := 0
	s.AddSlidedHandler(func(args *SliderSlidedEventArgs) {
		fired++
	})

	s.Set(10)
	s.eventManager.HandleFired()

	s.FireEvents()

	is.Equal(fired, 1)
}
