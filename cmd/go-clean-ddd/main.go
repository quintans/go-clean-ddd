package main

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/quintans/go-clean-ddd/internal/controller/web"
	"github.com/quintans/go-clean-ddd/internal/gateway/postgres"
	pg "github.com/quintans/go-clean-ddd/internal/infra/postgres"
	iWeb "github.com/quintans/go-clean-ddd/internal/infra/web"
	"github.com/quintans/go-clean-ddd/internal/usecase/command"
	ucEvent "github.com/quintans/go-clean-ddd/internal/usecase/event"
	"github.com/quintans/go-clean-ddd/internal/usecase/query"
	"github.com/quintans/go-clean-ddd/lib/event"
	"github.com/quintans/go-clean-ddd/lib/transaction"
)

const dbDriver = "postgres"

func main() {
	db, err := pg.New()
	if err != nil {
		log.Fatal(err)
	}

	eb := event.NewEventBus()
	dbx := sqlx.NewDb(db, dbDriver)
	trans := transaction.New[*sqlx.Tx](
		func(ctx context.Context) (transaction.Tx, error) {
			return dbx.BeginTxx(ctx, nil)
		},
		eb,
	)
	customerWrite := postgres.NewCustomerRepository(trans)
	customerRead := postgres.NewCustomerViewRepository(dbx)

	updateCustomer := command.NewUpdateCustomer(customerWrite, customerRead)
	allCustomers := query.NewAllCustomers(customerRead)
	c := web.NewCustomerController(
		updateCustomer,
		allCustomers,
	)

	registrationWrite := postgres.NewRegistrationRepository(trans)
	createRegistration := command.NewCreateRegistration(registrationWrite, customerRead)

	emailVerifiedHandler := ucEvent.NewEmailVerifiedHandler(customerWrite, customerRead)
	eb.AddHandler(emailVerifiedHandler)

	r := web.NewRegistrationController(createRegistration)

	if err := iWeb.StartWebServer(c, r); err != nil {
		log.Fatal(err)
	}
}
