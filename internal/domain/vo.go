package domain

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// FullNameVO is a Value Object representing a first and last names
// gog:record
type FullName struct {
	// gog:@required
	firstName string
	// gog:@required
	lastName string
}

func (f FullName) String() string {
	return fmt.Sprintf("%s %s", f.firstName, f.lastName)
}

// gog:record
type Email struct {
	// gog:@required
	email string
}

func (e Email) String() string {
	return fmt.Sprintf("Email{email: %+v}", e.email)
}

var emailRe = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (e Email) validate() error {
	if !emailRe.MatchString(e.email) {
		return errors.Errorf("%s is not a valid email", e.email)
	}
	return nil
}

type CustomerID struct {
	id uuid.UUID
}

func ParseCustomerID(s string) (CustomerID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return CustomerID{}, err
	}
	c := CustomerID{
		id: id,
	}

	return c, nil
}

func NewCustomerID() CustomerID {
	return CustomerID{
		id: uuid.New(),
	}
}

func MustParseCustomerID(
	id string,
) CustomerID {
	c, err := ParseCustomerID(id)
	if err != nil {
		panic(err)
	}
	return c
}

func (c CustomerID) Id() uuid.UUID {
	return c.id
}

func (c CustomerID) IsZero() bool {
	return c == CustomerID{}
}

func (c CustomerID) String() string {
	return c.id.String()
}
