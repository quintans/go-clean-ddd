package eventbus

import "context"

type DomainEvent interface {
	Kind() string
}

type Handler func(context.Context, DomainEvent) error

type Subscriber interface {
	Subscribe(string, Handler)
}
type Publisher interface {
	Publish(context.Context, ...DomainEvent) error
}

type EventBus struct {
	handlers map[string][]Handler
}

func New() *EventBus {
	return &EventBus{
		handlers: map[string][]Handler{},
	}
}

func (m *EventBus) Subscribe(kind string, handler Handler) {
	handlers := m.handlers[kind]
	m.handlers[kind] = append(handlers, handler)
}

func (m EventBus) Publish(ctx context.Context, events ...DomainEvent) error {
	for _, e := range events {
		handlers := m.handlers[e.Kind()]
		for _, h := range handlers {
			if err := h(ctx, e); err != nil {
				return nil
			}
		}
	}
	return nil
}
