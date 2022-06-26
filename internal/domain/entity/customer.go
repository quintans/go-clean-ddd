package entity

import (
	"context"

	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/internal/domain/vo"
)

// Customer entity (and also root aggregate).
type Customer struct {
	id vo.CustomerID

	fullName vo.FullName
	email    vo.Email
}

type UniqueEmailPolicy interface {
	IsUnique(context.Context, vo.Email) (bool, error)
}

// NewCustomer creates a customer
func NewCustomer(ctx context.Context, email vo.Email, policy UniqueEmailPolicy) (Customer, error) {
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
		id:    vo.NewCustomerID(),
		email: email,
	}, nil
}

// RestoreCustomer instantiates customer from a previous stored state
func RestoreCustomer(id vo.CustomerID, fullName vo.FullName, email vo.Email) Customer {
	return Customer{
		id:       id,
		email:    email,
		fullName: fullName,
	}
}

func (c *Customer) UpdateInfo(fullName vo.FullName) {
	c.fullName = fullName
}

func (c Customer) ID() vo.CustomerID {
	return c.id
}

func (c Customer) FullName() vo.FullName {
	return c.fullName
}

func (c Customer) Email() vo.Email {
	return c.email
}
