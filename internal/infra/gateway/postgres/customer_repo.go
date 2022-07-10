package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/internal/app"
	"github.com/quintans/go-clean-ddd/internal/domain"
	"github.com/quintans/go-clean-ddd/internal/domain/customer"
	"github.com/quintans/go-clean-ddd/lib/transaction"
)

type Customer struct {
	ID        string `db:"id"`
	Version   int
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
}

type CustomerRepository struct {
	trans *transaction.Transaction[*sqlx.Tx]
}

func NewCustomerRepository(trans *transaction.Transaction[*sqlx.Tx]) CustomerRepository {
	return CustomerRepository{
		trans: trans,
	}
}

func (r CustomerRepository) Create(ctx context.Context, c customer.Customer) error {
	err := r.trans.Current(ctx, func(ctx context.Context, tx *sqlx.Tx) (transaction.EventPopper, error) {
		_, err := tx.ExecContext(
			ctx,
			"INSERT INTO customers (id, first_name, last_name, email, version) VALUES ($1, $2, $3, $4, 1)",
			c.ID().String(),
			c.FullName().FirstName(),
			c.FullName().LastName(),
			c.Email().String(),
		)
		return nil, faults.Wrap(err)
	})

	return errorMap(err)
}

func (r CustomerRepository) Update(ctx context.Context, id customer.CustomerID, apply func(context.Context, *customer.Customer) error) error {
	err := r.trans.Current(ctx, func(ctx context.Context, tx *sqlx.Tx) (transaction.EventPopper, error) {
		c, err := r.getByID(ctx, tx, id.String())
		if err != nil {
			return nil, faults.Wrap(err)
		}

		cust, err := toDomainCustomer(c)
		if err != nil {
			return nil, faults.Wrap(err)
		}

		err = apply(ctx, &cust)
		if err != nil {
			return nil, faults.Wrap(err)
		}

		// uses optimistic locking
		res, err := tx.ExecContext(
			ctx,
			"UPDATE customers SET first_name=$1, last_name=$2, email=$3, version=version+1 WHERE id=$4 AND version=$5",
			c.FirstName,
			c.LastName,
			c.Email,
			c.ID,
			c.Version,
		)
		if err != nil {
			return nil, faults.Wrap(err)
		}
		n, err := res.RowsAffected()
		if err != nil {
			return nil, faults.Wrap(err)
		}
		if n == 0 {
			return nil, app.ErrOptimisticLocking
		}
		return nil, nil
	})

	return errorMap(err)
}

func (r CustomerRepository) getByID(ctx context.Context, tx *sqlx.Tx, id string) (Customer, error) {
	c := Customer{}
	err := tx.GetContext(
		ctx,
		&c,
		"SELECT * FROM customers WHERE id=$1",
		id,
	)
	return c, errorMap(err)
}

func errorMap(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return faults.Wrap(app.ErrNotFound)
	}
	return faults.Wrap(err)
}

func toDomainCustomers(cs []Customer) ([]customer.Customer, error) {
	dcs := make([]customer.Customer, len(cs))
	for k, v := range cs {
		dc, err := toDomainCustomer(v)
		if err != nil {
			return nil, faults.Wrap(err)
		}
		dcs[k] = dc
	}
	return dcs, nil
}

func toDomainCustomer(c Customer) (customer.Customer, error) {
	id, err := customer.ParseCustomerID(c.ID)
	if err != nil {
		return customer.Customer{}, faults.Wrap(err)
	}
	fullName, err := domain.NewFullName(c.FirstName, c.LastName)
	if err != nil {
		return customer.Customer{}, faults.Wrap(err)
	}
	email, err := domain.NewEmail(c.Email)
	if err != nil {
		return customer.Customer{}, faults.Wrap(err)
	}
	return customer.RestoreCustomer(id, fullName, email), nil
}
