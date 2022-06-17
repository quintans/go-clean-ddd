package postgres

import (
	"context"
	"errors"

	"github.com/quintans/go-clean-ddd/internal/domain/entity"
	"github.com/quintans/go-clean-ddd/internal/domain/usecase"
	"github.com/quintans/go-clean-ddd/internal/domain/vo"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres/ent"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres/ent/customer"
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
	trans transaction.Transactioner[*ent.Tx]
}

func NewCustomerRepository(trans transaction.Transactioner[*ent.Tx]) CustomerRepository {
	return CustomerRepository{
		trans: trans,
	}
}

func (r CustomerRepository) Save(ctx context.Context, c entity.Customer) error {
	err := r.trans.Current(ctx, func(ctx context.Context, tx *ent.Tx) ([]event.DomainEvent, error) {
		_, err := tx.Customer.
			Create().
			SetID(c.ID().String()).
			SetFirstName(c.FullName().FirstName()).
			SetLastName(c.FullName().LastName()).
			SetEmail(c.Email().Email()).
			Save(ctx)
		return nil, err
	})

	return errorMap(err)
}

func (r CustomerRepository) Apply(ctx context.Context, id vo.CustomerID, apply func(context.Context, *entity.Customer) error) error {
	err := r.trans.Current(ctx, func(ctx context.Context, tx *ent.Tx) ([]event.DomainEvent, error) {
		c, err := r.getByID(ctx, tx, id)
		if err != nil {
			return nil, err
		}

		cust, err := toDomainCustomer(c)
		if err != nil {
			return nil, err
		}

		err = apply(ctx, &cust)
		if err != nil {
			return nil, err
		}

		n, err := tx.Customer.
			Update().
			Where(
				customer.And(
					customer.IDEQ(cust.ID().String()),
					customer.VersionEQ(c.Version),
				),
			).
			SetFirstName(cust.FullName().FirstName()).
			SetLastName(cust.FullName().LastName()).
			SetEmail(cust.Email().Email()).
			Save(ctx)
		if err != nil {
			return nil, err
		}
		if n == 0 {
			return nil, usecase.ErrOptimisticLocking
		}
		return nil, nil
	})

	return errorMap(err)
}

func (r CustomerRepository) getByID(ctx context.Context, tx *ent.Tx, id vo.CustomerID) (*ent.Customer, error) {
	c, err := tx.Customer.Query().Where(customer.ID(id.String())).Only(ctx)
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
		return usecase.ErrNotFound
	}
	return err
}

func toDomainCustomers(cs []*ent.Customer) ([]entity.Customer, error) {
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

func toDomainCustomer(c *ent.Customer) (entity.Customer, error) {
	id, err := vo.ParseCustomerID(c.ID)
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
