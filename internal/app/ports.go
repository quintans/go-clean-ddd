package app

import (
	"context"
	"errors"

	"github.com/quintans/go-clean-ddd/internal/domain"
	"github.com/quintans/go-clean-ddd/internal/domain/customer"
	"github.com/quintans/go-clean-ddd/internal/domain/registration"
)

type Publisher interface {
	Publish(ctx context.Context, event Event) error
}

type Event struct {
	Kind    string
	Payload []byte
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

type OutboxRepository interface {
	Consume(ctx context.Context, handler func([]*Outbox) error) error
}

type Outbox struct {
	Kind    string
	Payload []byte
}

func RestoreOutbox(kind string, payload []byte) *Outbox {
	return &Outbox{
		Kind:    kind,
		Payload: payload,
	}
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
