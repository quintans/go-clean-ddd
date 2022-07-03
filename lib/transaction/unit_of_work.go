package transaction

import (
	"context"
	"fmt"

	"github.com/quintans/faults"
)

type UoWFunc func(context.Context) error

type UnitOfWorkManager interface {
	Current(context.Context, UoWFunc) error
	New(context.Context, UoWFunc) error
}

// UnitOfWork defines a database transactional boundary.
// This can be use to avoid using domain events bus for database transactional boundaries
type UnitOfWork[T Tx] struct {
	txFactory func(context.Context) (Tx, error)
}

func NewUnitOfWorkManager[T Tx](txFactory func(context.Context) (Tx, error)) UnitOfWork[T] {
	return UnitOfWork[T]{
		txFactory: txFactory,
	}
}
func (tm UnitOfWork[T]) Current(ctx context.Context, fn UoWFunc) error {
	t := ctx.Value(txID)
	if t == nil {
		return tm.makeTxHandler(ctx, fn)
	}

	err := fn(ctx)
	return faults.Wrap(err)
}

func (tm UnitOfWork[T]) New(ctx context.Context, fn UoWFunc) error {
	return tm.makeTxHandler(ctx, fn)
}

func (tm UnitOfWork[T]) makeTxHandler(ctx context.Context, fn UoWFunc) error {
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

	err = fn(c)
	if err != nil {
		return faults.Wrap(err)
	}

	_ = tx.Commit()
	return nil
}
