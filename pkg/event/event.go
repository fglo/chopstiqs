package event

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

// A HandlerFunc is a function that receives and handles an event. When firing an event using
// EventManager.Fire, arbitrary event arguments may be passed that are in turn passed on to the handler function.
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

// AddOneTimeHandler registers event handler with event. When event fires, handler is removed from it immediately.
func (e *Event) AddOneTimeHandler(handler HandlerFunc) {
	var removeHandler RemoveHandlerFunc

	oneShotHandlerWrapperFunc := func(args interface{}) {
		removeHandler()
		handler(args)
	}

	removeHandler = e.AddHandler(oneShotHandlerWrapperFunc)
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
