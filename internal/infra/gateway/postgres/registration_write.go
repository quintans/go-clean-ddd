package postgres

import (
	"context"

	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/internal/domain/entity"
	"github.com/quintans/go-clean-ddd/internal/domain/vo"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres/ent"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres/ent/registration"
	"github.com/quintans/go-clean-ddd/lib/transaction"
)

type RegistrationRepository struct {
	trans transaction.Transactioner[*ent.Tx]
}

func NewRegistrationRepository(trans transaction.Transactioner[*ent.Tx]) RegistrationRepository {
	return RegistrationRepository{
		trans: trans,
	}
}

func (r RegistrationRepository) Create(ctx context.Context, c entity.Registration) error {
	err := r.trans.Current(ctx, func(ctx context.Context, tx *ent.Tx) (transaction.EventPopper, error) {
		_, err := tx.ExecContext(
			ctx,
			"INSERT INTO registration(id, email, verified) SET VALUES($1, $2, $3)",
			c.ID(), c.Email(), c.Verified(),
		)
		return c, faults.Wrap(err)
	})

	return errorMap(err)
}

func (r RegistrationRepository) Update(ctx context.Context, id string, apply func(context.Context, *entity.Registration) error) error {
	err := r.trans.Current(ctx, func(ctx context.Context, tx *ent.Tx) (transaction.EventPopper, error) {
		reg, err := r.getByID(ctx, tx, id)
		if err != nil {
			return nil, faults.Wrap(err)
		}

		registration, err := toDomainRegistration(reg)
		if err != nil {
			return nil, faults.Wrap(err)
		}

		err = apply(ctx, &registration)
		if err != nil {
			return nil, faults.Wrap(err)
		}

		_, err = tx.Registration.
			UpdateOne(reg).
			SetVerified(true).
			Save(ctx)
		return registration, faults.Wrap(err)
	})

	return errorMap(err)
}

func (r RegistrationRepository) getByID(ctx context.Context, tx *ent.Tx, id string) (*ent.Registration, error) {
	e, err := tx.Registration.Query().Where(registration.ID(id)).Only(ctx)
	if err != nil {
		return nil, errorMap(err)
	}
	return e, nil
}

func toDomainRegistration(e *ent.Registration) (entity.Registration, error) {
	email, err := vo.NewEmail(e.Email)
	if err != nil {
		return entity.Registration{}, faults.Wrap(err)
	}
	return entity.RestoreRegistration(e.ID, email, e.Verified)
}
