package postgres

import (
	"context"
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

func (r OutboxRepository) Consume(ctx context.Context, policy func([]entity.Outbox) error) error {
	return r.trans.Current(ctx, func(ctx context.Context, tx *sqlx.Tx) ([]event.DomainEvent, error) {
		var entities []entity.Outbox
		now := time.Now().UTC()
		until := now.Add(10 * time.Second)
		outbox := []Outbox{}
		err := tx.SelectContext(
			ctx, &outbox,
			`UPDATE outbox SET consumed=TRUE WHERE consumed=FALSE 
			ORDER BY id	ASC LIMIT $3
			FOR UPDATE SKIP LOCKED
			RETURNING id, kind, payload`,
			until, now, r.batchSize,
		)
		if err != nil {
			return nil, errorMap(err)
		}
		for _, o := range outbox {
			entities = append(entities, entity.RestoreOutbox(o.Id, o.Kind, o.Payload))
		}
		return nil, policy(entities)
	})
}
