package infra

import (
	"context"
	"time"

	"github.com/quintans/go-clean-ddd/fake"
	"github.com/quintans/go-clean-ddd/internal/app/command"
	"github.com/quintans/go-clean-ddd/internal/app/query"
	"github.com/quintans/go-clean-ddd/internal/domain/registration"
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

func Start(ctx context.Context, lock *latch.CountDownLatch, cfg Config) {
	db := NewDB(cfg.DbConfig)
	defer db.Close()

	bus := event.NewEventBus()
	trans := transaction.New[*ent.Tx](
		bus,
		func(ctx context.Context) (transaction.Tx, error) {
			return db.Tx(ctx)
		},
	)
	customerWrite := postgres.NewCustomerRepository(trans)
	customerRead := postgres.NewCustomerViewRepository(db)

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

	registrationController := web.NewRegistrationController(createRegistration, confirmRegistration)
	StartWebServer(ctx, lock, cfg.WebConfig, customerController, registrationController)

	emailClient := fake.NewEmailClient()
	emailGateway := fakeemail.NewClient(emailClient)
	sendEmail := command.NewSendEmail("http://localhost:"+cfg.Port+"/registrations/", emailGateway)
	registrationHandler := fakesub.NewRegistrationController(sendEmail)
	mq := StartMQ(ctx, lock, registrationHandler)

	pub := fakepub.NewFakePublisher(mq)
	outboxMan := outbox.New(trans, 5, pub)
	outboxMan.Start(ctx, lock, 5*time.Second)
	bus.AddHandlerF(registration.EventRegistrationCreated, func(ctx context.Context, de event.DomainEvent) error {
		// DEMO: transform the incoming domain event into an integration event if there is a need to.
		// In this case there is no need
		return outboxMan.Create(ctx, de)
	})
}
