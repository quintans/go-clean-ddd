package fakepub

import (
	"context"

	"github.com/quintans/go-clean-ddd/lib/fakemq"
)

type FakePublisher struct {
	mq *fakemq.FakeMQ
}

func NewFakePublisher(mq *fakemq.FakeMQ) *FakePublisher {
	return &FakePublisher{
		mq: mq,
	}
}

func (f *FakePublisher) Publish(_ context.Context, event fakemq.Event) error {
	f.mq.Publish(event)
	return nil
}
