package service

import (
	"context"

	"github.com/quintans/go-clean-ddd/internal/domain/model"
)

// CustomerService allow us to do things that could not be done inside the entity.
// CustomerService is a domain service
type CustomerService interface {
	FindDuplicate(context.Context, string) (*model.Customer, error)
}

type CustomerServiceImpl struct {
	CustomerRepository model.CustomerRepository
}

// FindDuplicate finds a customer with the same email
// The service can be implemented strictly using the domain layer,
// so both the interface and the implementation are part of the domain layer.
// If the service requires to access esternal resources, the implementation could part of the infrastructure layer.
func (s CustomerServiceImpl) FindDuplicate(ctx context.Context, email string) (*model.Customer, error) {
	return s.CustomerRepository.FindByEmail(ctx, email)
}
