package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/quintans/go-clean-ddd/internal/domain/usecase/command"
	"github.com/quintans/go-clean-ddd/internal/domain/usecase/query"
	"github.com/quintans/go-clean-ddd/internal/infra"
	"github.com/quintans/go-clean-ddd/internal/infra/controller/scheduler"
	"github.com/quintans/go-clean-ddd/internal/infra/controller/web"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres/ent"
	"github.com/quintans/go-clean-ddd/lib/event"
	"github.com/quintans/go-clean-ddd/lib/transaction"
	"github.com/quintans/toolkit/latch"
)

func main() {
	lock := latch.NewCountDownLatch()

	cfg := infra.LoadEnvVars()

	client := infra.NewDB(cfg.DbConfig)
	defer client.Close()

	eb := event.NewEventBus()
	trans := transaction.New[*ent.Tx](
		eb,
		func(ctx context.Context) (transaction.Tx, error) {
			return client.Tx(ctx)
		},
	)
	customerWrite := postgres.NewCustomerRepository(trans)
	customerRead := postgres.NewCustomerViewRepository(client)

	updateCustomer := command.NewUpdateCustomer(customerWrite, customerRead)
	allCustomers := query.NewAllCustomers(customerRead)
	customerController := web.NewCustomerController(
		updateCustomer,
		allCustomers,
	)

	registrationWrite := postgres.NewRegistrationRepository(trans)
	createRegistration := command.NewCreateRegistration(registrationWrite, customerRead)
	confirmRegistration := command.NewConfirmRegistration(registrationWrite)

	eb.AddHandler(command.NewRegistrationHandler(cfg.Port))
	eb.AddHandler(command.NewEmailVerifiedHandler(customerWrite, customerRead))

	outboxRepository := postgres.NewOutboxRepository(trans, 5)
	outboxUC := command.NewFlushOutbox(outboxRepository)

	ctx, cancel := context.WithCancel(context.Background())
	scheduler.StartOutboxScheduler(ctx, lock, 5*time.Second, outboxUC)

	registrationController := web.NewRegistrationController(createRegistration, confirmRegistration)
	infra.StartWebServer(ctx, lock, cfg.WebConfig, customerController, registrationController)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit
	cancel()
	lock.WaitWithTimeout(3 * time.Second)
}
