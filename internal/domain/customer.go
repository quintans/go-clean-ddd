package domain

import (
	"context"

	"github.com/pkg/errors"
)

// NewCustomer creates a customer
func NewCustomer(id CustomerID, fullName FullName, email Email) (Customer, error) {
	if id.IsZero() {
		return Customer{}, errors.New("customer ID is undefined")
	}
	if email.IsZero() {
		return Customer{}, errors.New("email is undefined")
	}
	return Customer{
		id:       id,
		email:    email,
		fullName: fullName,
	}, nil
}

// Customer entity (and also root aggregate).
type Customer struct {
	id CustomerID

	fullName FullName
	email    Email
}

type UniquenessPolicy interface {
	IsUnique(context.Context, Email) (bool, error)
}

// AssignEmail updates email if is unique
func (c *Customer) AssignEmail(ctx context.Context, email Email, policy UniquenessPolicy) error {
	ok, err := policy.IsUnique(ctx, email)
	if err != nil {
		return errors.Wrap(err, "failed to check uniqueness on AssignEmail")
	}
	if !ok {
		return errors.Errorf("the provided e-mail %s is not unique", email)
	}
	c.email = email
	return nil
}

func (c *Customer) UpdateInfo(fullName FullName) {
	c.fullName = fullName
}

func (c Customer) ID() CustomerID {
	return c.id
}

func (c Customer) FullName() FullName {
	return c.fullName
}

func (c Customer) Email() Email {
	return c.email
}
