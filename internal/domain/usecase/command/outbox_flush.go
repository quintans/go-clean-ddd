package command

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/quintans/go-clean-ddd/internal/domain/entity"
	"github.com/quintans/go-clean-ddd/internal/domain/event"
	"github.com/quintans/go-clean-ddd/internal/domain/usecase"
)

type FlushOutboxHandler interface {
	Handle(context.Context) error
}

type FlushOutbox struct {
	outboxRepository usecase.OutboxRepository
	notifier         usecase.Notifier
}

func NewFlushOutbox(outboxRepository usecase.OutboxRepository) FlushOutbox {
	return FlushOutbox{
		outboxRepository: outboxRepository,
	}
}

func (f FlushOutbox) Handle(ctx context.Context) error {
	outboxes, err := f.outboxRepository.LockAndLoad(ctx)
	if errors.Is(err, usecase.ErrModelNotFound) {
		return nil
	}
	for _, o := range outboxes {
		switch o.Kind() {
		case event.EventNewRegistration:
			if err := f.handleEventNewRegistration(ctx, o); err != nil {
				return err
			}
		default:
			return errors.New("unknown event in outbox: " + o.Kind())
		}
	}
	return nil
}

func (f FlushOutbox) handleEventNewRegistration(ctx context.Context, o entity.Outbox) error {
	event := event.NewRegistration{}
	err := json.Unmarshal(o.Payload(), &event)
	if err != nil {
		return err
	}
	return f.notifier.Confirm(ctx, event.Email.String(), event.Id)
}
