package outbox

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres/ent"
	"github.com/quintans/go-clean-ddd/lib/event"
	"github.com/quintans/go-clean-ddd/lib/transaction"
	"github.com/quintans/toolkit/latch"
)

var errLocking = errors.New("failed to get advisory lock")

type Publisher interface {
	Publish(ctx context.Context, event Event) error
}

type Event struct {
	Kind    string
	Payload []byte
}

type outbox struct {
	Kind    string
	Payload []byte
}

func restoreOutbox(kind string, payload []byte) *outbox {
	return &outbox{
		Kind:    kind,
		Payload: payload,
	}
}

type OutboxManager struct {
	trans     transaction.Transactioner[*ent.Tx]
	batchSize uint
	publisher Publisher
}

func New(
	trans transaction.Transactioner[*ent.Tx],
	batchSize uint,
	publisher Publisher,
) OutboxManager {
	return OutboxManager{
		trans:     trans,
		batchSize: batchSize,
		publisher: publisher,
	}
}

func (r OutboxManager) Handle(ctx context.Context, e event.DomainEvent) error {
	b, err := json.Marshal(e)
	if err != nil {
		return faults.Wrap(err)
	}
	return r.trans.Current(ctx, func(ctx context.Context, tx *ent.Tx) (transaction.EventPopper, error) {
		_, err := tx.ExecContext(
			ctx,
			"INSERT INTO outbox(kind, payload, consumed) VALUES($1, $2, FALSE)",
			e.Kind(), b,
		)
		return nil, faults.Wrap(err)
	})
}

func (c OutboxManager) Start(ctx context.Context, lock *latch.CountDownLatch, heartbeat time.Duration) {
	lock.Add(1)
	go func() {
		defer lock.Done()
		ticker := time.NewTicker(heartbeat)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := c.consumeAll(ctx)
				if err != nil {
					log.Printf("ERROR: failed to execute flush outbox: %s\n", err)
				}
			}
		}
	}()
}

func (f OutboxManager) consumeAll(ctx context.Context) error {
	for {
		err := f.consumeBatch(ctx, func(events []*outbox) error {
			for _, e := range events {
				err := f.publisher.Publish(ctx, Event{
					Kind:    e.Kind,
					Payload: e.Payload,
				})
				if err != nil {
					return faults.Wrap(err)
				}
			}
			return nil
		})

		if err != nil {
			if errors.Is(err, errLocking) {
				return nil
			}
			return faults.Wrap(err)
		}
	}
}

func (r OutboxManager) consumeBatch(ctx context.Context, fn func([]*outbox) error) error {
	for {
		err := r.trans.Current(ctx, func(ctx context.Context, tx *ent.Tx) (transaction.EventPopper, error) {
			ok, err := getAdvisoryLock(ctx, tx)
			if err != nil {
				return nil, faults.Errorf("getting advisory lock: %w", err)
			}
			if !ok {
				return nil, errLocking
			}

			rows, err := tx.QueryContext(
				ctx,
				`UPDATE outbox SET consumed=TRUE WHERE consumed=FALSE 
				ORDER BY id	ASC LIMIT $1
				RETURNING kind, payload`,
				r.batchSize,
			)
			if err != nil {
				return nil, faults.Errorf("fetching events batch: %w", err)
			}
			var entities []*outbox
			var kind string
			var payload []byte
			forEachRow(rows, func() {
				entities = append(entities, restoreOutbox(kind, payload))
				payload = nil
			}, &kind, &payload)

			return nil, fn(entities)
		})
		if err != nil {
			return faults.Wrap(err)
		}
	}
}

func getAdvisoryLock(ctx context.Context, tx *ent.Tx) (bool, error) {
	// 123 is the lock key
	rows, err := tx.QueryContext(ctx, "SELECT pg_try_advisory_xact_lock(123)")
	if err != nil {
		return false, faults.Wrap(err)
	}

	var ok bool
	err = forEachRow(rows, nil, &ok)
	if err != nil {
		return false, faults.Wrap(err)
	}
	return ok, nil
}

func forEachRow(rows *sql.Rows, fn func(), v ...interface{}) error {
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(v...)
		if err != nil {
			return faults.Wrap(err)
		}
		if fn != nil {
			fn()
		} else {
			break
		}
	}

	if err := rows.Err(); err != nil {
		if err != nil {
			return faults.Wrap(err)
		}
	}
	return nil
}
