package command

import (
	"context"

	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/internal/app"
	"github.com/quintans/go-clean-ddd/internal/domain/customer"
	"github.com/quintans/go-clean-ddd/internal/domain/registration"
	"github.com/quintans/go-clean-ddd/lib/event"
)

type EmailVerifiedHandler struct {
	customerRepository app.CustomerRepository
	customerView       app.CustomerViewRepository
}

func NewEmailVerifiedHandler(customerRepository app.CustomerRepository, customerView app.CustomerViewRepository) EmailVerifiedHandler {
	return EmailVerifiedHandler{
		customerRepository: customerRepository,
		customerView:       customerView,
	}
}

func (h EmailVerifiedHandler) Handle(ctx context.Context, e event.DomainEvent) error {
	switch t := e.(type) {
	case registration.EmailVerified:
		return h.handleEmailVerified(ctx, t)
	default:
		return faults.Errorf("EmailVerifiedHandler cannot handle event of type %T", e)
	}
}

func (h EmailVerifiedHandler) handleEmailVerified(ctx context.Context, e registration.EmailVerified) error {
	customer, err := customer.NewCustomer(ctx, e.Email, app.NewUniquenessPolicy(h.customerView))
	if err != nil {
		return faults.Wrap(err)
	}

	if err := h.customerRepository.Create(ctx, customer); err != nil {
		return faults.Wrap(err)
	}
	return nil
}
