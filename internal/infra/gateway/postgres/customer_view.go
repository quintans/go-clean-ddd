package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/quintans/go-clean-ddd/internal/domain/entity"
	"github.com/quintans/go-clean-ddd/internal/domain/vo"
)

type CustomerViewRepository struct {
	db *sqlx.DB
}

func NewCustomerViewRepository(db *sqlx.DB) CustomerViewRepository {
	return CustomerViewRepository{
		db: db,
	}
}

func (r CustomerViewRepository) GetAll(ctx context.Context) ([]entity.Customer, error) {
	customers := []Customer{}
	err := r.db.Select(&customers, "SELECT * FROM customer")
	if err != nil {
		return nil, errorMap(err)
	}
	return toDomainCustomers(customers)
}

func (r CustomerViewRepository) GetByEmail(ctx context.Context, email vo.Email) (entity.Customer, error) {
	customer := Customer{}
	err := r.db.Get(&customer, "SELECT * FROM customer WHERE email=$1", email.Email())
	if err != nil {
		return entity.Customer{}, errorMap(err)
	}
	return toDomainCustomer(customer)
}
