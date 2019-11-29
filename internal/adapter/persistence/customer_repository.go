package persistence

import (
	"context"

	"github.com/quintans/go-clean-ddd/internal/domain/model"
	"github.com/google/uuid"
	"github.com/quintans/goSQL/db"
	"github.com/quintans/toolkit"
)

var (
	CUSTOMER             = db.TABLE("CUSTOMER")
	CUSTOMER_C_UUID      = CUSTOMER.KEY("UUID")          // implicit map to field Id
	CUSTOMER_C_FIRSTNAME = CUSTOMER.COLUMN("FIRST_NAME") // implicit map to field Firstname
	CUSTOMER_C_LASTNAME  = CUSTOMER.COLUMN("LAST_NAME")  // implicit map to field Lastname
	CUSTOMER_C_EMAIL     = CUSTOMER.COLUMN("EMAIL")      // implicit map to field Email
)

type Customer struct {
	Uuid     string
	FullName model.FullNameVO `sql:"embeded"`
	Email    string
}

func (c *Customer) Equals(e interface{}) bool {
	if c == e {
		return true
	}

	switch t := e.(type) {
	case *Customer:
		return c.Uuid == t.Uuid
	}
	return false
}

func (c *Customer) HashCode() int {
	result := toolkit.HashType(toolkit.HASH_SEED, c)
	result = toolkit.HashString(result, c.Uuid)
	return result
}

func toCustomers(db []*Customer) []*model.Customer {
	customers := make([]*model.Customer, len(db))
	for k, v := range db {
		customers[k] = toCustomer(v)
	}
	return customers
}

func toCustomer(db *Customer) *model.Customer {
	return model.NewCustomerAll(
		uuid.MustParse(db.Uuid),
		db.FullName,
		db.Email,
	)
}

func fromCustomer(c *model.Customer) *Customer {
	return &Customer{
		Uuid:     c.Uuid().String(),
		FullName: c.Fullname(),
		Email:    c.Email(),
	}
}

type CustomerRepositoryImpl struct {
	translator db.Translator
}

func NewCustomerRepositoryImpl(translator db.Translator) CustomerRepositoryImpl {
	return CustomerRepositoryImpl{
		translator: translator,
	}
}

func (r CustomerRepositoryImpl) newDb(ctx context.Context) *db.Db {
	tx := GetConn(ctx)
	return db.NewDb(tx, r.translator)
}

func (r CustomerRepositoryImpl) FindAll(ctx context.Context) ([]*model.Customer, error) {
	sql := r.newDb(ctx)
	var customers []*Customer
	err := sql.Query(CUSTOMER).List(&customers)
	if err != nil {
		return nil, err
	}
	return toCustomers(customers), nil
}

func (r CustomerRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.Customer, error) {
	sql := r.newDb(ctx)

	customer := &Customer{}
	ok, err := sql.Query(CUSTOMER).
		All().
		Where(CUSTOMER_C_EMAIL.Matches(email)).
		SelectTo(customer)
	if err != nil {
		return nil, err
	}
	if ok {
		return toCustomer(customer), nil
	}
	return nil, nil
}

func (r CustomerRepositoryImpl) Store(ctx context.Context, c *model.Customer) error {
	sql := r.newDb(ctx)

	cdb := fromCustomer(c)
	err := sql.Create(cdb)
	if err != nil {
		return err
	}

	return nil
}
