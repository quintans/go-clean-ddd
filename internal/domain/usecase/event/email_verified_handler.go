package event

import (
	"context"

	"github.com/pkg/errors"
	"github.com/quintans/go-clean-ddd/internal/domain/entity"
	devent "github.com/quintans/go-clean-ddd/internal/domain/event"
	"github.com/quintans/go-clean-ddd/internal/usecase"
	"github.com/quintans/go-clean-ddd/lib/event"
)

type EmailVerifiedHandler struct {
	customerRepository usecase.CustomerRepository
	customerView       usecase.CustomerViewRepository
}

func NewEmailVerifiedHandler(customerRepository usecase.CustomerRepository, customerView usecase.CustomerViewRepository) EmailVerifiedHandler {
	return EmailVerifiedHandler{
		customerRepository: customerRepository,
		customerView:       customerView,
	}
}

func (h EmailVerifiedHandler) Accept(e event.DomainEvent) bool {
	return e.Kind() == devent.EventEmailVerified
}

func (h EmailVerifiedHandler) Handle(ctx context.Context, e event.DomainEvent) error {
	switch t := e.(type) {
	case devent.EmailVerified:
		return h.handleEmailVerified(ctx, t)
	default:
		return errors.Errorf("EmailVerifiedHandler cannot handle event of type %T", e)
	}
}

func (h EmailVerifiedHandler) handleEmailVerified(ctx context.Context, e devent.EmailVerified) error {
	customer, err := entity.NewCustomer(ctx, e.Email, usecase.NewUniquenessPolicy(h.customerView))
	if err != nil {
		return err
	}

	if err := h.customerRepository.Save(ctx, customer); err != nil {
		return err
	}
	return nil
}
