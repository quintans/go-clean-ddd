package query

import (
	"context"

	"github.com/quintans/go-clean-ddd/internal/domain/entity"
	"github.com/quintans/go-clean-ddd/internal/usecase"
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
	repo usecase.CustomerViewRepository
}

func NewAllCustomers(repo usecase.CustomerViewRepository) AllCustomers {
	return AllCustomers{
		repo: repo,
	}
}

func (r AllCustomers) Handle(ctx context.Context) ([]CustomerDTO, error) {
	customers, err := r.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return toCustomerDTOs(customers), nil
}

func toCustomerDTOs(in []entity.Customer) []CustomerDTO {
	out := make([]CustomerDTO, len(in))
	for k, v := range in {
		out[k] = toCustomerDTO(v)
	}
	return out
}

func toCustomerDTO(c entity.Customer) CustomerDTO {
	return CustomerDTO{
		Id:        c.ID().String(),
		Email:     c.Email().String(),
		FirstName: c.FullName().FirstName(),
		LastName:  c.FullName().LastName(),
	}
}
