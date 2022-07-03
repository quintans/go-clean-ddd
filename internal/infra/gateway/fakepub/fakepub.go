package fakepub

import (
	"context"

	"github.com/quintans/go-clean-ddd/fake"
	"github.com/quintans/go-clean-ddd/lib/outbox"
)

type FakePublisher struct {
	mq *fake.FakeMQ
}

func NewFakePublisher(mq *fake.FakeMQ) *FakePublisher {
	return &FakePublisher{
		mq: mq,
	}
}

func (f *FakePublisher) Publish(_ context.Context, event outbox.Event) error {
	f.mq.Publish(fake.MQEvent{
		Kind:    event.Kind,
		Payload: event.Payload,
	})
	return nil
}
