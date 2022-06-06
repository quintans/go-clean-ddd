package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/quintans/go-clean-ddd/internal/domain/usecase/command"
	ucEvent "github.com/quintans/go-clean-ddd/internal/domain/usecase/event"
	"github.com/quintans/go-clean-ddd/internal/domain/usecase/query"
	"github.com/quintans/go-clean-ddd/internal/infra"
	"github.com/quintans/go-clean-ddd/internal/infra/controller/web"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres"
	"github.com/quintans/go-clean-ddd/lib/event"
	"github.com/quintans/go-clean-ddd/lib/transaction"
	"github.com/quintans/toolkit/latch"
)

const dbDriver = "postgres"

func main() {
	lock := latch.NewCountDownLatch()

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
	confirmRegistration := command.NewConfirmRegistration(registrationWrite)

	emailVerifiedHandler := ucEvent.NewEmailVerifiedHandler(customerWrite, customerRead)
	eb.AddHandler(emailVerifiedHandler)

	outboxRepository := postgres.NewOutboxRepository(trans, 5)
	outboxUC := command.NewFlushOutbox(outboxRepository)

	ctx, cancel := context.WithCancel(context.Background())
	infra.StartOutboxScheduler(ctx, lock, 5*time.Second, outboxUC)

	registrationController := web.NewRegistrationController(createRegistration, confirmRegistration)
	infra.StartWebServer(ctx, lock, ":8080", customerController, registrationController)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit
	cancel()
	lock.WaitWithTimeout(3 * time.Second)
}
