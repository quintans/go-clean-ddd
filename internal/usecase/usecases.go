package usecase

import (
	"context"
)

type CustomerDTO struct {
	Id        string
	Email     string
	FirstName string
	LastName  string
}

type Querier interface {
	GetAllCustomers(context.Context) ([]CustomerDTO, error)
}

type Commander interface {
	Register(context.Context, RegistrationCommand) (string, error)
	Update(context.Context, UpdateCommand) error
}

type RegistrationCommand struct {
	Email string
}

type UpdateCommand struct {
	Id        string
	FirstName string
	LastName  string
}
