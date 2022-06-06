package infra

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	// initializes the postgres driver
	_ "github.com/lib/pq"
)

// NewDB creates a new postgres database connection.
// It should receive database connection configuration but for the demo purposes we will ignore it
func NewDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "dbname=postgres user=postgres password=secret port=5432 sslmode=disable")
	if err != nil {
		return nil, errors.Wrap(err, "Unable to open connection")
	}
	// wake up the database pool
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to ping database")
	}

	err = migration()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func migration() error {
	options := []string{"sslmode=disable"}
	addr := fmt.Sprintf("postgres://postgres:secret@localhost:5432/postgres?%s", strings.Join(options, "&"))

	p := &postgres.Postgres{}
	d, err := p.Open(addr)
	if err != nil {
		return err
	}
	defer func() {
		if err := d.Close(); err != nil {
			log.Error(err)
		}
	}()
	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", d)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "failed to migrate database")
	}

	return nil
}
