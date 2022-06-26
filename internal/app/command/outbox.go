package command

import (
	"context"
	"errors"

	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/internal/app"
)

type FlushOutboxHandler interface {
	Handle(context.Context) error
}

type FlushOutbox struct {
	outboxRepository app.OutboxRepository
	publisher        app.Publisher
}

func NewFlushOutbox(outboxRepository app.OutboxRepository) FlushOutbox {
	return FlushOutbox{
		outboxRepository: outboxRepository,
	}
}

// Handle handles the events in the outbox left to be consumed by publishing them
func (f FlushOutbox) Handle(ctx context.Context) error {
	for {
		err := f.outboxRepository.Consume(ctx, func(events []*app.Outbox) error {
			for _, e := range events {
				err := f.publisher.Publish(ctx, app.Event{
					Kind:    e.Kind,
					Payload: e.Payload,
				})
				if err != nil {
					return faults.Wrap(err)
				}
			}
			return nil
		})
		if errors.Is(err, app.ErrNotFound) {
			return nil
		}
	}
}
