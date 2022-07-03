package postgres

import (
	"context"

	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/internal/domain"
	"github.com/quintans/go-clean-ddd/internal/domain/registration"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres/ent"
	entreg "github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres/ent/registration"
	"github.com/quintans/go-clean-ddd/lib/transaction"
)

type RegistrationRepository struct {
	trans *transaction.Transaction[*ent.Tx]
}

func NewRegistrationRepository(trans *transaction.Transaction[*ent.Tx]) RegistrationRepository {
	return RegistrationRepository{
		trans: trans,
	}
}

func (r RegistrationRepository) Create(ctx context.Context, c registration.Registration) error {
	err := r.trans.Current(ctx, func(ctx context.Context, tx *ent.Tx) (transaction.EventPopper, error) {
		_, err := tx.Registration.Create().
			SetID(c.ID()).
			SetEmail(c.Email().String()).
			SetVerified(c.Verified()).
			Save(ctx)

		return c, faults.Wrap(err)
	})

	return errorMap(err)
}

func (r RegistrationRepository) Update(ctx context.Context, id string, apply func(context.Context, *registration.Registration) error) error {
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
	e, err := tx.Registration.Query().Where(entreg.ID(id)).Only(ctx)
	if err != nil {
		return nil, errorMap(err)
	}
	return e, nil
}

func toDomainRegistration(e *ent.Registration) (registration.Registration, error) {
	email, err := domain.NewEmail(e.Email)
	if err != nil {
		return registration.Registration{}, faults.Wrap(err)
	}
	return registration.RestoreRegistration(e.ID, email, e.Verified)
}
