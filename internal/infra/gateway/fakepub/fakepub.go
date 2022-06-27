package fakepub

import (
	"context"

	"github.com/quintans/go-clean-ddd/fake"
)

type FakePublisher struct {
	mq *fake.FakeMQ
}

func NewFakePublisher(mq *fake.FakeMQ) *FakePublisher {
	return &FakePublisher{
		mq: mq,
	}
}

func (f *FakePublisher) Publish(_ context.Context, event fake.MQEvent) error {
	f.mq.Publish(event)
	return nil
}
