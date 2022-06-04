package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/quintans/go-clean-ddd/lib/event"
)

var ErrTxNotFound = errors.New("no transaction found in context")

type txKey struct{}

var txID txKey

type Tx interface {
	Rollback() error
	Commit() error
}

type TxFunc[T Tx] func(context.Context, T) ([]event.DomainEvent, error)

func setTxToContext(ctx context.Context, tx Tx) context.Context {
	return context.WithValue(ctx, txID, tx)
}

func GetTxFromContext[T Tx](ctx context.Context) (T, error) {
	tx, ok := ctx.Value(txID).(T)
	if !ok {
		var zero T
		return zero, ErrTxNotFound
	}
	return tx, nil
}

type Transactioner[T Tx] interface {
	Current(ctx context.Context, fn TxFunc[T]) error
	New(ctx context.Context, fn TxFunc[T]) error
}

type Transaction[T Tx] struct {
	txFactory func(context.Context) (Tx, error)
	eventBus  event.EventBuser
}

func Default(db *sql.DB, eventBus event.EventBuser) *Transaction[*sql.Tx] {
	tm := &Transaction[*sql.Tx]{
		txFactory: func(ctx context.Context) (Tx, error) {
			return db.BeginTx(ctx, nil)
		},
		eventBus: eventBus,
	}
	return tm
}

func New[T Tx](txFactory func(context.Context) (Tx, error), eventBus event.EventBuser) *Transaction[T] {
	return &Transaction[T]{
		txFactory: txFactory,
		eventBus:  eventBus,
	}
}

func (tm *Transaction[T]) Current(ctx context.Context, fn TxFunc[T]) error {
	t := ctx.Value(txID)
	if t == nil {
		return tm.makeTxHandler(ctx, fn)
	}

	tx, ok := t.(T)
	if !ok {
		return ErrTxNotFound
	}

	events, err := fn(ctx, tx)
	if err != nil {
		return err
	}
	if events != nil && tm.eventBus != nil {
		return tm.eventBus.Fire(ctx, events...)
	}
	return nil
}

func (tm *Transaction[T]) New(ctx context.Context, fn TxFunc[T]) error {
	return tm.makeTxHandler(ctx, fn)
}

func (tm *Transaction[T]) makeTxHandler(ctx context.Context, fn TxFunc[T]) error {
	// Begin Transaction
	t, err := tm.txFactory(ctx)
	if err != nil {
		return err
	}
	tx, ok := t.(T)
	if !ok {
		var zero T
		return fmt.Errorf("factory transaction produced wrong type: want %T, got %T", zero, t)
	}
	c := setTxToContext(ctx, tx)
	defer func() {
		err := recover()
		if err != nil {
			_ = tx.Rollback()
			panic(err) // up you go
		}
	}()

	events, err := fn(c, tx)
	if err != nil {
		return err
	}
	if events != nil && tm.eventBus != nil {
		err := tm.eventBus.Fire(ctx, events...)
		if err != nil {
			return err
		}
	}

	if err == nil {
		_ = tx.Commit()
	} else {
		_ = tx.Rollback()
	}

	return err
}
