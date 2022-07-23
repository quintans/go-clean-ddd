package outbox

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/lib/transaction"
	"github.com/quintans/toolkit/latch"
)

var errLocking = errors.New("failed to get advisory lock")

type Publisher interface {
	Publish(ctx context.Context, event Message) error
}

type Message struct {
	Kind    string
	Payload []byte
}

type Event interface {
	Kind() string
}

type outbox struct {
	Id       int64
	Kind     string
	Payload  []byte
	Consumed bool
}

type OutboxManager struct {
	trans     *transaction.Transaction[*sqlx.Tx]
	batchSize uint
	publisher Publisher
}

func New(
	trans *transaction.Transaction[*sqlx.Tx],
	batchSize uint,
	publisher Publisher,
) OutboxManager {
	return OutboxManager{
		trans:     trans,
		batchSize: batchSize,
		publisher: publisher,
	}
}

func (r OutboxManager) Create(ctx context.Context, e Event) error {
	b, err := json.Marshal(e)
	if err != nil {
		return faults.Wrap(err)
	}
	return r.trans.Current(ctx, func(ctx context.Context, tx *sqlx.Tx) (transaction.EventPopper, error) {
		_, err := tx.ExecContext(
			ctx,
			"INSERT INTO outbox(kind, payload, consumed) VALUES($1, $2, FALSE)",
			e.Kind(), b,
		)
		return nil, faults.Wrap(err)
	})
}

func (c OutboxManager) Start(ctx context.Context, lock *latch.CountDownLatch, heartbeat time.Duration) {
	fmt.Println("===> heartbeat:", heartbeat)
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
				fmt.Println("===> tick:", time.Now())
				err := c.consumeAll(ctx)
				if err != nil {
					log.Printf("ERROR: failed to execute flush outbox: %+v\n", err)
				}
			}
		}
	}()
}

func (f OutboxManager) consumeAll(ctx context.Context) error {
	done := false
	for !done {
		var err error
		done, err = f.consumeBatch(ctx, func(events []outbox) error {
			for _, e := range events {
				err := f.publisher.Publish(ctx, Message{
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
	return nil
}

func (r OutboxManager) consumeBatch(ctx context.Context, fn func([]outbox) error) (done bool, err error) {
	err = r.trans.New(ctx, func(ctx context.Context, tx *sqlx.Tx) (transaction.EventPopper, error) {
		locked, err := getAdvisoryLock(ctx, tx)
		if err != nil {
			return nil, faults.Errorf("getting advisory lock: %w", err)
		}
		if !locked {
			return nil, errLocking
		}

		// It should give an error sql.ErrNoRows
		outboxes := []outbox{}
		err = tx.SelectContext(
			ctx,
			&outboxes,
			`SELECT * FROM outbox WHERE consumed=FALSE ORDER BY id ASC LIMIT $1`,
			r.batchSize,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				done = true
				return nil, nil
			}
			return nil, faults.Errorf("fetching events batch: %w", err)
		}

		fmt.Println("===> outboxes:", outboxes)

		done = len(outboxes) < int(r.batchSize)
		if len(outboxes) == 0 {
			return nil, nil
		}

		ids := make([]int64, len(outboxes))
		for k, v := range outboxes {
			ids[k] = v.Id
		}
		query, args, err := sqlx.In("UPDATE outbox SET consumed=TRUE WHERE id IN(?)", ids)
		if err != nil {
			return nil, faults.Errorf("expanding ids to consume: %w", err)
		}
		query = tx.Rebind(query)
		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			return nil, faults.Errorf("consuming events batch: %w", err)
		}

		return nil, fn(outboxes)
	})
	return done, faults.Wrap(err)
}

func getAdvisoryLock(ctx context.Context, tx *sqlx.Tx) (bool, error) {
	ok := false
	// 123 is the lock key
	err := tx.GetContext(ctx, &ok, "SELECT pg_try_advisory_xact_lock(123)")
	if err != nil {
		return false, faults.Wrap(err)
	}

	return ok, nil
}
