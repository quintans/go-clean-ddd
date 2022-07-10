package entity

import "github.com/quintans/go-clean-ddd/lib/eventbus"

type Core struct {
	events []eventbus.DomainEvent
}

func (c *Core) AddEvent(e eventbus.DomainEvent) {
	c.events = append(c.events, e)
}

func (c *Core) PopEvents() []eventbus.DomainEvent {
	events := c.events
	c.events = nil
	return events
}
