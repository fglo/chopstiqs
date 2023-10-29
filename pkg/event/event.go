package event

import "fmt"

// Event represents an event that can be fired.
type Event struct {
	idCounter uint32
	handlers  []handler
}

// handler represents a handler that is registered with an event. It contains the handler function and the
// unique id of the handler. This id is used to remove the handler from the event.
type handler struct {
	id     uint32
	handle HandlerFunc
}

// Fire fires an event to all registered handlers. Arbitrary event arguments may be passed
// which are in turn passed on to event handlers.
//
// Events are not fired directly, but are put into a deferred queue. This queue is then
// processed by the GUI.
func (e *Event) Fire(args interface{}) {
	firedEvents = append(firedEvents, &FiredEvent{
		event: e,
		args:  args,
	})
}

// A HandlerFunc is a function that receives and handles an event. When firing an event using
// Event.Fire, arbitrary event arguments may be passed that are in turn passed on to the handler function.
type HandlerFunc func(args interface{})

// RemoveHandlerFunc is a function that removes a handler from an event.
type RemoveHandlerFunc func()

// AddHandler registers event handler with event. It returns a function to remove handler from event.
func (e *Event) AddHandler(h HandlerFunc) RemoveHandlerFunc {
	e.idCounter++

	id := e.idCounter

	e.handlers = append(e.handlers, handler{
		id:     id,
		handle: h,
	})

	return func() {
		e.removeHandler(id)
	}
}

func (e *Event) removeHandler(id uint32) {
	index := -1
	for i, handler := range e.handlers {
		if handler.id == id {
			index = i
			break
		}
	}

	if index < 0 {
		return
	}

	e.handlers = append(e.handlers[:index], e.handlers[index+1:]...)
}

var firedEvents []*FiredEvent

// FiredEvent represents an event that has been fired.
type FiredEvent struct {
	event *Event
	args  interface{}
}

// AddEventHandlerOneShot registers event handler with event. When event fires, handler is removed from it immediately.
func AddEventHandlerOneShot(event *Event, handler HandlerFunc) {
	var removeHandler RemoveHandlerFunc

	oneShotHandlerWrapperFunc := func(args interface{}) {
		removeHandler()
		handler(args)
	}

	removeHandler = event.AddHandler(oneShotHandlerWrapperFunc)
}

// HandleFired processes the queue of fired events and calls their handlers.
func HandleFired() {
	for len(firedEvents) > 0 {
		fired := firedEvents[0]
		firedEvents = firedEvents[1:]

		if fired == nil {
			fmt.Println("test")
		}

		if fired.event == nil {
			fmt.Println("test2")
		}

		for _, handler := range fired.event.handlers {
			handler.handle(fired.args)
		}
	}

	// resetting the deferredActions slice
	firedEvents = firedEvents[:0]
}
