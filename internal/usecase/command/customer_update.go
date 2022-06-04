package command

import (
	"context"

	"github.com/quintans/go-clean-ddd/internal/domain/entity"
	"github.com/quintans/go-clean-ddd/internal/domain/vo"
	"github.com/quintans/go-clean-ddd/internal/usecase"
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
	customerRepository usecase.CustomerRepository
}

func NewUpdateCustomer(customerRepository usecase.CustomerRepository, customerView usecase.CustomerViewRepository) UpdateCustomer {
	return UpdateCustomer{
		customerRepository: customerRepository,
	}
}

func (r UpdateCustomer) Handle(ctx context.Context, cmd UpdateCustomerCommand) error {
	id, err := vo.ParseCustomerID(cmd.Id)
	if err != nil {
		return err
	}

	fullName, err := vo.NewFullName(cmd.FirstName, cmd.LastName)
	if err != nil {
		return err
	}

	return r.customerRepository.Apply(ctx, id, func(ctx context.Context, c *entity.Customer) error {
		c.UpdateInfo(fullName)
		return nil
	})
}
