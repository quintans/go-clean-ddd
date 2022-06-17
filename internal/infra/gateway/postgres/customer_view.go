package postgres

import (
	"context"

	"github.com/quintans/go-clean-ddd/internal/domain/entity"
	"github.com/quintans/go-clean-ddd/internal/domain/vo"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres/ent"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres/ent/customer"
)

type CustomerViewRepository struct {
	client *ent.Client
}

func NewCustomerViewRepository(db *ent.Client) CustomerViewRepository {
	return CustomerViewRepository{
		client: db,
	}
}

func (r CustomerViewRepository) GetAll(ctx context.Context) ([]entity.Customer, error) {
	all, err := r.client.Customer.Query().All(ctx)
	if err != nil {
		return nil, errorMap(err)
	}
	return toDomainCustomers(all)
}

func (r CustomerViewRepository) GetByEmail(ctx context.Context, email vo.Email) (entity.Customer, error) {
	c, err := r.client.Customer.Query().Where(customer.EmailEQ(email.String())).First(ctx)
	if err != nil {
		return entity.Customer{}, errorMap(err)
	}
	return toDomainCustomer(c)
}
