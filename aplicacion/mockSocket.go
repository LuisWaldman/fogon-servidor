package aplicacion

type EmittedEvent struct {
	Event string
	Args  []any
}

type MockSocket struct {
	Events []EmittedEvent
}

func (m *MockSocket) Emit(ev string, args ...any) error {
	m.Events = append(m.Events, EmittedEvent{
		Event: ev,
		Args:  args,
	})
	return nil
}

func (m *MockSocket) UltimoEmitted() *EmittedEvent {
	if len(m.Events) == 0 {
		return nil
	}
	return &m.Events[len(m.Events)-1]
}

func (m *MockSocket) TieneMensaje(event string) bool {
	for _, e := range m.Events {
		if e.Event == event {
			return true
		}
	}
	return false
}

func (m *MockSocket) ListaEvents() []EmittedEvent {
	return m.Events
}
