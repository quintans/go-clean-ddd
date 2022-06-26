package command

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/internal/app"
	"github.com/quintans/go-clean-ddd/internal/domain/outbox"
	"github.com/quintans/go-clean-ddd/internal/domain/registration"
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

func (f FlushOutbox) Handle(ctx context.Context) error {
	for {
		err := f.outboxRepository.Consume(ctx, func(events []outbox.Outbox) error {
			for _, o := range events {
				switch o.Kind() {
				case registration.EventRegistered:
					if err := f.handleEventNewRegistration(ctx, o); err != nil {
						return faults.Wrap(err)
					}
				default:
					return errors.New("unknown event in outbox: " + o.Kind())
				}
			}
			return nil
		})
		if errors.Is(err, app.ErrNotFound) {
			return nil
		}
	}
}

func (f FlushOutbox) handleEventNewRegistration(ctx context.Context, o outbox.Outbox) error {
	event := registration.RegisteredEvent{}
	err := json.Unmarshal(o.Payload(), &event)
	if err != nil {
		return faults.Wrap(err)
	}
	return f.publisher.Publish(ctx, app.NewRegistration{Id: event.Id, Email: event.Email})
}
