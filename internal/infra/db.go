package infra

import (
	"fmt"
	"log"
	"strings"

	"github.com/quintans/faults"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewDB creates a new postgres database connection.
// It should receive database connection configuration but for the demo purposes we will ignore it
func NewDB(cfg DbConfig) *sqlx.DB {
	options := []string{"sslmode=disable"}
	addr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName, strings.Join(options, "&"))

	db, err := sqlx.Connect("postgres", addr)
	if err != nil {
		log.Fatalf("Unable to open connection: %s", err)
	}

	err = migration(addr)
	if err != nil {
		log.Fatalf("Unable to migrate database: %s", err)
	}

	return db
}

func migration(addr string) error {
	p := &postgres.Postgres{}
	d, err := p.Open(addr)
	if err != nil {
		return faults.Wrap(err)
	}
	defer func() {
		if err := d.Close(); err != nil {
			log.Printf("[ERROR] %+v", err)
		}
	}()
	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", d)
	if err != nil {
		return faults.Wrap(err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return faults.Wrapf(err, "failed to migrate database")
	}

	return nil
}
