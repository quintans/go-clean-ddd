package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/quintans/go-clean-ddd/internal/domain/entity"
	"github.com/quintans/go-clean-ddd/internal/domain/usecase"
	"github.com/quintans/go-clean-ddd/internal/domain/vo"
	"github.com/quintans/go-clean-ddd/lib/event"
	"github.com/quintans/go-clean-ddd/lib/transaction"
)

type Customer struct {
	Id        string
	Version   int
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
}

type CustomerRepository struct {
	trans transaction.Transactioner[*sqlx.Tx]
}

func NewCustomerRepository(trans transaction.Transactioner[*sqlx.Tx]) CustomerRepository {
	return CustomerRepository{
		trans: trans,
	}
}

func (r CustomerRepository) Save(ctx context.Context, c entity.Customer) error {
	err := r.trans.Current(ctx, func(ctx context.Context, tx *sqlx.Tx) ([]event.DomainEvent, error) {
		_, err := tx.ExecContext(
			ctx,
			"INSERT INTO customer(id, version, first_name, last_name, email) SET VALUES($1, $2, $3, $4, $5)",
			c.ID(), 0, c.FullName().FirstName(), c.FullName().LastName(), c.Email(),
		)
		return nil, err
	})

	return errorMap(err)
}

func (r CustomerRepository) Apply(ctx context.Context, id vo.CustomerID, apply func(context.Context, *entity.Customer) error) error {
	err := r.trans.Current(ctx, func(ctx context.Context, tx *sqlx.Tx) ([]event.DomainEvent, error) {
		c, err := r.getByID(ctx, tx, id)
		if err != nil {
			return nil, err
		}

		customer, err := toDomainCustomer(c)
		if err != nil {
			return nil, err
		}

		err = apply(ctx, &customer)
		if err != nil {
			return nil, err
		}

		// optimistic locking is used
		_, err = tx.ExecContext(
			ctx,
			"UPDATE customer SET first_name=$1, last_name=$2, email=$3, version=version+1 WHERE id=$4 AND version=$5",
			customer.FullName().FirstName(), customer.FullName().LastName(), customer.Email(), customer.ID(), c.Version,
		)
		return nil, err
	})

	return errorMap(err)
}

func (r CustomerRepository) getByID(ctx context.Context, tx *sqlx.Tx, id vo.CustomerID) (Customer, error) {
	customer := Customer{}
	err := tx.Get(&customer, "SELECT * FROM customer WHERE id=$1", id.String())
	if err != nil {
		return Customer{}, errorMap(err)
	}
	return customer, nil
}

func errorMap(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return usecase.ErrModelNotFound
	}
	return err
}

func toDomainCustomers(cs []Customer) ([]entity.Customer, error) {
	dcs := make([]entity.Customer, len(cs))
	for k, v := range cs {
		dc, err := toDomainCustomer(v)
		if err != nil {
			return nil, err
		}
		dcs[k] = dc
	}
	return dcs, nil
}

func toDomainCustomer(c Customer) (entity.Customer, error) {
	id, err := vo.ParseCustomerID(c.Id)
	if err != nil {
		return entity.Customer{}, err
	}
	fullName, err := vo.NewFullName(c.FirstName, c.LastName)
	if err != nil {
		return entity.Customer{}, err
	}
	email, err := vo.NewEmail(c.Email)
	if err != nil {
		return entity.Customer{}, err
	}
	return entity.RestoreCustomer(id, fullName, email), nil
}
