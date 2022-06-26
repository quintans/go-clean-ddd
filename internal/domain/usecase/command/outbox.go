package command

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/internal/domain/entity"
	"github.com/quintans/go-clean-ddd/internal/domain/event"
	"github.com/quintans/go-clean-ddd/internal/domain/usecase"
)

type FlushOutboxHandler interface {
	Handle(context.Context) error
}

type FlushOutbox struct {
	outboxRepository usecase.OutboxRepository
	publisher        usecase.Publisher
}

func NewFlushOutbox(outboxRepository usecase.OutboxRepository) FlushOutbox {
	return FlushOutbox{
		outboxRepository: outboxRepository,
	}
}

func (f FlushOutbox) Handle(ctx context.Context) error {
	for {
		err := f.outboxRepository.Consume(ctx, func(events []entity.Outbox) error {
			for _, o := range events {
				switch o.Kind() {
				case event.EventNewRegistration:
					if err := f.handleEventNewRegistration(ctx, o); err != nil {
						return faults.Wrap(err)
					}
				default:
					return errors.New("unknown event in outbox: " + o.Kind())
				}
			}
			return nil
		})
		if errors.Is(err, usecase.ErrNotFound) {
			return nil
		}
	}
}

func (f FlushOutbox) handleEventNewRegistration(ctx context.Context, o entity.Outbox) error {
	event := event.NewRegistration{}
	err := json.Unmarshal(o.Payload(), &event)
	if err != nil {
		return faults.Wrap(err)
	}
	return f.publisher.Publish(ctx, usecase.NewRegistration{Id: event.Id, Email: event.Email})
}
