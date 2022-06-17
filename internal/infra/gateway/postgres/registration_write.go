package postgres

import (
	"context"

	"github.com/quintans/go-clean-ddd/internal/domain/entity"
	"github.com/quintans/go-clean-ddd/internal/domain/vo"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres/ent"
	"github.com/quintans/go-clean-ddd/lib/event"
	"github.com/quintans/go-clean-ddd/lib/transaction"
)

type Registration struct {
	Id       string
	Email    string
	Verified bool
}

type RegistrationRepository struct {
	trans transaction.Transactioner[*ent.Tx]
}

func NewRegistrationRepository(trans transaction.Transactioner[*ent.Tx]) RegistrationRepository {
	return RegistrationRepository{
		trans: trans,
	}
}

func (r RegistrationRepository) Save(ctx context.Context, c entity.Registration) error {
	err := r.trans.Current(ctx, func(ctx context.Context, tx *ent.Tx) ([]event.DomainEvent, error) {
		_, err := tx.ExecContext(
			ctx,
			"INSERT INTO registration(id, email, verified) SET VALUES($1, $2, $3)",
			c.ID(), c.Email(), c.Verified(),
		)
		return c.PopEvents(), err
	})

	return errorMap(err)
}

func (r RegistrationRepository) Apply(ctx context.Context, id string, apply func(context.Context, *entity.Registration) error) error {
	err := r.trans.Current(ctx, func(ctx context.Context, tx *ent.Tx) ([]event.DomainEvent, error) {
		reg, err := r.getByID(ctx, tx, id)
		if err != nil {
			return nil, err
		}

		registration, err := toDomainRegistration(reg)
		if err != nil {
			return nil, err
		}

		err = apply(ctx, &registration)
		if err != nil {
			return nil, err
		}

		// optimistic locking is used
		_, err = tx.ExecContext(
			ctx,
			"UPDATE customer SET email=$1, verified=$2  WHERE id=$3",
			registration.Email(), registration.Verified(), registration.ID(),
		)
		return registration.PopEvents(), err
	})

	return errorMap(err)
}

func (r RegistrationRepository) getByID(ctx context.Context, tx *ent.Tx, id string) (Registration, error) {
	customer := Registration{}
	err := tx.Get(&customer, "SELECT * FROM registration WHERE id=$1", id)
	if err != nil {
		return Registration{}, errorMap(err)
	}
	return customer, nil
}

func toDomainRegistration(c Registration) (entity.Registration, error) {
	email, err := vo.NewEmail(c.Email)
	if err != nil {
		return entity.Registration{}, err
	}
	return entity.RestoreRegistration(c.Id, email, c.Verified)
}
