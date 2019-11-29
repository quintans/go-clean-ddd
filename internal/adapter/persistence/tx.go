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

func (tm *TxManager) Transaction(ctx context.Context, handler func(context.Context) error) error {
	return tm.transaction(ctx, true, handler)
}

func (tm *TxManager) NewTransaction(ctx context.Context, handler func(context.Context) error) error {
	return tm.transaction(ctx, false, handler)
}

func (tm *TxManager) transaction(ctx context.Context, useCurrent bool, handler func(context.Context) error) error {
	var isNew bool

	conn := GetConn(ctx)
	tx, ok := conn.(Tx)

	var txCtx context.Context
	if !ok || !useCurrent {
		// No ongoing transaction
		// Begin Transaction
		var err error
		tx, err = tm.db.Begin()
		if err != nil {
			return err
		}
		txCtx = setConn(ctx, tx)
		isNew = true
	}

	defer func() {
		err := recover()
		if err != nil {
			tx.Rollback()
			panic(err) // up you go
		}
	}()

	err := handler(txCtx)

	if isNew {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}
	return err
}

func (tm *TxManager) Transactionless(ctx context.Context, handler func(context.Context) error) error {
	var conn = GetConn(ctx)

	var txCtx context.Context
	if conn == nil {
		// No ongoing connection
		// Use the connection pool
		txCtx = setConn(ctx, tm.db)
	}

	return handler(txCtx)
}

func (tm *TxManager) Pool() *sql.DB {
	return tm.db
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
