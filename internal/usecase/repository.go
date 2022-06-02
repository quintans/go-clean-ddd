package usecase

import (
	"context"
	"errors"

	"github.com/quintans/go-clean-ddd/internal/domain"
)

type Apply func(context.Context, domain.Customer) (domain.Customer, error)

// CustomerRepository interface for handling the persistence of the Customer
type CustomerRepository interface {
	Save(context.Context, domain.Customer) error
	Apply(context.Context, domain.CustomerID, Apply) error
}

var ErrReadModelNotFound = errors.New("read model was not found")

// CustomerViewRepository interface for customer reads (queries)
type CustomerViewRepository interface {
	GetAll(context.Context) ([]domain.Customer, error)
	GetByEmail(ctx context.Context, email domain.Email) (domain.Customer, error)
}
