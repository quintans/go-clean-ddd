package app

import (
	"context"
	"errors"

	"github.com/quintans/go-clean-ddd/internal/domain"
	"github.com/quintans/go-clean-ddd/internal/domain/customer"
	"github.com/quintans/go-clean-ddd/internal/domain/registration"
)

type EmailSender interface {
	Send(ctx context.Context, destination domain.Email, body string) error
}

type NewRegistration struct {
	Id    string
	Email domain.Email
}

// RegistrationRepository interface for handling the persistence of RegistrationRepository
type RegistrationRepository interface {
	Create(context.Context, registration.Registration) error
	Update(context.Context, string, func(context.Context, *registration.Registration) error) error
}

var (
	ErrNotFound          = errors.New("model was not found")
	ErrOptimisticLocking = errors.New("optimistic locking failure")
)

// CustomerRepository interface for handling the persistence of Customer
type CustomerRepository interface {
	Create(context.Context, customer.Customer) error
	Update(context.Context, customer.CustomerID, func(context.Context, *customer.Customer) error) error
}

// CustomerViewRepository interface for customer reads (queries)
type CustomerViewRepository interface {
	GetAll(context.Context) ([]customer.Customer, error)
	GetByEmail(ctx context.Context, email domain.Email) (customer.Customer, error)
}
