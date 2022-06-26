package event

import "context"

type DomainEvent interface {
	Kind() string
}

type Handler interface {
	Handle(context.Context, DomainEvent) error
}

type EventBuser interface {
	AddHandler(string, Handler)
	Fire(context.Context, ...DomainEvent) error
}

type EventBus struct {
	handlers map[string][]Handler
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: map[string][]Handler{},
	}
}

func (m *EventBus) AddHandler(kind string, handler Handler) {
	handlers := m.handlers[kind]
	m.handlers[kind] = append(handlers, handler)
}

func (m EventBus) Fire(ctx context.Context, events ...DomainEvent) error {
	for _, e := range events {
		handlers := m.handlers[e.Kind()]
		for _, h := range handlers {
			if err := h.Handle(ctx, e); err != nil {
				return nil
			}
		}
	}
	return nil
}
