package postgres

import (
	"context"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/quintans/go-clean-ddd/internal/domain/entity"
	"github.com/quintans/go-clean-ddd/lib/event"
	"github.com/quintans/go-clean-ddd/lib/transaction"
)

type Outbox struct {
	Id      int
	Kind    string
	Payload []byte
}

type OutboxRepository struct {
	trans     transaction.Transactioner[*sqlx.Tx]
	batchSize uint
}

func NewOutboxRepository(trans transaction.Transactioner[*sqlx.Tx], batchSize uint) OutboxRepository {
	return OutboxRepository{
		trans:     trans,
		batchSize: batchSize,
	}
}

func (r OutboxRepository) Save(ctx context.Context, ob entity.Outbox) error {
	return r.trans.Current(ctx, func(ctx context.Context, tx *sqlx.Tx) ([]event.DomainEvent, error) {
		_, err := tx.ExecContext(
			ctx,
			"INSERT INTO outbox(kind, payload, consumed) SET VALUES($1, $2, false)",
			ob.Kind(), ob.Payload(),
		)
		return nil, errorMap(err)
	})
}

func (r OutboxRepository) LockAndLoad(ctx context.Context) ([]entity.Outbox, error) {
	var entities []entity.Outbox
	err := r.trans.Current(ctx, func(ctx context.Context, tx *sqlx.Tx) ([]event.DomainEvent, error) {
		now := time.Now().UTC()
		until := now.Add(10 * time.Second)
		outbox := []Outbox{}
		err := tx.SelectContext(
			ctx, &outbox,
			`UPDATE outbox SET locked_until = $1
			WHERE consumed=FALSE AND locked_until IS NULL OR locked_until < $2
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
		return nil, nil
	})

	return entities, err
}

func (r OutboxRepository) Consume(ctx context.Context, outboxes []entity.Outbox) error {
	return r.consume(ctx, outboxes, true)
}

func (r OutboxRepository) Release(ctx context.Context, outboxes []entity.Outbox) error {
	return r.consume(ctx, outboxes, false)
}

func (r OutboxRepository) consume(ctx context.Context, outboxes []entity.Outbox, yes bool) error {
	return r.trans.Current(ctx, func(ctx context.Context, tx *sqlx.Tx) ([]event.DomainEvent, error) {
		sb := strings.Builder{}
		sb.WriteString("UPDATE outbox SET locked_until = NULL")
		if yes {
			sb.WriteString(", consumed=TRUE")
		}
		sb.WriteString(" WHERE id IN($1)")

		ids := make([]int, len(outboxes))
		for k, v := range outboxes {
			ids[k] = v.Id()
		}
		upd, args, err := sqlx.In(sb.String(), ids)
		if err != nil {
			return nil, err
		}
		upd = tx.Rebind(upd)
		_, err = tx.ExecContext(ctx, upd, args...)
		return nil, err
	})
}
