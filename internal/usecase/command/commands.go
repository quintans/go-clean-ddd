package command

import (
	"context"
	"errors"

	"github.com/quintans/go-clean-ddd/internal/domain"
	"github.com/quintans/go-clean-ddd/internal/usecase"
)

type Commander struct {
	customerRepository usecase.CustomerRepository
	customerView       usecase.CustomerViewRepository
}

func NewCommander(customerRepository usecase.CustomerRepository, customerView usecase.CustomerViewRepository) Commander {
	return Commander{
		customerRepository: customerRepository,
		customerView:       customerView,
	}
}

type UniquenessPolicy func(ctx context.Context, email domain.Email) (bool, error)

func (p UniquenessPolicy) IsUnique(ctx context.Context, email domain.Email) (bool, error) {
	return p(ctx, email)
}

func (r Commander) Register(ctx context.Context, cmd usecase.RegistrationCommand) (string, error) {
	id := domain.NewCustomerID()

	email, err := domain.NewEmail(cmd.Email)
	if err != nil {
		return "", err
	}

	customer, err := domain.NewCustomer(id, domain.FullName{}, email)
	if err != nil {
		return "", err
	}

	p := func(ctx context.Context, email domain.Email) (bool, error) {
		_, err := r.customerView.GetByEmail(ctx, email)
		if errors.Is(err, usecase.ErrReadModelNotFound) {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		return true, nil
	}

	if err := customer.AssignEmail(ctx, email, UniquenessPolicy(p)); err != nil {
		return "", err
	}

	if err := r.customerRepository.Save(ctx, customer); err != nil {
		return "", err
	}
	return id.String(), nil
}

func (r Commander) Update(ctx context.Context, cmd usecase.UpdateCommand) error {
	id, err := domain.ParseCustomerID(cmd.Id)
	if err != nil {
		return err
	}

	fullName, err := domain.NewFullName(cmd.FirstName, cmd.LastName)
	if err != nil {
		return err
	}

	return r.customerRepository.Apply(ctx, id, func(ctx context.Context, c domain.Customer) (domain.Customer, error) {
		c.UpdateInfo(fullName)
		return c, nil
	})
}
