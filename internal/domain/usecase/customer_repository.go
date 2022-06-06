package usecase

import (
	"context"
	"errors"

	"github.com/quintans/go-clean-ddd/internal/domain/entity"
	"github.com/quintans/go-clean-ddd/internal/domain/vo"
)

var ErrModelNotFound = errors.New("model was not found")

// CustomerRepository interface for handling the persistence of Customer
type CustomerRepository interface {
	Save(context.Context, entity.Customer) error
	Apply(context.Context, vo.CustomerID, func(context.Context, *entity.Customer) error) error
}

// CustomerViewRepository interface for customer reads (queries)
type CustomerViewRepository interface {
	GetAll(context.Context) ([]entity.Customer, error)
	GetByEmail(ctx context.Context, email vo.Email) (entity.Customer, error)
}
