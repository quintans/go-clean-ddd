package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/quintans/go-clean-ddd/internal/domain/entity"
	"github.com/quintans/go-clean-ddd/internal/domain/usecase"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres/ent"
	"github.com/quintans/go-clean-ddd/lib/event"
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

func (r OutboxRepository) Save(ctx context.Context, ob entity.Outbox) error {
	return r.trans.Current(ctx, func(ctx context.Context, tx *ent.Tx) ([]event.DomainEvent, error) {
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
		err := r.trans.Current(ctx, func(ctx context.Context, tx *ent.Tx) ([]event.DomainEvent, error) {
			ok, err := getAdvisoryLock(ctx, tx)
			if err != nil {
				return nil, err
			}
			if !ok {
				return nil, usecase.ErrNotFound
			}

			var entities []entity.Outbox
			now := time.Now().UTC()
			until := now.Add(10 * time.Second)
			outbox := []Outbox{}
			err = tx.SelectContext(
				ctx, &outbox,
				`UPDATE outbox SET consumed=TRUE WHERE consumed=FALSE 
			ORDER BY id	ASC LIMIT $3
			RETURNING id, kind, payload`,
				until, now, r.batchSize,
			)
			if err != nil {
				return nil, errorMap(err)
			}
			for _, o := range outbox {
				entities = append(entities, entity.RestoreOutbox(o.Id, o.Kind, o.Payload))
			}
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
	err = getOneFromRow(rows, &ok)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func getOneFromRow(rows *sql.Rows, v interface{}) error {
	if rows.Next() {
		err := rows.Scan(v)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	if closeErr := rows.Close(); closeErr != nil {
		if closeErr != nil {
			return errors.WithStack(closeErr)
		}
	}

	if err := rows.Err(); err != nil {
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
