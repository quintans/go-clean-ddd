package model

import (
	"regexp"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var emailRe = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// NewCustomerAll creates a customer
func NewCustomerAll(uuid uuid.UUID, fullname FullNameVO, email string) *Customer {
	return &Customer{
		uuid:     uuid,
		email:    email,
		fullname: fullname,
	}
}

// Customer entity (and also root agregate).
type Customer struct {
	uuid uuid.UUID

	fullname FullNameVO
	email    string
}

// SetEmail applies the rule of e-mail validation.
// It could be a more complex rule
func (c *Customer) SetEmail(email string) error {
	if !emailRe.MatchString(email) {
		return errors.Errorf("%s is not a valid email", email)
	}
	c.email = email
	return nil
}

func (c *Customer) Uuid() uuid.UUID {
	return c.uuid
}

func (c *Customer) Fullname() FullNameVO {
	return c.fullname
}

func (c *Customer) Email() string {
	return c.email
}

// FullNameVO is a Value Object representing a first and last names
//gog:value
type FullNameVO struct {
	//gog:@wither
	firstName string
	//gog:@wither
	lastName string
}
