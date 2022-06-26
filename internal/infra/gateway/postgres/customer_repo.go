package postgres

import (
	"context"
	"errors"

	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/internal/app"
	"github.com/quintans/go-clean-ddd/internal/domain"
	"github.com/quintans/go-clean-ddd/internal/domain/customer"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres/ent"
	entcust "github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres/ent/customer"
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
	trans transaction.Transactioner[*ent.Tx]
}

func NewCustomerRepository(trans transaction.Transactioner[*ent.Tx]) CustomerRepository {
	return CustomerRepository{
		trans: trans,
	}
}

func (r CustomerRepository) Create(ctx context.Context, c customer.Customer) error {
	err := r.trans.Current(ctx, func(ctx context.Context, tx *ent.Tx) (transaction.EventPopper, error) {
		_, err := tx.Customer.
			Create().
			SetID(c.ID().String()).
			SetFirstName(c.FullName().FirstName()).
			SetLastName(c.FullName().LastName()).
			SetEmail(c.Email().Email()).
			Save(ctx)
		return nil, faults.Wrap(err)
	})

	return errorMap(err)
}

func (r CustomerRepository) Update(ctx context.Context, id customer.CustomerID, apply func(context.Context, *customer.Customer) error) error {
	err := r.trans.Current(ctx, func(ctx context.Context, tx *ent.Tx) (transaction.EventPopper, error) {
		c, err := r.getByID(ctx, tx, id)
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
		n, err := tx.Customer.
			Update().
			Where(
				entcust.And(
					entcust.IDEQ(cust.ID().String()),
					entcust.VersionEQ(c.Version),
				),
			).
			SetFirstName(cust.FullName().FirstName()).
			SetLastName(cust.FullName().LastName()).
			SetEmail(cust.Email().Email()).
			Save(ctx)
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

func (r CustomerRepository) getByID(ctx context.Context, tx *ent.Tx, id customer.CustomerID) (*ent.Customer, error) {
	c, err := tx.Customer.Query().Where(entcust.ID(id.String())).Only(ctx)
	if err != nil {
		return nil, errorMap(err)
	}
	return c, nil
}

func errorMap(err error) error {
	if err == nil {
		return nil
	}

	var target *ent.NotFoundError
	if errors.As(err, &target) {
		return faults.Wrap(app.ErrNotFound)
	}
	return faults.Wrap(err)
}

func toDomainCustomers(cs []*ent.Customer) ([]customer.Customer, error) {
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

func toDomainCustomer(c *ent.Customer) (customer.Customer, error) {
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
