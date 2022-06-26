package infra

import (
	"context"

	"github.com/quintans/go-clean-ddd/internal/domain/registration"
	"github.com/quintans/go-clean-ddd/internal/infra/controller/fakesub"
	"github.com/quintans/go-clean-ddd/lib/fakemq"
)

func StartMQ(
	ctx context.Context,
	regSub fakesub.RegistrationController,
) {
	mq := fakemq.New()

	mq.Subscribe(registration.EventRegistrationCreated, regSub)

	go func() {
		<-ctx.Done()
		mq.Close()
	}()
}
