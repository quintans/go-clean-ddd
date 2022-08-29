package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/internal/domain"
	"github.com/quintans/go-clean-ddd/internal/domain/registration"
	"github.com/quintans/go-clean-ddd/lib/transaction"
)

type Registration struct {
	ID       string `db:"id"`
	Email    string
	Verified bool
}

type RegistrationRepository struct {
	trans *transaction.Transaction[*sqlx.Tx]
}

func NewRegistrationRepository(trans *transaction.Transaction[*sqlx.Tx]) RegistrationRepository {
	return RegistrationRepository{
		trans: trans,
	}
}

func (r RegistrationRepository) Create(ctx context.Context, c *registration.Registration) error {
	err := r.trans.Current(ctx, func(ctx context.Context, tx *sqlx.Tx) (transaction.EventPopper, error) {
		_, err := tx.ExecContext(
			ctx,
			"INSERT INTO registrations(id, email, verified) VALUES ($1, $2, $3)",
			c.ID(),
			c.Email().String(),
			false,
		)
		return c, faults.Wrap(err)
	})

	return errorMap(err)
}

func (r RegistrationRepository) Update(ctx context.Context, id string, apply func(context.Context, *registration.Registration) error) (*registration.Registration, error) {
	var registration *registration.Registration
	err := r.trans.Current(ctx, func(ctx context.Context, tx *sqlx.Tx) (transaction.EventPopper, error) {
		reg, err := r.getByID(ctx, tx, id)
		if err != nil {
			return nil, faults.Wrap(err)
		}

		registration, err = toDomainRegistration(reg)
		if err != nil {
			return nil, faults.Wrap(err)
		}

		err = apply(ctx, registration)
		if err != nil {
			return nil, faults.Wrap(err)
		}

		_, err = tx.ExecContext(
			ctx,
			"UPDATE registrations SET verified=$1 WHERE id=$2",
			true,
			id,
		)
		if err != nil {
			return nil, faults.Wrap(err)
		}
		return registration, nil
	})
	if err != nil {
		return nil, errorMap(err)
	}
	return registration, nil

}

func (r RegistrationRepository) getByID(ctx context.Context, tx *sqlx.Tx, id string) (Registration, error) {
	reg := Registration{}
	err := tx.GetContext(
		ctx,
		&reg,
		"SELECT * FROM registrations WHERE id=$1",
		id,
	)
	return reg, errorMap(err)
}

func toDomainRegistration(e Registration) (*registration.Registration, error) {
	email, err := domain.NewEmail(e.Email)
	if err != nil {
		return nil, faults.Wrap(err)
	}
	return registration.RestoreRegistration(e.ID, email, e.Verified)
}
