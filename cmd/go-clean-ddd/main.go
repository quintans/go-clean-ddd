package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/quintans/go-clean-ddd/internal/domain/usecase/command"
	ucEvent "github.com/quintans/go-clean-ddd/internal/domain/usecase/event"
	"github.com/quintans/go-clean-ddd/internal/domain/usecase/query"
	"github.com/quintans/go-clean-ddd/internal/infra"
	"github.com/quintans/go-clean-ddd/internal/infra/controller/outbox"
	"github.com/quintans/go-clean-ddd/internal/infra/controller/web"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres"
	"github.com/quintans/go-clean-ddd/lib/event"
	"github.com/quintans/go-clean-ddd/lib/transaction"
)

const dbDriver = "postgres"

func main() {
	var wg sync.WaitGroup

	db, err := infra.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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
	customerController := web.NewCustomerController(
		updateCustomer,
		allCustomers,
	)

	registrationWrite := postgres.NewRegistrationRepository(trans)
	createRegistration := command.NewCreateRegistration(registrationWrite, customerRead)

	emailVerifiedHandler := ucEvent.NewEmailVerifiedHandler(customerWrite, customerRead)
	eb.AddHandler(emailVerifiedHandler)

	outboxRepository := postgres.NewOutboxRepository(trans, 5)
	outboxUC := command.NewFlushOutbox(outboxRepository)
	outboxController := outbox.NewOutboxController(outboxUC)

	ctx, cancel := context.WithCancel(context.Background())
	infra.StartOutboxRelay(ctx, 5*time.Second, outboxController)

	registrationController := web.NewRegistrationController(createRegistration)
	infra.StartWebServer(ctx, &wg, customerController, registrationController)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit
	cancel()
	latch.WaitWithTimeout(3 * time.Second)
}
