package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrTxNotFound = errors.New("no transaction found in context")

type txKey struct{}

var txID txKey

type Tx interface {
	Rollback() error
	Commit() error
}

type TxFunc[T Tx] func(context.Context, T) error

func setTxToContext(ctx context.Context, tx Tx) context.Context {
	return context.WithValue(ctx, txID, tx)
}

func getTxFromContext(ctx context.Context) interface{} {
	return ctx.Value(txID)
}

type Transactioner[T Tx] interface {
	Current(ctx context.Context, fn TxFunc[T]) error
	New(ctx context.Context, fn TxFunc[T]) error
	DB() *sql.DB
}

type Transaction[T Tx] struct {
	db        *sql.DB
	txFactory func(context.Context, *sql.DB) (Tx, error)
}

type TxOption[T Tx] func(*Transaction[T])

func WithTxFactory[T Tx](fn func(context.Context, *sql.DB) (Tx, error)) TxOption[T] {
	return func(tx *Transaction[T]) {
		tx.txFactory = fn
	}
}

func New[T Tx](db *sql.DB, options ...TxOption[T]) *Transaction[T] {
	tm := &Transaction[T]{
		db: db,
		txFactory: func(ctx context.Context, d *sql.DB) (Tx, error) {
			return d.BeginTx(ctx, nil)
		},
	}
	for _, o := range options {
		o(tm)
	}
	return tm
}

func (tm Transaction[T]) DB() *sql.DB {
	return tm.db
}

func (tm *Transaction[T]) Current(ctx context.Context, fn TxFunc[T]) error {
	t := getTxFromContext(ctx)
	if t == nil {
		return tm.makeTxHandler(ctx, fn)
	}

	tx, ok := t.(T)
	if !ok {
		return ErrTxNotFound
	}

	return fn(ctx, tx)
}

func (tm *Transaction[T]) New(ctx context.Context, fn TxFunc[T]) error {
	return tm.makeTxHandler(ctx, fn)
}

func (tm *Transaction[T]) makeTxHandler(ctx context.Context, fn TxFunc[T]) error {
	// Begin Transaction
	t, err := tm.txFactory(ctx, tm.db)
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

	err = fn(c, tx)

	if err == nil {
		_ = tx.Commit()
	} else {
		_ = tx.Rollback()
	}

	return err
}
