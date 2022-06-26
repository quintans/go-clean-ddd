package command

import (
	"context"

	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/internal/app"
	"github.com/quintans/go-clean-ddd/internal/domain"
	"github.com/quintans/go-clean-ddd/internal/domain/customer"
)

type UpdateCustomerCommand struct {
	Id        string
	FirstName string
	LastName  string
}

type UpdateCustomerHandler interface {
	Handle(context.Context, UpdateCustomerCommand) error
}

type UpdateCustomer struct {
	customerRepository app.CustomerRepository
}

func NewUpdateCustomer(customerRepository app.CustomerRepository, customerView app.CustomerViewRepository) UpdateCustomer {
	return UpdateCustomer{
		customerRepository: customerRepository,
	}
}

func (r UpdateCustomer) Handle(ctx context.Context, cmd UpdateCustomerCommand) error {
	id, err := customer.ParseCustomerID(cmd.Id)
	if err != nil {
		return faults.Wrap(err)
	}

	fullName, err := domain.NewFullName(cmd.FirstName, cmd.LastName)
	if err != nil {
		return faults.Wrap(err)
	}

	return r.customerRepository.Update(ctx, id, func(ctx context.Context, c *customer.Customer) error {
		c.UpdateInfo(fullName)
		return nil
	})
}
