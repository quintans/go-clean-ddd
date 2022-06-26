package command

import (
	"context"

	"github.com/quintans/go-clean-ddd/internal/domain/entity"
	"github.com/quintans/go-clean-ddd/internal/domain/usecase"
	"github.com/quintans/go-clean-ddd/internal/domain/vo"
)

type CreateRegistrationCommand struct {
	Email string
}

type CreateRegistrationHandler interface {
	Handle(context.Context, CreateRegistrationCommand) (string, error)
}

type CreateRegistration struct {
	registrationRepository usecase.RegistrationRepository
	customerView           usecase.CustomerViewRepository
}

func NewCreateRegistration(registrationRepository usecase.RegistrationRepository, customerView usecase.CustomerViewRepository) CreateRegistration {
	return CreateRegistration{
		registrationRepository: registrationRepository,
		customerView:           customerView,
	}
}

func (c CreateRegistration) Handle(ctx context.Context, cmd CreateRegistrationCommand) (string, error) {
	email, err := vo.NewEmail(cmd.Email)
	if err != nil {
		return "", err
	}
	r, err := entity.NewRegistration(ctx, email, usecase.NewUniquenessPolicy(c.customerView))
	if err != nil {
		return "", err
	}

	err = c.registrationRepository.Create(ctx, r)
	if err != nil {
		return "", err
	}
	return r.ID(), nil
}
