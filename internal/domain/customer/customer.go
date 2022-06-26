package customer

import (
	"context"

	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/internal/domain"
)

// Customer entity (and also root aggregate).
type Customer struct {
	id CustomerID

	fullName domain.FullName
	email    domain.Email
}

// NewCustomer creates a customer
func NewCustomer(ctx context.Context, email domain.Email, policy domain.UniqueEmailPolicy) (Customer, error) {
	if email.IsZero() {
		return Customer{}, faults.New("email is undefined")
	}
	ok, err := policy.IsUnique(ctx, email)
	if err != nil {
		return Customer{}, faults.Wrapf(err, "failed to check uniqueness of email")
	}
	if !ok {
		return Customer{}, faults.Errorf("the provided e-mail %s is already taken unique", email)
	}
	return Customer{
		id:    NewCustomerID(),
		email: email,
	}, nil
}

// RestoreCustomer instantiates customer from a previous stored state
func RestoreCustomer(id CustomerID, fullName domain.FullName, email domain.Email) Customer {
	return Customer{
		id:       id,
		email:    email,
		fullName: fullName,
	}
}

func (c *Customer) UpdateInfo(fullName domain.FullName) {
	c.fullName = fullName
}

func (c Customer) ID() CustomerID {
	return c.id
}

func (c Customer) FullName() domain.FullName {
	return c.fullName
}

func (c Customer) Email() domain.Email {
	return c.email
}
