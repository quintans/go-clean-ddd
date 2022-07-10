package infra

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/quintans/go-clean-ddd/fake"
	"github.com/quintans/go-clean-ddd/internal/app/command"
	"github.com/quintans/go-clean-ddd/internal/app/event"
	"github.com/quintans/go-clean-ddd/internal/app/query"
	"github.com/quintans/go-clean-ddd/internal/domain/registration"
	"github.com/quintans/go-clean-ddd/internal/infra/controller/fakesub"
	"github.com/quintans/go-clean-ddd/internal/infra/controller/web"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/fakeemail"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/fakepub"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres"
	"github.com/quintans/go-clean-ddd/lib/eventbus"
	"github.com/quintans/go-clean-ddd/lib/outbox"
	"github.com/quintans/go-clean-ddd/lib/transaction"
	"github.com/quintans/toolkit/latch"
)

func Start(ctx context.Context, lock *latch.CountDownLatch, cfg Config) {
	db := NewDB(cfg.DbConfig)
	lock.Add(1)
	go func() {
		<-ctx.Done()
		db.Close()
		lock.Done()
	}()

	bus := eventbus.New()
	trans := transaction.New[*sqlx.Tx](
		bus,
		func(ctx context.Context) (transaction.Tx, error) {
			return db.Beginx()
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

	bus.Subscribe(registration.EventEmailVerified, event.NewEmailVerifiedHandler(customerWrite, customerRead).Handle)

	registrationController := web.NewRegistrationController(createRegistration, confirmRegistration)
	StartWebServer(ctx, lock, cfg.WebConfig, customerController, registrationController)

	emailClient := fake.NewEmailClient()
	emailGateway := fakeemail.NewClient(emailClient)
	sendEmail := command.NewSendEmail("http://localhost:"+cfg.Port+"/registrations/", emailGateway)
	registrationHandler := fakesub.NewRegistrationController(sendEmail)
	mq := StartMQ(ctx, lock, registrationHandler)

	pub := fakepub.NewFakePublisher(mq)
	outboxMan := outbox.New(trans, 5, pub)
	bus.Subscribe(registration.EventRegistrationCreated, func(ctx context.Context, de eventbus.DomainEvent) error {
		// DEMO: transform the incoming domain event into an integration event if there is a need to.
		// In this case there is no need to.
		return outboxMan.Create(ctx, de)
	})
	outboxMan.Start(ctx, lock, cfg.OutboxHeartbeat)
}
