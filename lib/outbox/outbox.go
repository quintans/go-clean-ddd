package outbox

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres/ent"
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
	trans     *transaction.Transaction[*ent.Tx]
	batchSize uint
	publisher Publisher
}

func New(
	trans *transaction.Transaction[*ent.Tx],
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
	return r.trans.Current(ctx, func(ctx context.Context, tx *ent.Tx) (transaction.EventPopper, error) {
		fmt.Println("===> inserting into outbox", e)
		_, err := tx.ExecContext(
			ctx,
			"INSERT INTO outbox(kind, payload, consumed) VALUES($1, $2, FALSE)",
			e.Kind, b,
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
		done, err = f.consumeBatch(ctx, func(events []*outbox) error {
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

func (r OutboxManager) consumeBatch(ctx context.Context, fn func([]*outbox) error) (done bool, err error) {
	err = r.trans.New(ctx, func(ctx context.Context, tx *ent.Tx) (transaction.EventPopper, error) {
		locked, err := getAdvisoryLock(ctx, tx)
		if err != nil {
			return nil, faults.Errorf("getting advisory lock: %w", err)
		}
		if !locked {
			return nil, errLocking
		}

		rows, err := tx.QueryContext(
			ctx,
			`SELECT id, kind, payload FROM outbox WHERE kind = 'XXX' ORDER BY id ASC LIMIT $1`,
			r.batchSize,
		)
		if err != nil {
			fmt.Println("===> no rows?", err)
			return nil, faults.Errorf("fetching events batch: %w", err)
		}
		fmt.Println("===> fetching rows")
		x := 0
		for rows.Next() {
			fmt.Println("===> select outbox")
			if err := rows.Scan(&x); err != nil {
				return nil, faults.Wrap(err)
			}
			fmt.Println("===> X", x)
		}

		var entities []*outbox
		var id int64
		var kind string
		var payload []byte

		var in strings.Builder
		var args []any
		err = forEachRow(rows, func() {
			entities = append(entities, restoreOutbox(kind, payload))
			fmt.Println("===>", id, ", ", kind, ", ", payload)

			args = append(args, id)
			if len(args) > 1 {
				in.WriteString(", ")
			}
			in.WriteString("$" + strconv.Itoa(len(args)))

			// reset
			id = 0
			kind = ""
			payload = nil
		}, &id, &kind, &payload)
		if err != nil {
			return nil, faults.Errorf("iterating over events rows: %w", err)
		}

		if len(args) == 0 {
			done = true
			return nil, nil
		}

		s := fmt.Sprintf("UPDATE outbox SET consumed=TRUE WHERE id IN(%s)", in.String())
		fmt.Println("===> sql:", s)
		_, err = tx.ExecContext(ctx, s, args...)
		if err != nil {
			return nil, faults.Errorf("consuming events batch: %w", err)
		}

		return nil, fn(entities)
	})
	return done, faults.Wrap(err)
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

func forEachRow(rows *sql.Rows, fn func(), v ...any) error {
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(v...); err != nil {
			return faults.Wrap(err)
		}
		if fn != nil {
			fn()
		} else {
			break
		}
	}

	return faults.Wrap(rows.Err())
}
