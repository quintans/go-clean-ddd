package postgres

import (
	"context"

	"github.com/quintans/go-clean-ddd/internal/domain"
	"github.com/quintans/go-clean-ddd/internal/domain/customer"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres/ent"
	entcust "github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres/ent/customer"
)

type CustomerViewRepository struct {
	client *ent.Client
}

func NewCustomerViewRepository(db *ent.Client) CustomerViewRepository {
	return CustomerViewRepository{
		client: db,
	}
}

func (r CustomerViewRepository) GetAll(ctx context.Context) ([]customer.Customer, error) {
	all, err := r.client.Customer.Query().All(ctx)
	if err != nil {
		return nil, errorMap(err)
	}
	return toDomainCustomers(all)
}

func (r CustomerViewRepository) GetByEmail(ctx context.Context, email domain.Email) (customer.Customer, error) {
	c, err := r.client.Customer.Query().Where(entcust.EmailEQ(email.String())).First(ctx)
	if err != nil {
		return customer.Customer{}, errorMap(err)
	}
	return toDomainCustomer(c)
}
