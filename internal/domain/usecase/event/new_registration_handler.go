package event

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
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
		return errors.Errorf("NewRegistrationHandler cannot handle event of type %T", e)
	}
}

func (h NewRegistrationHandler) handleNewRegistration(ctx context.Context, e devent.NewRegistration) error {
	fmt.Println("===> fake send email")
	return nil
}
