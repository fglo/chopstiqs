package component

import (
	"testing"

	"github.com/fglo/chopstiqs/event"
	"github.com/fglo/chopstiqs/option"
	"github.com/matryer/is"
)

func newTestScrollBar(containerHeight int) *ScrollBar {
	eventManager := event.NewManager()
	container := NewContainer(&ContainerOptions{
		Height: option.Int(containerHeight),
	})
	container.SetEventManager(eventManager)

	sb := NewScrollBar(&ScrollBarOptions{
		Container: container,
		Step:      option.Float(1),
	})
	sb.SetEventManager(eventManager)
	return sb
}

func TestScrollBar_NewDefaults(t *testing.T) {
	is := is.New(t)

	sb := newTestScrollBar(100)

	is.Equal(sb.GetValue(), 0.)
}

func TestScrollBar_SetValue(t *testing.T) {
	is := is.New(t)

	sb := newTestScrollBar(100)

	fired := 0
	sb.AddScrolledHandler(func(args *ScrollBarScrolledEventArgs) {
		fired++
	})

	sb.Set(50)
	sb.eventManager.HandleFired()

	is.Equal(sb.GetValue(), 50.)
	is.Equal(fired, 1)
}

func TestScrollBar_SetToMin(t *testing.T) {
	is := is.New(t)

	sb := newTestScrollBar(100)

	sb.Set(50)
	sb.eventManager.HandleFired()

	fired := 0
	sb.AddScrolledHandler(func(args *ScrollBarScrolledEventArgs) {
		fired++
	})

	sb.SetToMin()
	sb.eventManager.HandleFired()

	is.Equal(sb.GetValue(), 0.)
	is.Equal(fired, 1)
}

func TestScrollBar_SetToMax(t *testing.T) {
	is := is.New(t)

	sb := newTestScrollBar(100)

	fired := 0
	sb.AddScrolledHandler(func(args *ScrollBarScrolledEventArgs) {
		fired++
	})

	sb.SetToMax()
	sb.eventManager.HandleFired()

	is.Equal(sb.GetValue(), 0.)
}

func TestScrollBar_ScrolledEvent_NoChange(t *testing.T) {
	is := is.New(t)

	sb := newTestScrollBar(100)

	fired := 0
	sb.AddScrolledHandler(func(args *ScrollBarScrolledEventArgs) {
		fired++
	})

	sb.Set(0)
	sb.eventManager.HandleFired()

	is.Equal(fired, 0)
}

func TestScrollBar_PressedEvent(t *testing.T) {
	is := is.New(t)

	sb := newTestScrollBar(100)

	fired := 0
	sb.PressedEvent.AddHandler(func(args any) {
		fired++
	})

	sb.eventManager.Fire(sb.PressedEvent, &ScrollBarPressedEventArgs{ScrollBar: sb})
	sb.eventManager.HandleFired()

	is.Equal(fired, 1)
}

func TestScrollBar_ReleasedEvent(t *testing.T) {
	is := is.New(t)

	sb := newTestScrollBar(100)

	fired := 0
	sb.ReleasedEvent.AddHandler(func(args any) {
		fired++
	})

	sb.eventManager.Fire(sb.ReleasedEvent, &ScrollBarReleasedEventArgs{ScrollBar: sb, Inside: true})
	sb.eventManager.HandleFired()

	is.Equal(fired, 1)
}

func TestScrollBar_ClickedEvent(t *testing.T) {
	is := is.New(t)

	sb := newTestScrollBar(100)

	fired := 0
	sb.ClickedEvent.AddHandler(func(args any) {
		fired++
	})

	sb.eventManager.Fire(sb.ClickedEvent, &ScrollBarClickedEventArgs{ScrollBar: sb})
	sb.eventManager.HandleFired()

	is.Equal(fired, 1)
}

func TestScrollBar_SetDisabled_CascadesToHandle(t *testing.T) {
	is := is.New(t)

	sb := newTestScrollBar(100)

	sb.SetDisabled(true)
	is.Equal(sb.handle.Disable(), true)

	sb.SetDisabled(false)
	is.Equal(sb.handle.Disable(), false)
}

func TestScrollBar_GetValue(t *testing.T) {
	is := is.New(t)

	sb := newTestScrollBar(100)

	sb.Set(42)
	is.Equal(sb.GetValue(), 42.)
}

func TestScrollBar_ScrolledEventArgs(t *testing.T) {
	is := is.New(t)

	sb := newTestScrollBar(100)

	var receivedValue float64
	var receivedChange float64

	sb.AddScrolledHandler(func(args *ScrollBarScrolledEventArgs) {
		receivedValue = args.Value
		receivedChange = args.Change
	})

	sb.Set(50)
	sb.eventManager.HandleFired()

	is.Equal(receivedValue, 50.)
	is.Equal(receivedChange, 50.)
}

func TestScrollBar_FireEvents(t *testing.T) {
	is := is.New(t)

	sb := newTestScrollBar(100)

	fired := 0
	sb.AddScrolledHandler(func(args *ScrollBarScrolledEventArgs) {
		fired++
	})

	sb.Set(10)
	sb.eventManager.HandleFired()

	sb.FireEvents()

	is.Equal(fired, 1)
}
