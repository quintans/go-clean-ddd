package event

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/quintans/go-clean-ddd/internal/domain/entity"
	devent "github.com/quintans/go-clean-ddd/internal/domain/event"
	"github.com/quintans/go-clean-ddd/internal/domain/usecase"
	"github.com/quintans/go-clean-ddd/lib/event"
)

// NewRegistration

type NewRegistrationHandler struct {
	outboxRepository usecase.OutboxRepository
}

func NewNewRegistrationHandler(outboxRepository usecase.OutboxRepository) NewRegistrationHandler {
	return NewRegistrationHandler{
		outboxRepository: outboxRepository,
	}
}

func (h NewRegistrationHandler) Accept(e event.DomainEvent) bool {
	return e.Kind() == devent.EventNewRegistration
}

func (h NewRegistrationHandler) Handle(ctx context.Context, e event.DomainEvent) error {
	switch t := e.(type) {
	case devent.NewRegistration:
		return h.handleNewRegistration(ctx, t)
	default:
		return errors.Errorf("EmailVerifiedHandler cannot handle event of type %T", e)
	}
}

func (h NewRegistrationHandler) handleNewRegistration(ctx context.Context, e devent.NewRegistration) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	ob, err := entity.NewOutbox(e.Kind(), b)
	if err != nil {
		return err
	}

	return h.outboxRepository.Save(ctx, ob)
}
