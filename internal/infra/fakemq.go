package infra

import (
	"context"

	"github.com/quintans/go-clean-ddd/fake"
	"github.com/quintans/go-clean-ddd/internal/domain/registration"
	"github.com/quintans/go-clean-ddd/internal/infra/controller/fakesub"
	"github.com/quintans/toolkit/latch"
)

func StartMQ(
	ctx context.Context,
	lock *latch.CountDownLatch,
	regSub fakesub.RegistrationController,
) {
	mq := fake.NewMQ()

	mq.Subscribe(registration.EventRegistrationCreated, regSub)

	go func() {
		lock.Add(1)
		defer lock.Done()
		<-ctx.Done()
		mq.Close()
	}()
}
