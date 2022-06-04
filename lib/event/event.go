package event

import "context"

type DomainEvent interface {
	Kind() string
}

type Handler interface {
	Accept(DomainEvent) bool
	Handle(context.Context, DomainEvent) error
}

type EventBuser interface {
	AddHandler(Handler)
	Fire(context.Context, ...DomainEvent) error
}

type EventBus struct {
	handlers []Handler
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: []Handler{},
	}
}

func (m *EventBus) AddHandler(handler Handler) {
	m.handlers = append(m.handlers, handler)
}

func (m EventBus) Fire(ctx context.Context, events ...DomainEvent) error {
	for _, e := range events {
		for _, h := range m.handlers {
			if !h.Accept(e) {
				continue
			}

			if err := h.Handle(ctx, e); err != nil {
				return nil
			}
		}
	}
	return nil
}
