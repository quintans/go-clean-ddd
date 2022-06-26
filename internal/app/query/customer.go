package query

import (
	"context"

	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/internal/app"
	"github.com/quintans/go-clean-ddd/internal/domain/customer"
)

type CustomerDTO struct {
	Id        string
	Email     string
	FirstName string
	LastName  string
}

type AllCustomersHandler interface {
	Handle(context.Context) ([]CustomerDTO, error)
}

type AllCustomers struct {
	repo app.CustomerViewRepository
}

func NewAllCustomers(repo app.CustomerViewRepository) AllCustomers {
	return AllCustomers{
		repo: repo,
	}
}

func (r AllCustomers) Handle(ctx context.Context) ([]CustomerDTO, error) {
	customers, err := r.repo.GetAll(ctx)
	if err != nil {
		return nil, faults.Wrap(err)
	}

	return toCustomerDTOs(customers), nil
}

func toCustomerDTOs(in []customer.Customer) []CustomerDTO {
	out := make([]CustomerDTO, len(in))
	for k, v := range in {
		out[k] = toCustomerDTO(v)
	}
	return out
}

func toCustomerDTO(c customer.Customer) CustomerDTO {
	return CustomerDTO{
		Id:        c.ID().String(),
		Email:     c.Email().String(),
		FirstName: c.FullName().FirstName(),
		LastName:  c.FullName().LastName(),
	}
}
