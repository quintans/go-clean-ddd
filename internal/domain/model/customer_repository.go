package model

import "context"

// CustomerRepository interface for handling the persistence of the Customer
type CustomerRepository interface {
	FindAll(context.Context) ([]*Customer, error)
	FindByEmail(context.Context, string) (*Customer, error)
	Store(context.Context, *Customer) error
}
