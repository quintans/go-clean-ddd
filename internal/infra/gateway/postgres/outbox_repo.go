package postgres

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/quintans/go-clean-ddd/internal/domain/entity"
	"github.com/quintans/go-clean-ddd/internal/domain/usecase"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres/ent"
	"github.com/quintans/go-clean-ddd/lib/transaction"
)

type Outbox struct {
	Id      int
	Kind    string
	Payload []byte
}

type OutboxRepository struct {
	trans     transaction.Transactioner[*ent.Tx]
	batchSize uint
}

func NewOutboxRepository(trans transaction.Transactioner[*ent.Tx], batchSize uint) OutboxRepository {
	return OutboxRepository{
		trans:     trans,
		batchSize: batchSize,
	}
}

func (r OutboxRepository) Create(ctx context.Context, ob entity.Outbox) error {
	return r.trans.Current(ctx, func(ctx context.Context, tx *ent.Tx) (transaction.EventPopper, error) {
		_, err := tx.Outbox.
			Create().
			SetKind(ob.Kind()).
			SetPayload(ob.Payload()).
			Save(ctx)
		return nil, errorMap(err)
	})
}

func (r OutboxRepository) Consume(ctx context.Context, fn func([]entity.Outbox) error) error {
	for {
		err := r.trans.Current(ctx, func(ctx context.Context, tx *ent.Tx) (transaction.EventPopper, error) {
			ok, err := getAdvisoryLock(ctx, tx)
			if err != nil {
				return nil, errorMap(err)
			}
			if !ok {
				return nil, usecase.ErrNotFound
			}

			rows, err := tx.QueryContext(
				ctx,
				`UPDATE outbox SET consumed=TRUE WHERE consumed=FALSE 
				ORDER BY id	ASC LIMIT $1
				RETURNING id, kind, payload`,
				r.batchSize,
			)
			if err != nil {
				return nil, errorMap(err)
			}
			var entities []entity.Outbox
			o := Outbox{}
			forEachRow(rows, func() {
				entities = append(entities, entity.RestoreOutbox(o.Id, o.Kind, o.Payload))
			}, &o.Id, &o.Kind, &o.Payload)

			return nil, fn(entities)
		})
		if err != nil {
			return err
		}
	}
}

func getAdvisoryLock(ctx context.Context, tx *ent.Tx) (bool, error) {
	// 123 is the lock key
	rows, err := tx.QueryContext(ctx, "SELECT pg_try_advisory_xact_lock(123)")
	if err != nil {
		return false, err
	}

	var ok bool
	err = forEachRow(rows, nil, &ok)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func forEachRow(rows *sql.Rows, fn func(), v ...interface{}) error {
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(v...)
		if err != nil {
			return errors.WithStack(err)
		}
		if fn != nil {
			fn()
		} else {
			break
		}
	}

	if err := rows.Err(); err != nil {
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
