package event

// Manager contains queue of fired events.
// Its role is to handle fire events.
type Manager struct {
	firedEvents []*Fired
}

// Fired represents an event that has been fired.
type Fired struct {
	event *Event
	args  interface{}
}

func NewManager() *Manager {
	return &Manager{
		firedEvents: make([]*Fired, 0),
	}
}

// Fire fires an event to all registered handlers. Arbitrary event arguments may be passed
// which are in turn passed on to event handlers.
//
// Events are not fired directly, but are put into a deferred queue. This queue is then
// processed by the GUI.
func (m *Manager) Fire(e *Event, args any) {
	if m == nil {
		return
	}

	m.firedEvents = append(m.firedEvents, &Fired{
		event: e,
		args:  args,
	})
}

// HandleFired processes the queue of fired events and calls their handlers.
func (m *Manager) HandleFired() {
	if m == nil {
		return
	}

	for len(m.firedEvents) > 0 {
		fired := m.firedEvents[0]
		m.firedEvents = m.firedEvents[1:]

		for _, handler := range fired.event.handlers {
			handler.handle(fired.args)
		}
	}

	// resetting the deferredActions slice
	m.firedEvents = m.firedEvents[:0]
}
