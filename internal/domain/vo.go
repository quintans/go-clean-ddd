package domain

import (
	"fmt"
	"regexp"

	"github.com/quintans/faults"
)

// FullName is a Value Object representing a first and last names
type FullName struct {
	firstName string
	lastName  string
}

func (f FullName) String() string {
	return fmt.Sprintf("%s %s", f.firstName, f.lastName)
}

func NewFullName(
	firstName string,
	lastName string,
) (FullName, error) {
	if firstName == "" {
		return FullName{}, faults.New("FullName.firstName cannot be empty")
	}
	if lastName == "" {
		return FullName{}, faults.New("FullName.lastName cannot be empty")
	}
	f := FullName{
		firstName: firstName,
		lastName:  lastName,
	}

	return f, nil
}

func MustNewFullName(
	firstName string,
	lastName string,
) FullName {
	f, err := NewFullName(
		firstName,
		lastName,
	)
	if err != nil {
		panic(err)
	}
	return f
}

func (f FullName) FirstName() string {
	return f.firstName
}

func (f FullName) LastName() string {
	return f.lastName
}

func (f FullName) IsZero() bool {
	return f == FullName{}
}

type Email struct {
	email string
}

func (e Email) String() string {
	return e.email
}

var emailRe = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (e Email) validate() error {
	if !emailRe.MatchString(e.email) {
		return faults.Errorf("%s is not a valid email", e.email)
	}
	return nil
}

func NewEmail(
	email string,
) (Email, error) {
	if email == "" {
		return Email{}, faults.New("Email.email cannot be empty")
	}
	e := Email{
		email: email,
	}
	if err := e.validate(); err != nil {
		return Email{}, faults.Wrap(err)
	}

	return e, nil
}

func MustNewEmail(
	email string,
) Email {
	e, err := NewEmail(
		email,
	)
	if err != nil {
		panic(err)
	}
	return e
}

func (e Email) IsZero() bool {
	return e == Email{}
}

func (e Email) MarshalText() (text []byte, err error) {
	return []byte(e.email), nil
}

func (e *Email) UnmarshalText(text []byte) error {
	e.email = string(text)
	return e.validate()
}
