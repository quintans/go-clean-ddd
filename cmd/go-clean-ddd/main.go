package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/quintans/go-clean-ddd/fake"
	"github.com/quintans/go-clean-ddd/internal/app/command"
	"github.com/quintans/go-clean-ddd/internal/app/query"
	"github.com/quintans/go-clean-ddd/internal/domain/registration"
	"github.com/quintans/go-clean-ddd/internal/infra"
	"github.com/quintans/go-clean-ddd/internal/infra/controller/fakesub"
	"github.com/quintans/go-clean-ddd/internal/infra/controller/web"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/fakeemail"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/fakepub"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres/ent"
	"github.com/quintans/go-clean-ddd/lib/event"
	"github.com/quintans/go-clean-ddd/lib/outbox"
	"github.com/quintans/go-clean-ddd/lib/transaction"
	"github.com/quintans/toolkit/latch"
)

func main() {
	lock := latch.NewCountDownLatch()

	cfg := infra.LoadEnvVars()

	client := infra.NewDB(cfg.DbConfig)
	defer client.Close()

	bus := event.NewEventBus()
	trans := transaction.New[*ent.Tx](
		bus,
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

	bus.AddHandler(registration.EventEmailVerified, command.NewEmailVerifiedHandler(customerWrite, customerRead))

	ctx, cancel := context.WithCancel(context.Background())
	mq := fake.NewMQ()
	pub := fakepub.NewFakePublisher(mq)
	outboxMan := outbox.New(trans, 5, pub)
	outboxMan.Start(ctx, lock, 5*time.Second)
	bus.AddHandler(registration.EventRegistrationCreated, outboxMan)

	registrationController := web.NewRegistrationController(createRegistration, confirmRegistration)
	infra.StartWebServer(ctx, lock, cfg.WebConfig, customerController, registrationController)

	emailClient := fake.NewEmailClient()
	emailGateway := fakeemail.NewClient(emailClient)
	sendEmail := command.NewSendEmail("http://localhost:"+cfg.Port+"/registrations/", emailGateway)
	registrationHandler := fakesub.NewRegistrationController(sendEmail)
	infra.StartMQ(ctx, lock, registrationHandler)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit
	cancel()
	lock.WaitWithTimeout(3 * time.Second)
}
