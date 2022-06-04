package usecase

import (
	"context"

	"github.com/quintans/go-clean-ddd/internal/domain/entity"
)

// RegistrationRepository interface for handling the persistence of RegistrationRepository
type RegistrationRepository interface {
	Save(context.Context, entity.Registration) error
	Apply(context.Context, string, func(context.Context, *entity.Registration) error) error
}
