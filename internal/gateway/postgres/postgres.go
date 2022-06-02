package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/quintans/go-clean-ddd/internal/domain"
	"github.com/quintans/go-clean-ddd/internal/usecase"
	"github.com/quintans/go-clean-ddd/lib/transaction"
)

const driver = "postgres"

type Customer struct {
	Id        string
	Version   int
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
}

type CustomerRepository struct {
	trans transaction.Transactioner[*sqlx.Tx]
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	// trans is instantiated here because it is only used inside this repository,
	// but it could be used in the use case layer if would like to have a database transaction span several repositories calls
	trans := transaction.New(db, transaction.WithTxFactory[*sqlx.Tx](func(ctx context.Context, db *sql.DB) (transaction.Tx, error) {
		return sqlx.NewDb(db, driver).BeginTxx(ctx, nil)
	}))

	return CustomerRepository{
		trans: trans,
	}
}

func (r CustomerRepository) Save(ctx context.Context, c domain.Customer) error {
	err := r.trans.Current(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		_, err := tx.ExecContext(
			ctx,
			"INSERT INTO customer(id, version, first_name, last_name, email) SET VALUES($1, $2, $3, $4, $5)",
			c.ID(), 0, c.FullName().FirstName(), c.FullName().LastName(), c.Email(),
		)
		return err
	})

	return errorMap(err)
}

func (r CustomerRepository) Apply(ctx context.Context, id domain.CustomerID, apply usecase.Apply) error {
	c, err := r.getByID(ctx, id)
	if err != nil {
		return err
	}
	err = r.trans.Current(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		customer, err := toDomainCustomer(c)
		if err != nil {
			return err
		}

		customer, err = apply(ctx, customer)
		if err != nil {
			return err
		}

		// optimistic locking is used
		_, err = tx.ExecContext(
			ctx,
			"UPDATE customer SET first_name=$1, last_name=$2, email=$3, version=version+1 WHERE id=$4 AND version=$5",
			customer.FullName().FirstName(), customer.FullName().LastName(), customer.Email(), customer.ID(), c.Version,
		)
		return err
	})

	return errorMap(err)
}

func (r CustomerRepository) getByID(ctx context.Context, id domain.CustomerID) (Customer, error) {
	db := r.getConn()
	customer := Customer{}
	err := db.Get(&customer, "SELECT * FROM customer WHERE id=$1", id.String())
	if err != nil {
		return Customer{}, errorMap(err)
	}
	return customer, nil
}

func (r CustomerRepository) GetByID(ctx context.Context, id domain.CustomerID) (domain.Customer, error) {
	c, err := r.getByID(ctx, id)
	if err != nil {
		return domain.Customer{}, err
	}
	return toDomainCustomer(c)
}

func (r CustomerRepository) GetAll(ctx context.Context) ([]domain.Customer, error) {
	db := r.getConn()
	customers := []Customer{}
	err := db.Select(&customers, "SELECT * FROM customer")
	if err != nil {
		return nil, errorMap(err)
	}
	return toDomainCustomers(customers)
}

func (r CustomerRepository) GetByEmail(ctx context.Context, email domain.Email) (domain.Customer, error) {
	db := r.getConn()
	customer := Customer{}
	err := db.Get(&customer, "SELECT * FROM customer WHERE email=$1", email.Email())
	if err != nil {
		return domain.Customer{}, errorMap(err)
	}
	return toDomainCustomer(customer)
}

func (r CustomerRepository) getConn() *sqlx.DB {
	return sqlx.NewDb(r.trans.DB(), driver)
}

func errorMap(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return usecase.ErrReadModelNotFound
	}
	return err
}

func toDomainCustomers(cs []Customer) ([]domain.Customer, error) {
	dcs := make([]domain.Customer, len(cs))
	for k, v := range cs {
		dc, err := toDomainCustomer(v)
		if err != nil {
			return nil, err
		}
		dcs[k] = dc
	}
	return dcs, nil
}

func toDomainCustomer(c Customer) (domain.Customer, error) {
	id, err := domain.ParseCustomerID(c.Id)
	if err != nil {
		return domain.Customer{}, err
	}
	fullName, err := domain.NewFullName(c.FirstName, c.LastName)
	if err != nil {
		return domain.Customer{}, err
	}
	email, err := domain.NewEmail(c.Email)
	if err != nil {
		return domain.Customer{}, err
	}
	return domain.NewCustomer(id, fullName, email)
}
