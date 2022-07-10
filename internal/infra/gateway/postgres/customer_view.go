package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/quintans/go-clean-ddd/internal/domain"
	"github.com/quintans/go-clean-ddd/internal/domain/customer"
)

type CustomerViewRepository struct {
	client *sqlx.DB
}

func NewCustomerViewRepository(db *sqlx.DB) CustomerViewRepository {
	return CustomerViewRepository{
		client: db,
	}
}

func (r CustomerViewRepository) GetAll(ctx context.Context) ([]customer.Customer, error) {
	customers := []Customer{}
	err := r.client.SelectContext(ctx, &customers, "SELECT * FROM customers")
	if err != nil {
		return nil, errorMap(err)
	}
	return toDomainCustomers(customers)
}

func (r CustomerViewRepository) GetByEmail(ctx context.Context, email domain.Email) (customer.Customer, error) {
	c := Customer{}
	err := r.client.GetContext(ctx, &c, "SELECT * FROM customers WHERE email=$1", email.String())
	if err != nil {
		return customer.Customer{}, errorMap(err)
	}
	return toDomainCustomer(c)
}
