package command

import (
	"context"

	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/internal/app"
	"github.com/quintans/go-clean-ddd/internal/domain"
	"github.com/quintans/go-clean-ddd/internal/domain/registration"
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
	registrationRepository app.RegistrationRepository
}

func NewConfirmRegistration(registrationRepository app.RegistrationRepository) ConfirmRegistration {
	return ConfirmRegistration{
		registrationRepository: registrationRepository,
	}
}

func (h ConfirmRegistration) Handle(ctx context.Context, cmd ConfirmRegistrationCommand) error {
	return h.registrationRepository.Update(ctx, cmd.Id, func(ctx context.Context, r *registration.Registration) error {
		r.Verify()
		return nil
	})
}
