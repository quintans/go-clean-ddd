package usecase

import (
	"context"
	"errors"

	"github.com/quintans/go-clean-ddd/internal/domain/entity"
	"github.com/quintans/go-clean-ddd/internal/domain/vo"
)

type Publisher interface {
	Publish(ctx context.Context, event NewRegistration) error
}

type NewRegistration struct {
	Id    string
	Email vo.Email
}

// RegistrationRepository interface for handling the persistence of RegistrationRepository
type RegistrationRepository interface {
	Create(context.Context, entity.Registration) error
	Update(context.Context, string, func(context.Context, *entity.Registration) error) error
}

type OutboxRepository interface {
	Create(ctx context.Context, ob entity.Outbox) error
	Consume(ctx context.Context, handler func([]entity.Outbox) error) error
}

var (
	ErrNotFound          = errors.New("model was not found")
	ErrOptimisticLocking = errors.New("optimistic locking failure")
)

// CustomerRepository interface for handling the persistence of Customer
type CustomerRepository interface {
	Create(context.Context, entity.Customer) error
	Update(context.Context, vo.CustomerID, func(context.Context, *entity.Customer) error) error
}

// CustomerViewRepository interface for customer reads (queries)
type CustomerViewRepository interface {
	GetAll(context.Context) ([]entity.Customer, error)
	GetByEmail(ctx context.Context, email vo.Email) (entity.Customer, error)
}
