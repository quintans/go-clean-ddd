package command

import (
	"context"

	"github.com/quintans/go-clean-ddd/internal/domain/entity"
	"github.com/quintans/go-clean-ddd/internal/domain/usecase"
)

type ConfirmRegistrationCommand struct {
	Id string
}

type ConfirmRegistrationHandler interface {
	Handle(context.Context, ConfirmRegistrationCommand) error
}

type ConfirmRegistration struct {
	registrationRepository usecase.RegistrationRepository
}

func NewConfirmRegistration(registrationRepository usecase.RegistrationRepository) ConfirmRegistration {
	return ConfirmRegistration{
		registrationRepository: registrationRepository,
	}
}

func (h ConfirmRegistration) Handle(ctx context.Context, cmd ConfirmRegistrationCommand) error {
	return h.registrationRepository.Apply(ctx, cmd.Id, func(ctx context.Context, r *entity.Registration) error {
		r.Verify()
		return nil
	})
}
