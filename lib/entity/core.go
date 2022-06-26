package entity

import "github.com/quintans/go-clean-ddd/lib/event"

type Core struct {
	events []event.DomainEvent
}

func (c *Core) AddEvent(e event.DomainEvent) {
	c.events = append(c.events, e)
}

func (c *Core) PopEvents() []event.DomainEvent {
	events := c.events
	c.events = nil
	return events
}
