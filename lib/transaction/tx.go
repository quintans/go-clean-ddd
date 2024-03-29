package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/lib/eventbus"
)

var ErrTxNotFound = errors.New("no transaction found in context")

type txKey struct{}

var txID txKey

type Tx interface {
	Rollback() error
	Commit() error
}

type EventPopper interface {
	PopEvents() []eventbus.DomainEvent
}

type TxFunc[T Tx] func(context.Context, T) (EventPopper, error)

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

type Transaction[T Tx] struct {
	txFactory func(context.Context) (Tx, error)
	eventBus  eventbus.Publisher
}

func New[T Tx](eventBus eventbus.Publisher, txFactory func(context.Context) (Tx, error)) *Transaction[T] {
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

	return tm.apply(ctx, tx, fn)
}

func (tm *Transaction[T]) apply(ctx context.Context, tx T, fn TxFunc[T]) error {
	popper, err := fn(ctx, tx)
	if err != nil {
		return faults.Wrap(err)
	}
	if popper != nil && tm.eventBus != nil {
		err := tm.eventBus.Publish(ctx, popper.PopEvents()...)
		if err != nil {
			return faults.Wrap(err)
		}
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
		return faults.Wrap(err)
	}
	tx, ok := t.(T)
	if !ok {
		var zero T
		return fmt.Errorf("factory transaction produced wrong type: want %T, got %T", zero, t)
	}
	c := setTxToContext(ctx, tx)
	defer func() {
		_ = tx.Rollback()
	}()

	err = tm.apply(c, tx, fn)
	if err != nil {
		return faults.Wrap(err)
	}

	_ = tx.Commit()
	return nil
}
