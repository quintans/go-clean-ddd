package query

import (
	"context"

	"github.com/quintans/go-clean-ddd/internal/domain"
	"github.com/quintans/go-clean-ddd/internal/usecase"
)

type Querier struct {
	repo usecase.CustomerViewRepository
}

func NewQuerier(repo usecase.CustomerViewRepository) Querier {
	return Querier{
		repo: repo,
	}
}

func (r Querier) GetAllCustomers(ctx context.Context) ([]usecase.CustomerDTO, error) {
	customers, err := r.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return toCustomerDTOs(customers), nil
}

func toCustomerDTOs(in []domain.Customer) []usecase.CustomerDTO {
	out := make([]usecase.CustomerDTO, len(in))
	for k, v := range in {
		out[k] = toCustomerDTO(v)
	}
	return out
}

func toCustomerDTO(c domain.Customer) usecase.CustomerDTO {
	return usecase.CustomerDTO{
		Id:        c.ID().String(),
		Email:     c.Email().String(),
		FirstName: c.FullName().FirstName(),
		LastName:  c.FullName().LastName(),
	}
}
