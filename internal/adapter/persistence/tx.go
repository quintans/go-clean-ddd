package persistence

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

type key struct{}

var id key

type Conn interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type Tx interface {
	Conn
	Rollback() error
	Commit() error
}

type TxFunc func(context.Context) error

type Transactioner interface {
	Do(TxFunc) error
}

type TxHandler func(TxFunc) error

func (h TxHandler) Do(fn TxFunc) error {
	return h(fn)
}

func setConn(ctx context.Context, tx Conn) context.Context {
	return context.WithValue(ctx, id, tx)
}

func GetConn(ctx context.Context) Conn {
	c, _ := ctx.Value(id).(Conn)
	return c
}

func NewTxManager(db *sql.DB) *TxManager {
	tm := &TxManager{
		db: db,
	}
	return tm
}

type TxManager struct {
	db *sql.DB
}

func (tm *TxManager) Current(ctx context.Context) Transactioner {
	conn := GetConn(ctx)
	if conn == nil {
		return tm.makeTxHandler(ctx)
	}

	return tm.makeNoTxHandler(ctx)
}

func (tm *TxManager) New(ctx context.Context) Transactioner {
	return tm.makeTxHandler(ctx)
}

func (tm *TxManager) None(ctx context.Context) Transactioner {
	return tm.makeNoTxHandler(ctx)
}

func (tm *TxManager) makeNoTxHandler(ctx context.Context) Transactioner {
	h := func(fn TxFunc) error {
		return fn(ctx)
	}
	return TxHandler(h)
}

func (tm *TxManager) makeTxHandler(ctx context.Context) Transactioner {

	h := func(fn TxFunc) error {
		// Begin Transaction
		tx, err := tm.db.Begin()
		if err != nil {
			return err
		}
		txCtx := setConn(ctx, tx)
		defer func() {
			err := recover()
			if err != nil {
				tx.Rollback()
				panic(err) // up you go
			}
		}()

		err = fn(txCtx)

		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}

		return err
	}
	return TxHandler(h)
}

func Pool(driverName string, dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to open connection")
	}
	// wake up the database pool
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to ping database")
	}
	return db, nil
}
