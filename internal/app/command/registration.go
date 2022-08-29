package command

import (
	"context"
	"errors"

	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/internal/app"
	"github.com/quintans/go-clean-ddd/internal/domain"
	"github.com/quintans/go-clean-ddd/internal/domain/customer"
	"github.com/quintans/go-clean-ddd/internal/domain/registration"
	"github.com/quintans/go-clean-ddd/lib/transaction"
)

type CreateRegistrationHandler interface {
	Handle(context.Context, CreateRegistrationCommand) (string, error)
}

type CreateRegistrationCommand struct {
	Email string
}

type CreateRegistration struct {
	registrationRepository app.RegistrationRepository
	customerView           app.CustomerViewRepository
}

func NewCreateRegistration(registrationRepository app.RegistrationRepository, customerView app.CustomerViewRepository) CreateRegistration {
	return CreateRegistration{
		registrationRepository: registrationRepository,
		customerView:           customerView,
	}
}

func (c CreateRegistration) Handle(ctx context.Context, cmd CreateRegistrationCommand) (string, error) {
	email, err := domain.NewEmail(cmd.Email)
	if err != nil {
		return "", faults.Wrap(err)
	}
	r, err := registration.NewRegistration(ctx, email, app.NewUniquenessPolicy(c.customerView))
	if err != nil {
		return "", faults.Wrap(err)
	}

	err = c.registrationRepository.Create(ctx, r)
	if err != nil {
		return "", faults.Wrap(err)
	}
	return r.ID(), nil
}

type ConfirmRegistrationCommand struct {
	Id string
}

type ConfirmRegistrationHandler interface {
	Handle(context.Context, ConfirmRegistrationCommand) error
}

type ConfirmRegistration struct {
	uowManager transaction.UnitOfWorkManager

	registrationRepository app.RegistrationRepository

	customerRepository app.CustomerRepository
	customerView       app.CustomerViewRepository
}

func NewConfirmRegistration(
	uowManager transaction.UnitOfWorkManager,
	registrationRepository app.RegistrationRepository,
	customerRepository app.CustomerRepository,
	customerView app.CustomerViewRepository,
) ConfirmRegistration {
	return ConfirmRegistration{
		uowManager:             uowManager,
		registrationRepository: registrationRepository,
		customerRepository:     customerRepository,
		customerView:           customerView,
	}
}

func (h ConfirmRegistration) Handle(ctx context.Context, cmd ConfirmRegistrationCommand) error {
	err := h.uowManager.Current(ctx, func(ctx context.Context) error {
		return h.handle(ctx, cmd)
	})
	if errors.Is(err, domain.ErrNoChange) {
		return nil
	}

	return faults.Wrap(err)
}

func (h ConfirmRegistration) handle(ctx context.Context, cmd ConfirmRegistrationCommand) error {
	r, err := h.registrationRepository.Update(ctx, cmd.Id, func(ctx context.Context, r *registration.Registration) error {
		return r.Verify()
	})
	if err != nil {
		return faults.Wrap(err)
	}

	customer, err := customer.NewCustomer(ctx, r.Email(), app.NewUniquenessPolicy(h.customerView))
	if err != nil {
		return faults.Wrap(err)
	}

	h.customerRepository.Create(ctx, customer)
	return faults.Wrap(err)
}
